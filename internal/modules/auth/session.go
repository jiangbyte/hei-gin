package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
)

// SessionAnalysis åœ¨çº¿ä¼šè¯åˆ†æžç»Ÿè®¡ã€‚
//
// Author: Charlie
type SessionAnalysis struct {
	OnlineAccountCount int `json:"online_account_count"`
	OnlineTokenCount   int `json:"online_token_count"`
	AdminAccountCount  int `json:"admin_account_count"`
	PortalAccountCount int `json:"portal_account_count"`
	OneHourNewCount    int `json:"one_hour_new_count"`
	MaxTokenCount      int `json:"max_token_count"`
}

// SessionAccount å•ä¸ªåœ¨çº¿è´¦å·çš„ä¼šè¯æ±‡æ€»ã€‚
//
// Author: Charlie
type SessionAccount struct {
	AccountID       string             `json:"account_id"`
	Account         string             `json:"account"`
	AccountType     string             `json:"account_type"`
	Name            *string            `json:"name"`
	Nickname        *string            `json:"nickname"`
	Avatar          *string            `json:"avatar"`
	LatestLoginIP   *string            `json:"latest_login_ip"`
	LatestLoginTime *time.Time         `json:"latest_login_time"`
	TokenCount      int                `json:"token_count"`
	FirstLoginAt    *time.Time         `json:"first_login_at"`
	LatestActiveAt  *time.Time         `json:"latest_active_at"`
	ClientIP        *string            `json:"client_ip"`
	DeviceLabel     *string            `json:"device_label"`
	Tokens          []SessionTokenInfo `json:"tokens"`
}

