// internal/modules/auth/session.go 会话管理（对齐 hei-fastapi session_admin_service）。
//
// Author: Charlie

package auth

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/modules/iam/view"
)

const analysisCacheTTL = 30 * time.Second

var (
	analysisCacheMu   sync.Mutex
	analysisCacheAt   time.Time
	analysisCacheData *SessionAnalysis
)

// SessionAnalysis 在线会话分析统计。
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

// SessionAccount 单个在线账号的会话汇总。
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

// SessionTokenInfo 单个 Token 会话详情。
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

// SessionPageParam 在线会话分页查询参数。
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

// Normalize 会话分页参数（size 上限 200，对齐 hei-boot AdminSessionController.page）。
func (q SessionPageParam) Normalize() (current, size int) {
	current = q.Current
	size = q.Size
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 200 {
		size = 200
	}
	return current, size
}

// SessionExitTarget 强制下线目标。
//
// Author: Charlie
type SessionExitTarget struct {
	AccountID   string `json:"account_id" binding:"required"`
	AccountType string `json:"account_type"`
}

// SessionExitParam 批量强制下线请求。
//
// Author: Charlie
type SessionExitParam struct {
	Targets []SessionExitTarget `json:"targets" binding:"required"`
}

// SessionTokenExitParam 按 Token 强制下线请求。
//
// Author: Charlie
type SessionTokenExitParam struct {
	Tokens []string `json:"tokens" binding:"required"`
}

type sessionGroupKey struct {
	AccountType string
	AccountID   string
}

// registerSessionRoutes 挂载管理端在线会话 API。
func (s *Service) registerSessionRoutes(api *gin.RouterGroup) {
	g := api.Group("/v1/admin/auth/sessions", middleware.RequireAccountType(security.AccountAdmin))
	g.GET("/analysis", middleware.RequirePermission(s.perms, "auth:session:analysis", "会话分析"), s.sessionAnalysis)
	g.GET("/page", middleware.RequirePermission(s.perms, "auth:session:page", "会话分页"), s.sessionPage)
	g.GET("/tokens", middleware.RequirePermission(s.perms, "auth:session:tokenlist", "会话 Token 列表"), s.sessionTokens)
	g.POST("/exit", middleware.RequirePermission(s.perms, "auth:session:exit", "会话强制下线"), middleware.OperationAudit(s.audit, "auth_session", "exit"), s.sessionExit)
	g.POST("/token/exit", middleware.RequirePermission(s.perms, "auth:session:tokenexit", "Token 强制下线"), middleware.OperationAudit(s.audit, "auth_session", "token_exit"), s.sessionTokenExit)
}

func invalidateAnalysisCache() {
	analysisCacheMu.Lock()
	analysisCacheData = nil
	analysisCacheAt = time.Time{}
	analysisCacheMu.Unlock()
}

