// internal/framework/platform/audit/registry_test.go 审计注册表单测。
//
// Author: Charlie

package audit

import "testing"

func TestRegistryMatch(t *testing.T) {
	reg := NewRegistry()
	reg.Register("POST", "/api/v1/admin/sys/banners/create", "sys_banner", "create")
	reg.Register("POST", "/api/v1/admin/sys/banners/update", "sys_banner", "update")
	reg.Register("POST", "/api/v1/admin/oauth/*/bind/authorize", "auth", "oauth_bind_authorize")
	reg.Register("POST", "/api/v1/portal/sys/feedbacks/submit", "sys_feedback", "submit")

	cases := []struct {
		name   string
		method string
		path   string
		want   AuditSpec
		ok     bool
	}{
		{"literal hit", "POST", "/api/v1/admin/sys/banners/create", AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/banners/create", ResourceType: "sys_banner", Action: "create"}, true},
		{"lowercase method normalized", "post", "/api/v1/admin/sys/banners/update", AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/banners/update", ResourceType: "sys_banner", Action: "update"}, true},
		{"path param glob", "POST", "/api/v1/admin/oauth/wechat/bind/authorize", AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/oauth/*/bind/authorize", ResourceType: "auth", Action: "oauth_bind_authorize"}, true},
		{"portal path", "POST", "/api/v1/portal/sys/feedbacks/submit", AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/sys/feedbacks/submit", ResourceType: "sys_feedback", Action: "submit"}, true},
		{"wrong method", "GET", "/api/v1/admin/sys/banners/create", AuditSpec{}, false},
		{"wrong path", "POST", "/api/v1/admin/sys/banners/list", AuditSpec{}, false},
		{"unregistered", "DELETE", "/api/v1/admin/sys/banners/delete", AuditSpec{}, false},
		{"empty registry", "POST", "/api/v1/x", AuditSpec{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := reg.Match(tc.method, tc.path)
			if ok != tc.ok {
				t.Fatalf("Match(%s %s) ok = %v, want %v", tc.method, tc.path, ok, tc.ok)
			}
			if !ok {
				return
			}
			if got != tc.want {
				t.Fatalf("Match(%s %s) = %+v, want %+v", tc.method, tc.path, got, tc.want)
			}
		})
	}
}

func TestRegistryNilSafe(t *testing.T) {
	var reg *Registry
	reg.Register("POST", "/x", "a", "b") // must not panic
	if _, ok := reg.Match("POST", "/x"); ok {
		t.Fatal("nil registry should not match")
	}
	if got := reg.All(); got != nil {
		t.Fatalf("nil registry All() = %v, want nil", got)
	}
}

func TestBuildModule(t *testing.T) {
	cases := map[string]string{
		"resources":      "resource",
		"sys_banner":     "iam",
		"iam_account":    "iam",
		"auth":           "iam",
		"profile_center": "iam",
		"":               "iam",
	}
	for in, want := range cases {
		if got := BuildModule(in); got != want {
			t.Fatalf("BuildModule(%q) = %q, want %q", in, got, want)
		}
	}
}
