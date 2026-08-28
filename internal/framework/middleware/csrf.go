// csrf.go Cookie 双提交 CSRF（HEI_CSRF + X-HEI-CSRF）。
//
// Author: Charlie
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/response"
)

const (
	csrfCookieName = "HEI_CSRF"
	csrfHeaderName = "X-HEI-CSRF"
)

// CSRFProtect 当请求携带会话 Cookie 时，对非安全方法校验双提交 CSRF。
// 仅 Header/Bearer 鉴权（无会话 Cookie）时跳过。
func CSRFProtect(cfg config.AuthConfig) gin.HandlerFunc {
	cookieName := cfg.SessionCookieName
	if cookieName == "" {
		cookieName = cfg.TokenName
	}
	if cookieName == "" {
		cookieName = "Authorization"
	}
	return func(c *gin.Context) {
		if !cfg.SessionCookieEnabled {
			c.Next()
			return
		}
		method := c.Request.Method
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions || method == http.MethodTrace {
			c.Next()
			return
		}
		path := c.Request.URL.Path
		if isCSRFExempt(path) {
			c.Next()
			return
		}
		sessCookie, err := c.Cookie(cookieName)
		if err != nil || strings.TrimSpace(sessCookie) == "" {
			c.Next()
			return
		}
		csrfCookie, err := c.Cookie(csrfCookieName)
		if err != nil || csrfCookie == "" {
			response.Fail(c, http.StatusForbidden, 403, "CSRF token missing")
			c.Abort()
			return
		}
		header := c.GetHeader(csrfHeaderName)
		if header == "" || header != csrfCookie {
			response.Fail(c, http.StatusForbidden, 403, "CSRF token mismatch")
			c.Abort()
			return
		}
		c.Next()
	}
}

func isCSRFExempt(path string) bool {
	p := strings.ToLower(path)
	if strings.Contains(p, "/oauth/") && strings.Contains(p, "/callback") {
		return true
	}
	if strings.HasSuffix(p, "/health") || strings.HasSuffix(p, "/ready") {
		return true
	}
	return false
}

// IssueCSRFCookie 登录成功后下发可读 CSRF Cookie（Path=/ 便于 SPA 读取）。
func IssueCSRFCookie(c *gin.Context, cfg config.AuthConfig, maxAge int) {
	if !cfg.SessionCookieEnabled {
		return
	}
	token, err := randomCSRFToken()
	if err != nil {
		return
	}
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(cfg.SessionCookieSameSite) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: false,
		Secure:   cfg.SessionCookieSecure,
		SameSite: sameSite,
	})
}

// ClearCSRFCookie 登出时清除 CSRF Cookie。
func ClearCSRFCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     csrfCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
	})
}

func randomCSRFToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
