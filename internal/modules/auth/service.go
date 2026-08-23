// internal/modules/auth/service.go 业务服务。
//
// Author: Charlie

package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/config"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/notify"
	"hei-gin/internal/framework/platform/runtimecfg"
	"hei-gin/internal/modules/auth/oauth"
)

// Service 认证服务。
//
// Author: Charlie
type Service struct {
	cfg            *config.Config
	db             *gorm.DB
	repo           *Repo
	sessions       *security.SessionStore
	accounts       AccountFinder
	notify         *notify.Facade
	audit          *audit.Queue
	runtime        *runtimecfg.Settings
	passwordPolicy *security.PasswordPolicy
	oauth          *oauth.Service
	perms          *security.PermissionRegistry
}

// NewService 构造认证服务。
func NewService(d *module.Deps, accounts AccountFinder) *Service {
	s := &Service{
		cfg:            d.Cfg,
		db:             d.DB,
		repo:           NewRepo(d.Redis),
		sessions:       d.Sessions,
		accounts:       accounts,
		notify:         d.Notify,
		audit:          d.Audit,
		runtime:        d.Runtime,
		passwordPolicy: security.NewPasswordPolicy(d.DB, d.Runtime),
		perms:          d.Perms,
	}
	s.oauth = oauth.NewService(d, func(ctx context.Context, accountType security.AccountType, accountID, clientIP, userAgent string, rememberMe bool) (oauth.LoginResult, error) {
		out, err := s.issueSession(ctx, accountType, accountID, clientIP, userAgent, rememberMe)
		if err != nil {
			return oauth.LoginResult{}, err
		}
		return oauth.LoginResult{
			Token: out.Token, AccountID: out.AccountID, AccountType: out.AccountType,
			PasswordExpired: out.PasswordExpired, ForceBindEmail: out.ForceBindEmail,
			ForceBindPhone: out.ForceBindPhone, PasswordExpiryWarningDays: out.PasswordExpiryWarningDays,
			ExpiresIn: out.ExpiresIn,
		}, nil
	}, oauth.OAuthHooks{
		AssignRegisterDefaults: func(ctx context.Context, accountID string) error {
			type defaults interface {
				AssignRegisterDefaults(ctx context.Context, accountID string, accountType security.AccountType) error
			}
			if reg, ok := accounts.(defaults); ok {
				return reg.AssignRegisterDefaults(ctx, accountID, security.AccountPortal)
			}
			return nil
		},
		RecordPasswordHistory: func(ctx context.Context, accountID, rawPassword, changedBy, reason string) error {
			return s.passwordPolicy.RecordHistory(ctx, accountID, rawPassword, changedBy, reason)
		},
	})
	return s
}

// New 构建认证模块。
func New(d *module.Deps, accounts AccountFinder) module.Module {
	s := NewService(d, accounts)
	models := []any{oauth.AccountOAuthBinding{}}
	return module.Module{
		Name:   "auth",
		Order:  10,
		Models: models,
		Routes: []module.RouteRegistrar{s.registerRoutes},
	}
}

