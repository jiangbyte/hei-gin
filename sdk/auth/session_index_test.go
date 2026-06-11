package auth

import (
	"testing"
)

func TestParseTokenData(t *testing.T) {
	raw := `{"user_id":"u1","type":"BUSINESS","created_at":"2026-06-11 10:20:30","extra":{"username":"alice","device_type":"web","device_id":"d1"}}`

	data, ok := parseTokenData(raw)
	if !ok {
		t.Fatal("parseTokenData returned false")
	}
	if data.CreatedAtString != "2026-06-11 10:20:30" {
		t.Fatalf("CreatedAtString = %q", data.CreatedAtString)
	}
	if data.Username != "alice" || data.DeviceType != "web" || data.DeviceID != "d1" {
		t.Fatalf("unexpected extra fields: %#v", data)
	}
}

func TestDiffTokens(t *testing.T) {
	stale := diffTokens([]string{"a", "b", "c"}, []string{"b"})
	if len(stale) != 2 || stale[0] != "a" || stale[1] != "c" {
		t.Fatalf("diffTokens = %#v", stale)
	}
}
