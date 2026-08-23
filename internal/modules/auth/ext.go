// internal/modules/auth/ext.go 认证扩展。
//
// Author: Charlie

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/modules/auth/oauth"
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
	RegisterRequireEmail       bool                   `json:"register_require_email"`
	RegisterRequirePhone       bool                   `json:"register_require_phone"`
	ForceBindEmail             bool                   `json:"force_bind_email"`
	ForceBindPhone             bool                   `json:"force_bind_phone"`
	OAuthProviders             []oauth.ProviderOption `json:"oauth_providers"`
	PasswordChangeVerifyMethod string                 `json:"password_change_verify_method"`
	CopyrightText              string                 `json:"copyright_text"`
	CopyrightURL               string                 `json:"copyright_url"`
	SiteFooter                 *SiteFooterResult      `json:"site_footer"`
}

// CancelParam 注销账号入参。
//
// Author: Charlie
type CancelParam struct {
	CancelReason *string `json:"cancel_reason"`
}

// AuthOptions 读取登录页公开配置。
func (s *Service) AuthOptions(ctx context.Context, accountType security.AccountType) *AuthOptions {
	typeName := strings.ToUpper(string(accountType))
	defaultRegister := accountType == security.AccountPortal
	o := &AuthOptions{
		AccountType:                string(accountType),
		AllowAccount:               true,
		AllowEmail:                 s.runtimeBool(ctx, "AUTH_LOGIN_"+typeName+"_ALLOW_EMAIL", true),
		AllowPhone:                 s.runtimeBool(ctx, "AUTH_LOGIN_"+typeName+"_ALLOW_PHONE", true),
		AllowOTP:                   s.runtimeBool(ctx, "AUTH_LOGIN_"+typeName+"_ALLOW_OTP", true),
		RegisterEnabled:            s.runtimeBool(ctx, "AUTH_REGISTER_"+typeName+"_ENABLED", defaultRegister),
		RegisterAllowAccount:       accountType == security.AccountPortal && s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_ALLOW_ACCOUNT", true),
		RegisterAllowEmail:         accountType == security.AccountPortal && s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_ALLOW_EMAIL", true),
		RegisterAllowPhone:         accountType == security.AccountPortal && s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_ALLOW_PHONE", false),
		RegisterRequireEmail:       s.runtimeBool(ctx, "AUTH_REGISTER_"+typeName+"_REQUIRE_EMAIL", accountType == security.AccountPortal),
		RegisterRequirePhone:       s.runtimeBool(ctx, "AUTH_REGISTER_"+typeName+"_REQUIRE_PHONE", false),
		ForceBindEmail:             s.runtimeBool(ctx, "AUTH_FORCE_BIND_"+typeName+"_EMAIL", false),
		ForceBindPhone:             s.runtimeBool(ctx, "AUTH_FORCE_BIND_"+typeName+"_PHONE", false),
		PasswordChangeVerifyMethod: s.runtimeString(ctx, "PASSWORD_CHANGE_VERIFY_METHOD", "OLD_PASSWORD"),
		CopyrightText:              s.runtimeString(ctx, "COPYRIGHT_TEXT", ""),
		CopyrightURL:               s.runtimeString(ctx, "COPYRIGHT_URL", ""),
		SiteFooter:                 resolveSiteFooter(ctx, s),
	}
	if s.oauth != nil {
		o.OAuthProviders = s.oauth.ProviderOptions(ctx, accountType)
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
	ttlSec := s.runtimeInt(ctx, "AUTH_TOKEN_TTL_SECONDS", s.cfg.Auth.TokenTTLSeconds)
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
	passwordExpired, warningDays := s.passwordPolicy.PasswordExpired(ctx, sess.AccountID)
	sess.PasswordExpired = passwordExpired
	forceEmail, forcePhone := s.forceBindFlags(ctx, sess.AccountType, sess.AccountID)
	return &LoginResult{
		Token:                     sess.Token,
		AccountID:                 sess.AccountID,
		AccountType:               sess.AccountType,
		PasswordExpired:           passwordExpired,
		ForceBindEmail:            forceEmail,
		ForceBindPhone:            forcePhone,
		PasswordExpiryWarningDays: warningDays,
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
	_ = s.sessions.DeleteAllForAccountAnyType(ctx, accountID)
	return nil
}

// sendRegisterCode 门户注册发送验证码（EMAIL/PHONE 通道）。
func (s *Service) sendRegisterCode(ctx context.Context, req SendLoginCodeParam) error {
	if !s.registerEnabled(ctx) {
		return errRegisterDisabled
	}
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return err
	}
	channel := strings.ToUpper(strings.TrimSpace(req.Channel))
	if channel == "" {
		channel, _ = resolveOTPTarget(req)
	}
	if channel != "EMAIL" && channel != "PHONE" {
		return fmt.Errorf("不支持的注册通道")
	}
	target := strings.TrimSpace(req.Target)
	if target == "" {
		if channel == "EMAIL" {
			target = strings.TrimSpace(req.Email)
		} else {
			target = strings.TrimSpace(req.Phone)
		}
	}
	if target == "" {
		return errOTPTargetRequired
	}
	if err := s.ensureRegisterChannelAllowed(ctx, channel); err != nil {
		return err
	}
	normTarget := normalizeRegisterTarget(channel, target)
	identityType := IdentityEmailConst
	if channel == "PHONE" {
		identityType = IdentityPhoneConst
	}
	if s.identityExists(ctx, identityType, normTarget) {
		if channel == "EMAIL" {
			return fmt.Errorf("邮箱已被使用")
		}
		return fmt.Errorf("手机号已被使用")
	}
	code, err := sixDigitCode()
	if err != nil {
		return err
	}
	ttl := s.otpTTL(ctx)
	if err := s.repo.StoreRegisterOTP(ctx, channel, normTarget, code, ttl); err != nil {
		return err
	}
	if s.notify != nil {
		expireMin := max(1, int(ttl/time.Minute))
		vars := map[string]any{"app_name": s.cfg.App.Name, "code": code, "expire_minutes": expireMin}
		if channel == "PHONE" {
			_ = s.notify.SendTemplated(ctx, "LOGIN_CODE", normTarget, vars)
		} else {
			_ = s.notify.SendTemplated(ctx, "LOGIN_CODE_MAIL", normTarget, vars)
		}
	}
	return nil
}
