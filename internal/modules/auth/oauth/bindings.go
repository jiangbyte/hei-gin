// internal/modules/auth/oauth/bindings.go 三方绑定。
//
// Author: Charlie

package oauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
)

// BindingResult 当前用户三方绑定列表项。
//
// Author: Charlie
type BindingResult struct {
	Provider     string    `json:"provider"`
	Label        string    `json:"label"`
	OpenIDMasked string    `json:"open_id_masked"`
	Nickname     *string   `json:"nickname"`
	Avatar       *string   `json:"avatar"`
	BoundAt      time.Time `json:"bound_at"`
}

// AdminUnbindParam 管理端强制解绑入参。
//
// Author: Charlie
type AdminUnbindParam struct {
	AccountID string `json:"account_id" binding:"required"`
	Provider  string `json:"provider" binding:"required"`
}

// RegisterBindingRoutes 挂载绑定相关路由（admin/portal OAuth 前缀下）。
func (s *Service) RegisterBindingRoutes(api *gin.RouterGroup, perms *security.PermissionRegistry) {
	for _, prefix := range []struct {
		base string
		typ  security.AccountType
	}{
		{"/v1/admin/oauth", security.AccountAdmin},
		{"/v1/portal/oauth", security.AccountPortal},
	} {
		p := prefix
		g := api.Group(p.base, middleware.RequireAccountType(p.typ))
		g.GET("/bindings", s.bindings(p.typ))
		g.POST("/:provider/bind/authorize", middleware.OperationAudit(s.audit, "auth", "oauth_bind_authorize"), s.bindAuthorize(p.typ))
		g.POST("/:provider/unbind", middleware.OperationAudit(s.audit, "auth", "oauth_unbind"), s.unbind(p.typ))
	}
	// 管理端强制解绑
	api.POST("/v1/admin/sys/accounts/oauth/unbind",
		middleware.RequireAccountType(security.AccountAdmin),
		middleware.RequirePermission(perms, "iam:account:update", "账户更新"),
		middleware.OperationAudit(s.audit, "iam_account", "oauth_unbind"),
		s.adminUnbind)
}

func (s *Service) bindings(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := contextx.Session(c.Request.Context())
		if sess == nil {
			response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
			return
		}
		rows, err := s.listBindings(c.Request.Context(), sess.AccountID)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, rows)
	}
}

func (s *Service) bindAuthorize(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := contextx.Session(c.Request.Context())
		if sess == nil {
			response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
			return
		}
		provider, err := parseProvider(c.Param("provider"))
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
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
		payload := fmt.Sprintf("{\"account_type\":\"%s\",\"provider\":\"%s\",\"intent\":\"BIND\",\"account_id\":\"%s\"}",
			string(accountType), string(provider), sess.AccountID)
		_ = s.rdb.Set(c.Request.Context(), stateKeyPrefix+state, payload, 10*time.Minute).Err()
		u, err := buildAuthorizeURL(provider, pc, state)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, AuthorizeResult{AuthorizeURL: u, State: state})
	}
}

