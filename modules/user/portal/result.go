package portal

import (
	"hei-gin/framework/core/security"
)

// MeResult 当前登录用户概览。
//
// Author: Charlie
type MeResult struct {
	AccountID      string               `json:"account_id"`
	AccountType    security.AccountType `json:"account_type"`
	Name           *string              `json:"name"`
	Nickname       *string              `json:"nickname"`
	Avatar         *string              `json:"avatar"`
	RoleIDs        []string             `json:"role_ids"`
	DeptIDs        []string             `json:"dept_ids"`
	GroupIDs       []string             `json:"group_ids"`
	PermissionKeys []string             `json:"permission_keys"`
	Profile        *Profile             `json:"profile"`
}

// AvatarResult 头像上传结果。
//
// Author: Charlie
type AvatarResult struct {
	Avatar string `json:"avatar"`
}
