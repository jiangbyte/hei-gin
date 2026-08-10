package weakpassword

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 弱密码业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造弱密码服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.weak_password 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.weak_password",
		Models: []any{&WeakPassword{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建弱密码。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := WeakPassword{ID: idgen.Next(), Password: req.Password}
	return s.repo.Create(ctx, &row)
}

// Update 更新弱密码。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	return s.repo.UpdatePassword(ctx, req.ID, req.Password)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*WeakPassword, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []WeakPassword, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// List 列表。
func (s *Service) List(ctx context.Context, q ListParam) ([]WeakPassword, error) {
	return s.repo.List(ctx, q)
}