// Captcha 生成验证码。
func (s *Service) Captcha(ctx context.Context) (*CaptchaResult, error) {
	ttl := time.Duration(s.cfg.Auth.CaptchaTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return s.repo.CreateCaptcha(ctx, ttl)
}

// PasswordKey 生成密码加密密钥。
func (s *Service) PasswordKey(ctx context.Context) (*PasswordKeyResult, error) {
	ttl := time.Duration(s.cfg.Auth.PasswordCryptoKeyTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return s.repo.CreatePasswordKey(ctx, ttl)
}

// Login 登录（密码或 OTP）。
func (s *Service) Login(ctx context.Context, accountType security.AccountType, req LoginParam, clientIP, userAgent string) (*LoginResult, error) {
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return nil, err
	}
	if err := s.repo.EnsureLoginAllowed(ctx, accountType, req.Account, clientIP); err != nil {
		return nil, err
	}
	identityType := strings.TrimSpace(req.IdentityType)
	if identityType == "" {
		identityType = "ACCOUNT"
	}
	loginMode := strings.ToUpper(strings.TrimSpace(req.LoginMode))
	if loginMode == "" {
		if strings.TrimSpace(req.OTPCode) != "" {
			loginMode = "OTP"
		} else {
			loginMode = "PASSWORD"
		}
	}
	if loginMode == "OTP" || strings.TrimSpace(req.OTPCode) != "" {
		switch strings.ToUpper(identityType) {
		case "EMAIL", "PHONE":
		default:
			if strings.Contains(req.Account, "@") {
				identityType = "EMAIL"
			} else {
				identityType = "PHONE"
			}
		}
	}

	var (
		accountID string
		err       error
	)
	defer func() {
		if err != nil {
			s.repo.RecordLoginFailure(ctx, s.protectCfg(ctx), accountType, req.Account, clientIP)
		}
	}()

	if loginMode == "OTP" || strings.TrimSpace(req.OTPCode) != "" {
		channel := identityChannel(identityType)
		target := normalizeAccount(req.Account)
		if !s.repo.ConsumeLoginOTP(ctx, string(accountType), channel, target, req.OTPCode) {
			err = errInvalidOTP
			return nil, err
		}
		if s.accounts == nil {
			err = errAccountFinder
			return nil, err
		}
		accountID, _, err = s.accounts.FindEnabledByIdentity(ctx, accountType, identityType, strings.TrimSpace(req.Account))
		if err != nil || accountID == "" {
			err = errInvalidCredentials
			return nil, err
		}
	} else {
		password, derr := s.repo.DecryptPassword(ctx, req.PasswordKeyID, req.Password)
		if derr != nil {
			err = derr
			return nil, err
		}
		if password == "" {
			err = errEmptyPassword
			return nil, err
		}
		if s.accounts == nil {
			err = errAccountFinder
			return nil, err
		}
		var hash string
		accountID, hash, err = s.accounts.FindEnabledByIdentity(ctx, accountType, identityType, strings.TrimSpace(req.Account))
		if err != nil || accountID == "" {
			err = errInvalidCredentials
			return nil, err
		}
		if !security.CheckPassword(hash, password) {
			err = errInvalidCredentials
			return nil, err
		}
	}

	out, serr := s.issueSession(ctx, accountType, accountID, clientIP, userAgent, rememberMeOrDefault(req.RememberMe))
	if serr != nil {
		err = serr
		return nil, err
	}
	err = nil
	s.repo.ClearLoginFailures(ctx, accountType, req.Account, clientIP)
	return out, nil
}

func (s *Service) issueSession(ctx context.Context, accountType security.AccountType, accountID, clientIP, userAgent string, rememberMe bool) (*LoginResult, error) {
	if s.accounts == nil {
		return nil, errAccountFinder
	}
	authz, err := s.accounts.GetSessionAuthorization(ctx, accountID)
	if err != nil {
		return nil, err
	}
	token, err := security.NewToken()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	ttlSec := s.runtimeInt(ctx, "AUTH_TOKEN_TTL_SECONDS", s.cfg.Auth.TokenTTLSeconds)
	if ttlSec <= 0 {
		ttlSec = 14400
	}
	ttl := time.Duration(ttlSec) * time.Second
	passwordExpired, warningDays := s.passwordPolicy.PasswordExpired(ctx, accountID)
	forceEmail, forcePhone := s.forceBindFlags(ctx, accountType, accountID)
	payload := &security.SessionPayload{
		Token:                token,
		AccountID:            accountID,
		AccountType:          accountType,
		RoleIDs:              authz.RoleIDs,
		DeptIDs:              authz.DeptIDs,
		GroupIDs:             authz.GroupIDs,
		ResourceIDs:          nil, // 菜单/按钮资源不进会话（对齐 hei-boot issueSession）
		PermissionKeys:       authz.PermissionKeys,
		PermissionGrants:     authz.PermissionGrants,
		ClientResourceIDs:    authz.ClientResourceIDs,
		ClientPermissionKeys: authz.ClientPermissionKeys,
		ClientIP:             &clientIP,
		UserAgent:            &userAgent,
		DeviceLabel:          security.DeviceLabelFromUserAgent(userAgent),
		RememberMe:           rememberMe,
		PasswordExpired:      passwordExpired,
		LoginAt:              now,
		LastActiveAt:         now,
		ExpiresAt:            now.Add(ttl),
	}
	if err := s.sessions.Save(ctx, payload, ttl); err != nil {
		return nil, err
	}
	s.updateLoginMeta(ctx, accountID, clientIP, userAgent)
	s.maybeNotifyPasswordExpiring(ctx, accountID)
	return &LoginResult{
		Token:                     token,
		AccountID:                 accountID,
		AccountType:               accountType,
		PasswordExpired:           passwordExpired,
		ForceBindEmail:            forceEmail,
		ForceBindPhone:            forcePhone,
		PasswordExpiryWarningDays: warningDays,
		ExpiresIn:                 ttlSec,
	}, nil
}

func (s *Service) updateLoginMeta(ctx context.Context, accountID, clientIP, userAgent string) {
	device := security.DeviceLabelFromUserAgent(userAgent)
	now := time.Now().UTC()
	_ = s.db.WithContext(ctx).Table("sys_account").Where("id = ?", accountID).Updates(map[string]any{
		"latest_login_ip":     clientIP,
		"latest_login_time":   now,
		"latest_login_device": device,
	}).Error
}

func (s *Service) maybeNotifyPasswordExpiring(ctx context.Context, accountID string) {
	if s.notify == nil || s.accounts == nil {
		return
	}
	_, warningDays := s.passwordPolicy.PasswordExpired(ctx, accountID)
	if warningDays <= 0 {
		return
	}
	if !s.repo.TryMarkPasswordExpiryNotified(ctx, accountID) {
		return
	}
	accountName := accountID
	var emailIdent, phoneIdent string
	_ = s.db.WithContext(ctx).Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ? AND bind_status = ?", accountID, "EMAIL", "BOUND").
		Select("identifier").Scan(&emailIdent)
	_ = s.db.WithContext(ctx).Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ? AND bind_status = ?", accountID, "ACCOUNT", "BOUND").
		Select("identifier").Scan(&accountName)
	_ = s.db.WithContext(ctx).Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ? AND bind_status = ?", accountID, "PHONE", "BOUND").
		Select("identifier").Scan(&phoneIdent)
	vars := map[string]any{
		"app_name":       s.cfg.App.Name,
		"account":        accountName,
		"remaining_days": fmt.Sprintf("%d", warningDays),
	}
	if emailIdent != "" {
		_ = s.notify.SendTemplated(ctx, "PASSWORD_EXPIRING", emailIdent, vars)
	}
	if phoneIdent != "" {
		_ = s.notify.SendTemplated(ctx, "PASSWORD_EXPIRING", phoneIdent, vars)
	}
}

func rememberMeOrDefault(v *bool) bool {
	if v == nil {
		return true
	}
	return *v
}

func (s *Service) forceBindFlags(ctx context.Context, accountType security.AccountType, accountID string) (email, phone bool) {
	if !s.forceBind(ctx, accountType, "EMAIL") && !s.forceBind(ctx, accountType, "PHONE") {
		return false, false
	}
	if s.accounts == nil {
		return s.forceBind(ctx, accountType, "EMAIL"), s.forceBind(ctx, accountType, "PHONE")
	}
	hasEmail, hasPhone := s.accounts.HasBoundIdentity(ctx, accountID, "EMAIL"), s.accounts.HasBoundIdentity(ctx, accountID, "PHONE")
	return s.forceBind(ctx, accountType, "EMAIL") && !hasEmail, s.forceBind(ctx, accountType, "PHONE") && !hasPhone
}

// registerEnabled 门户注册开关：优先运行时配置 AUTH_REGISTER_PORTAL_ENABLED，回退 yaml。
func (s *Service) registerEnabled(ctx context.Context) bool {
	return s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_ENABLED", s.cfg.Auth.PortalRegisterEnabled)
}

// forceBind 读取 AUTH_FORCE_BIND_{TYPE}_{CHANNEL} 运行时配置。
func (s *Service) forceBind(ctx context.Context, accountType security.AccountType, channel string) bool {
	key := "AUTH_FORCE_BIND_" + strings.ToUpper(string(accountType)) + "_" + channel
	return s.runtimeBool(ctx, key, false)
}

// Logout 登出（审计由 middleware 按注册表记录）。
func (s *Service) Logout(ctx context.Context, token, accountID, accountType, clientIP, userAgent string) error {
	var err error
	if token != "" {
		err = s.sessions.Delete(ctx, token)
	}
	return err
}

// SendLoginCode 发送登录 OTP。
func (s *Service) SendLoginCode(ctx context.Context, accountType security.AccountType, req SendLoginCodeParam) error {
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return err
	}
	channel, target := resolveOTPTarget(req)
	if target == "" {
		return errOTPTargetRequired
	}
	identityType := "EMAIL"
	if channel == "PHONE" {
		identityType = "PHONE"
	}
	if s.accounts != nil {
		_, _, err := s.accounts.FindEnabledByIdentity(ctx, accountType, identityType, target)
		if err != nil {
			// 静默返回，避免枚举账号
			return nil
		}
	}
	code, err := sixDigitCode()
	if err != nil {
		return err
	}
	ttl := 5 * time.Minute
	if err := s.repo.StoreLoginOTP(ctx, string(accountType), channel, target, code, ttl); err != nil {
		return err
	}
	if s.notify != nil {
		vars := map[string]any{
			"app_name":       s.cfg.App.Name,
			"code":           code,
			"expire_minutes": "5",
		}
		_ = s.notify.SendTemplated(ctx, "LOGIN_CODE", target, vars)
	}
	return nil
}

