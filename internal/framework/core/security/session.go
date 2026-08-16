// internal/framework/core/security/session.go 会话管理（对齐 hei-fastapi login:* Redis 键）。
//
// Author: Charlie

package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// PermissionGrant 为权限键及其数据范围。
//
// Author: Charlie
type PermissionGrant struct {
	PermissionKey      string    `json:"permission_key"`
	DataScope          DataScope `json:"data_scope"`
	CustomScopeDeptIDs []string  `json:"custom_scope_dept_ids"`
	SourceType         string    `json:"source_type"`
	SourceID           string    `json:"source_id"`
}

// Session 是会话载荷别名（API / datascope 常用简写）。
type Session = SessionPayload

// SessionPayload 存于 Redis（JSON）。
//
// Author: Charlie
type SessionPayload struct {
	Token                string            `json:"token"`
	AccountID            string            `json:"account_id"`
	AccountType          AccountType       `json:"account_type"`
	RoleIDs              []string          `json:"role_ids"`
	DeptIDs              []string          `json:"dept_ids"`
	GroupIDs             []string          `json:"group_ids"`
	ResourceIDs          []string          `json:"resource_ids"`
	PermissionKeys       []string          `json:"permission_keys"`
	PermissionGrants     []PermissionGrant `json:"permission_grants"`
	ClientResourceIDs    []string          `json:"client_resource_ids"`
	ClientPermissionKeys []string          `json:"client_permission_keys"`
	ClientIP             *string           `json:"client_ip"`
	UserAgent            *string           `json:"user_agent"`
	RememberMe           bool              `json:"remember_me"`
	PasswordExpired      bool              `json:"password_expired"`
	DeviceLabel          *string           `json:"device_label"`
	LoginAt              time.Time         `json:"login_at"`
	LastActiveAt         time.Time         `json:"last_active_at"`
	ExpiresAt            time.Time         `json:"expires_at"`
}

// SessionStore 管理不透明 Redis 会话（键对齐 fastapi app/core/cache/keys.py）。
//
// Author: Charlie
type SessionStore struct {
	rdb *redis.Client
}

// NewSessionStore 创建会话存储。
func NewSessionStore(rdb *redis.Client) *SessionStore {
	return &SessionStore{rdb: rdb}
}

// NewToken 生成随机不透明会话 token。
func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func loginTokenKey(token string) string { return "login:token:" + token }

func loginAccountTokensKey(accountType AccountType, accountID string) string {
	return "login:account:" + string(accountType) + ":" + accountID
}

func loginTokensKey() string { return "login:tokens" }

// Save 写入会话并维护账号索引与全局 token 集。
func (s *SessionStore) Save(ctx context.Context, p *SessionPayload, ttl time.Duration) error {
	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}
	acctKey := loginAccountTokensKey(p.AccountType, p.AccountID)
	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, loginTokenKey(p.Token), raw, ttl)
	pipe.SAdd(ctx, acctKey, p.Token)
	pipe.Expire(ctx, acctKey, ttl)
	pipe.SAdd(ctx, loginTokensKey(), p.Token)
	_, err = pipe.Exec(ctx)
	return err
}

