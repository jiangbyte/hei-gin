// Package storage 抽象对象存储（local 默认；S3 兼容可后续接入）。
//
// Author: Charlie
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/platform/runtimecfg"
)

// Provider 抽象对象存储（local / S3 兼容）。
//
// Author: Charlie
type Provider interface {
	Put(ctx context.Context, objectName string, r io.Reader, size int64, contentType string) (url string, err error)
	Get(ctx context.Context, objectName string) (io.ReadCloser, error)
	Delete(ctx context.Context, objectName string) error
	PublicURL(objectName string) string
}

// Manager 持有可热切换的存储 Provider。
//
// Author: Charlie
type Manager struct {
	mu       sync.RWMutex
	provider Provider
	cfg      config.StorageConfig
	runtime  *runtimecfg.Settings
}

// Presigner 可选接口：支持预签名 URL 的 Provider（本地返回公开 URL）。
type Presigner interface {
	PresignedURL(ctx context.Context, objectName string, expire time.Duration) (string, error)
}

// SetRuntime 注入运行时配置读取器（用于 STORAGE_PRESIGN_EXPIRE_SECONDS 等）。
func (m *Manager) SetRuntime(s *runtimecfg.Settings) {
	m.mu.Lock()
	m.runtime = s
	m.mu.Unlock()
}

// PresignExpireSeconds 预签名有效期（秒）：运行时 STORAGE_PRESIGN_EXPIRE_SECONDS 优先，回退 yaml。
func (m *Manager) PresignExpireSeconds(ctx context.Context) int {
	m.mu.RLock()
	rt := m.runtime
	def := m.cfg.PresignExpireSeconds
	m.mu.RUnlock()
	if rt != nil {
		if v := rt.GetInt(ctx, "STORAGE_PRESIGN_EXPIRE_SECONDS", def); v > 0 {
			return v
		}
	}
	if def <= 0 {
		return 3600
	}
	return def
}

// NewManager 按配置创建存储管理器。
func NewManager(cfg config.StorageConfig) (*Manager, error) {
	p, err := newProvider(cfg)
	if err != nil {
		return nil, err
	}
	return &Manager{provider: p, cfg: cfg}, nil
}

// Provider 返回当前 Provider。
func (m *Manager) Provider() Provider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.provider
}

// Reconfigure 用新配置重建并替换 Provider。
func (m *Manager) Reconfigure(cfg config.StorageConfig) error {
	p, err := newProvider(cfg)
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.provider = p
	m.cfg = cfg
	m.mu.Unlock()
	return nil
}

func newProvider(cfg config.StorageConfig) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "", "local":
		return NewLocal(cfg.LocalRoot, cfg.PublicPath, cfg.BaseURL), nil
	case "s3", "minio", "oss":
		return NewS3(cfg)
	default:
		return nil, fmt.Errorf("storage: unsupported provider %q (use local|s3|minio|oss)", cfg.Provider)
	}
}

// Local 基于本地目录的对象存储实现。
//
// Author: Charlie
type Local struct {
	root       string
	publicPath string
	baseURL    string
}

// NewLocal 创建本地存储并确保根目录存在。
func NewLocal(root, publicPath, baseURL string) *Local {
	_ = os.MkdirAll(root, 0o755)
	return &Local{root: root, publicPath: publicPath, baseURL: strings.TrimRight(baseURL, "/")}
}

// Put 写入对象并返回公开 URL。
func (l *Local) Put(_ context.Context, objectName string, r io.Reader, _ int64, _ string) (string, error) {
	path := filepath.Join(l.root, filepath.FromSlash(objectName))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	return l.PublicURL(objectName), nil
}

// Get 打开对象只读流。
func (l *Local) Get(_ context.Context, objectName string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.root, filepath.FromSlash(objectName)))
}

// Delete 删除本地对象文件。
func (l *Local) Delete(_ context.Context, objectName string) error {
	return os.Remove(filepath.Join(l.root, filepath.FromSlash(objectName)))
}

// PublicURL 拼出对象公开访问路径（可带 baseURL）。
func (l *Local) PublicURL(objectName string) string {
	p := strings.TrimRight(l.publicPath, "/") + "/" + strings.TrimLeft(objectName, "/")
	if l.baseURL != "" {
		return l.baseURL + p
	}
	return p
}

// PresignedURL 本地存储直接返回公开 URL（预签名无意义）。
func (l *Local) PresignedURL(_ context.Context, objectName string, _ time.Duration) (string, error) {
	return l.PublicURL(objectName), nil
}

// ObjectKey 用前缀与文件名拼对象键。
func ObjectKey(prefix, name string) string {
	return fmt.Sprintf("%s/%s", strings.Trim(prefix, "/"), name)
}
