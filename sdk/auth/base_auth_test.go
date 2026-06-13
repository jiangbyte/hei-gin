package auth

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/shared/contracts"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type stubPermissionAPI struct {
	perms map[string][]string
	roles map[string][]string
	scope map[string]map[string]contracts.ScopeInfo
}

func (s *stubPermissionAPI) GetPermissionList(ctx context.Context, realmID contracts.RealmID, userID string) ([]string, error) {
	_ = ctx
	_ = realmID
	return append([]string{}, s.perms[userID]...), nil
}

func (s *stubPermissionAPI) GetRoleList(ctx context.Context, realmID contracts.RealmID, userID string) ([]string, error) {
	_ = ctx
	_ = realmID
	return append([]string{}, s.roles[userID]...), nil
}

func (s *stubPermissionAPI) GetPermissionScopeMap(ctx context.Context, realmID contracts.RealmID, userID string) (map[string]contracts.ScopeInfo, error) {
	_ = ctx
	_ = realmID
	result := make(map[string]contracts.ScopeInfo)
	for code, info := range s.scope[userID] {
		result[code] = info
	}
	return result, nil
}

func setupAuthTest(t *testing.T) (*miniredis.Miniredis, *stubPermissionAPI, *Realm) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	mr := miniredis.RunT(t)
	db.Redis = redis.NewClient(&redis.Options{Addr: mr.Addr()})

	stub := &stubPermissionAPI{
		perms: map[string][]string{
			"u1": {"sys:user:view"},
		},
		roles: map[string][]string{
			"u1": {"ADMIN"},
		},
		scope: map[string]map[string]contracts.ScopeInfo{
			"u1": {
				"sys:user:view": {
					GroupScope: "SELF",
					OrgScope:   "SELF",
				},
			},
		},
	}
	PermissionDelegate = stub

	realm := &Realm{
		ID: BusinessID,
		tool: &baseAuthTool{
			realmID:   BusinessID,
			expire:    120,
			tokenName: "Authorization",
		},
	}

	t.Cleanup(func() {
		PermissionDelegate = nil
		if db.Redis != nil {
			_ = db.Redis.Close()
			db.Redis = nil
		}
	})

	return mr, stub, realm
}

func newAuthContext(tokenName, token string) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	if token != "" {
		req.Header.Set(tokenName, token)
	}
	c.Request = req
	return c
}

func TestRefreshUserSessionsACLUpdatesAllLiveTokens(t *testing.T) {
	_, stub, realm := setupAuthTest(t)

	token1, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login token1: %v", err)
	}
	token2, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login token2: %v", err)
	}

	stub.perms["u1"] = []string{"sys:user:view", "sys:user:edit"}
	stub.roles["u1"] = []string{"ADMIN", "EDITOR"}
	stub.scope["u1"] = map[string]contracts.ScopeInfo{
		"sys:user:view": {GroupScope: "SELF", OrgScope: "SELF"},
		"sys:user:edit": {GroupScope: "ORG", OrgScope: "ORG"},
	}

	if err := realm.RefreshUserSessionsACL(context.Background(), "u1"); err != nil {
		t.Fatalf("refresh acl: %v", err)
	}

	claims1, ok := realm.tool.getClaims(token1)
	if !ok {
		t.Fatal("claims1 not found")
	}
	claims2, ok := realm.tool.getClaims(token2)
	if !ok {
		t.Fatal("claims2 not found")
	}

	if len(claims1.ACL.Permissions) != 2 || len(claims2.ACL.Permissions) != 2 {
		t.Fatalf("unexpected permissions after refresh: %#v %#v", claims1.ACL.Permissions, claims2.ACL.Permissions)
	}
	if claims1.ACL.Permissions[0] != "sys:user:edit" || claims1.ACL.Permissions[1] != "sys:user:view" {
		t.Fatalf("permissions not sorted: %#v", claims1.ACL.Permissions)
	}
	if claims1.ACL.Roles[0] != "ADMIN" || claims1.ACL.Roles[1] != "EDITOR" {
		t.Fatalf("roles not sorted: %#v", claims1.ACL.Roles)
	}
	if _, ok := claims1.ACL.ScopeMap["sys:user:edit"]; !ok {
		t.Fatalf("scope map not refreshed: %#v", claims1.ACL.ScopeMap)
	}
}

func TestKickoutTokenWithContextRemovesOnlyTargetToken(t *testing.T) {
	_, _, realm := setupAuthTest(t)

	token1, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login token1: %v", err)
	}
	token2, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login token2: %v", err)
	}

	realm.KickoutTokenWithContext(context.Background(), "u1", token1)

	if _, ok := realm.tool.getClaims(token1); ok {
		t.Fatal("token1 should be removed")
	}
	if _, ok := realm.tool.getClaims(token2); !ok {
		t.Fatal("token2 should remain")
	}

	tokens, err := db.Redis.SMembers(context.Background(), realm.tool.getSessionKey("u1")).Result()
	if err != nil {
		t.Fatalf("read session tokens: %v", err)
	}
	if len(tokens) != 1 || tokens[0] != token2 {
		t.Fatalf("unexpected session tokens after kickout: %#v", tokens)
	}
}

func TestRenewTimeoutExtendsTokenAndSessionTTL(t *testing.T) {
	mr, _, realm := setupAuthTest(t)

	token, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	ctx := newAuthContext(realm.GetTokenName(), token)
	beforeTokenTTL := mr.TTL(realm.tool.getTokenKey(token))
	beforeSessionTTL := mr.TTL(realm.tool.getSessionKey("u1"))

	realm.RenewTimeout(ctx, 300)

	afterTokenTTL := mr.TTL(realm.tool.getTokenKey(token))
	afterSessionTTL := mr.TTL(realm.tool.getSessionKey("u1"))

	if !(afterTokenTTL > beforeTokenTTL) {
		t.Fatalf("token ttl not extended: before=%v after=%v", beforeTokenTTL, afterTokenTTL)
	}
	if !(afterSessionTTL > beforeSessionTTL) {
		t.Fatalf("session ttl not extended: before=%v after=%v", beforeSessionTTL, afterSessionTTL)
	}

	expiryScore, err := db.Redis.ZScore(context.Background(), realm.tool.getTokenExpiryIndexKey(), token).Result()
	if err != nil {
		t.Fatalf("read token expiry index: %v", err)
	}
	if expiryScore <= float64(time.Now().Unix()) {
		t.Fatalf("token expiry index not updated: %v", expiryScore)
	}
}
