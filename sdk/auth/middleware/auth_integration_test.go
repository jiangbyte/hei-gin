package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hei-gin/sdk/auth"
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

func setupMiddlewareAuthTest(t *testing.T) (*auth.Realm, *stubPermissionAPI) {
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
	auth.RegisterPermissionProvider(stub)
	realm := auth.NewRealm(auth.BusinessID, 120, "Authorization")

	t.Cleanup(func() {
		auth.RegisterPermissionProvider(nil)
		if db.Redis != nil {
			_ = db.Redis.Close()
			db.Redis = nil
		}
	})

	return realm, stub
}

func decodeResponseBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return payload
}

func TestRealmClaimsAndScopeFor(t *testing.T) {
	realm, _ := setupMiddlewareAuthTest(t)

	token, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	c := httptest.NewRecorder()
	gc, _ := gin.CreateTestContext(c)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(realm.GetTokenName(), token)
	gc.Request = req

	claims, ok := realm.Claims(gc)
	if !ok {
		t.Fatal("claims not found")
	}
	if claims.UserID != "u1" {
		t.Fatalf("unexpected claims user_id: %q", claims.UserID)
	}
	scope, ok := realm.ScopeFor(gc, "sys:user:view")
	if !ok {
		t.Fatal("scope not found")
	}
	if scope.GroupScope != "SELF" || scope.OrgScope != "SELF" {
		t.Fatalf("unexpected scope: %#v", scope)
	}
}

func TestCheckLoginMiddlewareAttachesContext(t *testing.T) {
	realm, _ := setupMiddlewareAuthTest(t)
	token, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	r := gin.New()
	r.GET("/protected", CheckLogin(realm), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"login_id":   c.GetString("login_id"),
			"login_type": c.GetString("login_type"),
			"login_user": c.MustGet("loginUser"),
			"ctx_login":  c.Request.Context().Value(db.CtxKeyLoginID{}),
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(realm.GetTokenName(), token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}

	body := decodeResponseBody(t, w.Body.Bytes())
	if body["login_id"] != "u1" {
		t.Fatalf("unexpected login_id: %#v", body)
	}
	if body["login_type"] != string(realm.ID) {
		t.Fatalf("unexpected login_type: %#v", body)
	}
	if body["login_user"] != "alice" {
		t.Fatalf("unexpected login_user: %#v", body)
	}
	if body["ctx_login"] != "u1" {
		t.Fatalf("unexpected ctx_login: %#v", body)
	}
}

func TestCheckPermissionMiddlewareRejectsAndAllows(t *testing.T) {
	realm, _ := setupMiddlewareAuthTest(t)
	token, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	r := gin.New()
	r.GET("/allow", CheckPermission(realm, []string{"sys:user:view"}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.GET("/deny", CheckPermission(realm, []string{"sys:user:delete"}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	reqAllow := httptest.NewRequest(http.MethodGet, "/allow", nil)
	reqAllow.Header.Set(realm.GetTokenName(), token)
	wAllow := httptest.NewRecorder()
	r.ServeHTTP(wAllow, reqAllow)
	if wAllow.Code != http.StatusOK {
		t.Fatalf("allow status = %d body=%s", wAllow.Code, wAllow.Body.String())
	}

	reqDeny := httptest.NewRequest(http.MethodGet, "/deny", nil)
	reqDeny.Header.Set(realm.GetTokenName(), token)
	wDeny := httptest.NewRecorder()
	r.ServeHTTP(wDeny, reqDeny)
	if wDeny.Code != http.StatusForbidden {
		t.Fatalf("deny status = %d body=%s", wDeny.Code, wDeny.Body.String())
	}
}

func TestCheckRoleMiddlewareRejectsAndAllows(t *testing.T) {
	realm, _ := setupMiddlewareAuthTest(t)
	token, err := realm.Login(nil, "u1", map[string]any{"username": "alice"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	r := gin.New()
	r.GET("/allow", CheckRole(realm, []string{"ADMIN"}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.GET("/deny", CheckRole(realm, []string{"AUDITOR"}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	reqAllow := httptest.NewRequest(http.MethodGet, "/allow", nil)
	reqAllow.Header.Set(realm.GetTokenName(), token)
	wAllow := httptest.NewRecorder()
	r.ServeHTTP(wAllow, reqAllow)
	if wAllow.Code != http.StatusOK {
		t.Fatalf("allow status = %d body=%s", wAllow.Code, wAllow.Body.String())
	}

	reqDeny := httptest.NewRequest(http.MethodGet, "/deny", nil)
	reqDeny.Header.Set(realm.GetTokenName(), token)
	wDeny := httptest.NewRecorder()
	r.ServeHTTP(wDeny, reqDeny)
	if wDeny.Code != http.StatusForbidden {
		t.Fatalf("deny status = %d body=%s", wDeny.Code, wDeny.Body.String())
	}
}
