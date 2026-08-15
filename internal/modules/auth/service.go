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
	"hei-gin/internal/modules/shared"
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
	passwordPolicy *shared.PasswordPolicy
	oauth          *oauth.Service
	perms          *security.PermissionRegistry
}

// NewService 构造认证服务。
func NewService(d *shared.Deps, accounts AccountFinder) *Service {
	s := &Service{
		cfg:            d.Cfg,
		db:             d.DB,
		repo:           NewRepo(d.Redis),
		sessions:       d.Sessions,
		accounts:       accounts,
		notify:         d.Notify,
		audit:          d.Audit,
		runtime:        d.Runtime,
		passwordPolicy: shared.NewPasswordPolicy(d.DB, d.Runtime),
		perms:          d.Perms,
	}
	s.oauth = oauth.NewService(d, func(ctx context.Context, accountType security.AccountType, accountID, clientIP, userAgent string, rememberMe bool) (string, error) {
		out, err := s.issueSession(ctx, accountType, accountID, clientIP, userAgent, rememberMe)
		if err != nil {
			return "", err
		}
		return out.Token, nil
	})
	return s
}

// New 构建认证模块。
func New(d *shared.Deps, accounts AccountFinder) module.Module {
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
			s.publishAudit(ctx, "login", false, "", string(accountType), clientIP, userAgent, err.Error())
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

	out, serr := s.issueSession(ctx, accountType, accountID, clientIP, userAgent, req.RememberMe)
	if serr != nil {
		err = serr
		return nil, err
	}
	err = nil
	s.repo.ClearLoginFailures(ctx, accountType, req.Account, clientIP)
	s.publishAudit(ctx, "login", true, accountID, string(accountType), clientIP, userAgent, "")
	return out, nil
}

func (s *Service) issueSession(ctx context.Context, accountType security.AccountType, accountID, clientIP, userAgent string, rememberMe bool) (*LoginResult, error) {
	if s.accounts == nil {
		return nil, errAccountFinder
	}
	keys, grants, err := s.accounts.EnsureSuperPermissions(ctx, accountID)
	if err != nil {
		return nil, err
	}
	token, err := security.NewToken()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	ttlSec := s.runtimeInt(ctx, "AUTH_TOKEN_TTL_SECONDS", s.cfg.Auth.TokenTTLSeconds)
	if !rememberMe && s.cfg.Auth.TokenTTLShortSeconds > 0 {
		ttlSec = s.cfg.Auth.TokenTTLShortSeconds
	}
	if ttlSec <= 0 {
		ttlSec = 14400
	}
	ttl := time.Duration(ttlSec) * time.Second
	passwordExpired, warningDays := s.passwordPolicy.PasswordExpired(ctx, accountID)
	payload := &security.SessionPayload{
		Token:            token,
		AccountID:        accountID,
		AccountType:      accountType,
		PermissionKeys:   keys,
		PermissionGrants: grants,
		ClientIP:         &clientIP,
		UserAgent:        &userAgent,
		RememberMe:       rememberMe,
		PasswordExpired:  passwordExpired,
		LoginAt:          now,
		LastActiveAt:     now,
		ExpiresAt:        now.Add(ttl),
	}
	if err := s.sessions.Save(ctx, payload, ttl); err != nil {
		return nil, err
	}
	return &LoginResult{
		Token:                     token,
		AccountID:                 accountID,
		AccountType:               accountType,
		PasswordExpired:           passwordExpired,
		ForceBindEmail:            s.forceBind(ctx, accountType, "EMAIL"),
		ForceBindPhone:            s.forceBind(ctx, accountType, "PHONE"),
		PasswordExpiryWarningDays: warningDays,
	}, nil
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

// Logout 登出。
func (s *Service) Logout(ctx context.Context, token, accountID, accountType, clientIP, userAgent string) error {
	var err error
	if token != "" {
		err = s.sessions.Delete(ctx, token)
	}
	s.publishAudit(ctx, "logout", err == nil, accountID, accountType, clientIP, userAgent, errString(err))
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

// Register 门户注册。
func (s *Service) Register(ctx context.Context, req RegisterParam) (*RegisterResult, error) {
	if !s.registerEnabled(ctx) {
		return nil, errRegisterDisabled
	}
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return nil, err
	}
	password, err := s.repo.DecryptPassword(ctx, req.PasswordKeyID, req.Password)
	if err != nil {
		return nil, err
	}
	if err := s.passwordPolicy.Validate(ctx, password, "", req.Account, req.Email, req.Phone); err != nil {
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
	var name, nickname, email, phone *string
	if req.Name != "" {
		name = &req.Name
	}
	if req.Nickname != "" {
		nickname = &req.Nickname
	}
	if req.Email != "" {
		email = &req.Email
	}
	if req.Phone != "" {
		phone = &req.Phone
	}
	accountID, err := reg.RegisterPortal(ctx, req.Account, hash, name, nickname, email, phone)
	if err != nil {
		return nil, err
	}
	_ = s.passwordPolicy.RecordHistory(ctx, accountID, password, accountID, "register")
	return &RegisterResult{
		AccountID:   accountID,
		Account:     req.Account,
		AccountType: security.AccountPortal,
	}, nil
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

func (s *Service) publishAudit(ctx context.Context, action string, success bool, accountID, accountType, ip, ua, errMsg string) {
	if s.audit == nil {
		return
	}
	s.audit.Publish(audit.Event{
		Module:       "auth",
		ResourceType: "auth",
		Action:       action,
		AccountID:    accountID,
		AccountType:  accountType,
		RequestID:    contextx.RequestID(ctx),
		IP:           ip,
		UserAgent:    ua,
		Success:      success,
		ErrorMessage: errMsg,
	})
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
)

type authErr struct{ msg string }

func (e authErr) Error() string { return e.msg }
