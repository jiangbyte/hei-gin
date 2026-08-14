package dashboard

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service ä»ªè¡¨ç›˜æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct {
	repo *Repo
}

// NewService æž„é€ æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º dashboard æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "dashboard",
		Order:  50,
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Overview æ¦‚è§ˆç»Ÿè®¡ã€‚
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
