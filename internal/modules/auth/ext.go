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

// AuthOptions 登录页公开配置（对齐 hei-boot AuthOptionsResult）。
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

// OauthProviderOption 三方登录入口选项。
//
// Author: Charlie
// RefreshResult 刷新会话结果（与 LoginResult 一致）。
type RefreshResult = LoginResult

// CancelParam 注销账号入参。
//
// Author: Charlie
type CancelParam struct {
	CancelReason *string `json:"cancel_reason"`
}

// AuthOptions 读取登录页公开配置。
func (s *Service) AuthOptions(ctx context.Context, accountType security.AccountType) *AuthOptions {
	o := &AuthOptions{
		AccountType:                string(accountType),
		AllowAccount:               true,
		AllowEmail:                 true,
		AllowPhone:                 true,
		AllowOTP:                   true,
		RegisterEnabled:            s.registerEnabled(ctx) && accountType == security.AccountPortal,
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

// RefreshSession 续期当前会话 Token。
func (s *Service) RefreshSession(ctx context.Context, accountType security.AccountType, clientIP, userAgent string) (*LoginResult, error) {
	sess := contextx.Session(ctx)
	if sess == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	if sess.AccountType != accountType {
		return nil, fmt.Errorf("unauthorized")
	}
	// 重新计算权限（角色/资源变更后刷新生效）
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
		Token:                     sess.Token,
		AccountID:                 sess.AccountID,
		AccountType:               sess.AccountType,
		PasswordExpired:           sess.PasswordExpired,
		ForceBindEmail:            s.forceBind(ctx, sess.AccountType, "EMAIL"),
		ForceBindPhone:            s.forceBind(ctx, sess.AccountType, "PHONE"),
		PasswordExpiryWarningDays: 0,
	}, nil
}

// CancelAccount 注销（停用）当前账号并清理会话。
func (s *Service) CancelAccount(ctx context.Context, accountType security.AccountType, clientIP, userAgent string, reason *string) error {
	sess := contextx.Session(ctx)
	if sess == nil {
		return fmt.Errorf("unauthorized")
	}
	if s.accounts == nil {
		return errAccountFinder
	}
	accountID := sess.AccountID
	// 标记账号已取消（软注销）：更新 sys_account 状态
	now := time.Now().UTC()
	if err := s.db.WithContext(ctx).Table("sys_account").
		Where("id = ?", accountID).
		Update("account_status", security.AccountStatusCancelled).Error; err != nil {
		return err
	}
	_ = s.db.WithContext(ctx).Table("profile_user_admin").
		Where("account_id = ?", accountID).
		Update("remark", "cancelled at "+now.Format(time.RFC3339)).Error
	_ = s.db.WithContext(ctx).Table("profile_user_portal").
		Where("account_id = ?", accountID).
		Update("remark", "cancelled at "+now.Format(time.RFC3339)).Error
	_ = s.sessions.DeleteAllForAccount(ctx, accountID)
	s.publishAudit(ctx, "cancel", true, accountID, string(accountType), clientIP, userAgent, reasonString(reason))
	return nil
}

// sendRegisterCode 门户注册发送验证码。
func (s *Service) sendRegisterCode(ctx context.Context, req SendLoginCodeParam) error {
	if !s.registerEnabled(ctx) {
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
