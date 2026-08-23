// internal/framework/middleware/audit_test.go 操作审计中间件单测。
//
// Author: Charlie

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/audit"
)

type fakePublisher struct{ events []audit.Event }

func (f *fakePublisher) Publish(ev audit.Event) { f.events = append(f.events, ev) }

// withSession 把会话注入请求上下文（模拟 AuthContext 已通过）。
func withSession(sess *security.SessionPayload) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := contextx.WithSession(c.Request.Context(), sess)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func TestOperationAuditPublishOnSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pub := &fakePublisher{}

	r := gin.New()
	r.Use(withSession(&security.SessionPayload{AccountID: "1", AccountType: security.AccountAdmin}))
	r.POST("/api/v1/admin/sys/banners/create", OperationAudit(pub, "sys_banner", "create"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/admin/sys/banners/create", nil))

	if len(pub.events) != 1 {
		t.Fatalf("published %d events, want 1", len(pub.events))
	}
	ev := pub.events[0]
	if ev.Module != "iam" || ev.ResourceType != "sys_banner" || ev.Action != "create" {
		t.Fatalf("unexpected event identity: %+v", ev)
	}
	if ev.Summary != "POST /api/v1/admin/sys/banners/create" {
		t.Fatalf("summary = %q, want %q", ev.Summary, "POST /api/v1/admin/sys/banners/create")
	}
	if ev.Success != true {
		t.Fatalf("success = %v, want true", ev.Success)
	}
	if ev.AccountID != "1" || ev.AccountType != "admin" {
		t.Fatalf("account = %q/%q, want 1/admin", ev.AccountID, ev.AccountType)
	}
	if ev.OccurredAt.IsZero() {
		t.Fatal("OccurredAt should be set")
	}
}

func TestOperationAuditRecordsFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pub := &fakePublisher{}

	r := gin.New()
	r.Use(withSession(&security.SessionPayload{AccountID: "1", AccountType: security.AccountAdmin}))
	r.POST("/api/v1/admin/sys/banners/create", OperationAudit(pub, "sys_banner", "create"), func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/admin/sys/banners/create", nil))
	if len(pub.events) != 1 {
		t.Fatalf("published %d events for failed request, want 1 (boot 对齐：失败也记审计)", len(pub.events))
	}
	if pub.events[0].Success {
		t.Fatalf("success = true, want false for 4xx")
	}
}

func TestOperationAuditSkipsWithoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pub := &fakePublisher{}

	r := gin.New()
	r.Use(withSession(&security.SessionPayload{AccountID: "1", AccountType: security.AccountAdmin}))
	r.GET("/api/v1/admin/sys/banners/page", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.POST("/api/v1/admin/sys/other", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/v1/admin/sys/banners/page", nil))
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/admin/sys/other", nil))
	if len(pub.events) != 0 {
		t.Fatalf("published %d events for unaudited paths, want 0", len(pub.events))
	}
}

func TestOperationAuditAnonymous(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pub := &fakePublisher{}

	r := gin.New()
	r.POST("/api/v1/admin/login", OperationAudit(pub, "auth", "login"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/admin/login", nil))
	if len(pub.events) != 1 {
		t.Fatalf("published %d events, want 1", len(pub.events))
	}
	if pub.events[0].AccountID != "" || pub.events[0].AccountType != "" {
		t.Fatalf("anonymous request should record empty account, got %+v", pub.events[0])
	}
	if pub.events[0].Module != "iam" {
		t.Fatalf("module = %q, want iam (buildModule)", pub.events[0].Module)
	}
}

func TestOperationAuditNilSafe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/x", OperationAudit(nil, "x", "y"), func(c *gin.Context) { c.Status(http.StatusOK) })
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/x", nil))
}