func (s *Service) unbind(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := contextx.Session(c.Request.Context())
		if sess == nil {
			response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
			return
		}
		provider, err := parseProvider(c.Param("provider"))
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.assertCanUnbind(c.Request.Context(), sess.AccountID, string(provider)); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.db.WithContext(c.Request.Context()).
			Where("account_id = ? AND provider = ?", sess.AccountID, string(provider)).
			Delete(&AccountOAuthBinding{}).Error; err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) adminUnbind(c *gin.Context) {
	var req AdminUnbindParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	provider, err := parseProvider(req.Provider)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.assertCanUnbind(c.Request.Context(), req.AccountID, string(provider)); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.db.WithContext(c.Request.Context()).
		Where("account_id = ? AND provider = ?", req.AccountID, string(provider)).
		Delete(&AccountOAuthBinding{}).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) listBindings(ctx context.Context, accountID string) ([]BindingResult, error) {
	var rows []AccountOAuthBinding
	if err := s.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Order("provider asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]BindingResult, 0, len(rows))
	for _, r := range rows {
		p, err := parseProvider(r.Provider)
		label := r.Provider
		if err == nil {
			label = providerLabel(p)
		}
		out = append(out, BindingResult{
			Provider:     r.Provider,
			Label:        label,
			OpenIDMasked: maskOpenID(r.OpenID),
			Nickname:     r.Nickname,
			Avatar:       r.Avatar,
			BoundAt:      r.BoundAt,
		})
	}
	return out, nil
}

// createBinding 写入三方绑定关系（对齐 hei-boot upsertBinding）。
func (s *Service) createBinding(ctx context.Context, accountID, provider string, profile *oauthProfile) error {
	p, err := parseProvider(provider)
	if err != nil {
		return err
	}
	return s.upsertBinding(ctx, accountID, p, profile)
}

func isWechatFamily(provider string) bool {
	switch strings.ToUpper(strings.TrimSpace(provider)) {
	case "WECHAT", "WECHAT_OPEN", "WECHAT_MP":
		return true
	default:
		return false
	}
}

func (s *Service) resolveBinding(ctx context.Context, provider oauthProvider, profile *oauthProfile) (*AccountOAuthBinding, error) {
	providerName := string(provider)
	if profile.Provider != "" {
		providerName = profile.Provider
	}
	if isWechatFamily(providerName) && strings.TrimSpace(profile.UnionID) != "" {
		var byUnion AccountOAuthBinding
		err := s.db.WithContext(ctx).Where("union_id = ?", strings.TrimSpace(profile.UnionID)).First(&byUnion).Error
		if err == nil {
			return &byUnion, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	var binding AccountOAuthBinding
	err := s.db.WithContext(ctx).
		Where("provider = ? AND open_id = ?", providerName, profile.OpenID).
		First(&binding).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &binding, nil
}

func (s *Service) upsertBinding(ctx context.Context, accountID string, provider oauthProvider, profile *oauthProfile) error {
	return s.upsertBindingDB(s.db.WithContext(ctx), accountID, provider, profile)
}

func (s *Service) upsertBindingDB(db *gorm.DB, accountID string, provider oauthProvider, profile *oauthProfile) error {
	providerName := string(provider)
	if profile.Provider != "" {
		providerName = profile.Provider
	}
	now := time.Now().UTC()
	raw := profile.Raw
	if raw == "" {
		raw = "{}"
	}
	var unionID *string
	if strings.TrimSpace(profile.UnionID) != "" {
		unionID = strPtr(strings.TrimSpace(profile.UnionID))
	}
	var existing AccountOAuthBinding
	err := db.
		Where("account_id = ? AND provider = ?", accountID, providerName).
		First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		row := AccountOAuthBinding{
			ID:         idgen.Next(),
			AccountID:  accountID,
			Provider:   providerName,
			OpenID:     profile.OpenID,
			UnionID:    unionID,
			Nickname:   strPtr(profile.Nickname),
			Avatar:     strPtr(profile.Avatar),
			RawProfile: raw,
			BoundAt:    now,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		return db.Create(&row).Error
	}
	if err != nil {
		return err
	}
	updates := map[string]any{
		"open_id":     profile.OpenID,
		"nickname":    strPtr(profile.Nickname),
		"avatar":      strPtr(profile.Avatar),
		"raw_profile": raw,
		"bound_at":    now,
		"updated_at":  now,
	}
	if unionID != nil {
		updates["union_id"] = *unionID
	}
	return db.Model(&existing).Updates(updates).Error
}

func (s *Service) assertCanUnbind(ctx context.Context, accountID, provider string) error {
	var oauthCount int64
	if err := s.db.WithContext(ctx).Model(&AccountOAuthBinding{}).
		Where("account_id = ?", accountID).Count(&oauthCount).Error; err != nil {
		return err
	}
	var hasAccount, hasEmail, hasPhone int64
	_ = s.db.WithContext(ctx).Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ? AND bind_status = ?", accountID, "ACCOUNT", "BOUND").
		Count(&hasAccount).Error
	_ = s.db.WithContext(ctx).Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ? AND bind_status = ?", accountID, "EMAIL", "BOUND").
		Count(&hasEmail).Error
	_ = s.db.WithContext(ctx).Table("sys_account_identity").
		Where("account_id = ? AND identity_type = ? AND bind_status = ?", accountID, "PHONE", "BOUND").
		Count(&hasPhone).Error
	other := int(hasAccount+hasEmail+hasPhone) + int(oauthCount) - 1
	if other <= 0 {
		return fmt.Errorf("无法解绑：请至少保留一种登录方式")
	}
	return nil
}

func maskOpenID(openID string) string {
	if len(openID) <= 8 {
		return strings.Repeat("*", len(openID))
	}
	return openID[:4] + strings.Repeat("*", len(openID)-8) + openID[len(openID)-4:]
}
