// Package xxljob 在 API 进程内嵌入 XXL-JOB 执行器。
//
// Author: Charlie
package xxljob

import (
	"context"
	"fmt"
	"net/http"
	"time"

	xxl "github.com/xxl-job/xxl-job-executor-go"
	"go.uber.org/zap"

	"hei-gin/framework/core/config"
	"hei-gin/framework/core/logger"
	"hei-gin/framework/platform/module"
)

// Manager 管理 XXL-JOB 执行器生命周期。
//
// Author: Charlie
type Manager struct {
	exec   xxl.Executor
	server *http.Server
	cfg    config.XxlJobConfig
}

// NewManager 从配置与模块 Jobs 构建执行器（不启动）。
func NewManager(cfg config.XxlJobConfig, regs *module.Registry) *Manager {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(cfg.Admin.Addresses),
		xxl.AccessToken(cfg.AccessToken),
		xxl.ExecutorPort(fmt.Sprintf("%d", cfg.Executor.Port)),
		xxl.RegistryKey(cfg.Executor.AppName),
		xxl.SetLogger(&zapLogger{}),
	)
	exec.Init()
	m := &Manager{exec: exec, cfg: cfg}
	if regs != nil {
		for _, mod := range regs.Modules {
			for _, job := range mod.Jobs {
				job := job
				exec.RegTask(job.Name, func(ctx context.Context, param *xxl.RunReq) string {
					p := ""
					if param != nil {
						p = param.ExecutorParams
					}
					if err := job.Run(ctx, p); err != nil {
						logger.L.Warn("xxl-job 失败", zap.String("handler", job.Name), zap.Error(err))
						return err.Error()
					}
					return "ok"
				})
			}
		}
	}
	return m
}

// Start 启动执行器 HTTP 端口并开始向 Admin 注册（不占用主进程信号）。
func (m *Manager) Start() error {
	if m == nil || !m.cfg.Enabled {
		return nil
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/run", m.exec.RunTask)
	mux.HandleFunc("/kill", m.exec.KillTask)
	mux.HandleFunc("/log", m.exec.TaskLog)
	mux.HandleFunc("/beat", m.exec.Beat)
	mux.HandleFunc("/idleBeat", m.exec.IdleBeat)
	m.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", m.cfg.Executor.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		logger.L.Info("xxl-job 执行器监听",
			zap.String("appname", m.cfg.Executor.AppName),
			zap.Int("port", m.cfg.Executor.Port),
			zap.String("admin", m.cfg.Admin.Addresses),
		)
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Error("xxl-job 执行器退出", zap.Error(err))
		}
	}()
	return nil
}

// Stop 注销并关闭执行器端口。
func (m *Manager) Stop(ctx context.Context) error {
	if m == nil {
		return nil
	}
	if m.exec != nil {
		m.exec.Stop()
	}
	if m.server != nil {
		return m.server.Shutdown(ctx)
	}
	return nil
}

type zapLogger struct{}

func (z *zapLogger) Info(format string, a ...interface{}) {
	logger.L.Info(fmt.Sprintf(format, a...))
}

func (z *zapLogger) Error(format string, a ...interface{}) {
	logger.L.Error(fmt.Sprintf(format, a...))
}
