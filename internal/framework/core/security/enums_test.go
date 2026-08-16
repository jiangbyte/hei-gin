// internal/framework/core/security/enums_test.go
//
// Author: Charlie

package security

import "testing"

func TestDeviceLabelFromUserAgent(t *testing.T) {
	cases := []struct {
		ua   string
		want string
		nil  bool
	}{
		{ua: "", nil: true},
		{ua: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0)", want: "Mobile"},
		{ua: "Mozilla/5.0 (Linux; Android 14)", want: "Mobile"},
		{ua: "Mozilla/5.0 (iPad; CPU OS 17_0)", want: "Tablet"},
		{ua: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)", want: "Desktop"},
	}
	for _, tc := range cases {
		got := DeviceLabelFromUserAgent(tc.ua)
		if tc.nil {
			if got != nil {
				t.Fatalf("ua=%q want nil got %v", tc.ua, *got)
			}
			continue
		}
		if got == nil || *got != tc.want {
			t.Fatalf("ua=%q want %q got %v", tc.ua, tc.want, got)
		}
	}
}
