// Package middleware æä¾› Gin é€šç”¨ä¸­é—´ä»¶ï¼šæ¢å¤ã€è¿½è¸ªã€CORSã€é‰´æƒä¸Žé”™è¯¯æ˜ å°„ã€‚
package middleware

import (
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"hei-gin/internal/framework/core/config"
	contextx "hei-gin/internal/framework/core/context"
	apperr "hei-gin/internal/framework/core/errors"
	"hei-gin/internal/framework/core/logger"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
)

// Recovery æ•èŽ· panic å¹¶è¿”å›ž 500 ä¿¡å°ã€‚
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.L.Error("panic", zap.Any("recover", r))
				response.Fail(c, http.StatusInternalServerError, 500, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// Trace æ³¨å…¥ X-Request-Id ä¸Žå®¢æˆ·ç«¯ IP åˆ°è¯·æ±‚ä¸Šä¸‹æ–‡ã€‚
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-Id")
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Writer.Header().Set("X-Request-Id", rid)
		ctx := contextx.WithRequestID(c.Request.Context(), rid)
		ctx = contextx.WithClientIP(ctx, c.ClientIP())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// AccessLog è®°å½•è®¿é—®æ—¥å¿—ï¼ˆæ–¹æ³•ã€è·¯å¾„ã€çŠ¶æ€ã€è€—æ—¶ã€request_idï¼‰ã€‚
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.L.Info("access",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("request_id", contextx.RequestID(c.Request.Context())),
		)
	}
}

// CORS æŒ‰é…ç½®å†™å…¥è·¨åŸŸå“åº”å¤´å¹¶å¤„ç† OPTIONS é¢„æ£€ã€‚
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	origins := map[string]struct{}{}
	for _, o := range cfg.AllowOrigins {
		origins[o] = struct{}{}
	}
	allowAll := false
	for _, o := range cfg.AllowOrigins {
		if o == "*" {
			allowAll = true
		}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if allowAll {
				c.Header("Access-Control-Allow-Origin", "*")
			} else if _, ok := origins[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				if cfg.AllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
			}
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ","))
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ","))
			c.Header("Access-Control-Expose-Headers", "X-Request-Id")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// SecurityHeaders å†™å…¥åŸºç¡€å®‰å…¨å“åº”å¤´ï¼›å¯é€‰å¼€å¯ HSTSã€‚
func SecurityHeaders(cfg config.SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		if cfg.HSTSEnabled {
			maxAge := cfg.HSTSMaxAgeSeconds
			if maxAge <= 0 {
				maxAge = 31536000
			}
			c.Header("Strict-Transport-Security", "max-age="+strconv.Itoa(maxAge)+"; includeSubDomains")
		}
		c.Next()
	}
}

// AuthContext ä»Ž Cookie/Header å¯é€‰åŠ è½½ä¼šè¯ï¼ˆä¸é€æ˜Ž tokenï¼Œéž Bearerï¼‰ã€‚
func AuthContext(cfg config.AuthConfig, store *security.SessionStore) gin.HandlerFunc {
	name := cfg.TokenName
	if name == "" {
		name = "Authorization"
	}
	return func(c *gin.Context) {
		token := ""
		if cfg.SessionCookieEnabled {
			if ck, err := c.Cookie(name); err == nil {
				token = ck
			}
		}
		if token == "" {
			h := c.GetHeader(name)
			if h != "" && !strings.HasPrefix(strings.ToLower(h), "bearer ") {
				token = h
			}
		}
		if token != "" {
			sess, err := store.Get(c.Request.Context(), token)
			if err == nil && sess != nil {
				ctx := contextx.WithSession(c.Request.Context(), sess)
				ctx = contextx.WithAccount(ctx, sess.AccountID, sess.AccountType)
				c.Request = c.Request.WithContext(ctx)
				c.Set("session", sess)
			}
		}
		c.Next()
	}
}

// AuthWhitelist å¯¹éžç™½åå•æœªç™»å½•è¯·æ±‚ç›´æŽ¥æ‹’ç»ã€‚
func AuthWhitelist(extra []string) gin.HandlerFunc {
	builtin := []string{
		"/",
		"/metrics",
		"/api/v1/internal/health/*",
		"/api/v1/admin/login",
		"/api/v1/admin/captcha",
		"/api/v1/admin/password-key",
		"/api/v1/admin/public/auth-options",
		"/api/v1/admin/send-login-code",
		"/api/v1/admin/forgot-password",
		"/api/v1/admin/reset-password",
		"/api/v1/admin/oauth/*",
		"/api/v1/portal/login",
		"/api/v1/portal/register",
		"/api/v1/portal/register/send-code",
		"/api/v1/portal/public/auth-options",
		"/api/v1/portal/captcha",
		"/api/v1/portal/password-key",
		"/api/v1/portal/send-login-code",
		"/api/v1/portal/forgot-password",
		"/api/v1/portal/reset-password",
		"/api/v1/portal/oauth/*",
		"/api/v1/files/*",
	}
	patterns := append(builtin, extra...)
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if matchAny(p, patterns) {
			c.Next()
			return
		}
		if contextx.Session(c.Request.Context()) == nil {
			response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func matchAny(p string, patterns []string) bool {
	for _, pat := range patterns {
		if ok, _ := path.Match(pat, p); ok {
			return true
		}
		if strings.HasSuffix(pat, "/*") {
			prefix := strings.TrimSuffix(pat, "/*")
			if p == prefix || strings.HasPrefix(p, prefix+"/") {
				return true
			}
		}
		if p == pat {
			return true
		}
	}
	return false
}

// RequireAccountType è¦æ±‚ä¼šè¯è´¦å·ç±»åž‹åŒ¹é…ã€‚
func RequireAccountType(t security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := contextx.Session(c.Request.Context())
		if sess == nil || sess.AccountType != t {
			response.Fail(c, http.StatusForbidden, 403, "forbidden account type")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequirePermission ç™»è®°æƒé™é”®å¹¶è¦æ±‚ä¼šè¯æŒæœ‰è¯¥æƒé™ã€‚
func RequirePermission(reg *security.PermissionRegistry, key, name string) gin.HandlerFunc {
	reg.Register(key, name)
	return func(c *gin.Context) {
		sess := contextx.Session(c.Request.Context())
		if sess == nil || !security.HasPermission(sess.PermissionKeys, key) {
			response.Fail(c, http.StatusForbidden, 403, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

// ErrorHandler å°† gin.Context é”™è¯¯æ˜ å°„ä¸ºç»Ÿä¸€å¤±è´¥ä¿¡å°ã€‚
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		if ae, ok := err.(*apperr.AppError); ok {
			status := ae.Code
			if status < 400 {
				status = 400
			}
			response.Fail(c, status, ae.Code, ae.Message)
			return
		}
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
	}
}
