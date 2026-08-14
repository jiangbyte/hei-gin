package portal

import (
	"hei-gin/internal/framework/core/security"
)

// MeResult å½“å‰ç™»å½•ç”¨æˆ·æ¦‚è§ˆã€‚
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

// AvatarResult å¤´åƒä¸Šä¼ ç»“æžœã€‚
//
// Author: Charlie
type AvatarResult struct {
	Avatar string `json:"avatar"`
}
