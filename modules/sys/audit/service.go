package audit

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 审计日志业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造审计服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.audit 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.audit",
		Models: []any{&OperationLog{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Page 分页查询。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []OperationLog, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*OperationLog, error) {
	return s.repo.GetByID(ctx, id)
}
