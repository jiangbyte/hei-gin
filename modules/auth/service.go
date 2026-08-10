package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	contextx "hei-gin/framework/core/context"
	"hei-gin/framework/core/config"
	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 认证服务。
//
// Author: Charlie
type Service struct {
	cfg      *config.Config
	repo     *Repo
	sessions *security.SessionStore
	accounts AccountFinder
}

// NewService 构造认证服务。
func NewService(d *shared.Deps, accounts AccountFinder) *Service {
	return &Service{
		cfg:      d.Cfg,
		repo:     NewRepo(d.Redis),
		sessions: d.Sessions,
		accounts: accounts,
	}
}

// New 构建认证模块。
func New(d *shared.Deps, accounts AccountFinder) module.Module {
	s := NewService(d, accounts)
	return module.Module{
		Name:   "auth",
		Order:  10,
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

// Login 登录。
func (s *Service) Login(ctx context.Context, accountType security.AccountType, req LoginParam, clientIP, userAgent string) (*LoginResult, error) {
	if err := s.repo.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return nil, err
	}
	identityType := strings.TrimSpace(req.IdentityType)
	if identityType == "" {
		identityType = "ACCOUNT"
	}
	password, err := s.repo.DecryptPassword(ctx, req.PasswordKeyID, req.Password)
	if err != nil {
		return nil, err
	}
	if password == "" {
		return nil, errEmptyPassword
	}
	if s.accounts == nil {
		return nil, errAccountFinder
	}
	accountID, hash, err := s.accounts.FindEnabledByIdentity(ctx, accountType, identityType, strings.TrimSpace(req.Account))
	if err != nil || accountID == "" {
		return nil, errInvalidCredentials
	}
	if !security.CheckPassword(hash, password) {
		return nil, errInvalidCredentials
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
	if !req.RememberMe && s.cfg.Auth.TokenTTLShortSeconds > 0 {
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
		RememberMe:       req.RememberMe,
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

// Logout 登出。
func (s *Service) Logout(ctx context.Context, token string) error {
	if token != "" {
		return s.sessions.Delete(ctx, token)
	}
	return nil
}

// Register 门户注册。
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

var (
	errEmptyPassword      = authErr{"密码不能为空"}
	errAccountFinder      = authErr{"账号查找未配置"}
	errInvalidCredentials = authErr{"账号或密码错误"}
	errRegisterDisabled   = authErr{"门户注册已关闭"}
	errPortalRegistrar    = authErr{"portal registrar not configured"}
)

type authErr struct{ msg string }

func (e authErr) Error() string { return e.msg }
