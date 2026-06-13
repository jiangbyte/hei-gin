package middleware

import (
	"hei-gin/sdk/auth"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

func CheckLogin(realm *auth.Realm) gin.HandlerFunc {
	return func(c *gin.Context) {
		if realm == nil {
			realm = auth.Business
		}
		if !realm.IsLogin(c) {
			c.Abort()
			result.Failure(c, "未授权/未登录", 401)
			return
		}
		AttachLoginContext(c, realm, realm.GetLoginIDDefaultNull(c))

		c.Next()
	}
}
