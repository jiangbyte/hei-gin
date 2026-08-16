// Package storage 抽象对象存储（仅 S3 兼容：minio/rustfs/oss/s3，对齐 hei-boot）。
//
// 运行时权威配置来自 sys_config（DEFAULT_FILE_ENGINE + STORAGE_{ENGINE}_*），
// 与 hei-boot StorageSettingsResolver / hei-fastapi config_reader 一致；不再使用 yaml 存储段。
//
// Author: Charlie
package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"hei-gin/internal/framework/platform/runtimecfg"
)

// Provider 抽象对象存储（S3 兼容）。
//
// Author: Charlie
type Provider interface {
	Put(ctx context.Context, objectName string, r io.Reader, size int64, contentType string) (url string, err error)
	Get(ctx context.Context, objectName string) (io.ReadCloser, error)
	Delete(ctx context.Context, objectName string) error
	// PublicURL 返回对象可访问 URL：公开桶走直链，私有桶走预签名。
	PublicURL(ctx context.Context, objectName string) string
}

// BucketHolder 可选接口：暴露对象存储桶名（供元数据落库）。
//
// Author: Charlie
type BucketHolder interface {
	BucketName() string
}

// Presigner 可选接口：支持预签名 URL 的 Provider。
type Presigner interface {
	PresignedURL(ctx context.Context, objectName string, expire time.Duration) (string, error)
}

// ResolvedConfig 从 sys_config 解析出的对象存储配置（对齐 hei-boot ResolvedStorageConfig）。
//
// Author: Charlie
type ResolvedConfig struct {
	Provider             string
	Bucket               string
	Endpoint             string
	AccessKey            string
	SecretKey            string
	Region               string
	UseSSL               bool
	BaseURL              string
	BucketPublic         bool
	PresignExpireSeconds int
}

// Manager 持有可热切换的存储 Provider（按 sys_config 懒加载，对齐 hei-boot StorageEngineFactory）。
//
// Author: Charlie
type Manager struct {
	mu      sync.RWMutex
	byName  map[string]Provider
	runtime *runtimecfg.Settings
	version int64
}

// NewManager 创建空管理器（引擎在首次使用时按 RuntimeSettings 构建）。
func NewManager() *Manager {
	return &Manager{byName: map[string]Provider{}}
}

// SetRuntime 注入运行时配置读取器。
func (m *Manager) SetRuntime(s *runtimecfg.Settings) {
	m.mu.Lock()
	m.runtime = s
	m.mu.Unlock()
	m.Refresh()
}

// Refresh 清空引擎缓存（sys_config 变更后调用，对齐 hei-boot / hei-fastapi clear_storage_cache）。
func (m *Manager) Refresh() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byName = map[string]Provider{}
	m.version++
}

// DefaultProviderName 缺省上传引擎名（对齐 hei-boot StorageSettingsResolverImpl）：
// 运行时 DEFAULT_FILE_ENGINE，缺省 rustfs。
func (m *Manager) DefaultProviderName(ctx context.Context) string {
	m.mu.RLock()
	rt := m.runtime
	m.mu.RUnlock()
	if rt != nil {
		if eng := strings.TrimSpace(rt.GetString(ctx, "DEFAULT_FILE_ENGINE", "")); eng != "" {
			if p := EngineToProvider(eng); p != "" {
				return p
			}
		}
	}
	return ProviderRustFS
}

// ResolveURL 对象引用 → 访问 URL（外部 http(s) 原样返回；历史 /api/v1/files/ 前缀剥离后按默认引擎拼 URL）。
func (m *Manager) ResolveURL(ctx context.Context, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if u, err := url.Parse(value); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		return value
	}
	key := StripToObjectKey(value)
	if key == "" {
		return ""
	}
	return m.ProviderByName(ctx, m.DefaultProviderName(ctx)).PublicURL(ctx, key)
}

// StripToObjectKey 把任意对象引用（纯 key / /api/v1/files/... / 完整 URL path）转成存储纯 key。
func StripToObjectKey(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if u, err := url.Parse(value); err == nil && u.Scheme != "" {
		value = u.Path
	}
	pathOnly := strings.ReplaceAll(value, "\\", "/")
	pathOnly = strings.TrimLeft(pathOnly, "/")
	for _, prefix := range []string{"api/v1/files/", "v1/files/"} {
		if strings.HasPrefix(pathOnly, prefix) {
			pathOnly = pathOnly[len(prefix):]
			break
		}
	}
	// path-style URL 常为 /{bucket}/{key}：若第二段起像业务 key（uploads/），去掉首段 bucket。
	if slash := strings.Index(pathOnly, "/"); slash > 0 {
		rest := pathOnly[slash+1:]
		if strings.HasPrefix(rest, "uploads/") {
			pathOnly = rest
		}
	}
	return strings.TrimLeft(pathOnly, "/")
}