func (s *Service) sessionAnalysis(c *gin.Context) {
	now := time.Now()
	analysisCacheMu.Lock()
	if analysisCacheData != nil && now.Before(analysisCacheAt) {
		cached := *analysisCacheData
		analysisCacheMu.Unlock()
		response.OK(c, cached)
		return
	}
	analysisCacheMu.Unlock()

	grouped, err := s.groupOnlineSessions(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	oneHourAgo := now.Add(-time.Hour)
	result := SessionAnalysis{}
	maxToken := 0
	for key, sessions := range grouped {
		result.OnlineAccountCount++
		n := len(sessions)
		result.OnlineTokenCount += n
		if n > maxToken {
			maxToken = n
		}
		switch security.AccountType(key.AccountType) {
		case security.AccountAdmin:
			result.AdminAccountCount++
		case security.AccountPortal:
			result.PortalAccountCount++
		}
		for _, sess := range sessions {
			if !sess.LoginAt.IsZero() && !sess.LoginAt.Before(oneHourAgo) {
				result.OneHourNewCount++
			}
		}
	}
	result.MaxTokenCount = maxToken

	analysisCacheMu.Lock()
	analysisCacheData = &result
	analysisCacheAt = now.Add(analysisCacheTTL)
	analysisCacheMu.Unlock()
	response.OK(c, result)
}

func (s *Service) sessionPage(c *gin.Context) {
	var q SessionPageParam
	_ = c.ShouldBindQuery(&q)
	cur, size := q.Normalize()
	grouped, err := s.groupOnlineSessions(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	grouped = filterGroupedSessions(grouped, q)
	needsProfile := q.Account != "" || q.Keyword != ""
	if needsProfile {
		items := s.buildSessionItems(c.Request.Context(), grouped)
		items = filterProfileSessionItems(items, q)
		sortSessionItems(items)
		total := int64(len(items))
		from := (cur - 1) * size
		to := from + size
		if from > len(items) {
			from = len(items)
		}
		if to > len(items) {
			to = len(items)
		}
		response.Page(c, int64(cur), int64(size), total, items[from:to])
		return
	}

	keys := make([]sessionGroupKey, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return groupSortLess(grouped[keys[i]], grouped[keys[j]], keys[i].AccountID, keys[j].AccountID)
	})
	total := int64(len(keys))
	from := (cur - 1) * size
	to := from + size
	if from > len(keys) {
		from = len(keys)
	}
	if to > len(keys) {
		to = len(keys)
	}
	pageGrouped := make(map[sessionGroupKey][]*security.SessionPayload, to-from)
	for _, k := range keys[from:to] {
		pageGrouped[k] = grouped[k]
	}
	items := s.buildSessionItems(c.Request.Context(), pageGrouped)
	sortSessionItems(items)
	response.Page(c, int64(cur), int64(size), total, items)
}

func (s *Service) sessionTokens(c *gin.Context) {
	accountID := strings.TrimSpace(c.Query("account_id"))
	if accountID == "" {
		response.Fail(c, http.StatusBadRequest, 400, "account_id required")
		return
	}
	accountType := security.AccountType(strings.TrimSpace(c.Query("account_type")))
	if accountType == "" {
		accountType = security.AccountAdmin
	}
	tokens, err := s.sessions.ListTokensForAccount(c.Request.Context(), accountType, accountID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	sessions, err := s.sessions.ListSessionsByTokens(c.Request.Context(), tokens)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out := make([]SessionTokenInfo, 0, len(sessions))
	for _, sess := range sessions {
		out = append(out, toTokenInfo(sess))
	}
	response.OK(c, out)
}

func (s *Service) sessionExit(c *gin.Context) {
	var req SessionExitParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	targets := make([]security.AccountSessionTarget, 0, len(req.Targets))
	for _, t := range req.Targets {
		accountType := security.AccountType(strings.TrimSpace(t.AccountType))
		if accountType == "" {
			accountType = security.AccountAdmin
		}
		targets = append(targets, security.AccountSessionTarget{
			AccountType: accountType,
			AccountID:   t.AccountID,
		})
	}
	if err := s.sessions.DeleteAccountsSessions(c.Request.Context(), targets); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	invalidateAnalysisCache()
	response.OK(c, nil)
}

func (s *Service) sessionTokenExit(c *gin.Context) {
	var req SessionTokenExitParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	seen := map[string]struct{}{}
	for _, tok := range req.Tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		if _, ok := seen[tok]; ok {
			continue
		}
		seen[tok] = struct{}{}
		_ = s.sessions.Delete(c.Request.Context(), tok)
	}
	invalidateAnalysisCache()
	response.OK(c, nil)
}

func (s *Service) groupOnlineSessions(ctx context.Context) (map[sessionGroupKey][]*security.SessionPayload, error) {
	tokens, err := s.sessions.ListTokens(ctx)
	if err != nil {
		return nil, err
	}
	sessions, err := s.sessions.ListSessionsByTokens(ctx, tokens)
	if err != nil {
		return nil, err
	}
	grouped := make(map[sessionGroupKey][]*security.SessionPayload)
	for _, sess := range sessions {
		if sess == nil {
			continue
		}
		key := sessionGroupKey{AccountType: string(sess.AccountType), AccountID: sess.AccountID}
		grouped[key] = append(grouped[key], sess)
	}
	return grouped, nil
}

func filterGroupedSessions(
	grouped map[sessionGroupKey][]*security.SessionPayload,
	q SessionPageParam,
) map[sessionGroupKey][]*security.SessionPayload {
	result := grouped
	if q.AccountType != "" {
		wanted := q.AccountType
		filtered := make(map[sessionGroupKey][]*security.SessionPayload)
		for k, sessions := range result {
			if k.AccountType == wanted {
				filtered[k] = sessions
			}
		}
		result = filtered
	}
	if q.AccountID != "" {
		filtered := make(map[sessionGroupKey][]*security.SessionPayload)
		for k, sessions := range result {
			if k.AccountID == q.AccountID {
				filtered[k] = sessions
			}
		}
		result = filtered
	}
	if q.IP != "" {
		filtered := make(map[sessionGroupKey][]*security.SessionPayload)
		for k, sessions := range result {
			for _, sess := range sessions {
				if sess.ClientIP != nil && strings.Contains(*sess.ClientIP, q.IP) {
					filtered[k] = sessions
					break
				}
			}
		}
		result = filtered
	}
	return result
}

func (s *Service) buildSessionItems(
	ctx context.Context,
	grouped map[sessionGroupKey][]*security.SessionPayload,
) []SessionAccount {
	if len(grouped) == 0 {
		return nil
	}
	ids := make([]string, 0, len(grouped))
	seen := map[string]struct{}{}
	for k := range grouped {
		if _, ok := seen[k.AccountID]; ok {
			continue
		}
		seen[k.AccountID] = struct{}{}
		ids = append(ids, k.AccountID)
	}
	views, _ := view.LoadAccountViews(ctx, s.db, ids)
	viewByID := make(map[string]view.AccountView, len(views))
	for _, v := range views {
		viewByID[v.ID] = v
	}

	items := make([]SessionAccount, 0, len(grouped))
	for key, sessions := range grouped {
		tokenInfos := make([]SessionTokenInfo, 0, len(sessions))
		for _, sess := range sessions {
			tokenInfos = append(tokenInfos, toTokenInfo(sess))
		}
		sort.Slice(tokenInfos, func(i, j int) bool {
			return tokenActiveAt(tokenInfos[i]).After(tokenActiveAt(tokenInfos[j]))
		})
		var firstLogin *time.Time
		var latestActive *time.Time
		for _, t := range tokenInfos {
			if t.LoginAt != nil && (firstLogin == nil || t.LoginAt.Before(*firstLogin)) {
				v := *t.LoginAt
				firstLogin = &v
			}
			if t.LastActiveAt != nil && (latestActive == nil || t.LastActiveAt.After(*latestActive)) {
				v := *t.LastActiveAt
				latestActive = &v
			}
		}
		item := SessionAccount{
			AccountID:      key.AccountID,
			AccountType:    key.AccountType,
			TokenCount:     len(tokenInfos),
			FirstLoginAt:   firstLogin,
			LatestActiveAt: latestActive,
			Tokens:         tokenInfos,
		}
		if len(tokenInfos) > 0 {
			newest := tokenInfos[0]
			item.ClientIP = newest.ClientIP
			item.DeviceLabel = newest.DeviceLabel
		}
		if v, ok := viewByID[key.AccountID]; ok {
			item.Account = v.Account
			if item.Account == "" {
				item.Account = key.AccountID
			}
			item.Name = v.Name
			item.Nickname = v.Nickname
			item.Avatar = v.Avatar
			item.LatestLoginIP = v.LatestLoginIP
			item.LatestLoginTime = v.LatestLoginTime
		} else {
			item.Account = key.AccountID
		}
		items = append(items, item)
	}
	return items
}

func filterProfileSessionItems(items []SessionAccount, q SessionPageParam) []SessionAccount {
	result := items
	if q.Account != "" {
		keyword := strings.ToLower(q.Account)
		filtered := make([]SessionAccount, 0, len(result))
		for _, item := range result {
			if containsFold(item.Account, keyword) ||
				containsFold(deref(item.Name), keyword) ||
				containsFold(deref(item.Nickname), keyword) {
				filtered = append(filtered, item)
			}
		}
		result = filtered
	}
	if q.Keyword != "" {
		keyword := strings.ToLower(q.Keyword)
		filtered := make([]SessionAccount, 0, len(result))
		for _, item := range result {
			if containsFold(item.Account, keyword) ||
				containsFold(item.AccountID, keyword) {
				filtered = append(filtered, item)
			}
		}
		result = filtered
	}
	return result
}

func sortSessionItems(items []SessionAccount) {
	sort.Slice(items, func(i, j int) bool {
		ai := itemSortActive(items[i])
		aj := itemSortActive(items[j])
		if !ai.Equal(aj) {
			return ai.After(aj)
		}
		return items[i].AccountID > items[j].AccountID
	})
}

func groupSortLess(a, b []*security.SessionPayload, accountIDA, accountIDB string) bool {
	ta := groupLatestActive(a)
	tb := groupLatestActive(b)
	if !ta.Equal(tb) {
		return ta.After(tb)
	}
	return accountIDA > accountIDB
}

func groupLatestActive(sessions []*security.SessionPayload) time.Time {
	var latest time.Time
	for _, sess := range sessions {
		active := sess.LastActiveAt
		if active.IsZero() {
			active = sess.LoginAt
		}
		if active.After(latest) {
			latest = active
		}
	}
	return latest
}

func itemSortActive(item SessionAccount) time.Time {
	if item.LatestActiveAt != nil {
		return *item.LatestActiveAt
	}
	if item.FirstLoginAt != nil {
		return *item.FirstLoginAt
	}
	return time.Time{}
}

func tokenActiveAt(t SessionTokenInfo) time.Time {
	if t.LastActiveAt != nil {
		return *t.LastActiveAt
	}
	if t.LoginAt != nil {
		return *t.LoginAt
	}
	return time.Time{}
}

func toTokenInfo(sess *security.SessionPayload) SessionTokenInfo {
	loginAt := sess.LoginAt
	lastActive := sess.LastActiveAt
	expires := sess.ExpiresAt
	return SessionTokenInfo{
		Token:        sess.Token,
		AccountID:    sess.AccountID,
		AccountType:  string(sess.AccountType),
		LoginAt:      &loginAt,
		LastActiveAt: &lastActive,
		ExpiresAt:    &expires,
		ClientIP:     sess.ClientIP,
		DeviceLabel:  sess.DeviceLabel,
		UserAgent:    sess.UserAgent,
		RememberMe:   sess.RememberMe,
	}
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

// SessionFromAuthContext 复用 auth 包 handler 用的会话读取。
func SessionFromAuthContext(c *gin.Context) *security.SessionPayload {
	return contextx.Session(c.Request.Context())
}
