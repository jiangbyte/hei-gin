// Package oauth 提供三方登录（GitHub 完整；Gitee/微信等桩）。
package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/framework/core/bind"
	"hei-gin/framework/core/config"
	"hei-gin/framework/core/response"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
	"hei-gin/framework/platform/idgen"
	"hei-gin/modules/shared"
)

const (
	stateKeyPrefix    = "oauth:state:"
	exchangeKeyPrefix = "oauth:exchange:"
)

// IssueSessionFunc 由 auth 注入，避免循环依赖。
//
// Author: Charlie
type IssueSessionFunc func(ctx context.Context, accountType security.AccountType, accountID, clientIP, userAgent string, rememberMe bool) (token string, err error)

// Service OAuth 服务。
//
// Author: Charlie
type Service struct {
	cfg   *config.Config
	db    *gorm.DB
	rdb   *redis.Client
	issue IssueSessionFunc
}

// NewService 构造 OAuth 服务。
func NewService(d *shared.Deps, issue IssueSessionFunc) *Service {
	return &Service{
		cfg:   d.Cfg,
		db:    d.DB,
		rdb:   d.Redis,
		issue: issue,
	}
}

// ProviderOption 三方登录入口选项（供登录页公开配置）。
type ProviderOption struct {
	Provider string `json:"provider"`
	Label    string `json:"label"`
	Enabled  bool   `json:"enabled"`
	WebOAuth bool   `json:"web_oauth"`
}

// ProviderOptions 三方登录入口选项。
func (s *Service) ProviderOptions() []ProviderOption {
	out := []ProviderOption{}
	for _, p := range []string{"github", "gitee", "wechat", "wechat_mp", "qq"} {
		_, err := s.providerConfig(p)
		out = append(out, ProviderOption{
			Provider: p,
			Label:    providerLabel(p),
			Enabled:  err == nil,
			WebOAuth: p == "github" || p == "gitee" || p == "qq",
		})
	}
	return out
}

func providerLabel(provider string) string {
	switch provider {
	case "github":
		return "GitHub"
	case "gitee":
		return "Gitee"
	case "wechat":
		return "微信"
	case "wechat_mp":
		return "微信公众号"
	case "qq":
		return "QQ"
	default:
		return provider
	}
}

// RegisterRoutes 挂载 admin/portal OAuth 路由。
func (s *Service) RegisterRoutes(api *gin.RouterGroup) {
	rdb := s.rdb
	for _, prefix := range []struct {
		base string
		typ  security.AccountType
		rl   string
	}{
		{"/v1/admin/oauth", security.AccountAdmin, "admin:oauth"},
		{"/v1/portal/oauth", security.AccountPortal, "portal:oauth"},
	} {
		p := prefix
		api.GET(p.base+"/:provider/authorize", middleware.RateLimit(rdb, p.rl+"-authorize", 30, 60), s.authorize(p.typ))
		api.GET(p.base+"/:provider/callback", middleware.RateLimit(rdb, p.rl+"-callback", 30, 60), s.callback(p.typ))
		api.POST(p.base+"/exchange", middleware.RateLimit(rdb, p.rl+"-exchange", 30, 60), s.exchange(p.typ))
	}
}

// AuthorizeResult 授权跳转结果。
type AuthorizeResult struct {
	AuthorizeURL string `json:"authorize_url"`
	State        string `json:"state"`
}

// ExchangeParam 兑换登录码。
type ExchangeParam struct {
	Code string `json:"code" binding:"required"`
}

// LoginResult OAuth 登录结果（字段与 auth 对齐）。
type LoginResult struct {
	Token           string               `json:"token"`
	AccountID       string               `json:"account_id"`
	AccountType     security.AccountType `json:"account_type"`
	PasswordExpired bool                 `json:"password_expired"`
}

func (s *Service) authorize(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := strings.ToLower(c.Param("provider"))
		pc, err := s.providerConfig(provider)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		state, err := randomHex(16)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, 500, err.Error())
			return
		}
		payload, _ := json.Marshal(map[string]string{
			"account_type": string(accountType),
			"provider":     provider,
			"redirect":     c.Query("redirect"),
			"intent":       c.Query("intent"),
			"account_id":   c.Query("account_id"),
		})
		_ = s.rdb.Set(c.Request.Context(), stateKeyPrefix+state, string(payload), 10*time.Minute).Err()
		u, err := buildAuthorizeURL(provider, pc, state)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, AuthorizeResult{AuthorizeURL: u, State: state})
	}
}

