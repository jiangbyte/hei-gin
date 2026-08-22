// Package oauth 提供三方登录（GitHub 完整；Gitee/微信等桩）。
//
// Author: Charlie
package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/bind"
	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/runtimecfg"
	"hei-gin/internal/framework/platform/module"
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
	cfg     *config.Config
	db      *gorm.DB
	rdb     *redis.Client
	runtime *runtimecfg.Settings
	audit   *audit.Queue
	issue   IssueSessionFunc
}

// NewService 构造 OAuth 服务。
func NewService(d *module.Deps, issue IssueSessionFunc) *Service {
	return &Service{
		cfg:     d.Cfg,
		db:      d.DB,
		rdb:     d.Redis,
		runtime: d.Runtime,
		audit:   d.Audit,
		issue:   issue,
	}
}

// ProviderOption 三方登录入口选项（供登录页公开配置）。
type ProviderOption struct {
	Provider string `json:"provider"`
	Label    string `json:"label"`
	Enabled  bool   `json:"enabled"`
	WebOAuth bool   `json:"web_oauth"`
}

// ProviderOptions 三方登录入口选项（对齐 hei-boot buildOauthProviderOptions）。
func (s *Service) ProviderOptions(ctx context.Context, accountType security.AccountType) []ProviderOption {
	out := make([]ProviderOption, 0, len(allProviders))
	for _, p := range allProviders {
		if accountType == security.AccountAdmin && p == providerWeChatMP {
			continue
		}
		out = append(out, ProviderOption{
			Provider: string(p),
			Label:    providerLabel(p),
			Enabled:  s.oauthEnabled(ctx, accountType, p),
			WebOAuth: isWebOAuth(p),
		})
	}
	return out
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
	// 微信小程序 code2session 登录（对齐 hei-boot PortalOauthController.wechatMpLogin）
	api.POST("/v1/portal/oauth/wechat-mp/login",
		middleware.RateLimit(rdb, "portal:oauth-wechat-mp", 30, 60),
		middleware.OperationAudit(s.audit, "auth", "oauth_wechat_mp_login"),
		s.wechatMpLogin)
}

