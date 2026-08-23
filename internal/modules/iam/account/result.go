// internal/modules/iam/account/result.go 出参定义（对齐 hei-boot SysAccountResult + 身份/三方绑定行）。
//
// Author: Charlie

package account

import (
	"time"

	"hei-gin/internal/modules/iam/client"
	"hei-gin/internal/modules/iam/group"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/resource"
	"hei-gin/internal/modules/iam/role"
	"hei-gin/internal/modules/profile/identity"
)

// IdentityResult 账号身份行（对齐 hei-boot AccountIdentityResult）。
//
// Author: Charlie
type IdentityResult struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	IdentityType string    `json:"identity_type"`
	Identifier   string    `json:"identifier"`
	Verified     bool      `json:"verified"`
	IsPrimary    bool      `json:"is_primary"`
	BindStatus   string    `json:"bind_status"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    *string   `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    *string   `json:"updated_by"`
}

// OAuthBindingResult 账号三方绑定行（对齐 hei-boot AccountOauthBindingResult）。
//
// Author: Charlie
type OAuthBindingResult struct {
	ID       string     `json:"id"`
	Provider string     `json:"provider"`
	OpenID   string     `json:"open_id"`
	UnionID  *string    `json:"union_id"`
	Nickname *string    `json:"nickname"`
	Avatar   *string    `json:"avatar"`
	BoundAt  *time.Time `json:"bound_at"`
}

// AccountListResult 账号分页列表行（仅 sys_account + profile 展示字段）。
//
// Author: Charlie
type AccountListResult struct {
	ID              string     `json:"id"`
	Account         string     `json:"account"`
	AccountType     string     `json:"account_type"`
	AccountStatus   string     `json:"account_status"`
	Nickname        *string    `json:"nickname"`
	Avatar          *string    `json:"avatar"`
	Phone           *string    `json:"phone"`
	Email           *string    `json:"email"`
	Remark          *string    `json:"remark"`
	LatestLoginTime *time.Time `json:"latest_login_time"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// AccountResult 账号详情。
//
// Author: Charlie
type AccountResult struct {
	ID                      string               `json:"id"`
	Account                 string               `json:"account"`
	AccountType             string               `json:"account_type"`
	AccountStatus           string               `json:"account_status"`
	Name                    *string              `json:"name"`
	Nickname                *string              `json:"nickname"`
	Avatar                  *string              `json:"avatar"`
	Signature               *string              `json:"signature"`
	Phone                   *string              `json:"phone"`
	Email                   *string              `json:"email"`
	EmailLoginEnabled       bool                 `json:"email_login_enabled"`
	PhoneLoginEnabled       bool                 `json:"phone_login_enabled"`
	EmailIdentity           *string              `json:"email_identity"`
	PhoneIdentity           *string              `json:"phone_identity"`
	EmailIdentityVerified   bool                 `json:"email_identity_verified"`
	PhoneIdentityVerified   bool                 `json:"phone_identity_verified"`
	EmailIdentityBindStatus *string              `json:"email_identity_bind_status"`
	PhoneIdentityBindStatus *string              `json:"phone_identity_bind_status"`
	Identities              []IdentityResult     `json:"identities"`
	OAuthBindings           []OAuthBindingResult `json:"oauth_bindings"`
	IdentityStatus          *identity.IdentityStatusResult `json:"identity_status"`
	Remark                  *string              `json:"remark"`
	CancelledAt             *time.Time           `json:"cancelled_at"`
	CancelledBy             *string              `json:"cancelled_by"`
	CancelReason            *string              `json:"cancel_reason"`
	LastLoginIP             *string              `json:"last_login_ip"`
	LastLoginAddress        *string              `json:"last_login_address"`
	LastLoginTime           *time.Time           `json:"last_login_time"`
	LastLoginDevice         *string              `json:"last_login_device"`
	LatestLoginIP           *string              `json:"latest_login_ip"`
	LatestLoginAddress      *string              `json:"latest_login_address"`
	LatestLoginTime         *time.Time           `json:"latest_login_time"`
	LatestLoginDevice       *string              `json:"latest_login_device"`
	CreatedAt               time.Time            `json:"created_at"`
	CreatedBy               *string              `json:"created_by"`
	UpdatedAt               time.Time            `json:"updated_at"`
	UpdatedBy               *string              `json:"updated_by"`
}

// OwnRoleResult 账号已拥有角色结果。
//
// Author: Charlie
type OwnRoleResult struct {
	ID      string      `json:"id"`
	Roles   []role.Role `json:"roles"`
	RoleIDs []string    `json:"role_ids"`
}

// OwnGroupResult 账号已拥有用户组结果。
//
// Author: Charlie
type OwnGroupResult struct {
	ID       string        `json:"id"`
	Groups   []group.Group `json:"groups"`
	GroupIDs []string      `json:"group_ids"`
}

// OwnDeptResult 账号已拥有部门授权结果。
//
// Author: Charlie
type OwnDeptResult struct {
	ID            string                   `json:"id"`
	GrantInfoList []relation.DeptGrantInfo `json:"grant_info_list"`
}

// OwnResourceResult 账号已拥有管理端资源授权结果。
//
// Author: Charlie
type OwnResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []resource.GrantModule       `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

// OwnClientResourceResult 账号已拥有客户端资源授权结果。
//
// Author: Charlie
type OwnClientResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []client.GrantModule         `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
