// internal/framework/platform/audit/module_test.go BuildModule 单测。
//
// Author: Charlie

package audit

import "testing"

func TestBuildModule(t *testing.T) {
	cases := map[string]string{
		"resources":   "resource",
		"sys_banner":  "iam",
		"iam_account": "iam",
		"auth":        "iam",
	}
	for in, want := range cases {
		if got := BuildModule(in); got != want {
			t.Fatalf("BuildModule(%q) = %q, want %q", in, got, want)
		}
	}
}