func (s *Service) callback(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := strings.ToLower(c.Param("provider"))
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			response.Fail(c, http.StatusBadRequest, 400, "缺少 code 或 state")
			return
		}
		raw, err := s.rdb.Get(c.Request.Context(), stateKeyPrefix+state).Result()
		_ = s.rdb.Del(c.Request.Context(), stateKeyPrefix+state)
		if err != nil || raw == "" {
			response.Fail(c, http.StatusBadRequest, 400, "无效或过期的 state")
			return
		}
		var st map[string]string
		_ = json.Unmarshal([]byte(raw), &st)
		if st["account_type"] != "" && st["account_type"] != string(accountType) {
			response.Fail(c, http.StatusBadRequest, 400, "账号类型不匹配")
			return
		}
		pc, err := s.providerConfig(provider)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		profile, err := exchangeCode(c.Request.Context(), provider, pc, code)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		var accountID string
		if st["intent"] == "BIND" && st["account_id"] != "" {
			// 绑定模式：写入绑定关系，不创建账号
			if err := s.createBinding(c.Request.Context(), st["account_id"], provider, profile); err != nil {
				response.Fail(c, http.StatusBadRequest, 400, err.Error())
				return
			}
			accountID = st["account_id"]
		} else {
			accountID, err = s.resolveOrBindAccount(c.Request.Context(), accountType, provider, profile)
			if err != nil {
				response.Fail(c, http.StatusBadRequest, 400, err.Error())
				return
			}
		}
		exCode, err := randomHex(16)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, 500, err.Error())
			return
		}
		exPayload, _ := json.Marshal(map[string]string{
			"account_id":   accountID,
			"account_type": string(accountType),
		})
		_ = s.rdb.Set(c.Request.Context(), exchangeKeyPrefix+exCode, string(exPayload), 5*time.Minute).Err()
		redirect := st["redirect"]
		if redirect == "" {
			response.OK(c, gin.H{"code": exCode})
			return
		}
		sep := "?"
		if strings.Contains(redirect, "?") {
			sep = "&"
		}
		c.Redirect(http.StatusFound, redirect+sep+"code="+url.QueryEscape(exCode))
	}
}

func (s *Service) exchange(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ExchangeParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		raw, err := s.rdb.Get(c.Request.Context(), exchangeKeyPrefix+req.Code).Result()
		_ = s.rdb.Del(c.Request.Context(), exchangeKeyPrefix+req.Code)
		if err != nil || raw == "" {
			response.Fail(c, http.StatusBadRequest, 400, "兑换码无效或已过期")
			return
		}
		var payload map[string]string
		_ = json.Unmarshal([]byte(raw), &payload)
		if payload["account_type"] != string(accountType) {
			response.Fail(c, http.StatusBadRequest, 400, "账号类型不匹配")
			return
		}
		if s.issue == nil {
			response.Fail(c, http.StatusInternalServerError, 500, "session issuer not configured")
			return
		}
		token, err := s.issue(c.Request.Context(), accountType, payload["account_id"], c.ClientIP(), c.Request.UserAgent(), true)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, LoginResult{
			Token:       token,
			AccountID:   payload["account_id"],
			AccountType: accountType,
		})
	}
}

type providerCfg struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (s *Service) providerConfig(provider string) (providerCfg, error) {
	o := s.cfg.OAuth
	switch strings.ToLower(provider) {
	case "github":
		if o.GitHub.ClientID == "" || o.GitHub.ClientSecret == "" {
			return providerCfg{}, fmt.Errorf("OAuth 提供商 github 未配置")
		}
		return providerCfg{
			ClientID:     o.GitHub.ClientID,
			ClientSecret: o.GitHub.ClientSecret,
			RedirectURL:  o.GitHub.RedirectURL,
		}, nil
	case "gitee", "wechat", "wechat_open", "wechat_mp", "qq":
		return providerCfg{}, fmt.Errorf("OAuth 提供商 %s 未配置", provider)
	default:
		return providerCfg{}, fmt.Errorf("不支持的 OAuth 提供商: %s", provider)
	}
}

func buildAuthorizeURL(provider string, pc providerCfg, state string) (string, error) {
	switch provider {
	case "github":
		q := url.Values{}
		q.Set("client_id", pc.ClientID)
		q.Set("redirect_uri", pc.RedirectURL)
		q.Set("scope", "read:user user:email")
		q.Set("state", state)
		return "https://github.com/login/oauth/authorize?" + q.Encode(), nil
	default:
		return "", fmt.Errorf("OAuth 提供商 %s 未配置", provider)
	}
}

