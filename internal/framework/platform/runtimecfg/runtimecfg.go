// Package runtimecfg 运行时配置读取：sys_config 表为权威，缺省回退默认值。
//
// 对齐 hei-boot RuntimeSettings / ConfigApi：管理端配置页（hei-admin sys/config）保存的
// 每个键（AUTH_LOGIN_* / AUTH_PASSWORD_* / MAIL_* / SMS_* / PUSH_* / STORAGE_* ...）
// 都应由业务侧在运行期消费，而非仅落库。
//
// Author: Charlie
package runtimecfg

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/crypto"
)

// Settings 基于 sys_config 的运行时配置读取器（敏感值自动解密，对齐 hei-boot getValue）。
//
// Author: Charlie
type Settings struct {
	db    *gorm.DB
	codec *crypto.Codec

	mu    sync.RWMutex
	cache map[string]string // 进程级缓存；config 变更后 Invalidate
}

// New 构造读取器。
func New(db *gorm.DB) *Settings {
	return &Settings{db: db, cache: map[string]string{}}
}

// WithCodec 注入配置加解密器（敏感键运行时读取需解密）。
func (s *Settings) WithCodec(codec *crypto.Codec) *Settings {
	s.codec = codec
	return s
}

// Invalidate 清空缓存（sys_config 变更后调用，对齐 hei-fastapi config_reader.reload）。
func (s *Settings) Invalidate() {
	if s == nil {
		return
	}
	s.mu.Lock()
	s.cache = map[string]string{}
	s.mu.Unlock()
}

// GetString 读取字符串配置；空值/缺失返回 def；密文自动解密。
func (s *Settings) GetString(ctx context.Context, key, def string) string {
	if s == nil || s.db == nil || key == "" {
		return def
	}
	s.mu.RLock()
	if v, ok := s.cache[key]; ok {
		s.mu.RUnlock()
		if v == "" {
			return def
		}
		return v
	}
	s.mu.RUnlock()

	var raw string
	if err := s.db.WithContext(ctx).Table("sys_config").Select("config_value").
		Where("config_key = ?", key).Limit(1).Scan(&raw).Error; err != nil || raw == "" {
		s.storeCache(key, "")
		return def
	}
	v := raw
	if s.codec != nil && crypto.LooksEncrypted(raw) {
		if plain, err := s.codec.Decrypt(raw); err == nil && plain != "" {
			v = plain
		}
	}
	s.storeCache(key, v)
	return v
}

func (s *Settings) storeCache(key, v string) {
	s.mu.Lock()
	if s.cache == nil {
		s.cache = map[string]string{}
	}
	s.cache[key] = v
	s.mu.Unlock()
}

// GetBool 读取布尔配置（TRUE/1/yes/on 视为真）。
func (s *Settings) GetBool(ctx context.Context, key string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(s.GetString(ctx, key, "")))
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

// GetInt 读取整数配置；非法值回退默认。
func (s *Settings) GetInt(ctx context.Context, key string, def int) int {
	v := strings.TrimSpace(s.GetString(ctx, key, ""))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
