package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newHashContext(method, target, contentType, body string, contentLength int64) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = contentLength
	c.Request = req
	return c
}

func TestParamsHashRestoresSmallJSONBody(t *testing.T) {
	body := `{"name":"alice"}`
	c := newHashContext(http.MethodPost, "/submit?b=2&a=1", "application/json", body, int64(len(body)))

	hash := paramsHash(c)
	if hash == "" {
		t.Fatal("hash is empty")
	}

	restored, err := io.ReadAll(c.Request.Body)
	if err != nil {
		t.Fatalf("read restored body: %v", err)
	}
	if string(restored) != body {
		t.Fatalf("restored body = %q, want %q", string(restored), body)
	}
}

func TestParamsHashSkipsUnknownOrLargeBody(t *testing.T) {
	body := `{"name":"alice"}`
	c := newHashContext(http.MethodPost, "/submit?a=1", "application/json", body, -1)

	hash := paramsHash(c)
	if hash == "" {
		t.Fatal("hash is empty")
	}

	restored, err := io.ReadAll(c.Request.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if string(restored) != body {
		t.Fatalf("body should not be consumed for unknown length, got %q", string(restored))
	}

	large := newHashContext(http.MethodPost, "/submit?a=1", "application/json", body, maxNoRepeatBodyBytes+1)
	largeHash := paramsHash(large)
	if largeHash != hash {
		t.Fatalf("large and unknown body hashes should both use query params only: %q != %q", largeHash, hash)
	}
}