// ForgotPassword 发送重置邮件。
func (s *Service) ForgotPassword(ctx context.Context, accountType security.AccountType, req ForgotPasswordParam) error {
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return err
	}
	email := normalizeAccount(req.Email)
	if s.accounts == nil {
		return nil
	}
	accountID, _, err := s.accounts.FindEnabledByIdentity(ctx, accountType, "EMAIL", email)
	if err != nil || accountID == "" {
		return nil
	}
	token, err := newResetToken()
	if err != nil {
		return err
	}
	ttlSeconds := s.runtimeInt(ctx, "AUTH_PASSWORD_RESET_TOKEN_TTL_SECONDS", 600)
	if ttlSeconds < 60 {
		ttlSeconds = 60
	}
	ttl := time.Duration(ttlSeconds) * time.Second
	if err := s.repo.StoreResetToken(ctx, token, accountID, ttl); err != nil {
		return err
	}
	resetLink, err := s.buildPasswordResetLink(ctx, token, accountType)
	if err != nil {
		return err
	}
	if s.notify != nil {
		vars := map[string]any{
			"app_name":       s.runtimeString(ctx, "COPYRIGHT_TEXT", "HEI"),
			"reset_link":     resetLink,
			"email":          email,
			"expire_minutes": strconv.Itoa(max(1, ttlSeconds/60)),
		}
		_ = s.notify.SendTemplated(ctx, "RESET_PASSWORD_CODE", email, vars)
	}
	return nil
}