// SessionTokenInfo å•ä¸ª Token ä¼šè¯è¯¦æƒ…ã€‚
//
// Author: Charlie
type SessionTokenInfo struct {
	Token        string     `json:"token"`
	AccountID    string     `json:"account_id"`
	AccountType  string     `json:"account_type"`
	LoginAt      *time.Time `json:"login_at"`
	LastActiveAt *time.Time `json:"last_active_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	ClientIP     *string    `json:"client_ip"`
	DeviceLabel  *string    `json:"device_label"`
	UserAgent    *string    `json:"user_agent"`
	RememberMe   bool       `json:"remember_me"`
}

// SessionPageParam åœ¨çº¿ä¼šè¯åˆ†é¡µæŸ¥è¯¢å‚æ•°ã€‚
//
// Author: Charlie
type SessionPageParam struct {
	schema.PageQuery
	AccountType string `form:"account_type"`
	AccountID   string `form:"account_id"`
	Account     string `form:"account"`
	IP          string `form:"ip"`
	Keyword     string `form:"keyword"`
}

// SessionExitTarget å¼ºåˆ¶ä¸‹çº¿ç›®æ ‡ã€‚
//
// Author: Charlie
type SessionExitTarget struct {
	AccountID   string `json:"account_id" binding:"required"`
	AccountType string `json:"account_type"`
}

// SessionExitParam æ‰¹é‡å¼ºåˆ¶ä¸‹çº¿è¯·æ±‚ã€‚
//
// Author: Charlie
type SessionExitParam struct {
	Targets []SessionExitTarget `json:"targets" binding:"required"`
}

// SessionTokenExitParam æŒ‰ Token å¼ºåˆ¶ä¸‹çº¿è¯·æ±‚ã€‚
//
// Author: Charlie
type SessionTokenExitParam struct {
	Tokens []string `json:"tokens" binding:"required"`
}

// registerSessionRoutes æŒ‚è½½ç®¡ç†ç«¯åœ¨çº¿ä¼šè¯ APIã€‚
func (s *Service) registerSessionRoutes(api *gin.RouterGroup) {
	g := api.Group("/v1/admin/auth/sessions", middleware.RequireAccountType(security.AccountAdmin))
	g.GET("/analysis", middleware.RequirePermission(s.perms, "auth:session:analysis", "ä¼šè¯åˆ†æž"), s.sessionAnalysis)
	g.GET("/page", middleware.RequirePermission(s.perms, "auth:session:page", "ä¼šè¯åˆ†é¡µ"), s.sessionPage)
	g.GET("/tokens", middleware.RequirePermission(s.perms, "auth:session:tokenlist", "ä¼šè¯ Token åˆ—è¡¨"), s.sessionTokens)
	g.POST("/exit", middleware.RequirePermission(s.perms, "auth:session:exit", "ä¼šè¯å¼ºåˆ¶ä¸‹çº¿"), s.sessionExit)
	g.POST("/token/exit", middleware.RequirePermission(s.perms, "auth:session:tokenexit", "Token å¼ºåˆ¶ä¸‹çº¿"), s.sessionTokenExit)
}

func (s *Service) sessionAnalysis(c *gin.Context) {
	admin := s.analyzeSessions(c, security.AccountAdmin)
	portal := s.analyzeSessions(c, security.AccountPortal)
	response.OK(c, SessionAnalysis{
		OnlineAccountCount: admin.accountCount + portal.accountCount,
		OnlineTokenCount:   admin.tokenCount + portal.tokenCount,
		AdminAccountCount:  admin.accountCount,
		PortalAccountCount: portal.accountCount,
		OneHourNewCount:    admin.oneHourNew + portal.oneHourNew,
		MaxTokenCount:      maxInt(admin.maxToken, portal.maxToken),
	})
}

type sessionAccumulator struct {
	accountCount int
	tokenCount   int
	oneHourNew   int
	maxToken     int
}

func (s *Service) analyzeSessions(c *gin.Context, accountType security.AccountType) sessionAccumulator {
	acc := sessionAccumulator{}
	accountIDs, err := s.sessions.ListAccountIDs(c.Request.Context())
	if err != nil {
		return acc
	}
	for _, accountID := range accountIDs {
		sess, _ := s.sessions.Get(c.Request.Context(), s.firstToken(c, accountID))
		if sess != nil && sess.AccountType != accountType {
			continue
		}
		tokens, _ := s.sessions.ListTokensForAccount(c.Request.Context(), accountID)
		if len(tokens) == 0 {
			continue
		}
		// æ ¡éªŒè¯¥è´¦å·å­˜åœ¨å¯¹åº”ç±»åž‹çš„ä¼šè¯
		matched := 0
		for _, tok := range tokens {
			one, _ := s.sessions.Get(c.Request.Context(), tok)
			if one == nil || one.AccountType != accountType {
				continue
			}
			matched++
			if time.Since(one.LoginAt) < time.Hour {
				acc.oneHourNew++
			}
		}
		if matched == 0 {
			continue
		}
		acc.accountCount++
		acc.tokenCount += matched
		if matched > acc.maxToken {
			acc.maxToken = matched
		}
	}
	return acc
}

func (s *Service) firstToken(c *gin.Context, accountID string) string {
	tokens, _ := s.sessions.ListTokensForAccount(c.Request.Context(), accountID)
	if len(tokens) == 0 {
		return ""
	}
	return tokens[0]
}

func (s *Service) sessionPage(c *gin.Context) {
	var q SessionPageParam
	_ = c.ShouldBindQuery(&q)
	cur, size := q.Normalize()
	accountIDs, err := s.sessions.ListAccountIDs(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	var matched []SessionAccount
	for _, accountID := range accountIDs {
		item := s.hydrateSession(c, accountID, q.AccountType)
		if item == nil {
			continue
		}
		if !matchesSessionFilter(item, q) {
			continue
		}
		matched = append(matched, *item)
	}
	total := int64(len(matched))
	from := (cur - 1) * size
	to := from + size
	if from >= len(matched) {
		from = len(matched)
	}
	if to > len(matched) {
		to = len(matched)
	}
	response.Page(c, int64(cur), int64(size), total, matched[from:to])
}

func (s *Service) sessionTokens(c *gin.Context) {
	accountID := c.Query("account_id")
	if accountID == "" {
		response.Fail(c, http.StatusBadRequest, 400, "account_id required")
		return
	}
	accountType := c.Query("account_type")
	var tokens []SessionTokenInfo
	for _, tok := range s.tokenList(c, accountID) {
		sess, _ := s.sessions.Get(c.Request.Context(), tok)
		if sess == nil {
			continue
		}
		if accountType != "" && string(sess.AccountType) != accountType {
			continue
		}
		tokens = append(tokens, toTokenInfo(sess))
	}
	response.OK(c, tokens)
}

func (s *Service) sessionExit(c *gin.Context) {
	var req SessionExitParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	for _, target := range req.Targets {
		_ = s.sessions.DeleteAllForAccount(c.Request.Context(), target.AccountID)
	}
	response.OK(c, nil)
}

func (s *Service) sessionTokenExit(c *gin.Context) {
	var req SessionTokenExitParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	for _, tok := range req.Tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		_ = s.sessions.Delete(c.Request.Context(), tok)
	}
	response.OK(c, nil)
}

func (s *Service) hydrateSession(c *gin.Context, accountID, accountTypeFilter string) *SessionAccount {
	var tokens []SessionTokenInfo
	var latest *security.SessionPayload
	firstLogin := time.Time{}
	for _, tok := range s.tokenList(c, accountID) {
		sess, _ := s.sessions.Get(c.Request.Context(), tok)
		if sess == nil {
			continue
		}
		if accountTypeFilter != "" && string(sess.AccountType) != accountTypeFilter {
			continue
		}
		if firstLogin.IsZero() || sess.LoginAt.Before(firstLogin) {
			firstLogin = sess.LoginAt
		}
		if latest == nil || sess.LastActiveAt.After(latest.LastActiveAt) {
			latest = sess
		}
		tokens = append(tokens, toTokenInfo(sess))
	}
	if len(tokens) == 0 {
		return nil
	}
	item := &SessionAccount{
		AccountID:   accountID,
		AccountType: string(accountTypeOf(tokens)),
		TokenCount:  len(tokens),
		Tokens:      tokens,
	}
	if !firstLogin.IsZero() {
		item.FirstLoginAt = &firstLogin
	}
	if latest != nil {
		item.LatestLoginIP = latest.ClientIP
		item.LatestLoginTime = &latest.LoginAt
		item.LatestActiveAt = &latest.LastActiveAt
		item.ClientIP = latest.ClientIP
		item.DeviceLabel = latest.DeviceLabel
		item.AccountType = string(latest.AccountType)
	}
	item.Account = s.accountName(c, accountID)
	name, nickname, avatar := s.profileNames(c, accountID)
	item.Name = name
	item.Nickname = nickname
	item.Avatar = avatar
	return item
}

func accountTypeOf(tokens []SessionTokenInfo) security.AccountType {
	if len(tokens) == 0 {
		return ""
	}
	return security.AccountType(tokens[0].AccountType)
}

func (s *Service) tokenList(c *gin.Context, accountID string) []string {
	tokens, _ := s.sessions.ListTokensForAccount(c.Request.Context(), accountID)
	return tokens
}

func toTokenInfo(sess *security.SessionPayload) SessionTokenInfo {
	return SessionTokenInfo{
		Token:        sess.Token,
		AccountID:    sess.AccountID,
		AccountType:  string(sess.AccountType),
		LoginAt:      &sess.LoginAt,
		LastActiveAt: &sess.LastActiveAt,
		ExpiresAt:    &sess.ExpiresAt,
		ClientIP:     sess.ClientIP,
		DeviceLabel:  sess.DeviceLabel,
		UserAgent:    sess.UserAgent,
		RememberMe:   sess.RememberMe,
	}
}

func (s *Service) accountName(c *gin.Context, accountID string) string {
	if s.db == nil {
		return accountID
	}
	var account string
	_ = s.db.WithContext(c.Request.Context()).Table("sys_account").
		Select("account").
		Where("id = ?", accountID).Limit(1).Scan(&account).Error
	if account == "" {
		return accountID
	}
	return account
}

func (s *Service) profileNames(c *gin.Context, accountID string) (*string, *string, *string) {
	if s.db == nil {
		return nil, nil, nil
	}
	var name, nickname, avatar *string
	_ = s.db.WithContext(c.Request.Context()).Table("admin_user_profile").
		Select("name, nickname, avatar").
		Where("account_id = ?", accountID).Limit(1).
		Scan(&struct {
			Name     *string `gorm:"column:name"`
			Nickname *string `gorm:"column:nickname"`
			Avatar   *string `gorm:"column:avatar"`
		}{name, nickname, avatar}).Error
	return name, nickname, avatar
}

func matchesSessionFilter(item *SessionAccount, q SessionPageParam) bool {
	if q.AccountID != "" && q.AccountID != item.AccountID {
		return false
	}
	if q.Account != "" {
		k := strings.ToLower(q.Account)
		if !containsFold(item.Account, k) && !containsFold(deref(item.Name), k) && !containsFold(deref(item.Nickname), k) {
			return false
		}
	}
	if q.IP != "" {
		hit := false
		for _, t := range item.Tokens {
			if t.ClientIP != nil && strings.Contains(*t.ClientIP, q.IP) {
				hit = true
				break
			}
		}
		if !hit {
			return false
		}
	}
	if q.Keyword != "" {
		k := strings.ToLower(q.Keyword)
		if !containsFold(item.Account, k) && !containsFold(item.AccountID, k) {
			return false
		}
	}
	return true
}

func containsFold(value, keyword string) bool {
	return strings.Contains(strings.ToLower(value), keyword)
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SessionFromAuthContext å¤ç”¨ auth åŒ… handler ç”¨çš„ä¼šè¯è¯»å–ã€‚
func SessionFromAuthContext(c *gin.Context) *security.SessionPayload {
	return contextx.Session(c.Request.Context())
}
