package message

import (
	"net/http/httptest"
	"testing"

	"hei-gin/sdk/auth"

	"github.com/gin-gonic/gin"
)

func newMessageContext(path string) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", path, nil)
	return c
}

func TestGetRealmIDPrefersContextValue(t *testing.T) {
	c := newMessageContext("/api/v1/c/im/message/page")
	c.Set("login_realm_id", string(auth.BusinessID))

	if got := getRealmID(c); got != string(auth.BusinessID) {
		t.Fatalf("getRealmID = %q", got)
	}
}

func TestGetRealmIDDetectsConsumerAPIPath(t *testing.T) {
	tests := map[string]string{
		"/api/v1/c/im/message/page":  string(auth.ConsumerID),
		"/api/v12/c/im/message/page": string(auth.ConsumerID),
		"/api/v1/b/im/message/page":  string(auth.BusinessID),
		"/c/im/message/page":         string(auth.BusinessID),
	}

	for path, want := range tests {
		t.Run(path, func(t *testing.T) {
			if got := getRealmID(newMessageContext(path)); got != want {
				t.Fatalf("getRealmID = %q, want %q", got, want)
			}
		})
	}
}

func TestNormalizeReceiverIDs(t *testing.T) {
	got := normalizeReceiverIDs([]string{" u1 ", "", "u2", "u1", "  "})
	want := []string{"u1", "u2"}
	if len(got) != len(want) {
		t.Fatalf("normalizeReceiverIDs len = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("normalizeReceiverIDs[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
