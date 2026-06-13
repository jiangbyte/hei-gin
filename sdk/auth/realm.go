package auth

import (
	"context"

	"hei-gin/sdk/shared/contracts"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type RealmID = contracts.RealmID

const (
	BusinessID RealmID = "BUSINESS"
	ConsumerID RealmID = "CONSUMER"
)

type ACLSnapshot struct {
	Permissions []string             `json:"permissions"`
	Roles       []string             `json:"roles"`
	ScopeMap    map[string]ScopeInfo `json:"scope_map"`
}

type SessionClaims struct {
	UserID    string         `json:"user_id"`
	RealmID   RealmID        `json:"realm_id"`
	CreatedAt string         `json:"created_at"`
	Extra     map[string]any `json:"extra"`
	ACL       ACLSnapshot    `json:"acl"`
}

type Realm struct {
	ID   RealmID
	tool *baseAuthTool
}

func NewRealm(id RealmID, expire int, tokenName string) *Realm {
	realm := &Realm{
		ID:   id,
		tool: newBaseAuthTool(id),
	}
	realm.Init(expire, tokenName)
	return realm
}

func (r *Realm) Init(expire int, tokenName string) {
	if r == nil {
		return
	}
	r.tool.Init(expire, tokenName)
}

func (r *Realm) Login(c *gin.Context, id string, extra map[string]any) (string, error) {
	return r.tool.Login(c, id, extra)
}

func (r *Realm) Logout(c *gin.Context, loginID ...string) {
	r.tool.Logout(c, loginID...)
}

func (r *Realm) IsLogin(c *gin.Context) bool {
	return r.tool.IsLogin(c)
}

func (r *Realm) CheckLogin(c *gin.Context) error {
	return r.tool.CheckLogin(c)
}

func (r *Realm) GetLoginID(c *gin.Context) string {
	return r.tool.GetLoginID(c)
}

func (r *Realm) GetLoginIDDefaultNull(c *gin.Context) string {
	return r.tool.GetLoginIDDefaultNull(c)
}

func (r *Realm) GetLoginIDByToken(token string) string {
	return r.tool.GetLoginIDByToken(token)
}

func (r *Realm) GetTokenName() string {
	return r.tool.GetTokenName()
}

func (r *Realm) GetTokenValue(c *gin.Context) string {
	return r.tool.GetTokenValue(c)
}

func (r *Realm) GetExtra(c *gin.Context, key string) any {
	return r.tool.GetExtra(c, key)
}

func (r *Realm) GetSession(c *gin.Context) map[string]any {
	return r.tool.GetSession(c)
}

func (r *Realm) RenewTimeout(c *gin.Context, timeout ...int) {
	r.tool.RenewTimeout(c, timeout...)
}

func (r *Realm) GetTokenTimeout(c *gin.Context) int {
	return r.tool.GetTokenTimeout(c)
}

func (r *Realm) GetSessionTimeout(c *gin.Context) int {
	return r.tool.GetSessionTimeout(c)
}

func (r *Realm) Kickout(loginID string) {
	r.tool.Kickout(loginID)
}

func (r *Realm) KickoutWithContext(ctx context.Context, loginID string) {
	r.tool.KickoutWithContext(ctx, loginID)
}

func (r *Realm) KickoutToken(loginID, token string) {
	r.tool.KickoutToken(loginID, token)
}

func (r *Realm) KickoutTokenWithContext(ctx context.Context, loginID, token string) {
	r.tool.KickoutTokenWithContext(ctx, loginID, token)
}

func (r *Realm) Disable(loginID string, timeSeconds int) {
	r.tool.Disable(loginID, timeSeconds)
}

func (r *Realm) IsDisable(loginID string) bool {
	return r.tool.IsDisable(loginID)
}

func (r *Realm) CheckDisable(loginID string) error {
	return r.tool.CheckDisable(loginID)
}

func (r *Realm) GetDisableTime(loginID string) int {
	return r.tool.GetDisableTime(loginID)
}

func (r *Realm) UntieDisable(loginID string) {
	r.tool.UntieDisable(loginID)
}

func (r *Realm) HasPermission(c *gin.Context, permission string) bool {
	permissions, err := r.PermissionList(c)
	if err != nil {
		return false
	}
	return MatchPermission(permission, permissions)
}

func (r *Realm) HasPermissionAnd(c *gin.Context, permissions ...string) bool {
	perms, err := r.PermissionList(c)
	if err != nil {
		return false
	}
	return MatchPermissionsAnd(permissions, perms)
}

func (r *Realm) HasPermissionOr(c *gin.Context, permissions ...string) bool {
	perms, err := r.PermissionList(c)
	if err != nil {
		return false
	}
	return MatchPermissionsOr(permissions, perms)
}

func (r *Realm) HasRole(c *gin.Context, role string) bool {
	roles, err := r.RoleList(c)
	if err != nil {
		return false
	}
	for _, item := range roles {
		if item == role {
			return true
		}
	}
	return false
}

func (r *Realm) HasRoleAnd(c *gin.Context, roles ...string) bool {
	userRoles, err := r.RoleList(c)
	if err != nil {
		return false
	}
	roleSet := make(map[string]struct{}, len(userRoles))
	for _, item := range userRoles {
		roleSet[item] = struct{}{}
	}
	for _, item := range roles {
		if _, ok := roleSet[item]; !ok {
			return false
		}
	}
	return true
}

func (r *Realm) HasRoleOr(c *gin.Context, roles ...string) bool {
	userRoles, err := r.RoleList(c)
	if err != nil {
		return false
	}
	roleSet := make(map[string]struct{}, len(userRoles))
	for _, item := range userRoles {
		roleSet[item] = struct{}{}
	}
	for _, item := range roles {
		if _, ok := roleSet[item]; ok {
			return true
		}
	}
	return false
}

func (r *Realm) PermissionList(c *gin.Context) ([]string, error) {
	if r == nil {
		return []string{}, nil
	}
	if claims, ok := r.Claims(c); ok {
		return claims.ACL.Permissions, nil
	}
	return []string{}, nil
}

func (r *Realm) RoleList(c *gin.Context) ([]string, error) {
	if r == nil {
		return []string{}, nil
	}
	if claims, ok := r.Claims(c); ok {
		return claims.ACL.Roles, nil
	}
	return []string{}, nil
}

func (r *Realm) ScopeFor(c *gin.Context, permission string) (ScopeInfo, bool) {
	if r == nil || permission == "" {
		return ScopeInfo{}, false
	}
	claims, ok := r.Claims(c)
	if !ok {
		return ScopeInfo{}, false
	}
	scope, exists := claims.ACL.ScopeMap[permission]
	return scope, exists
}

func (r *Realm) Claims(c *gin.Context) (*SessionClaims, bool) {
	return r.tool.GetClaims(c)
}

func (r *Realm) Sessions() *RealmSessionService {
	return &RealmSessionService{realm: r}
}

func (r *Realm) RefreshUserSessionsACL(ctx context.Context, userID string) error {
	if r == nil {
		return nil
	}
	return r.tool.refreshUserSessionsACL(ctx, userID)
}

func (r *Realm) RefreshACL(ctx context.Context, userID string) (ACLSnapshot, error) {
	if r == nil {
		return ACLSnapshot{
			Permissions: []string{},
			Roles:       []string{},
			ScopeMap:    map[string]ScopeInfo{},
		}, nil
	}
	return r.tool.loadACL(ctx, userID)
}

func (r *Realm) CheckPermission(c *gin.Context, permission string) {
	if !r.HasPermission(c, permission) {
		result.Failure(c, "缺少权限: "+permission, 403)
	}
}

func (r *Realm) CheckRole(c *gin.Context, role string) {
	if !r.HasRole(c, role) {
		result.Failure(c, "缺少角色: "+role, 403)
	}
}
