// internal/modules/dashboard/service.go 业务服务。
//
// Author: Charlie

package dashboard

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service 仪表盘服务。
//
// Author: Charlie
type Service struct {
	repo *Repo
}

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 dashboard 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "dashboard",
		Order:  50,
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Overview 概览统计。
func (s *Service) Overview(ctx context.Context) OverviewResult {
	return OverviewResult{
		AccountTotal:    s.repo.Count(ctx, "sys_account"),
		NoticeTotal:     s.repo.Count(ctx, "msg_notice"),
		FeedbackTotal:   s.repo.Count(ctx, "msg_feedback"),
		FeedbackPending: s.repo.CountWhere(ctx, "msg_feedback", "status = ?", "PENDING"),
		FileTotal:       s.repo.Count(ctx, "sys_file"),
		RoleTotal:       s.repo.Count(ctx, "sys_role"),
		DeptTotal:       s.repo.Count(ctx, "sys_dept"),
	}
}
