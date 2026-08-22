// internal/modules/profile/result.go 用户中心出参定义。
//
// Author: Charlie

package profile

import (
	"time"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/modules/profile/identity"
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
	AccountID         string               `json:"account_id"`
	AccountType       security.AccountType `json:"account_type"`
	Account           string               `json:"account"`
	Nickname          *string              `json:"nickname"`
	Avatar            *string              `json:"avatar"`
	Identity          *identity.IdentityStatusResult `json:"identity"`
	RoleIDs           []string             `json:"role_ids"`
	DeptIDs           []string             `json:"dept_ids"`
	GroupIDs          []string             `json:"group_ids"`
	RoleIDNames       []IDName             `json:"role_id_names"`
	DeptIDNames       []IDName             `json:"dept_id_names"`
	GroupIDNames      []IDName             `json:"group_id_names"`
	PermissionKeys    []string               `json:"permission_keys"`
	Profile           *UserProfileResult     `json:"profile"`
	PasswordExpired   bool                 `json:"password_expired"`
	ForceBindEmail    bool                 `json:"force_bind_email"`
	ForceBindPhone    bool                 `json:"force_bind_phone"`
	ForceBindIdentity bool                 `json:"force_bind_identity"`
}

// UserProfileResult 用户资料 DTO（对齐 hei-boot UserProfileResult）。
type UserProfileResult struct {
	AccountID          string     `json:"account_id"`
	Nickname           *string    `json:"nickname"`
	Avatar             *string    `json:"avatar"`
	Signature          *string    `json:"signature"`
	Phone              *string    `json:"phone"`
	Email              *string    `json:"email"`
	PhoneLoginEnabled  bool       `json:"phone_login_enabled"`
	EmailLoginEnabled  bool       `json:"email_login_enabled"`
	Remark             *string    `json:"remark"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// AvatarResult 头像更新出参。
//
// Author: Charlie
type AvatarResult struct {
	Avatar string `json:"avatar"`
}

// OrgInfoResult 组织关联信息（对齐 hei-boot OrgInfoResult）。
type OrgInfoResult struct {
	RoleIDNames  []IDName `json:"role_id_names"`
	DeptIDNames  []IDName `json:"dept_id_names"`
	GroupIDNames []IDName `json:"group_id_names"`
}
