// Package middleware 提供 Gin 通用中间件：恢复、追踪、CORS、鉴权与错误映射。
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

	"hei-gin/framework/core/config"
	contextx "hei-gin/framework/core/context"
	apperr "hei-gin/framework/core/errors"
	"hei-gin/framework/core/logger"
	"hei-gin/framework/core/response"
	"hei-gin/framework/core/security"
)

// Recovery 捕获 panic 并返回 500 信封。
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

// Trace 注入 X-Request-Id 与客户端 IP 到请求上下文。
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

// AccessLog 记录访问日志（方法、路径、状态、耗时、request_id）。
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

// CORS 按配置写入跨域响应头并处理 OPTIONS 预检。
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

// SecurityHeaders 写入基础安全响应头；可选开启 HSTS。
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

// AuthContext 从 Cookie/Header 可选加载会话（不透明 token，非 Bearer）。
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

// AuthWhitelist 对非白名单未登录请求直接拒绝。
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

// RequireAccountType 要求会话账号类型匹配。
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

// RequirePermission 登记权限键并要求会话持有该权限。
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

// ErrorHandler 将 gin.Context 错误映射为统一失败信封。
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