package middleware

import (
	"strings"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

func CheckRole(realm *auth.Realm, roles []string, mode ...string) gin.HandlerFunc {
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
			if !realm.HasRoleOr(c, roles...) {
				c.Abort()
				result.Failure(c, "缺少角色: "+strings.Join(roles, ","), 403)
				return
			}
		} else {
			if !realm.HasRoleAnd(c, roles...) {
				c.Abort()
				result.Failure(c, "缺少角色: "+strings.Join(roles, ","), 403)
				return
			}
		}
		c.Next()
	}
}
