package middleware

import (
	"strings"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

func CheckPermission(realm *auth.Realm, permissions []string, mode ...string) gin.HandlerFunc {
	checkMode := "AND"
	if len(mode) > 0 {
		checkMode = mode[0]
	}
	return func(c *gin.Context) {
		if realm == nil {
			realm = auth.Business
		}
		if !realm.IsLogin(c) {
			c.Abort()
			result.Failure(c, "未授权/未登录", 401)
			return
		}

		if checkMode == "OR" {
			if !realm.HasPermissionOr(c, permissions...) {
				c.Abort()
				result.Failure(c, "缺少权限: "+strings.Join(permissions, ","), 403)
				return
			}
		} else {
			if !realm.HasPermissionAnd(c, permissions...) {
				c.Abort()
				result.Failure(c, "缺少权限: "+strings.Join(permissions, ","), 403)
				return
			}
		}
		c.Next()
	}
}
