// internal/modules/auth/oauth/providers.go 三方登录提供商（对齐 hei-boot JustAuth）。
//
// Author: Charlie

package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/security"
)

type oauthProvider string

const (
	providerGitHub     oauthProvider = "GITHUB"
	providerGitee      oauthProvider = "GITEE"
	providerQQ         oauthProvider = "QQ"
	providerWeChatOpen oauthProvider = "WECHAT_OPEN"
	providerWeChatMP   oauthProvider = "WECHAT_MP"
)

var allProviders = []oauthProvider{
	providerGitHub,
	providerGitee,
	providerQQ,
	providerWeChatOpen,
	providerWeChatMP,
}

var giteeEndpoint = oauth2.Endpoint{
	AuthURL:  "https://gitee.com/oauth/authorize",
	TokenURL: "https://gitee.com/oauth/token",
}

var jsonpPrefix = regexp.MustCompile(`^[^(]+\(`)

func parseProvider(raw string) (oauthProvider, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", fmt.Errorf("provider required")
	}
	normalized := strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
	switch normalized {
	case "GITHUB":
		return providerGitHub, nil
	case "GITEE":
		return providerGitee, nil
	case "QQ":
		return providerQQ, nil
	case "WECHAT", "WECHAT_OPEN":
		return providerWeChatOpen, nil
	case "WECHAT_MP":
		return providerWeChatMP, nil
	default:
		return "", fmt.Errorf("不支持的 OAuth 提供商: %s", raw)
	}
}

func providerLabel(p oauthProvider) string {
	switch p {
	case providerGitHub:
		return "GitHub"
	case providerGitee:
		return "Gitee"
	case providerQQ:
		return "QQ"
	case providerWeChatOpen:
		return "微信"
	case providerWeChatMP:
		return "微信小程序"
	default:
		return string(p)
	}
}

func isWebOAuth(p oauthProvider) bool {
	return p != providerWeChatMP
}

func oauthConfigKey(accountType security.AccountType, provider oauthProvider, field string) string {
	typ := "PORTAL"
	if accountType == security.AccountAdmin {
		typ = "ADMIN"
	}
	return fmt.Sprintf("AUTH_OAUTH_%s_%s_%s", typ, provider, field)
}

func (s *Service) oauthEnabled(ctx context.Context, accountType security.AccountType, provider oauthProvider) bool {
	if s.runtime == nil {
		return false
	}
	return s.runtime.GetBool(ctx, oauthConfigKey(accountType, provider, "ENABLED"), false)
}

func (s *Service) oauthConfigValue(ctx context.Context, accountType security.AccountType, provider oauthProvider, field string) string {
	if s.runtime == nil {
		return ""
	}
	return strings.TrimSpace(s.runtime.GetString(ctx, oauthConfigKey(accountType, provider, field), ""))
}

func (s *Service) yamlOAuthFallback(provider oauthProvider) config.OAuthProviderConfig {
	if s.cfg == nil {
		return config.OAuthProviderConfig{}
	}
	switch provider {
	case providerGitHub:
		return s.cfg.OAuth.GitHub
	case providerGitee:
		return s.cfg.OAuth.Gitee
	case providerQQ:
		return s.cfg.OAuth.QQ
	case providerWeChatOpen:
		return s.cfg.OAuth.WeChat
	case providerWeChatMP:
		return s.cfg.OAuth.WeChatMP
	default:
		return config.OAuthProviderConfig{}
	}
}

