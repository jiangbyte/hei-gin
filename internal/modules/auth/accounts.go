package auth

import (
	"context"

	"hei-gin/internal/framework/core/security"
)

// AccountFinder ç™»å½•ä½¿ç”¨çš„çª„ IAM æŽ¥å£ï¼ˆç”± iam/account å®žçŽ°ï¼‰ã€‚
//
// Author: Charlie
type AccountFinder interface {
	// FindEnabledByIdentity æŒ‰èº«ä»½æŸ¥æ‰¾å·²å¯ç”¨è´¦å·ï¼Œè¿”å›ž ID ä¸Žå¯†ç å“ˆå¸Œã€‚
	FindEnabledByIdentity(ctx context.Context, accountType security.AccountType, identityType, identifier string) (accountID, passwordHash string, err error)
	// EnsureSuperPermissions è§£æžè´¦å·æƒé™é”®ä¸ŽæŽˆæƒåˆ—è¡¨ã€‚
	EnsureSuperPermissions(ctx context.Context, accountID string) (keys []string, grants []security.PermissionGrant, err error)
	// GetEnabledAccount æŒ‰ ID å–å·²å¯ç”¨è´¦å·ç±»åž‹ã€‚
	GetEnabledAccount(ctx context.Context, accountID string) (accountType security.AccountType, err error)
	// UpdatePasswordHash æ›´æ–°å¯†ç å“ˆå¸Œã€‚
	UpdatePasswordHash(ctx context.Context, accountID, passwordHash string) error
}

// PortalRegistrar å¯é€‰é—¨æˆ·æ³¨å†Œèƒ½åŠ›ï¼›AccountFinder åŒæ—¶å®žçŽ°æ—¶å¯ç”¨ã€‚
//
// Author: Charlie
type PortalRegistrar interface {
	// RegisterPortal æ³¨å†Œé—¨æˆ·è´¦å·å¹¶è¿”å›žè´¦å· IDã€‚
	RegisterPortal(ctx context.Context, account, passwordHash string, name, nickname, email, phone *string) (accountID string, err error)
}
