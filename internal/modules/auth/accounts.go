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
	// GetEnabledAccount 按 ID 取已启用账号类型。
	GetEnabledAccount(ctx context.Context, accountID string) (accountType security.AccountType, err error)
	// UpdatePasswordHash 更新密码哈希。
	UpdatePasswordHash(ctx context.Context, accountID, passwordHash string) error
}

// PortalRegistrar 可选门户注册能力；AccountFinder 同时实现时可用。
//
// Author: Charlie
type PortalRegistrar interface {
	// RegisterPortal 注册门户账号并返回账号 ID。
	RegisterPortal(ctx context.Context, account, passwordHash string, name, nickname, email, phone *string) (accountID string, err error)
}