// buildPasswordResetLink 按账户类型读取重置页 URL 并拼接 token（对齐 hei-boot）。
func (s *Service) buildPasswordResetLink(ctx context.Context, token string, accountType security.AccountType) (string, error) {
	key := "AUTH_PASSWORD_RESET_URL_" + strings.ToUpper(string(accountType))
	base := strings.TrimSpace(s.runtimeString(ctx, key, ""))
	if base == "" {
		return "", fmt.Errorf("缺少系统配置: %s", key)
	}
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	return base + sep + "token=" + url.QueryEscape(token), nil
}

// ResetPassword 校验令牌并设置新密码。
func (s *Service) ResetPassword(ctx context.Context, accountType security.AccountType, req ResetPasswordParam) error {
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return err
	}
	password, err := s.repo.DecryptPassword(ctx, req.PasswordKeyID, req.Password)
	if err != nil {
		return err
	}
	if password == "" {
		return errEmptyPassword
	}
	accountID, err := s.repo.ConsumeResetToken(ctx, strings.TrimSpace(req.Token))
	if err != nil {
		return err
	}
	if s.accounts == nil {
		return errAccountFinder
	}
	gotType, err := s.accounts.GetEnabledAccount(ctx, accountID)
	if err != nil || gotType != accountType {
		return errResetTokenInvalid
	}
	if err := s.passwordPolicy.Validate(ctx, password, accountID, "", "", ""); err != nil {
		return err
	}
	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	if err := s.accounts.UpdatePasswordHash(ctx, accountID, hash); err != nil {
		return err
	}
	return s.passwordPolicy.RecordHistory(ctx, accountID, password, accountID, "self_reset")
}

