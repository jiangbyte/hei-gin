package auth

import (
	"context"
	"fmt"
	"hei-gin/sdk/web/result"
	"time"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/enums"
)

func getLoginID(c *gin.Context, loginType string) string {
	if loginType == string(enums.LoginTypeConsumer) {
		return Consumer.GetLoginID(c)
	}
	return GetLoginID(c)
}

func GetPermissionList(c *gin.Context, loginType string) ([]string, error) {
	cacheKey := "_perm_cache_" + loginType
	if cached, exists := c.Get(cacheKey); exists {
		if perms, ok := cached.([]string); ok {
			return perms, nil
		}
	}

	iface := PermissionDelegate
	if iface == nil {
		return []string{}, nil
	}

	loginID := getLoginID(c, loginType)
	if loginID == "" {
		return []string{}, nil
	}

	perms, err := iface.GetPermissionList(c.Request.Context(), loginID, loginType)
	if err != nil {
		return perms, err
	}
	c.Set(cacheKey, perms)
	return perms, nil
}

func GetRoleList(c *gin.Context, loginType string) ([]string, error) {
	cacheKey := "_role_cache_" + loginType
	if cached, exists := c.Get(cacheKey); exists {
		if roles, ok := cached.([]string); ok {
			return roles, nil
		}
	}

	iface := PermissionDelegate
	if iface == nil {
		return []string{}, nil
	}

	loginID := getLoginID(c, loginType)
	if loginID == "" {
		return []string{}, nil
	}

	roles, err := iface.GetRoleList(c.Request.Context(), loginID, loginType)
	if err != nil {
		return roles, err
	}
	c.Set(cacheKey, roles)
	return roles, nil
}

func GetPermissionListByLoginID(ctx context.Context, loginID, loginType string) ([]string, error) {
	iface := PermissionDelegate
	if iface == nil {
		return []string{}, nil
	}
	return iface.GetPermissionList(ctx, loginID, loginType)
}

func GetRoleListByLoginID(ctx context.Context, loginID, loginType string) ([]string, error) {
	iface := PermissionDelegate
	if iface == nil {
		return []string{}, nil
	}
	return iface.GetRoleList(ctx, loginID, loginType)
}

func HasPermission(c *gin.Context, permission, loginType string) bool {
	permissions, err := GetPermissionList(c, loginType)
	if err != nil {
		return false
	}
	return MatchPermission(permission, permissions)
}

func CheckPermission(c *gin.Context, permission, loginType string) {
	if !HasPermission(c, permission, loginType) {
		result.Failure(c, fmt.Sprintf("缺少权限: %s", permission), 403)
	}
}

func HasPermissionAnd(c *gin.Context, loginType string, permissions ...string) bool {
	perms, err := GetPermissionList(c, loginType)
	if err != nil {
		return false
	}
	return MatchPermissionsAnd(permissions, perms)
}

func CheckPermissionAnd(c *gin.Context, loginType string, permissions ...string) {
	for _, permission := range permissions {
		if !HasPermission(c, permission, loginType) {
			result.Failure(c, fmt.Sprintf("缺少权限: %s", permission), 403)
			return
		}
	}
}

func HasPermissionOr(c *gin.Context, loginType string, permissions ...string) bool {
	perms, err := GetPermissionList(c, loginType)
	if err != nil {
		return false
	}
	return MatchPermissionsOr(permissions, perms)
}

func CheckPermissionOr(c *gin.Context, loginType string, permissions ...string) {
	if !HasPermissionOr(c, loginType, permissions...) {
		result.Failure(c, fmt.Sprintf("缺少权限: %v", permissions), 403)
	}
}

func HasRole(c *gin.Context, role, loginType string) bool {
	roles, err := GetRoleList(c, loginType)
	if err != nil {
		return false
	}
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func CheckRole(c *gin.Context, role, loginType string) {
	if !HasRole(c, role, loginType) {
		result.Failure(c, fmt.Sprintf("缺少角色: %s", role), 403)
	}
}

func HasRoleAnd(c *gin.Context, loginType string, roles ...string) bool {
	userRoles, err := GetRoleList(c, loginType)
	if err != nil {
		return false
	}
	roleSet := make(map[string]struct{}, len(userRoles))
	for _, r := range userRoles {
		roleSet[r] = struct{}{}
	}
	for _, r := range roles {
		if _, ok := roleSet[r]; !ok {
			return false
		}
	}
	return true
}

func CheckRoleAnd(c *gin.Context, loginType string, roles ...string) {
	for _, role := range roles {
		if !HasRole(c, role, loginType) {
			result.Failure(c, fmt.Sprintf("缺少角色: %s", role), 403)
			return
		}
	}
}

func HasRoleOr(c *gin.Context, loginType string, roles ...string) bool {
	userRoles, err := GetRoleList(c, loginType)
	if err != nil {
		return false
	}
	roleSet := make(map[string]struct{}, len(userRoles))
	for _, r := range userRoles {
		roleSet[r] = struct{}{}
	}
	for _, r := range roles {
		if _, ok := roleSet[r]; ok {
			return true
		}
	}
	return false
}

func CheckRoleOr(c *gin.Context, loginType string, roles ...string) {
	if !HasRoleOr(c, loginType, roles...) {
		result.Failure(c, fmt.Sprintf("缺少角色: %v", roles), 403)
	}
}

var _ = context.Background
var _ = time.Time{}