func (s *Service) providerConfig(ctx context.Context, accountType security.AccountType, provider oauthProvider) (providerCfg, error) {
	if !isWebOAuth(provider) && provider != providerWeChatMP {
		return providerCfg{}, fmt.Errorf("请使用小程序登录接口")
	}
	if provider != providerWeChatMP {
		if !s.oauthEnabled(ctx, accountType, provider) {
			return providerCfg{}, fmt.Errorf("%s 登录未启用", providerLabel(provider))
		}
	} else if !s.oauthEnabled(ctx, accountType, provider) {
		return providerCfg{}, fmt.Errorf("%s 登录未启用", providerLabel(provider))
	}

	clientID := firstNonBlank(
		s.oauthConfigValue(ctx, accountType, provider, "CLIENT_ID"),
		s.oauthConfigValue(ctx, accountType, provider, "APP_ID"),
	)
	clientSecret := firstNonBlank(
		s.oauthConfigValue(ctx, accountType, provider, "CLIENT_SECRET"),
		s.oauthConfigValue(ctx, accountType, provider, "APP_SECRET"),
	)
	redirectURL := s.oauthConfigValue(ctx, accountType, provider, "REDIRECT_URI")

	if clientID == "" || clientSecret == "" {
		fb := s.yamlOAuthFallback(provider)
		if clientID == "" {
			clientID = strings.TrimSpace(fb.ClientID)
		}
		if clientSecret == "" {
			clientSecret = strings.TrimSpace(fb.ClientSecret)
		}
		if redirectURL == "" {
			redirectURL = strings.TrimSpace(fb.RedirectURL)
		}
	}
	if clientID == "" || clientSecret == "" {
		return providerCfg{}, fmt.Errorf("%s 未配置 ClientId/Secret", providerLabel(provider))
	}
	if provider != providerWeChatMP && redirectURL == "" {
		return providerCfg{}, fmt.Errorf("请配置 %s", oauthConfigKey(accountType, provider, "REDIRECT_URI"))
	}
	return providerCfg{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}, nil
}

func firstNonBlank(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func buildAuthorizeURL(provider oauthProvider, pc providerCfg, state string) (string, error) {
	switch provider {
	case providerGitHub:
		q := url.Values{}
		q.Set("client_id", pc.ClientID)
		q.Set("redirect_uri", pc.RedirectURL)
		q.Set("scope", "read:user user:email")
		q.Set("state", state)
		return "https://github.com/login/oauth/authorize?" + q.Encode(), nil
	case providerGitee:
		q := url.Values{}
		q.Set("client_id", pc.ClientID)
		q.Set("redirect_uri", pc.RedirectURL)
		q.Set("response_type", "code")
		q.Set("scope", "user_info")
		q.Set("state", state)
		return "https://gitee.com/oauth/authorize?" + q.Encode(), nil
	case providerQQ:
		q := url.Values{}
		q.Set("response_type", "code")
		q.Set("client_id", pc.ClientID)
		q.Set("redirect_uri", pc.RedirectURL)
		q.Set("state", state)
		q.Set("scope", "get_user_info")
		return "https://graph.qq.com/oauth2.0/authorize?" + q.Encode(), nil
	case providerWeChatOpen:
		q := url.Values{}
		q.Set("appid", pc.ClientID)
		q.Set("redirect_uri", pc.RedirectURL)
		q.Set("response_type", "code")
		q.Set("scope", "snsapi_login")
		q.Set("state", state)
		return "https://open.weixin.qq.com/connect/qrconnect?" + q.Encode() + "#wechat_redirect", nil
	default:
		return "", fmt.Errorf("该提供商不支持网页授权")
	}
}

func exchangeCode(ctx context.Context, provider oauthProvider, pc providerCfg, code string) (*oauthProfile, error) {
	switch provider {
	case providerGitHub:
		return exchangeGitHub(ctx, pc, code)
	case providerGitee:
		return exchangeGitee(ctx, pc, code)
	case providerQQ:
		return exchangeQQ(ctx, pc, code)
	case providerWeChatOpen:
		return exchangeWeChatOpen(ctx, pc, code)
	default:
		return nil, fmt.Errorf("该提供商不支持网页授权回调")
	}
}

func exchangeGitHub(ctx context.Context, pc providerCfg, code string) (*oauthProfile, error) {
	conf := &oauth2.Config{
		ClientID:     pc.ClientID,
		ClientSecret: pc.ClientSecret,
		RedirectURL:  pc.RedirectURL,
		Endpoint:     githuboauth.Endpoint,
		Scopes:       []string{"read:user", "user:email"},
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("三方登录失败: %v", err)
	}
	client := conf.Client(ctx, tok)
	uresp, err := client.Get("https://api.github.com/user")
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
		return nil, fmt.Errorf("三方登录失败: 无法获取 GitHub 用户信息")
	}
	nick := user.Name
	if nick == "" {
		nick = user.Login
	}
	return &oauthProfile{
		Provider: string(providerGitHub),
		OpenID:   fmt.Sprintf("%d", user.ID),
		Nickname: nick,
		Avatar:   user.AvatarURL,
		Raw:      string(ubody),
	}, nil
}

