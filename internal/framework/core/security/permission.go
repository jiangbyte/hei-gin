// internal/framework/core/security/permission.go 权限注册表。
//
// Author: Charlie

package security

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

// PermissionRegistry 启动时收集路由权限键并同步到 Redis。
//
// Author: Charlie
type PermissionRegistry struct {
	mu       sync.Mutex
	items    map[string]string // 权限键 -> 展示名
	rdb      *redis.Client
	redisKey string
}

// NewPermissionRegistry 创建权限注册表并绑定 Redis。
func NewPermissionRegistry(rdb *redis.Client) *PermissionRegistry {
	return &PermissionRegistry{
		items:    make(map[string]string),
		rdb:      rdb,
		redisKey: "hei:permission:registry",
	}
}

// Register 登记权限键（空键或通配符 *:*:* 忽略）。
func (r *PermissionRegistry) Register(key, name string) {
	if key == "" || key == "*:*:*" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if name == "" {
		name = key
	}
	if _, ok := r.items[key]; !ok {
		r.items[key] = name
	}
}

// All 返回已登记权限键到展示名的快照。
func (r *PermissionRegistry) All() map[string]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]string, len(r.items))
	for k, v := range r.items {
		out[k] = v
	}
	return out
}

// Sync 将注册表 JSON 写入 Redis。
func (r *PermissionRegistry) Sync(ctx context.Context) error {
	all := r.All()
	raw, err := json.Marshal(all)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, r.redisKey, raw, 0).Err()
}

// EnsureRegistered 校验权限键已在注册表中登记（对齐 hei-boot PermissionRegistryServiceImpl.ensureRegistered）。
func (r *PermissionRegistry) EnsureRegistered(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[key]; ok {
		return nil
	}
	return fmt.Errorf("permission key not registered: %s", key)
}

// HasPermission 判断权限键集合是否满足 need（空 need 或含 *:*:* 视为通过）。
func HasPermission(keys []string, need string) bool {
	if need == "" {
		return true
	}
	for _, k := range keys {
		if k == "*:*:*" || k == need {
			return true
		}
	}
	return false
}
