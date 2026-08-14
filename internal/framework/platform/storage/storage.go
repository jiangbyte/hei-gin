// Package storage æŠ½è±¡å¯¹è±¡å­˜å‚¨ï¼ˆlocal é»˜è®¤ï¼›S3 å…¼å®¹å¯åŽç»­æŽ¥å…¥ï¼‰ã€‚
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

	"hei-gin/internal/framework/core/config"
)

// Provider æŠ½è±¡å¯¹è±¡å­˜å‚¨ï¼ˆlocal / S3 å…¼å®¹ï¼‰ã€‚
//
// Author: Charlie
type Provider interface {
	Put(ctx context.Context, objectName string, r io.Reader, size int64, contentType string) (url string, err error)
	Get(ctx context.Context, objectName string) (io.ReadCloser, error)
	Delete(ctx context.Context, objectName string) error
	PublicURL(objectName string) string
}

// Manager æŒæœ‰å¯çƒ­åˆ‡æ¢çš„å­˜å‚¨ Providerã€‚
//
// Author: Charlie
type Manager struct {
	mu       sync.RWMutex
	provider Provider
	cfg      config.StorageConfig
}

// NewManager æŒ‰é…ç½®åˆ›å»ºå­˜å‚¨ç®¡ç†å™¨ã€‚
func NewManager(cfg config.StorageConfig) (*Manager, error) {
	p, err := newProvider(cfg)
	if err != nil {
		return nil, err
	}
	return &Manager{provider: p, cfg: cfg}, nil
}

// Provider è¿”å›žå½“å‰ Providerã€‚
func (m *Manager) Provider() Provider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.provider
}

// Reconfigure ç”¨æ–°é…ç½®é‡å»ºå¹¶æ›¿æ¢ Providerã€‚
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

// Local åŸºäºŽæœ¬åœ°ç›®å½•çš„å¯¹è±¡å­˜å‚¨å®žçŽ°ã€‚
//
// Author: Charlie
type Local struct {
	root       string
	publicPath string
	baseURL    string
}

// NewLocal åˆ›å»ºæœ¬åœ°å­˜å‚¨å¹¶ç¡®ä¿æ ¹ç›®å½•å­˜åœ¨ã€‚
func NewLocal(root, publicPath, baseURL string) *Local {
	_ = os.MkdirAll(root, 0o755)
	return &Local{root: root, publicPath: publicPath, baseURL: strings.TrimRight(baseURL, "/")}
}

// Put å†™å…¥å¯¹è±¡å¹¶è¿”å›žå…¬å¼€ URLã€‚
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

// Get æ‰“å¼€å¯¹è±¡åªè¯»æµã€‚
func (l *Local) Get(_ context.Context, objectName string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.root, filepath.FromSlash(objectName)))
}

// Delete åˆ é™¤æœ¬åœ°å¯¹è±¡æ–‡ä»¶ã€‚
func (l *Local) Delete(_ context.Context, objectName string) error {
	return os.Remove(filepath.Join(l.root, filepath.FromSlash(objectName)))
}

// PublicURL æ‹¼å‡ºå¯¹è±¡å…¬å¼€è®¿é—®è·¯å¾„ï¼ˆå¯å¸¦ baseURLï¼‰ã€‚
func (l *Local) PublicURL(objectName string) string {
	p := strings.TrimRight(l.publicPath, "/") + "/" + strings.TrimLeft(objectName, "/")
	if l.baseURL != "" {
		return l.baseURL + p
	}
	return p
}

// ObjectKey ç”¨å‰ç¼€ä¸Žæ–‡ä»¶åæ‹¼å¯¹è±¡é”®ã€‚
func ObjectKey(prefix, name string) string {
	return fmt.Sprintf("%s/%s", strings.Trim(prefix, "/"), name)
}
