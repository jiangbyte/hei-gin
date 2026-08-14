// internal/modules/sys/weak_password/service.go 业务服务。
//
// Author: Charlie

package weakpassword

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service å¼±å¯†ç ä¸šåŠ¡æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ å¼±å¯†ç æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º sys.weak_password æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.weak_password",
		Models: []any{&WeakPassword{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create åˆ›å»ºå¼±å¯†ç ã€‚
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := WeakPassword{ID: idgen.Next(), Password: req.Password}
	return s.repo.Create(ctx, &row)
}

// Update æ›´æ–°å¼±å¯†ç ã€‚
func (s *Service) Update(ctx context.Context, req EditParam) error {
	return s.repo.UpdatePassword(ctx, req.ID, req.Password)
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*WeakPassword, error) {
	return s.repo.GetByID(ctx, id)
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, q PageParam) (rows []WeakPassword, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// List åˆ—è¡¨ã€‚
func (s *Service) List(ctx context.Context, q ListParam) ([]WeakPassword, error) {
	return s.repo.List(ctx, q)
}