// Get 按 token 读取会话，不存在时返回 nil。
func (s *SessionStore) Get(ctx context.Context, token string) (*SessionPayload, error) {
	raw, err := s.rdb.Get(ctx, loginTokenKey(token)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var p SessionPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Delete 删除单个会话并更新账号索引与全局集。
func (s *SessionStore) Delete(ctx context.Context, token string) error {
	p, err := s.Get(ctx, token)
	if err != nil {
		return err
	}
	pipe := s.rdb.Pipeline()
	pipe.Del(ctx, loginTokenKey(token))
	pipe.SRem(ctx, loginTokensKey(), token)
	if p != nil {
		pipe.SRem(ctx, loginAccountTokensKey(p.AccountType, p.AccountID), token)
	}
	_, err = pipe.Exec(ctx)
	return err
}

// DeleteAllForAccount 删除指定账号类型下的全部会话。
func (s *SessionStore) DeleteAllForAccount(ctx context.Context, accountType AccountType, accountID string) error {
	if accountType == "" || accountID == "" {
		return nil
	}
	return s.DeleteAccountsSessions(ctx, []AccountSessionTarget{{AccountType: accountType, AccountID: accountID}})
}

// DeleteAllForAccountAnyType 删除某账号在 ADMIN/PORTAL 两端的全部会话（授权变更等仅有 accountID 时用）。
func (s *SessionStore) DeleteAllForAccountAnyType(ctx context.Context, accountID string) error {
	if accountID == "" {
		return nil
	}
	return s.DeleteAccountsSessions(ctx, []AccountSessionTarget{
		{AccountType: AccountAdmin, AccountID: accountID},
		{AccountType: AccountPortal, AccountID: accountID},
	})
}

// AccountSessionTarget 按类型+账号定位会话。
//
// Author: Charlie
type AccountSessionTarget struct {
	AccountType AccountType
	AccountID   string
}

// DeleteAccountsSessions 批量删除多个账户的在线会话（对齐 fastapi delete_accounts_sessions）。
func (s *SessionStore) DeleteAccountsSessions(ctx context.Context, targets []AccountSessionTarget) error {
	if len(targets) == 0 {
		return nil
	}
	seen := map[string]AccountSessionTarget{}
	for _, t := range targets {
		if t.AccountID == "" || t.AccountType == "" {
			continue
		}
		seen[string(t.AccountType)+":"+t.AccountID] = t
	}
	if len(seen) == 0 {
		return nil
	}
	pipe := s.rdb.Pipeline()
	smCmds := make([]*redis.StringSliceCmd, 0, len(seen))
	order := make([]AccountSessionTarget, 0, len(seen))
	for _, t := range seen {
		order = append(order, t)
		smCmds = append(smCmds, pipe.SMembers(ctx, loginAccountTokensKey(t.AccountType, t.AccountID)))
	}
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return err
	}
	delPipe := s.rdb.Pipeline()
	for i, t := range order {
		tokens, _ := smCmds[i].Result()
		for _, tok := range tokens {
			delPipe.Del(ctx, loginTokenKey(tok))
			delPipe.SRem(ctx, loginTokensKey(), tok)
		}
		delPipe.Del(ctx, loginAccountTokensKey(t.AccountType, t.AccountID))
	}
	_, err := delPipe.Exec(ctx)
	return err
}

// ListTokens 列出全局在线 token（可能含过期残留）。
func (s *SessionStore) ListTokens(ctx context.Context) ([]string, error) {
	return s.rdb.SMembers(ctx, loginTokensKey()).Result()
}

// CountTokens 全局在线 token 数（SCARD login:tokens）。
func (s *SessionStore) CountTokens(ctx context.Context) (int64, error) {
	return s.rdb.SCard(ctx, loginTokensKey()).Result()
}

// ListTokensForAccount 列出某类型账号下的全部 token。
func (s *SessionStore) ListTokensForAccount(ctx context.Context, accountType AccountType, accountID string) ([]string, error) {
	return s.rdb.SMembers(ctx, loginAccountTokensKey(accountType, accountID)).Result()
}

// ListSessionsByTokens 批量读取会话（MGET），并清理全局集中已过期的 token。
func (s *SessionStore) ListSessionsByTokens(ctx context.Context, tokens []string) ([]*SessionPayload, error) {
	uniq := make([]string, 0, len(tokens))
	seen := map[string]struct{}{}
	for _, t := range tokens {
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		uniq = append(uniq, t)
	}
	if len(uniq) == 0 {
		return nil, nil
	}
	keys := make([]string, len(uniq))
	for i, t := range uniq {
		keys[i] = loginTokenKey(t)
	}
	rawValues, err := s.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	out := make([]*SessionPayload, 0, len(uniq))
	var stale []string
	for i, raw := range rawValues {
		if raw == nil {
			stale = append(stale, uniq[i])
			continue
		}
		var text string
		switch v := raw.(type) {
		case string:
			text = v
		case []byte:
			text = string(v)
		default:
			stale = append(stale, uniq[i])
			continue
		}
		var p SessionPayload
		if err := json.Unmarshal([]byte(text), &p); err != nil {
			stale = append(stale, uniq[i])
			continue
		}
		out = append(out, &p)
	}
	if len(stale) > 0 {
		args := make([]any, len(stale))
		for i, t := range stale {
			args[i] = t
		}
		_ = s.rdb.SRem(ctx, loginTokensKey(), args...).Err()
	}
	return out, nil
}

// SessionCookiePath 按账号类型返回隔离的 Cookie Path。
func SessionCookiePath(accountType AccountType) string {
	return fmt.Sprintf("/api/v1/%s", accountType.URLSegment())
}