type oauthProfile struct {
	OpenID   string
	UnionID  string
	Nickname string
	Avatar   string
	Raw      string
}

func exchangeCode(ctx context.Context, provider string, pc providerCfg, code string) (*oauthProfile, error) {
	if provider != "github" {
		return nil, fmt.Errorf("OAuth 提供商 %s 未配置", provider)
	}
	form := url.Values{}
	form.Set("client_id", pc.ClientID)
	form.Set("client_secret", pc.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", pc.RedirectURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tok struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, err
	}
	if tok.AccessToken == "" {
		if tok.Error != "" {
			return nil, fmt.Errorf("github token: %s", tok.Error)
		}
		return nil, fmt.Errorf("github token 交换失败")
	}
	ureq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	ureq.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	ureq.Header.Set("Accept", "application/json")
	uresp, err := http.DefaultClient.Do(ureq)
	if err != nil {
		return nil, err
	}
	defer uresp.Body.Close()
	ubody, _ := io.ReadAll(uresp.Body)
	var user struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.Unmarshal(ubody, &user); err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("无法获取 GitHub 用户信息")
	}
	nick := user.Name
	if nick == "" {
		nick = user.Login
	}
	return &oauthProfile{
		OpenID:   fmt.Sprintf("%d", user.ID),
		Nickname: nick,
		Avatar:   user.AvatarURL,
		Raw:      string(ubody),
	}, nil
}

func (s *Service) resolveOrBindAccount(ctx context.Context, accountType security.AccountType, provider string, profile *oauthProfile) (string, error) {
	var binding AccountOAuthBinding
	err := s.db.WithContext(ctx).
		Where("provider = ? AND open_id = ?", strings.ToUpper(provider), profile.OpenID).
		First(&binding).Error
	if err == nil {
		var acc struct {
			ID            string
			AccountType   string
			AccountStatus string
		}
		if err := s.db.WithContext(ctx).Table("sys_account").
			Select("id, account_type, account_status").
			Where("id = ?", binding.AccountID).First(&acc).Error; err != nil {
			return "", fmt.Errorf("绑定账号不存在")
		}
		if acc.AccountType != string(accountType) || acc.AccountStatus != security.AccountStatusEnabled {
			return "", fmt.Errorf("账号不可用")
		}
		return acc.ID, nil
	}
	if err != gorm.ErrRecordNotFound {
		return "", err
	}
	// 首次绑定：自动创建对应类型账号（门户）；管理端要求已有绑定
	if accountType == security.AccountAdmin {
		return "", fmt.Errorf("未绑定的 GitHub 账号，请先在管理端绑定")
	}
	accountID := idgen.Next()
	now := time.Now().UTC()
	hash, _ := security.HashPassword(randomPassword())
	acc := map[string]any{
		"id":             accountID,
		"password_hash":  hash,
		"account_type":   string(accountType),
		"account_status": security.AccountStatusEnabled,
		"created_at":     now,
		"updated_at":     now,
	}
	ident := map[string]any{
		"id":            idgen.Next(),
		"account_id":    accountID,
		"identity_type": "ACCOUNT",
		"identifier":    "gh_" + profile.OpenID,
		"verified":      true,
		"is_primary":    true,
		"bind_status":   "BOUND",
		"created_at":    now,
		"updated_at":    now,
	}
	bind := AccountOAuthBinding{
		ID:         idgen.Next(),
		AccountID:  accountID,
		Provider:   strings.ToUpper(provider),
		OpenID:     profile.OpenID,
		Nickname:   strPtr(profile.Nickname),
		Avatar:     strPtr(profile.Avatar),
		RawProfile: profile.Raw,
		BoundAt:    now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if bind.RawProfile == "" {
		bind.RawProfile = "{}"
	}
	return accountID, s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("sys_account").Create(acc).Error; err != nil {
			return err
		}
		if err := tx.Table("sys_account_identity").Create(ident).Error; err != nil {
			return err
		}
		_ = tx.Table("portal_user_profile").Create(map[string]any{
			"account_id": accountID,
			"nickname":   profile.Nickname,
			"avatar":     profile.Avatar,
			"created_at": now,
			"updated_at": now,
		}).Error
		return tx.Create(&bind).Error
	})
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func randomPassword() string {
	s, _ := randomHex(16)
	return s
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