// ForgotPasswordByPhone 向绑定手机发送重置 OTP。
func (s *Service) ForgotPasswordByPhone(ctx context.Context, accountType security.AccountType, req ForgotPasswordByPhoneParam) error {
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return err
	}
	phone := normalizePhone(req.Phone)
	if phone == "" {
		return errOTPTargetRequired
	}
	if s.accounts == nil {
		return nil
	}
	accountID, _, err := s.accounts.FindEnabledByIdentity(ctx, accountType, "PHONE", phone)
	if err != nil || accountID == "" {
		return nil
	}
	code, err := sixDigitCode()
	if err != nil {
		return err
	}
	ttl := s.otpTTL(ctx)
	if err := s.repo.StoreResetPasswordOTP(ctx, string(accountType), phone, code, ttl); err != nil {
		return err
	}
	if s.notify != nil {
		vars := map[string]any{
			"app_name":       s.cfg.App.Name,
			"code":           code,
			"expire_minutes": strconv.Itoa(max(1, int(ttl.Seconds())/60)),
		}
		_ = s.notify.SendTemplated(ctx, "RESET_PASSWORD_CODE", phone, vars)
	}
	return nil
}

// ResetPasswordByPhone 通过手机 OTP 重置密码。
func (s *Service) ResetPasswordByPhone(ctx context.Context, accountType security.AccountType, req ResetPasswordByPhoneParam) error {
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return err
	}
	phone := normalizePhone(req.Phone)
	if !s.repo.ConsumeResetPasswordOTP(ctx, string(accountType), phone, req.OTPCode) {
		return errInvalidOTP
	}
	if s.accounts == nil {
		return errAccountFinder
	}
	accountID, _, err := s.accounts.FindEnabledByIdentity(ctx, accountType, "PHONE", phone)
	if err != nil || accountID == "" {
		return errAccountNotFound
	}
	password, err := s.repo.DecryptPassword(ctx, req.PasswordKeyID, req.Password)
	if err != nil {
		return err
	}
	if password == "" {
		return errEmptyPassword
	}
	if err := s.passwordPolicy.Validate(ctx, password, accountID, "", "", phone); err != nil {
		return err
	}
	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	if err := s.accounts.UpdatePasswordHash(ctx, accountID, hash); err != nil {
		return err
	}
	return s.passwordPolicy.RecordHistory(ctx, accountID, password, accountID, "self_reset_phone")
}

// SiteFooter 解析站点页脚公开配置。
func (s *Service) SiteFooter(ctx context.Context) *SiteFooterResult {
	return resolveSiteFooter(ctx, s)
}

