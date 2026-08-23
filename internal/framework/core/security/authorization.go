// Package security authorization snapshot for sessions.
//
// Author: Charlie
package security

// AuthorizationSnapshot 登录会话授权快照（对齐 hei-boot AccountAuthorizationInfo）。
type AuthorizationSnapshot struct {
	RoleIDs              []string
	RoleCodes            []string
	DeptIDs              []string
	GroupIDs             []string
	ResourceIDs          []string
	ClientResourceIDs    []string
	PermissionKeys       []string
	ClientPermissionKeys []string
	PermissionGrants     []PermissionGrant
}
