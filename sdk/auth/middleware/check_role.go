package middleware

import (
	"strings"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

// HeiCheckRole returns a middleware that checks the user has the required roles.
// mode defaults to "AND" (all roles required). Pass "OR" for any role.
// This middleware is for BUSINESS login type.
func HeiCheckRole(roles []string, mode ...string) gin.HandlerFunc {
	m := "AND"
	if len(mode) > 0 {
		m = mode[0]
	}
	return heiCheckRoleInner(string(enums.LoginTypeBusiness), roles, m)
}

// HeiClientCheckRole returns a middleware that checks the CONSUMER user has the required roles.
// mode defaults to "AND" (all roles required). Pass "OR" for any role.
func HeiClientCheckRole(roles []string, mode ...string) gin.HandlerFunc {
	m := "AND"
	if len(mode) > 0 {
		m = mode[0]
	}
	return heiCheckRoleInner(string(enums.LoginTypeConsumer), roles, m)
}

// heiCheckRoleInner is a shared implementation for both BUSINESS and CONSUMER role checks.
func heiCheckRoleInner(loginType string, roles []string, mode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check login first
		var isLogin bool
		if loginType == string(enums.LoginTypeConsumer) {
			tool := auth.Consumer
			isLogin = tool.IsLogin(c)
		} else {
			isLogin = auth.IsLogin(c)
		}
		if !isLogin {
			c.Abort()
			result.Failure(c, "未授权/未登录", 401)
			return
		}

		// Check role
		if mode == "OR" {
			if !auth.HasRoleOr(c, loginType, roles...) {
				c.Abort()
				result.Failure(c, "缺少角色: "+strings.Join(roles, ","), 403)
				return
			}
		} else {
			if !auth.HasRoleAnd(c, loginType, roles...) {
				c.Abort()
				result.Failure(c, "缺少角色: "+strings.Join(roles, ","), 403)
				return
			}
		}
		c.Next()
	}
}
