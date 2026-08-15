// internal/modules/profile/result.go 用户中心出参定义。
//
// Author: Charlie

package profile

import (
	"hei-gin/internal/framework/core/security"
)

// IDName 组织关联 ID 与名称。
//
// Author: Charlie
type IDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MeResult 当前登录用户概览（对齐 hei-boot MeResult）。
//
// Author: Charlie
type MeResult struct {
	AccountID       string               `json:"account_id"`
	AccountType     security.AccountType `json:"account_type"`
	Account         string               `json:"account"`
	Name            *string              `json:"name"`
	Nickname        *string              `json:"nickname"`
	Avatar          *string              `json:"avatar"`
	RoleIDs         []string             `json:"role_ids"`
	DeptIDs         []string             `json:"dept_ids"`
	GroupIDs        []string             `json:"group_ids"`
	RoleIDNames     []IDName             `json:"role_id_names"`
	DeptIDNames     []IDName             `json:"dept_id_names"`
	GroupIDNames    []IDName             `json:"group_id_names"`
	PermissionKeys  []string             `json:"permission_keys"`
	Profile         *Profile             `json:"profile"`
	PasswordExpired bool                 `json:"password_expired"`
	ForceBindEmail  bool                 `json:"force_bind_email"`
	ForceBindPhone  bool                 `json:"force_bind_phone"`
}

// AvatarResult 头像上传结果。
//
// Author: Charlie
type AvatarResult struct {
	Avatar string `json:"avatar"`
}

// OrgInfoResult 组织关联信息。
//
// Author: Charlie
type OrgInfoResult struct {
	RoleIDs  []string `json:"role_ids"`
	DeptIDs  []string `json:"dept_ids"`
	GroupIDs []string `json:"group_ids"`
}
