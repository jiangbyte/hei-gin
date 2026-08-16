// internal/framework/core/security/session_test.go
//
// Author: Charlie

package security

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestSessionStore(t *testing.T) (*SessionStore, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	t.Cleanup(mr.Close)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	return NewSessionStore(rdb), mr
}

func samplePayload(token string, accountType AccountType, accountID string) *SessionPayload {
	now := time.Now().UTC()
	return &SessionPayload{
		Token:        token,
		AccountID:    accountID,
		AccountType:  accountType,
		LoginAt:      now,
		LastActiveAt: now,
		ExpiresAt:    now.Add(time.Hour),
	}
}

func TestSessionStoreLoginKeys(t *testing.T) {
	store, mr := newTestSessionStore(t)
	ctx := context.Background()
	p := samplePayload("tok-a", AccountAdmin, "acc-1")
	if err := store.Save(ctx, p, time.Hour); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if !mr.Exists("login:token:tok-a") {
		t.Fatal("expected login:token:tok-a")
	}
	if !mr.Exists("login:account:ADMIN:acc-1") {
		t.Fatal("expected login:account:ADMIN:acc-1")
	}
	if !mr.Exists("login:tokens") {
		t.Fatal("expected login:tokens")
	}
	if mr.Exists("hei:session:token:tok-a") {
		t.Fatal("must not write legacy hei:session:* keys")
	}
}

func TestSessionStoreTypedIndexAndMGet(t *testing.T) {
	store, _ := newTestSessionStore(t)
	ctx := context.Background()
	admin := samplePayload("tok-admin", AccountAdmin, "same-id")
	portal := samplePayload("tok-portal", AccountPortal, "same-id")
	if err := store.Save(ctx, admin, time.Hour); err != nil {
		t.Fatalf("Save admin: %v", err)
	}
	if err := store.Save(ctx, portal, time.Hour); err != nil {
		t.Fatalf("Save portal: %v", err)
	}

	adminTokens, err := store.ListTokensForAccount(ctx, AccountAdmin, "same-id")
	if err != nil || len(adminTokens) != 1 || adminTokens[0] != "tok-admin" {
		t.Fatalf("admin tokens=%v err=%v", adminTokens, err)
	}
	portalTokens, err := store.ListTokensForAccount(ctx, AccountPortal, "same-id")
	if err != nil || len(portalTokens) != 1 || portalTokens[0] != "tok-portal" {
		t.Fatalf("portal tokens=%v err=%v", portalTokens, err)
	}

	n, err := store.CountTokens(ctx)
	if err != nil || n != 2 {
		t.Fatalf("CountTokens=%d err=%v", n, err)
	}
	all, err := store.ListTokens(ctx)
	if err != nil || len(all) != 2 {
		t.Fatalf("ListTokens=%v err=%v", all, err)
	}
	sessions, err := store.ListSessionsByTokens(ctx, all)
	if err != nil || len(sessions) != 2 {
		t.Fatalf("ListSessionsByTokens=%d err=%v", len(sessions), err)
	}
}

func TestSessionStoreDeleteByType(t *testing.T) {
	store, mr := newTestSessionStore(t)
	ctx := context.Background()
	_ = store.Save(ctx, samplePayload("tok-admin", AccountAdmin, "same-id"), time.Hour)
	_ = store.Save(ctx, samplePayload("tok-portal", AccountPortal, "same-id"), time.Hour)

	if err := store.DeleteAllForAccount(ctx, AccountAdmin, "same-id"); err != nil {
		t.Fatalf("DeleteAllForAccount: %v", err)
	}
	if mr.Exists("login:token:tok-admin") {
		t.Fatal("admin token should be deleted")
	}
	if !mr.Exists("login:token:tok-portal") {
		t.Fatal("portal token should remain")
	}
	n, _ := store.CountTokens(ctx)
	if n != 1 {
		t.Fatalf("CountTokens after typed delete want 1 got %d", n)
	}
}

func TestSessionStoreDeleteCleansIndexes(t *testing.T) {
	store, mr := newTestSessionStore(t)
	ctx := context.Background()
	_ = store.Save(ctx, samplePayload("tok-x", AccountAdmin, "acc-x"), time.Hour)
	if err := store.Delete(ctx, "tok-x"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if mr.Exists("login:token:tok-x") {
		t.Fatal("token key should be gone")
	}
	members, _ := mr.SMembers("login:tokens")
	if len(members) != 0 {
		t.Fatalf("global set should be empty, got %v", members)
	}
	members, _ = mr.SMembers("login:account:ADMIN:acc-x")
	if len(members) != 0 {
		t.Fatalf("account set should be empty, got %v", members)
	}
}

func TestSessionStoreListSessionsCleansStale(t *testing.T) {
	store, mr := newTestSessionStore(t)
	ctx := context.Background()
	_ = store.Save(ctx, samplePayload("alive", AccountAdmin, "a1"), time.Hour)
	if _, err := mr.SAdd("login:tokens", "stale-token"); err != nil {
		t.Fatalf("seed stale: %v", err)
	}

	sessions, err := store.ListSessionsByTokens(ctx, []string{"alive", "stale-token"})
	if err != nil {
		t.Fatalf("ListSessionsByTokens: %v", err)
	}
	if len(sessions) != 1 || sessions[0].Token != "alive" {
		t.Fatalf("sessions=%v", sessions)
	}
	members, _ := mr.SMembers("login:tokens")
	for _, m := range members {
		if m == "stale-token" {
			t.Fatal("stale token should be removed from login:tokens")
		}
	}
}