// PresignExpireSeconds 预签名有效期（秒）：运行时 STORAGE_PRESIGN_EXPIRE_SECONDS，缺省 3600。
func (m *Manager) PresignExpireSeconds(ctx context.Context) int {
	m.mu.RLock()
	rt := m.runtime
	m.mu.RUnlock()
	def := 3600
	if rt != nil {
		if v := rt.GetInt(ctx, "STORAGE_PRESIGN_EXPIRE_SECONDS", def); v > 0 {
			return v
		}
	}
	return def
}

// Provider 返回默认引擎（DEFAULT_FILE_ENGINE）。
func (m *Manager) Provider(ctx context.Context) Provider {
	return m.ProviderByName(ctx, m.DefaultProviderName(ctx))
}

// ProviderByName 按提供商名解析引擎（对齐 hei-boot StorageEngineFactory + RuntimeSettings）。
// name 可为 provider（minio/rustfs/oss/s3）或引擎（MINIO/…）；未知回退默认引擎。
func (m *Manager) ProviderByName(ctx context.Context, name string) Provider {
	resolved := ResolveProvider(name)
	if resolved == "" {
		resolved = m.DefaultProviderName(ctx)
	}
	m.mu.RLock()
	if p, ok := m.byName[resolved]; ok {
		m.mu.RUnlock()
		return p
	}
	rt := m.runtime
	m.mu.RUnlock()

	p, err := m.buildProvider(ctx, resolved, rt)
	if err != nil {
		return &errProvider{err: err}
	}
	m.mu.Lock()
	if m.byName == nil {
		m.byName = map[string]Provider{}
	}
	m.byName[resolved] = p
	m.mu.Unlock()
	return p
}

func (m *Manager) buildProvider(ctx context.Context, name string, rt *runtimecfg.Settings) (Provider, error) {
	cfg, err := resolveConfig(ctx, name, rt, m.PresignExpireSeconds(ctx))
	if err != nil {
		return nil, err
	}
	return NewS3(cfg)
}

// resolveConfig 仅从 sys_config 组装配置（对齐 hei-boot StorageSettingsResolverImpl.buildForProvider）。
func resolveConfig(ctx context.Context, name string, rt *runtimecfg.Settings, presignExpire int) (ResolvedConfig, error) {
	prefix := ProviderConfigKeyPrefix(name)
	if prefix == "" {
		return ResolvedConfig{}, fmt.Errorf("storage: unsupported provider %q", name)
	}
	get := func(suffix, def string) string {
		if rt != nil {
			if v := rt.GetString(ctx, prefix+"_"+suffix, ""); v != "" {
				return v
			}
		}
		return def
	}
	getBool := func(suffix string, def bool) bool {
		if rt != nil {
			if v := rt.GetString(ctx, prefix+"_"+suffix, ""); v != "" {
				return strings.EqualFold(v, "true") || v == "1" || strings.EqualFold(v, "yes")
			}
		}
		return def
	}
	region := get("REGION", "us-east-1")
	if region == "" {
		region = "us-east-1"
	}
	cfg := ResolvedConfig{
		Provider:             name,
		Bucket:               get("BUCKET", ""),
		Endpoint:             get("ENDPOINT", ""),
		AccessKey:            get("ACCESS_KEY", ""),
		SecretKey:            get("SECRET_KEY", ""),
		Region:               region,
		UseSSL:               getBool("USE_SSL", false),
		BaseURL:              get("BASE_URL", ""),
		BucketPublic:         getBool("BUCKET_PUBLIC", false),
		PresignExpireSeconds: presignExpire,
	}
	if cfg.Bucket == "" {
		return ResolvedConfig{}, fmt.Errorf("storage: %s_BUCKET is required（请在系统配置中填写）", prefix)
	}
	return cfg, nil
}

// ObjectKey 用前缀与文件名拼对象键。
func ObjectKey(prefix, name string) string {
	return fmt.Sprintf("%s/%s", strings.Trim(prefix, "/"), name)
}

// errProvider 配置缺失时的占位引擎（所有操作返回明确错误）。
type errProvider struct{ err error }

func (p *errProvider) Put(context.Context, string, io.Reader, int64, string) (string, error) {
	return "", p.err
}
func (p *errProvider) Get(context.Context, string) (io.ReadCloser, error) { return nil, p.err }
func (p *errProvider) Delete(context.Context, string) error               { return p.err }
func (p *errProvider) PublicURL(context.Context, string) string           { return "" }