func resolveSiteFooter(ctx context.Context, s *Service) *SiteFooterResult {
	return &SiteFooterResult{
		CopyrightText: s.runtimeString(ctx, "COPYRIGHT_TEXT", ""),
		CopyrightURL:  s.runtimeString(ctx, "COPYRIGHT_URL", ""),
		IcpNumber:     s.runtimeString(ctx, "SITE_ICP_NUMBER", ""),
		IcpURL:        s.runtimeString(ctx, "SITE_ICP_URL", ""),
		PsbNumber:     s.runtimeString(ctx, "SITE_PSB_NUMBER", ""),
		PsbURL:        s.runtimeString(ctx, "SITE_PSB_URL", ""),
	}
}

func (s *Service) otpTTL(ctx context.Context) time.Duration {
	sec := s.runtimeInt(ctx, "AUTH_OTP_TTL_SECONDS", 300)
	if sec <= 0 {
		sec = 300
	}
	return time.Duration(sec) * time.Second
}

// Register 门户注册（ACCOUNT / EMAIL / PHONE 三通道，对齐 hei-boot registerPortal）。
func (s *Service) Register(ctx context.Context, req RegisterParam) (*RegisterResult, error) {
	if !s.registerEnabled(ctx) {
		return nil, errRegisterDisabled
	}
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return nil, err
	}
	channel := strings.ToUpper(strings.TrimSpace(req.RegisterChannel))
	if channel == "" {
		return nil, fmt.Errorf("请选择注册通道")
	}
	if channel != "ACCOUNT" && channel != "EMAIL" && channel != "PHONE" {
		return nil, fmt.Errorf("不支持的注册通道")
	}
	if err := s.ensureRegisterChannelAllowed(ctx, channel); err != nil {
		return nil, err
	}

	var accountLogin string
	var email, phone *string
	switch channel {
	case "ACCOUNT":
		login, err := security.RequireAccountLogin(req.Account)
		if err != nil {
			return nil, err
		}
		accountLogin = login
		if s.identityExists(ctx, IdentityAccountConst, accountLogin) {
			return nil, fmt.Errorf("账号已存在")
		}
		if strings.TrimSpace(req.Email) != "" {
			e := strings.ToLower(strings.TrimSpace(req.Email))
			email = &e
		}
		if strings.TrimSpace(req.Phone) != "" {
			p := strings.TrimSpace(req.Phone)
			phone = &p
		}
		if s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_REQUIRE_EMAIL", true) && email == nil {
			return nil, fmt.Errorf("注册必填邮箱")
		}
		if s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_REQUIRE_PHONE", false) && phone == nil {
			return nil, fmt.Errorf("注册必填手机号")
		}
	case "EMAIL":
		e := strings.ToLower(strings.TrimSpace(req.Email))
		if e == "" || !strings.Contains(e, "@") {
			return nil, fmt.Errorf("邮箱格式不正确")
		}
		if err := s.consumeRegisterCode(ctx, channel, e, req.OTPCode); err != nil {
			return nil, err
		}
		if s.identityExists(ctx, IdentityEmailConst, e) {
			return nil, fmt.Errorf("邮箱已被使用")
		}
		email = &e
		base := allocateBaseFromEmail(e)
		allocated, err := s.allocateAccountLogin(ctx, base)
		if err != nil {
			return nil, err
		}
		accountLogin = allocated
	case "PHONE":
		p := strings.TrimSpace(req.Phone)
		if p == "" {
			return nil, fmt.Errorf("手机号不能为空")
		}
		if err := s.consumeRegisterCode(ctx, channel, p, req.OTPCode); err != nil {
			return nil, err
		}
		if s.identityExists(ctx, IdentityPhoneConst, p) {
			return nil, fmt.Errorf("手机号已被使用")
		}
		phone = &p
		base := allocateBaseFromPhone(p)
		allocated, err := s.allocateAccountLogin(ctx, base)
		if err != nil {
			return nil, err
		}
		accountLogin = allocated
	}

	password, err := s.repo.DecryptPassword(ctx, req.PasswordKeyID, req.Password)
	if err != nil {
		return nil, err
	}
	emailStr, phoneStr := "", ""
	if email != nil {
		emailStr = *email
	}
	if phone != nil {
		phoneStr = *phone
	}
	if err := s.passwordPolicy.Validate(ctx, password, "", accountLogin, emailStr, phoneStr); err != nil {
		return nil, err
	}
	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}
	reg, ok := s.accounts.(PortalRegistrar)
	if !ok || reg == nil {
		return nil, errPortalRegistrar
	}
	var nickname *string
	if strings.TrimSpace(req.Nickname) != "" {
		n := strings.TrimSpace(req.Nickname)
		nickname = &n
	}
	in := registerPortalInput(accountLogin, hash, email, phone, nickname)
	accountID, registeredLogin, err := reg.RegisterPortal(ctx, in)
	if err != nil {
		return nil, err
	}
	_ = s.passwordPolicy.RecordHistory(ctx, accountID, password, accountID, "register")
	s.sendRegisterSuccessNotify(ctx, registeredLogin, email, phone)
	return &RegisterResult{
		AccountID:   accountID,
		Account:     registeredLogin,
		AccountType: security.AccountPortal,
	}, nil
}

