package middleware

import (
	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

// HeiCheckLogin returns a middleware that checks if the user is logged in.
// loginType defaults to BUSINESS. Pass CONSUMER for client-side users.
// Sets "loginUser" in the Gin context for downstream audit logging.
// Injects the user ID into c.Request.Context() so GORM hooks can auto-fill CreatedBy/UpdatedBy.
func HeiCheckLogin(loginType ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lt := string(enums.LoginTypeBusiness)
		if len(loginType) > 0 {
			lt = loginType[0]
		}

		var uid string
		if lt == string(enums.LoginTypeConsumer) {
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

		AttachLoginContext(c, lt, uid)

		c.Next()
	}
}

// HeiClientCheckLogin returns a middleware that checks if the CONSUMER user is logged in.
func HeiClientCheckLogin() gin.HandlerFunc {
	return HeiCheckLogin(string(enums.LoginTypeConsumer))
}