func exchangeGitee(ctx context.Context, pc providerCfg, code string) (*oauthProfile, error) {
	conf := &oauth2.Config{
		ClientID:     pc.ClientID,
		ClientSecret: pc.ClientSecret,
		RedirectURL:  pc.RedirectURL,
		Endpoint:     giteeEndpoint,
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("三方登录失败: %v", err)
	}
	client := conf.Client(ctx, tok)
	uresp, err := client.Get("https://gitee.com/api/v5/user")
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
		return nil, fmt.Errorf("三方登录失败: 无法获取 Gitee 用户信息")
	}
	nick := user.Name
	if nick == "" {
		nick = user.Login
	}
	return &oauthProfile{
		Provider: string(providerGitee),
		OpenID:   fmt.Sprintf("%d", user.ID),
		Nickname: nick,
		Avatar:   user.AvatarURL,
		Raw:      string(ubody),
	}, nil
}

func exchangeQQ(ctx context.Context, pc providerCfg, code string) (*oauthProfile, error) {
	q := url.Values{}
	q.Set("grant_type", "authorization_code")
	q.Set("client_id", pc.ClientID)
	q.Set("client_secret", pc.ClientSecret)
	q.Set("code", code)
	q.Set("redirect_uri", pc.RedirectURL)
	tokenURL := "https://graph.qq.com/oauth2.0/token?" + q.Encode()
	body, err := httpGetText(ctx, tokenURL)
	if err != nil {
		return nil, err
	}
	body = stripJSONP(body)
	var tokenData struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Error       int    `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.Unmarshal([]byte(body), &tokenData); err != nil {
		return nil, fmt.Errorf("三方登录失败: token 解析错误")
	}
	if tokenData.AccessToken == "" {
		msg := tokenData.ErrorDesc
		if msg == "" {
			msg = body
		}
		return nil, fmt.Errorf("三方登录失败: %s", msg)
	}
	meBody, err := httpGetText(ctx, "https://graph.qq.com/oauth2.0/me?access_token="+url.QueryEscape(tokenData.AccessToken))
	if err != nil {
		return nil, err
	}
	meBody = stripJSONP(meBody)
	var me struct {
		ClientID string `json:"client_id"`
		OpenID   string `json:"openid"`
	}
	if err := json.Unmarshal([]byte(meBody), &me); err != nil || me.OpenID == "" {
		return nil, fmt.Errorf("三方登录失败: 无法获取 QQ openid")
	}
	userURL := "https://graph.qq.com/user/get_user_info?" + url.Values{
		"access_token":       {tokenData.AccessToken},
		"oauth_consumer_key": {pc.ClientID},
		"openid":             {me.OpenID},
	}.Encode()
	userBody, err := httpGetText(ctx, userURL)
	if err != nil {
		return nil, err
	}
	var user struct {
		Ret       int    `json:"ret"`
		Msg       string `json:"msg"`
		Nickname  string `json:"nickname"`
		FigureURL string `json:"figureurl_qq_2"`
	}
	_ = json.Unmarshal([]byte(userBody), &user)
	if user.Ret != 0 {
		return nil, fmt.Errorf("三方登录失败: %s", user.Msg)
	}
	avatar := user.FigureURL
	if avatar == "" {
		avatar = ""
	}
	return &oauthProfile{
		Provider: string(providerQQ),
		OpenID:   me.OpenID,
		Nickname: user.Nickname,
		Avatar:   avatar,
		Raw:      userBody,
	}, nil
}

func exchangeWeChatOpen(ctx context.Context, pc providerCfg, code string) (*oauthProfile, error) {
	q := url.Values{}
	q.Set("appid", pc.ClientID)
	q.Set("secret", pc.ClientSecret)
	q.Set("code", code)
	q.Set("grant_type", "authorization_code")
	tokenURL := "https://api.weixin.qq.com/sns/oauth2/access_token?" + q.Encode()
	body, err := httpGetText(ctx, tokenURL)
	if err != nil {
		return nil, err
	}
	var tokenData struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		UnionID     string `json:"unionid"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.Unmarshal([]byte(body), &tokenData); err != nil {
		return nil, fmt.Errorf("三方登录失败: token 解析错误")
	}
	if tokenData.ErrCode != 0 || tokenData.AccessToken == "" || tokenData.OpenID == "" {
		msg := tokenData.ErrMsg
		if msg == "" {
			msg = body
		}
		return nil, fmt.Errorf("三方登录失败: %s", msg)
	}
	userURL := "https://api.weixin.qq.com/sns/userinfo?" + url.Values{
		"access_token": {tokenData.AccessToken},
		"openid":       {tokenData.OpenID},
		"lang":         {"zh_CN"},
	}.Encode()
	userBody, err := httpGetText(ctx, userURL)
	if err != nil {
		return nil, err
	}
	var user struct {
		OpenID     string `json:"openid"`
		UnionID    string `json:"unionid"`
		Nickname   string `json:"nickname"`
		HeadImgURL string `json:"headimgurl"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err := json.Unmarshal([]byte(userBody), &user); err != nil {
		return nil, err
	}
	if user.ErrCode != 0 {
		return nil, fmt.Errorf("三方登录失败: %s", user.ErrMsg)
	}
	unionID := tokenData.UnionID
	if unionID == "" {
		unionID = user.UnionID
	}
	return &oauthProfile{
		Provider: string(providerWeChatOpen),
		OpenID:   tokenData.OpenID,
		UnionID:  unionID,
		Nickname: user.Nickname,
		Avatar:   user.HeadImgURL,
		Raw:      userBody,
	}, nil
}

func exchangeWechatMp(ctx context.Context, pc providerCfg, code string) (*oauthProfile, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("缺少 code")
	}
	q := url.Values{}
	q.Set("appid", pc.ClientID)
	q.Set("secret", pc.ClientSecret)
	q.Set("js_code", code)
	q.Set("grant_type", "authorization_code")
	u := "https://api.weixin.qq.com/sns/jscode2session?" + q.Encode()
	body, err := httpGetText(ctx, u)
	if err != nil {
		return nil, err
	}
	var data struct {
		OpenID     string `json:"openid"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
		SessionKey string `json:"session_key"`
	}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}
	if data.ErrCode != 0 {
		msg := data.ErrMsg
		if msg == "" {
			msg = fmt.Sprintf("%d", data.ErrCode)
		}
		return nil, fmt.Errorf("微信小程序登录失败: %s", msg)
	}
	if data.OpenID == "" {
		return nil, fmt.Errorf("微信小程序登录失败: 未返回 openid")
	}
	return &oauthProfile{
		Provider: string(providerWeChatMP),
		OpenID:   data.OpenID,
		UnionID:  data.UnionID,
		Raw:      body,
	}, nil
}

func httpGetText(ctx context.Context, rawURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func stripJSONP(body string) string {
	body = strings.TrimSpace(body)
	if idx := strings.Index(body, "{"); idx > 0 {
		body = body[idx:]
	}
	if idx := strings.LastIndex(body, "}"); idx >= 0 && idx < len(body)-1 {
		body = body[:idx+1]
	}
	body = jsonpPrefix.ReplaceAllString(body, "")
	return strings.TrimSpace(body)
}

func oauthAccountPrefix(provider oauthProvider) string {
	switch provider {
	case providerGitHub:
		return "gh_"
	case providerGitee:
		return "ge_"
	case providerQQ:
		return "qq_"
	case providerWeChatOpen, providerWeChatMP:
		return "wx_"
	default:
		return "oauth_"
	}
}

func allocateOauthAccountName(openID string, provider oauthProvider) string {
	prefix := oauthAccountPrefix(provider)
	suffix := regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(openID, "")
	if len(suffix) > 12 {
		suffix = suffix[:12]
	}
	if suffix == "" {
		suffix = fmt.Sprintf("%d", rand.Intn(1_000_000))
	}
	return strings.ToLower(prefix + suffix)
}