func (s *Service) sendRegisterSuccessNotify(ctx context.Context, accountLogin string, email, phone *string) {
	if s.notify == nil {
		return
	}
	vars := map[string]any{"app_name": s.cfg.App.Name, "account": accountLogin}
	if email != nil && strings.TrimSpace(*email) != "" {
		_ = s.notify.SendTemplated(ctx, "REGISTER_SUCCESS", strings.TrimSpace(*email), vars)
	}
	if phone != nil && strings.TrimSpace(*phone) != "" {
		_ = s.notify.SendTemplated(ctx, "REGISTER_SUCCESS", strings.TrimSpace(*phone), vars)
	}
}

const (
	IdentityAccountConst = "ACCOUNT"
	IdentityEmailConst   = "EMAIL"
	IdentityPhoneConst   = "PHONE"
)

func (s *Service) ensureRegisterChannelAllowed(ctx context.Context, channel string) error {
	switch channel {
	case "ACCOUNT":
		if !s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_ALLOW_ACCOUNT", true) {
			return fmt.Errorf("用户名注册已关闭")
		}
	case "EMAIL":
		if !s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_ALLOW_EMAIL", true) {
			return fmt.Errorf("邮箱注册已关闭")
		}
	case "PHONE":
		if !s.runtimeBool(ctx, "AUTH_REGISTER_PORTAL_ALLOW_PHONE", false) {
			return fmt.Errorf("手机注册已关闭")
		}
	}
	return nil
}

// ResolveLogoutToken 解析登出 token。
func (s *Service) ResolveLogoutToken(c *gin.Context) string {
	sess := contextx.Session(c.Request.Context())
	if sess != nil {
		return sess.Token
	}
	name := s.sessionCookieName()
	if ck, err := c.Cookie(name); err == nil {
		return ck
	}
	return c.GetHeader(name)
}

// SetSessionCookie 设置会话 Cookie。
func (s *Service) SetSessionCookie(c *gin.Context, token string, accountType security.AccountType, ttl time.Duration, remember bool) {
	if !s.cfg.Auth.SessionCookieEnabled {
		return
	}
	name := s.sessionCookieName()
	path := security.SessionCookiePath(accountType)
	maxAge := int(ttl.Seconds())
	if !remember {
		maxAge = 0
	}
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(s.cfg.Auth.SessionCookieSameSite) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   s.cfg.Auth.SessionCookieSecure,
		SameSite: sameSite,
	})
}

