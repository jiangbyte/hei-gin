package log

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

const maxLoggedBodyBytes int64 = 1 << 20

// ParseUserAgent extracts browser and OS from a User-Agent string.
// Delegates to utils.GetBrowser and utils.GetOS.
func ParseUserAgent(ua string) (browser, os string) {
	if ua == "" {
		return "-", "-"
	}
	return utils.GetBrowser(ua), utils.GetOS(ua)
}

// ExtractParamsJson extracts POST/PUT/PATCH body params as JSON,
// excluding infrastructure params (request, db, file). For GET/DELETE returns "".
func ExtractParamsJson(c *gin.Context) string {
	if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
		return ""
	}
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {
		return ""
	}
	if c.Request.ContentLength < 0 || c.Request.ContentLength > maxLoggedBodyBytes {
		return ""
	}

	bodyBytes, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, maxLoggedBodyBytes))
	if err != nil {
		return ""
	}
	// Restore body for downstream handlers (e.g. ShouldBindJSON)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var params map[string]any
	if err := json.Unmarshal(bodyBytes, &params); err != nil {
		return ""
	}

	excluded := map[string]bool{"request": true, "db": true, "file": true}
	filtered := make(map[string]any)
	for k, v := range params {
		if excluded[k] || v == nil {
			continue
		}
		filtered[k] = v
	}

	if len(filtered) == 0 {
		return ""
	}

	data, err := json.Marshal(filtered)
	if err != nil {
		return ""
	}
	return string(data)
}

// GetResultJson serializes a result value to JSON string.
// Returns "" if result is nil or serialization fails.
func GetResultJson(result any) string {
	if result == nil {
		return ""
	}
	data, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(data)
}

// GenerateLogSignature generates an SM3-based signature for log tamper-proofing.
// Uses the salt "hei-log-sign" matching the Python implementation.
func GenerateLogSignature(opData map[string]any) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(opData); err != nil {
		return ""
	}
	content := strings.TrimRight(buf.String(), "\n")
	return utils.HashWithSalt(content, "hei-log-sign")
}
