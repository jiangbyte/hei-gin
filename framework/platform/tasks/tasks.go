// Package tasks 在 worker 进程内运行模块周期任务调度。
package tasks

import (
	"context"
	"sync"
	"time"

	"hei-gin/framework/core/logger"
	"hei-gin/framework/platform/module"

	"go.uber.org/zap"
)

// Manager 在进程内运行模块定时任务（worker）。
//
// Author: Charlie
type Manager struct {
	schedules []module.Schedule
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// NewManager 从模块注册表收集全部 Schedule。
func NewManager(regs *module.Registry) *Manager {
	var ss []module.Schedule
	for _, m := range regs.Modules {
		ss = append(ss, m.Schedules...)
	}
	return &Manager{schedules: ss}
}

// Start 启动全部周期任务 goroutine。
func (m *Manager) Start(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	m.cancel = cancel
	for _, s := range m.schedules {
		s := s
		interval := parseEvery(s.Interval)
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			t := time.NewTicker(interval)
			defer t.Stop()
			logger.L.Info("schedule started", zap.String("name", s.Name), zap.Duration("interval", interval))
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					if err := s.Run(ctx); err != nil {
						logger.L.Warn("schedule failed", zap.String("name", s.Name), zap.Error(err))
					}
				}
			}
		}()
	}
}

// Stop 取消调度并等待全部任务退出。
func (m *Manager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
	m.wg.Wait()
}

func parseEvery(s string) time.Duration {
	// 支持 "@every 1m" / "@every 30s" / 纯 duration 字符串
	if len(s) > 7 && s[:7] == "@every " {
		d, err := time.ParseDuration(s[7:])
		if err == nil {
			return d
		}
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return time.Minute
	}
	return d
}