// ClearSessionCookie 清除会话 Cookie。
func (s *Service) ClearSessionCookie(c *gin.Context, accountType security.AccountType) {
	if !s.cfg.Auth.SessionCookieEnabled {
		return
	}
	name := s.sessionCookieName()
	path := "/"
	if accountType != "" {
		path = security.SessionCookiePath(accountType)
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func (s *Service) sessionCookieName() string {
	name := s.cfg.Auth.SessionCookieName
	if name == "" {
		name = s.cfg.Auth.TokenName
	}
	if name == "" {
		name = "Authorization"
	}
	return name
}

// protectCfg 登录保护参数：优先运行时配置 AUTH_LOGIN_*（配置页可改），回退 yaml。
func (s *Service) protectCfg(ctx context.Context) loginProtectCfg {
	return loginProtectCfg{
		WindowSeconds: s.runtimeInt(ctx, "AUTH_LOGIN_FAILURE_WINDOW_SECONDS", s.cfg.Auth.LoginFailureWindowSeconds),
		AccountMax:    s.runtimeInt(ctx, "AUTH_LOGIN_ACCOUNT_MAX_FAILURES", s.cfg.Auth.LoginAccountMaxFailures),
		IPMax:         s.runtimeInt(ctx, "AUTH_LOGIN_IP_MAX_FAILURES", s.cfg.Auth.LoginIPMaxFailures),
		LockSeconds:   s.runtimeInt(ctx, "AUTH_LOGIN_LOCK_SECONDS", s.cfg.Auth.LoginLockSeconds),
	}
}

// runtimeInt 运行时整型配置（sys_config 优先，yaml 回退）。
func (s *Service) runtimeInt(ctx context.Context, key string, def int) int {
	if s.runtime != nil {
		return s.runtime.GetInt(ctx, key, def)
	}
	return def
}

// runtimeString 运行时字符串配置（sys_config 优先，yaml 回退）。
func (s *Service) runtimeString(ctx context.Context, key, def string) string {
	if s.runtime != nil {
		return s.runtime.GetString(ctx, key, def)
	}
	return def
}

// runtimeBool 运行时布尔配置（sys_config 优先，yaml 回退）。
func (s *Service) runtimeBool(ctx context.Context, key string, def bool) bool {
	if s.runtime != nil {
		return s.runtime.GetBool(ctx, key, def)
	}
	return def
}

func resolveOTPTarget(req SendLoginCodeParam) (channel, target string) {
	channel = strings.ToUpper(strings.TrimSpace(req.Channel))
	target = strings.TrimSpace(req.Target)
	if target == "" {
		if e := strings.TrimSpace(req.Email); e != "" {
			target, channel = e, "EMAIL"
		} else if p := strings.TrimSpace(req.Phone); p != "" {
			target, channel = p, "PHONE"
		} else if a := strings.TrimSpace(req.Account); a != "" {
			target = a
			if strings.Contains(a, "@") {
				channel = "EMAIL"
			} else {
				channel = "PHONE"
			}
		}
	}
	if channel == "" {
		if strings.Contains(target, "@") {
			channel = "EMAIL"
		} else {
			channel = "PHONE"
		}
	}
	target = normalizeAccount(target)
	return channel, target
}

func identityChannel(identityType string) string {
	switch strings.ToUpper(strings.TrimSpace(identityType)) {
	case "PHONE":
		return "PHONE"
	default:
		return "EMAIL"
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

var (
	errEmptyPassword      = authErr{"密码不能为空"}
	errAccountFinder      = authErr{"账号查找未配置"}
	errInvalidCredentials = authErr{"账号或密码错误"}
	errInvalidOTP         = authErr{"验证码无效或已过期"}
	errRegisterDisabled   = authErr{"门户注册已关闭"}
	errPortalRegistrar    = authErr{"portal registrar not configured"}
	errAccountLocked      = authErr{"账号已临时锁定"}
	errIPLocked           = authErr{"该 IP 登录失败次数过多"}
	errOTPTargetRequired  = authErr{"请提供邮箱或手机号"}
	errResetTokenInvalid  = authErr{"重置令牌无效"}
	errAccountNotFound    = authErr{"账号不存在"}
)

type authErr struct{ msg string }

func (e authErr) Error() string { return e.msg }
