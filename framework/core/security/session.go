package security

import (
	"context"
	"strings"
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

// SessionStore 管理不透明 Redis 会话。
//
// Author: Charlie
type SessionStore struct {
	rdb       *redis.Client
	keyPrefix string
	indexPref string
}

// NewSessionStore 创建会话存储。
func NewSessionStore(rdb *redis.Client) *SessionStore {
	return &SessionStore{
		rdb:       rdb,
		keyPrefix: "hei:session:",
		indexPref: "hei:session:account:",
	}
}

// NewToken 生成随机不透明会话 token。
func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *SessionStore) key(token string) string { return s.keyPrefix + token }

// Save 写入会话并维护账号 token 索引。
func (s *SessionStore) Save(ctx context.Context, p *SessionPayload, ttl time.Duration) error {
	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}
	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, s.key(p.Token), raw, ttl)
	pipe.SAdd(ctx, s.indexPref+p.AccountID, p.Token)
	pipe.Expire(ctx, s.indexPref+p.AccountID, ttl)
	_, err = pipe.Exec(ctx)
	return err
}

// Get 按 token 读取会话，不存在时返回 nil。
func (s *SessionStore) Get(ctx context.Context, token string) (*SessionPayload, error) {
	raw, err := s.rdb.Get(ctx, s.key(token)).Bytes()
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

// Delete 删除单个会话并更新账号索引。
func (s *SessionStore) Delete(ctx context.Context, token string) error {
	p, err := s.Get(ctx, token)
	if err != nil {
		return err
	}
	pipe := s.rdb.Pipeline()
	pipe.Del(ctx, s.key(token))
	if p != nil {
		pipe.SRem(ctx, s.indexPref+p.AccountID, token)
	}
	_, err = pipe.Exec(ctx)
	return err
}

// DeleteAllForAccount 删除某账号下全部会话。
func (s *SessionStore) DeleteAllForAccount(ctx context.Context, accountID string) error {
	tokens, err := s.rdb.SMembers(ctx, s.indexPref+accountID).Result()
	if err != nil {
		return err
	}
	pipe := s.rdb.Pipeline()
	for _, t := range tokens {
		pipe.Del(ctx, s.key(t))
	}
	pipe.Del(ctx, s.indexPref+accountID)
	_, err = pipe.Exec(ctx)
	return err
}

// Touch 刷新最后活跃时间并续期。
func (s *SessionStore) Touch(ctx context.Context, p *SessionPayload, ttl time.Duration) error {
	p.LastActiveAt = time.Now().UTC()
	return s.Save(ctx, p, ttl)
}

// CountForAccount 返回账号当前会话数。
func (s *SessionStore) CountForAccount(ctx context.Context, accountID string) (int64, error) {
	return s.rdb.SCard(ctx, s.indexPref+accountID).Result()
}

// ListAccountIDs 列出全部有会话的账号 ID（按索引前缀 SCAN）。
func (s *SessionStore) ListAccountIDs(ctx context.Context) ([]string, error) {
	var out []string
	var cursor uint64
	for {
		keys, next, err := s.rdb.Scan(ctx, cursor, s.indexPref+"*", 200).Result()
		if err != nil {
			return nil, err
		}
		for _, k := range keys {
			out = append(out, strings.TrimPrefix(k, s.indexPref))
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return out, nil
}

// ListTokensForAccount 列出某账号全部 token。
func (s *SessionStore) ListTokensForAccount(ctx context.Context, accountID string) ([]string, error) {
	return s.rdb.SMembers(ctx, s.indexPref+accountID).Result()
}

// GetTokenSet 读取账号 token 索引集合（不存在返回空）。
func (s *SessionStore) GetTokenSet(ctx context.Context, accountID string) ([]string, error) {
	return s.ListTokensForAccount(ctx, accountID)
}

// SessionCookiePath 按账号类型返回隔离的 Cookie Path。
func SessionCookiePath(accountType AccountType) string {
	return fmt.Sprintf("/api/v1/%s", accountType.URLSegment())
}