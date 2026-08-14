// Package snailjob åœ¨ API è¿›ç¨‹å†…åµŒå…¥ SnailJob Go å®¢æˆ·ç«¯æ‰§è¡Œå™¨ã€‚
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

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/logger"
	"hei-gin/internal/framework/platform/module"
)

// Manager ç®¡ç† SnailJob å®¢æˆ·ç«¯ç”Ÿå‘½å‘¨æœŸã€‚
//
// Author: Charlie
type Manager struct {
	mgr *snailjobgo.SnailJobManager
	cfg config.SnailJobConfig
}

// NewManager ä»Žé…ç½®ä¸Žæ¨¡å— Jobs æž„å»ºå®¢æˆ·ç«¯ï¼ˆä¸å¯åŠ¨ï¼‰ã€‚
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

// Start åˆå§‹åŒ–å¹¶åœ¨åŽå°è¿è¡Œ gRPC æ‰§è¡Œå™¨ä¸Žå¿ƒè·³ã€‚
func (m *Manager) Start() error {
	if m == nil || !m.cfg.Enabled || m.mgr == nil {
		return nil
	}
	if err := m.mgr.Init(); err != nil {
		return fmt.Errorf("snailjob init: %w", err)
	}
	logger.L.Info("snailjob å®¢æˆ·ç«¯å¯åŠ¨",
		zap.String("group", m.cfg.GroupName),
		zap.String("namespace", m.cfg.Namespace),
		zap.String("server", m.cfg.ServerHost+":"+m.cfg.ServerPort),
		zap.String("host_port", m.cfg.HostPort),
	)
	go m.mgr.Run()
	return nil
}

// Stop å°è¯•åœæ­¢å®¢æˆ·ç«¯ï¼ˆä¸Šæ¸¸åº“æ— æ˜¾å¼ Shutdownï¼Œè¿›ç¨‹é€€å‡ºæ—¶éš gRPC ç»“æŸï¼‰ã€‚
func (m *Manager) Stop(_ context.Context) error {
	if m == nil || m.mgr == nil {
		return nil
	}
	logger.L.Info("snailjob å®¢æˆ·ç«¯åœæ­¢è¯·æ±‚å·²å‘å‡º")
	return nil
}

// adapterExecutor å°† module.Job.Run é€‚é…ä¸º SnailJob BaseJobExecutorã€‚
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
