// internal/modules/auth/service.go 业务服务。
//
// Author: Charlie

package auth

import (
	"context"
	"net/http"
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
	"hei-gin/internal/modules/auth/oauth"
	"hei-gin/internal/modules/shared"
)

// Service è®¤è¯æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct {
	cfg      *config.Config
	db       *gorm.DB
	repo     *Repo
	sessions *security.SessionStore
	accounts AccountFinder
	notify   *notify.Facade
	audit    *audit.Queue
	oauth    *oauth.Service
	perms    *security.PermissionRegistry
}

// NewService æž„é€ è®¤è¯æœåŠ¡ã€‚
func NewService(d *shared.Deps, accounts AccountFinder) *Service {
	s := &Service{
		cfg:      d.Cfg,
		db:       d.DB,
		repo:     NewRepo(d.Redis),
		sessions: d.Sessions,
		accounts: accounts,
		notify:   d.Notify,
		audit:    d.Audit,
		perms:    d.Perms,
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

// New æž„å»ºè®¤è¯æ¨¡å—ã€‚
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

// Captcha ç”ŸæˆéªŒè¯ç ã€‚
func (s *Service) Captcha(ctx context.Context) (*CaptchaResult, error) {
	ttl := time.Duration(s.cfg.Auth.CaptchaTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return s.repo.CreateCaptcha(ctx, ttl)
}

// PasswordKey ç”Ÿæˆå¯†ç åŠ å¯†å¯†é’¥ã€‚
func (s *Service) PasswordKey(ctx context.Context) (*PasswordKeyResult, error) {
	ttl := time.Duration(s.cfg.Auth.PasswordCryptoKeyTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return s.repo.CreatePasswordKey(ctx, ttl)
}

// Login ç™»å½•ï¼ˆå¯†ç æˆ– OTPï¼‰ã€‚
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
			s.repo.RecordLoginFailure(ctx, s.protectCfg(), accountType, req.Account, clientIP)
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
	ttlSec := s.cfg.Auth.TokenTTLSeconds
	if !rememberMe && s.cfg.Auth.TokenTTLShortSeconds > 0 {
		ttlSec = s.cfg.Auth.TokenTTLShortSeconds
	}
	if ttlSec <= 0 {
		ttlSec = 14400
	}
	ttl := time.Duration(ttlSec) * time.Second
	payload := &security.SessionPayload{
		Token:            token,
		AccountID:        accountID,
		AccountType:      accountType,
		PermissionKeys:   keys,
		PermissionGrants: grants,
		ClientIP:         &clientIP,
		UserAgent:        &userAgent,
		RememberMe:       rememberMe,
		PasswordExpired:  false,
		LoginAt:          now,
		LastActiveAt:     now,
		ExpiresAt:        now.Add(ttl),
	}
	if err := s.sessions.Save(ctx, payload, ttl); err != nil {
		return nil, err
	}
	return &LoginResult{
		Token:           token,
		AccountID:       accountID,
		AccountType:     accountType,
		PasswordExpired: false,
	}, nil
}

// Logout ç™»å‡ºã€‚
func (s *Service) Logout(ctx context.Context, token, accountID, accountType, clientIP, userAgent string) error {
	var err error
	if token != "" {
		err = s.sessions.Delete(ctx, token)
	}
	s.publishAudit(ctx, "logout", err == nil, accountID, accountType, clientIP, userAgent, errString(err))
	return err
}

// SendLoginCode å‘é€ç™»å½• OTPã€‚
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
			// é™é»˜è¿”å›žï¼Œé¿å…æžšä¸¾è´¦å·
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

// ForgotPassword å‘é€é‡ç½®é‚®ä»¶ã€‚
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
	ttl := 10 * time.Minute
	if err := s.repo.StoreResetToken(ctx, token, accountID, ttl); err != nil {
		return err
	}
	if s.notify != nil {
		vars := map[string]any{
			"app_name":       s.cfg.App.Name,
			"token":          token,
			"email":          email,
			"expire_minutes": "10",
			"reset_link":     "",
		}
		_ = s.notify.SendTemplated(ctx, "RESET_PASSWORD_CODE", email, vars)
	}
	return nil
}

// ResetPassword æ ¡éªŒä»¤ç‰Œå¹¶è®¾ç½®æ–°å¯†ç ã€‚
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
	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	return s.accounts.UpdatePasswordHash(ctx, accountID, hash)
}

// Register é—¨æˆ·æ³¨å†Œã€‚
func (s *Service) Register(ctx context.Context, req RegisterParam) (*RegisterResult, error) {
	if !s.cfg.Auth.PortalRegisterEnabled {
		return nil, errRegisterDisabled
	}
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return nil, err
	}
	password, err := s.repo.DecryptPassword(ctx, req.PasswordKeyID, req.Password)
	if err != nil {
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
	return &RegisterResult{
		AccountID:   accountID,
		Account:     req.Account,
		AccountType: security.AccountPortal,
	}, nil
}

// ResolveLogoutToken è§£æžç™»å‡º tokenã€‚
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

// SetSessionCookie è®¾ç½®ä¼šè¯ Cookieã€‚
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

// ClearSessionCookie æ¸…é™¤ä¼šè¯ Cookieã€‚
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

func (s *Service) protectCfg() loginProtectCfg {
	return loginProtectCfg{
		WindowSeconds: s.cfg.Auth.LoginFailureWindowSeconds,
		AccountMax:    s.cfg.Auth.LoginAccountMaxFailures,
		IPMax:         s.cfg.Auth.LoginIPMaxFailures,
		LockSeconds:   s.cfg.Auth.LoginLockSeconds,
	}
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
	errEmptyPassword      = authErr{"å¯†ç ä¸èƒ½ä¸ºç©º"}
	errAccountFinder      = authErr{"è´¦å·æŸ¥æ‰¾æœªé…ç½®"}
	errInvalidCredentials = authErr{"è´¦å·æˆ–å¯†ç é”™è¯¯"}
	errInvalidOTP         = authErr{"éªŒè¯ç æ— æ•ˆæˆ–å·²è¿‡æœŸ"}
	errRegisterDisabled   = authErr{"é—¨æˆ·æ³¨å†Œå·²å…³é—­"}
	errPortalRegistrar    = authErr{"portal registrar not configured"}
	errAccountLocked      = authErr{"è´¦å·å·²ä¸´æ—¶é”å®š"}
	errIPLocked           = authErr{"è¯¥ IP ç™»å½•å¤±è´¥æ¬¡æ•°è¿‡å¤š"}
	errOTPTargetRequired  = authErr{"è¯·æä¾›é‚®ç®±æˆ–æ‰‹æœºå·"}
	errResetTokenInvalid  = authErr{"é‡ç½®ä»¤ç‰Œæ— æ•ˆ"}
)

type authErr struct{ msg string }

func (e authErr) Error() string { return e.msg }
