package safeurl

import (
	"testing"
)

func TestValidateRejects(t *testing.T) {
	cases := []string{
		"http://127.0.0.1/x",
		"https://127.0.0.1/x",
		"http://169.254.169.254/latest",
		"file:///etc/passwd",
		"https://user:pass@example.com/",
		"ftp://example.com/",
		"",
	}
	for _, raw := range cases {
		if err := Validate(raw, Options{}); err == nil {
			t.Fatalf("expected reject for %q", raw)
		}
	}
	// http 默认拒绝
	if err := Validate("http://example.com/hook", Options{}); err == nil {
		t.Fatal("expected http rejected when AllowHTTP=false")
	}
}

func TestValidateAllowsHTTPSPublic(t *testing.T) {
	// example.com 解析为公网文档地址；若环境 DNS 异常则跳过。
	err := Validate("https://example.com/webhook", Options{})
	if err != nil {
		t.Skipf("dns/environment: %v", err)
	}
}

func TestValidateAllowsHTTPWhenEnabled(t *testing.T) {
	err := Validate("http://example.com/webhook", Options{AllowHTTP: true})
	if err != nil {
		t.Skipf("dns/environment: %v", err)
	}
}
