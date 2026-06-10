package middleware

import (
	"context"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/db"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"
)

// HeiCheckLogin returns a middleware that checks if the user is logged in.
// loginType defaults to "BUSINESS". Pass "CONSUMER" for client-side users.
// Sets "loginUser" in the Gin context for downstream audit logging.
// Injects the user ID into c.Request.Context() so GORM hooks can auto-fill CreatedBy/UpdatedBy.
func HeiCheckLogin(loginType ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lt := "BUSINESS"
		if len(loginType) > 0 {
			lt = loginType[0]
		}

		var uid string
		if lt == "CONSUMER" {
			if !auth.Consumer.IsLogin(c) {
				c.Abort()
				result.Failure(c, "未授权/未登录", 401)
				return
			}
			uid = auth.Consumer.GetLoginIDDefaultNull(c)
		} else {
			if !auth.IsLogin(c) {
				c.Abort()
				result.Failure(c, "未授权/未登录", 401)
				return
			}
			uid = auth.GetLoginIDDefaultNull(c)
		}

		// Set loginUser for downstream audit logging
		if username := auth.GetExtra(c, "username"); username != nil {
			if u, ok := username.(string); ok && u != "" {
				c.Set("loginUser", u)
			}
		}

		// Inject user ID into request context for GORM global callback
		if uid != "" {
			ctx := context.WithValue(c.Request.Context(), db.CtxKeyLoginID{}, uid)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// HeiClientCheckLogin returns a middleware that checks if the CONSUMER user is logged in.
func HeiClientCheckLogin() gin.HandlerFunc {
	return HeiCheckLogin("CONSUMER")
}
