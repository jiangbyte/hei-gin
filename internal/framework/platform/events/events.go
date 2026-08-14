// Package events 提供进程内生命周期 / 领域事件总线。
//
// Author: Charlie
package events

import (
	"context"
	"sync"
)

const (
	OnDBReady           = "on_db_ready"
	OnConfigLoaded      = "on_config_loaded"
	OnStorageConfigured = "on_storage_configured"
	OnPermissionsSynced = "on_permissions_synced"
	OnAuditEvent        = "on_audit_event"
)

// Handler 为事件处理函数。
//
// Author: Charlie
type Handler func(ctx context.Context, payload any) error

// Bus 为进程内生命周期副作用发布订阅。
//
// Author: Charlie
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]Handler
}

// NewBus 创建空事件总线。
func NewBus() *Bus {
	return &Bus{subs: make(map[string][]Handler)}
}

// Subscribe 订阅指定事件。
func (b *Bus) Subscribe(event string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[event] = append(b.subs[event], h)
}

// Emit 同步依次调用订阅者，任一失败即返回。
func (b *Bus) Emit(ctx context.Context, event string, payload any) error {
	b.mu.RLock()
	hs := append([]Handler{}, b.subs[event]...)
	b.mu.RUnlock()
	for _, h := range hs {
		if err := h(ctx, payload); err != nil {
			return err
		}
	}
	return nil
}
