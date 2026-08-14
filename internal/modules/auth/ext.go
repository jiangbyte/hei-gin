// internal/modules/auth/ext.go 认证扩展。
//
// Author: Charlie

package auth

import (
	"context"
	"fmt"
	"hei-gin/internal/modules/auth/oauth"
	"strings"
	"time"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/security"
)

// AuthOptions ç™»å½•é¡µå…¬å¼€é…ç½®ï¼ˆå¯¹é½ hei-boot AuthOptionsResultï¼‰ã€‚
//
// Author: Charlie
type AuthOptions struct {
	AccountType                string                 `json:"account_type"`
	AllowAccount               bool                   `json:"allow_account"`
	AllowEmail                 bool                   `json:"allow_email"`
	AllowPhone                 bool                   `json:"allow_phone"`
	AllowOTP                   bool                   `json:"allow_otp"`
	RegisterEnabled            bool                   `json:"register_enabled"`
	RegisterAllowAccount       bool                   `json:"register_allow_account"`
	RegisterAllowEmail         bool                   `json:"register_allow_email"`
	RegisterAllowPhone         bool                   `json:"register_allow_phone"`
	ForceBindEmail             bool                   `json:"force_bind_email"`
	ForceBindPhone             bool                   `json:"force_bind_phone"`
	OAuthProviders             []oauth.ProviderOption `json:"oauth_providers"`
	PasswordChangeVerifyMethod string                 `json:"password_change_verify_method"`
	CopyrightText              string                 `json:"copyright_text"`
	CopyrightURL               string                 `json:"copyright_url"`
}

// OauthProviderOption ä¸‰æ–¹ç™»å½•å…¥å£é€‰é¡¹ã€‚
//
// Author: Charlie
// RefreshResult åˆ·æ–°ä¼šè¯ç»“æžœï¼ˆä¸Ž LoginResult ä¸€è‡´ï¼‰ã€‚
type RefreshResult = LoginResult

// CancelParam æ³¨é”€è´¦å·å…¥å‚ã€‚
//
// Author: Charlie
type CancelParam struct {
	CancelReason *string `json:"cancel_reason"`
}

// AuthOptions è¯»å–ç™»å½•é¡µå…¬å¼€é…ç½®ã€‚
func (s *Service) AuthOptions(ctx context.Context, accountType security.AccountType) *AuthOptions {
	o := &AuthOptions{
		AccountType:                string(accountType),
		AllowAccount:               true,
		AllowEmail:                 true,
		AllowPhone:                 true,
		AllowOTP:                   true,
		RegisterEnabled:            s.cfg.Auth.PortalRegisterEnabled && accountType == security.AccountPortal,
		RegisterAllowAccount:       accountType == security.AccountPortal,
		RegisterAllowEmail:         accountType == security.AccountPortal,
		RegisterAllowPhone:         accountType == security.AccountPortal,
		PasswordChangeVerifyMethod: "OLD_PASSWORD",
	}
	if s.oauth != nil {
		o.OAuthProviders = s.oauth.ProviderOptions()
	}
	return o
}

// RefreshSession ç»­æœŸå½“å‰ä¼šè¯ Tokenã€‚
func (s *Service) RefreshSession(ctx context.Context, accountType security.AccountType, clientIP, userAgent string) (*LoginResult, error) {
	sess := contextx.Session(ctx)
	if sess == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	if sess.AccountType != accountType {
		return nil, fmt.Errorf("unauthorized")
	}
	// é‡æ–°è®¡ç®—æƒé™ï¼ˆè§’è‰²/èµ„æºå˜æ›´åŽåˆ·æ–°ç”Ÿæ•ˆï¼‰
	keys, grants, err := s.accounts.EnsureSuperPermissions(ctx, sess.AccountID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	ttlSec := s.cfg.Auth.TokenTTLSeconds
	if !sess.RememberMe && s.cfg.Auth.TokenTTLShortSeconds > 0 {
		ttlSec = s.cfg.Auth.TokenTTLShortSeconds
	}
	if ttlSec <= 0 {
		ttlSec = 14400
	}
	ttl := time.Duration(ttlSec) * time.Second
	sess.PermissionKeys = keys
	sess.PermissionGrants = grants
	sess.LastActiveAt = now
	sess.ExpiresAt = now.Add(ttl)
	if err := s.sessions.Save(ctx, sess, ttl); err != nil {
		return nil, err
	}
	return &LoginResult{
		Token:           sess.Token,
		AccountID:       sess.AccountID,
		AccountType:     sess.AccountType,
		PasswordExpired: sess.PasswordExpired,
	}, nil
}

// CancelAccount æ³¨é”€ï¼ˆåœç”¨ï¼‰å½“å‰è´¦å·å¹¶æ¸…ç†ä¼šè¯ã€‚
func (s *Service) CancelAccount(ctx context.Context, accountType security.AccountType, clientIP, userAgent string, reason *string) error {
	sess := contextx.Session(ctx)
	if sess == nil {
		return fmt.Errorf("unauthorized")
	}
	if s.accounts == nil {
		return errAccountFinder
	}
	accountID := sess.AccountID
	// æ ‡è®°è´¦å·å·²å–æ¶ˆï¼ˆè½¯æ³¨é”€ï¼‰ï¼šæ›´æ–° sys_account çŠ¶æ€
	now := time.Now().UTC()
	if err := s.db.WithContext(ctx).Table("sys_account").
		Where("id = ?", accountID).
		Update("account_status", security.AccountStatusCancelled).Error; err != nil {
		return err
	}
	_ = s.db.WithContext(ctx).Table("admin_user_profile").
		Where("account_id = ?", accountID).
		Update("remark", "cancelled at "+now.Format(time.RFC3339)).Error
	_ = s.db.WithContext(ctx).Table("portal_user_profile").
		Where("account_id = ?", accountID).
		Update("remark", "cancelled at "+now.Format(time.RFC3339)).Error
	_ = s.sessions.DeleteAllForAccount(ctx, accountID)
	s.publishAudit(ctx, "cancel", true, accountID, string(accountType), clientIP, userAgent, reasonString(reason))
	return nil
}

// sendRegisterCode é—¨æˆ·æ³¨å†Œå‘é€éªŒè¯ç ã€‚
func (s *Service) sendRegisterCode(ctx context.Context, req SendLoginCodeParam) error {
	if !s.cfg.Auth.PortalRegisterEnabled {
		return errRegisterDisabled
	}
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return err
	}
	channel, target := resolveOTPTarget(req)
	if target == "" {
		return errOTPTargetRequired
	}
	code, err := sixDigitCode()
	if err != nil {
		return err
	}
	ttl := 5 * time.Minute
	if err := s.repo.StoreLoginOTP(ctx, string(security.AccountPortal), channel, target, code, ttl); err != nil {
		return err
	}
	if s.notify != nil {
		vars := map[string]any{"app_name": s.cfg.App.Name, "code": code, "expire_minutes": 5}
		if channel == "PHONE" {
			_ = s.notify.SendTemplated(ctx, "LOGIN_CODE", target, vars)
		} else {
			_ = s.notify.SendTemplated(ctx, "LOGIN_CODE_MAIL", target, vars)
		}
	}
	return nil
}

func reasonString(reason *string) string {
	if reason == nil || strings.TrimSpace(*reason) == "" {
		return "cancel account"
	}
	return "cancel account: " + *reason
}
