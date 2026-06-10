package middleware

import (
	"log"
	"strings"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/constants"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"
)

// HeiCheckPermission returns a middleware that checks the user has the required permissions.
// mode defaults to "AND" (all permissions required). Pass "OR" for any permission.
// This middleware is for BUSINESS login type.
func HeiCheckPermission(permissions []string, mode ...string) gin.HandlerFunc {
	m := "AND"
	if len(mode) > 0 {
		m = mode[0]
	}
	return heiCheckPermissionInner("BUSINESS", permissions, m)
}

// HeiClientCheckPermission returns a middleware that checks the CONSUMER user has the required permissions.
// mode defaults to "AND" (all permissions required). Pass "OR" for any permission.
func HeiClientCheckPermission(permissions []string, mode ...string) gin.HandlerFunc {
	m := "AND"
	if len(mode) > 0 {
		m = mode[0]
	}
	return heiCheckPermissionInner("CONSUMER", permissions, m)
}

// heiCheckPermissionInner is a shared implementation for both BUSINESS and CONSUMER permission checks.
func heiCheckPermissionInner(loginType string, permissions []string, mode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check login first
		var isLogin bool
		if loginType == "CONSUMER" {
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

		// Super admin bypass: users with the SUPER_ADMIN role automatically pass all permission checks.
		// This avoids reliance on the Redis permission cache which may not be populated.
		roles, err := auth.GetRoleList(c, loginType)
		if err != nil {
			log.Printf("[Permission] Failed to get role list: %v", err)
		}
		for _, role := range roles {
			if role == constants.SUPER_ADMIN_CODE {
				c.Next()
				return
			}
		}

		// Check permission
		if mode == "OR" {
			if !auth.HasPermissionOr(c, loginType, permissions...) {
				c.Abort()
				result.Failure(c, "缺少权限: "+strings.Join(permissions, ","), 403)
				return
			}
		} else {
			if !auth.HasPermissionAnd(c, loginType, permissions...) {
				c.Abort()
				result.Failure(c, "缺少权限: "+strings.Join(permissions, ","), 403)
				return
			}
		}
		c.Next()
	}
}

// getModuleFromCode extracts the module segment from a permission code.
// Example: "user:add" -> "user", "sys:user:view" -> "sys:user"
func getModuleFromCode(code string) string {
	parts := strings.Split(code, ":")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ":")
	}
	return code
}
