// Package snailjob 在 API 进程内嵌入 SnailJob Go 客户端执行器。
//
// Author: Charlie
package snailjob

import (
	"context"
	"fmt"

	snailjobgo "github.com/open-snail/snail-job-go"
	"github.com/open-snail/snail-job-go/dto"
	"github.com/open-snail/snail-job-go/job"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"hei-gin/framework/core/config"
	"hei-gin/framework/core/logger"
	"hei-gin/framework/platform/module"
)

// Manager 管理 SnailJob 客户端生命周期。
//
// Author: Charlie
type Manager struct {
	mgr *snailjobgo.SnailJobManager
	cfg config.SnailJobConfig
}

// NewManager 从配置与模块 Jobs 构建客户端（不启动）。
func NewManager(cfg config.SnailJobConfig, regs *module.Registry) *Manager {
	m := &Manager{cfg: cfg}
	if !cfg.Enabled {
		return m
	}
	opts := &dto.Options{
		ServerHost:   cfg.ServerHost,
		ServerPort:   cfg.ServerPort,
		HostIP:       cfg.HostIP,
		HostPort:     cfg.HostPort,
		Namespace:    cfg.Namespace,
		GroupName:    cfg.GroupName,
		Token:        cfg.Token,
		Level:        logrus.InfoLevel,
		ReportCaller: false,
	}
	mgr := snailjobgo.NewSnailJobManager(opts)
	if regs != nil {
		for _, mod := range regs.Modules {
			for _, j := range mod.Jobs {
				j := j
				mgr.Register(j.Name, func() job.IJobExecutor {
					return &adapterExecutor{run: j.Run}
				})
			}
		}
	}
	m.mgr = mgr
	return m
}

// Start 初始化并在后台运行 gRPC 执行器与心跳。
func (m *Manager) Start() error {
	if m == nil || !m.cfg.Enabled || m.mgr == nil {
		return nil
	}
	if err := m.mgr.Init(); err != nil {
		return fmt.Errorf("snailjob init: %w", err)
	}
	logger.L.Info("snailjob 客户端启动",
		zap.String("group", m.cfg.GroupName),
		zap.String("namespace", m.cfg.Namespace),
		zap.String("server", m.cfg.ServerHost+":"+m.cfg.ServerPort),
		zap.String("host_port", m.cfg.HostPort),
	)
	go m.mgr.Run()
	return nil
}

// Stop 尝试停止客户端（上游库无显式 Shutdown，进程退出时随 gRPC 结束）。
func (m *Manager) Stop(_ context.Context) error {
	if m == nil || m.mgr == nil {
		return nil
	}
	logger.L.Info("snailjob 客户端停止请求已发出")
	return nil
}

// adapterExecutor 将 module.Job.Run 适配为 SnailJob BaseJobExecutor。
type adapterExecutor struct {
	job.BaseJobExecutor
	run func(ctx context.Context, param string) error
}

func (e *adapterExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	param := ""
	if jobArgs != nil {
		if p := jobArgs.GetJobParams(); p != nil {
			param = fmt.Sprint(p)
		}
	}
	ctx := e.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if err := e.run(ctx, param); err != nil {
		if e.LocalLogger != nil {
			e.LocalLogger.Errorf("job failed: %v", err)
		}
		return *dto.Failure().WithMessage(err.Error())
	}
	return *dto.Success()
}
