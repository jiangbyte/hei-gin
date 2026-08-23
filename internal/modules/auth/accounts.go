// internal/modules/auth/accounts.go 账号查找接口。
//
// Author: Charlie

package auth

import (
	"context"

	"hei-gin/internal/framework/core/security"
)

// AccountFinder 登录使用的窄 IAM 接口（由 iam/account 实现）。
//
// Author: Charlie
type AccountFinder interface {
	// FindEnabledByIdentity 按身份查找已启用账号，返回 ID 与密码哈希。
	FindEnabledByIdentity(ctx context.Context, accountType security.AccountType, identityType, identifier string) (accountID, passwordHash string, err error)
	// EnsureSuperPermissions 解析账号权限键与授权列表。
	EnsureSuperPermissions(ctx context.Context, accountID string) (keys []string, grants []security.PermissionGrant, err error)
	// GetSessionAuthorization 聚合会话授权快照（角色/组织/资源/权限）。
	GetSessionAuthorization(ctx context.Context, accountID string) (*security.AuthorizationSnapshot, error)
	// GetEnabledAccount 按 ID 取已启用账号类型。
	GetEnabledAccount(ctx context.Context, accountID string) (accountType security.AccountType, err error)
	// UpdatePasswordHash 更新密码哈希。
	UpdatePasswordHash(ctx context.Context, accountID, passwordHash string) error
	// HasBoundIdentity 账号是否已绑定指定身份类型。
	HasBoundIdentity(ctx context.Context, accountID, identityType string) bool
}

// PortalRegistrar 可选门户注册能力；AccountFinder 同时实现时可用。
//
// Author: Charlie
type PortalRegisterInput struct {
	AccountLogin   string
	PasswordHash   string
	Nickname       *string
	Email          *string
	Phone          *string
	EmailEnabled   bool
	PhoneEnabled   bool
	EmailVerified  bool
	PhoneVerified  bool
}

type PortalRegistrar interface {
	// RegisterPortal 注册门户账号并返回账号 ID 与登录名。
	RegisterPortal(ctx context.Context, in PortalRegisterInput) (accountID, accountLogin string, err error)
	// AllocateUniqueAccount 分配唯一登录账号名。
	AllocateUniqueAccount(ctx context.Context, base string) (string, error)
	// IdentityExists 身份标识是否已被占用。
	IdentityExists(ctx context.Context, identityType, identifier string) bool
}