// WechatMpLoginParam 微信小程序登录入参。
//
// Author: Charlie
type WechatMpLoginParam struct {
	Code string `json:"code" binding:"required"`
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
		provider, err := parseProvider(c.Param("provider"))
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if !isWebOAuth(provider) {
			response.Fail(c, http.StatusBadRequest, 400, "请使用小程序登录接口")
			return
		}
		pc, err := s.providerConfig(c.Request.Context(), accountType, provider)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		state, err := randomHex(16)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, 500, err.Error())
			return
		}
		redirect := c.Query("redirect")
		if redirect == "" {
			redirect = s.frontendCallback(c.Request.Context(), accountType)
		}
		intent := strings.TrimSpace(c.Query("intent"))
		if intent == "" {
			intent = "LOGIN"
		}
		intent = strings.ToUpper(intent)
		if intent != "LOGIN" && intent != "BIND" {
			response.Fail(c, http.StatusBadRequest, 400, "不支持的 OAuth intent")
			return
		}
		payload, _ := json.Marshal(map[string]string{
			"account_type": string(accountType),
			"provider":     string(provider),
			"redirect":     redirect,
			"intent":       intent,
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
		provider, err := parseProvider(c.Param("provider"))
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
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
		if st["provider"] != "" && !strings.EqualFold(st["provider"], string(provider)) {
			response.Fail(c, http.StatusBadRequest, 400, "授权状态不匹配")
			return
		}
		pc, err := s.providerConfig(c.Request.Context(), accountType, provider)
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
			if err := s.createBinding(c.Request.Context(), st["account_id"], string(provider), profile); err != nil {
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

// wechatMpLogin 微信小程序 code2session 登录并签发会话。
func (s *Service) wechatMpLogin(c *gin.Context) {
	var req WechatMpLoginParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	pc, err := s.providerConfig(c.Request.Context(), security.AccountPortal, providerWeChatMP)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	profile, err := exchangeWechatMp(c.Request.Context(), pc, req.Code)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	accountID, err := s.resolveOrBindAccount(c.Request.Context(), security.AccountPortal, providerWeChatMP, profile)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if s.issue == nil {
		response.Fail(c, http.StatusInternalServerError, 500, "session issuer not configured")
		return
	}
	token, err := s.issue(c.Request.Context(), security.AccountPortal, accountID, c.ClientIP(), c.Request.UserAgent(), true)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, LoginResult{
		Token:       token,
		AccountID:   accountID,
		AccountType: security.AccountPortal,
	})
}

// frontendCallback 前端 OAuth 回调页：优先运行时 AUTH_OAUTH_FRONTEND_CALLBACK_{TYPE}，缺省同源路径（对齐 hei-boot）。
func (s *Service) frontendCallback(ctx context.Context, accountType security.AccountType) string {
	key := "AUTH_OAUTH_FRONTEND_CALLBACK_ADMIN"
	if accountType != security.AccountAdmin {
		key = "AUTH_OAUTH_FRONTEND_CALLBACK_PORTAL"
	}
	if s.runtime != nil {
		if v := strings.TrimSpace(s.runtime.GetString(ctx, key, "")); v != "" {
			return v
		}
	}
	return "/auth/oauth/callback"
}

type providerCfg struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type oauthProfile struct {
	Provider string
	OpenID   string
	UnionID  string
	Nickname string
	Avatar   string
	Raw      string
}

func (s *Service) resolveOrBindAccount(ctx context.Context, accountType security.AccountType, provider oauthProvider, profile *oauthProfile) (string, error) {
	providerName := string(provider)
	if profile.Provider != "" {
		providerName = profile.Provider
	}
	var binding AccountOAuthBinding
	err := s.db.WithContext(ctx).
		Where("provider = ? AND open_id = ?", providerName, profile.OpenID).
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
	if accountType == security.AccountAdmin {
		return "", fmt.Errorf("请先使用账号密码登录后再绑定该三方账号")
	}
	accountID := idgen.Next()
	now := time.Now().UTC()
	hash, _ := security.HashPassword(randomPassword())
	identifier := allocateOauthAccountName(profile.OpenID, provider)
	for i := 0; ; i++ {
		var cnt int64
		_ = s.db.WithContext(ctx).Table("sys_account_identity").
			Where("identity_type = ? AND identifier = ?", "ACCOUNT", identifier).Count(&cnt).Error
		if cnt == 0 {
			break
		}
		if i == 0 {
			identifier = allocateOauthAccountName(profile.OpenID+fmt.Sprintf("%d", i), provider)
		} else {
			identifier = allocateOauthAccountName(profile.OpenID, provider) + fmt.Sprintf("%d", i)
		}
	}
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
		"identifier":    identifier,
		"verified":      true,
		"is_primary":    true,
		"bind_status":   "BOUND",
		"created_at":    now,
		"updated_at":    now,
	}
	bindRow := AccountOAuthBinding{
		ID:         idgen.Next(),
		AccountID:  accountID,
		Provider:   providerName,
		OpenID:     profile.OpenID,
		Nickname:   strPtr(profile.Nickname),
		Avatar:     strPtr(profile.Avatar),
		RawProfile: profile.Raw,
		BoundAt:    now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if bindRow.RawProfile == "" {
		bindRow.RawProfile = "{}"
	}
	return accountID, s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("sys_account").Create(acc).Error; err != nil {
			return err
		}
		if err := tx.Table("sys_account_identity").Create(ident).Error; err != nil {
			return err
		}
		_ = tx.Table("profile_user_portal").Create(map[string]any{
			"account_id": accountID,
			"nickname":   profile.Nickname,
			"avatar":     profile.Avatar,
			"created_at": now,
			"updated_at": now,
		}).Error
		return tx.Create(&bindRow).Error
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
