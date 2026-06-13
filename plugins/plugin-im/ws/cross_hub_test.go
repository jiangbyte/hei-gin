package ws

import (
	"testing"

	"hei-gin/sdk/auth"
)

func TestCrossHubKeyHelpers(t *testing.T) {
	ch := &CrossHub{instanceID: "inst-1"}

	if got := ch.userSetKey(auth.ConsumerID, "u1"); got != "ws:user:CONSUMER:u1" {
		t.Fatalf("userSetKey = %q", got)
	}
	if got := ch.userCountKey(auth.BusinessID, "u2"); got != "ws:usercnt:BUSINESS:u2" {
		t.Fatalf("userCountKey = %q", got)
	}
	if got := ch.dedupKey("m1"); got != "ws:dedup:inst-1:m1" {
		t.Fatalf("dedupKey = %q", got)
	}
	if got := userCountKeyFromSetKey("ws:user:CONSUMER:u1"); got != "ws:usercnt:CONSUMER:u1" {
		t.Fatalf("userCountKeyFromSetKey = %q", got)
	}
	if got := userCountKeyFromSetKey("ws:usercnt:CONSUMER:u1"); got != "" {
		t.Fatalf("unexpected count key for non-user-set key: %q", got)
	}
}
