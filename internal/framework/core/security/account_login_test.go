package security

import "testing"

func TestRequireAccountLogin(t *testing.T) {
	if _, err := RequireAccountLogin("admin-iam"); err == nil {
		t.Fatal("expected hyphen account to fail")
	}
	value, err := RequireAccountLogin(" admin_iam ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "admin_iam" {
		t.Fatalf("got %q", value)
	}
}

func TestSanitizeAccountBase(t *testing.T) {
	if got := SanitizeAccountBase("admin-iam"); got != "adminiam" {
		t.Fatalf("got %q", got)
	}
}
