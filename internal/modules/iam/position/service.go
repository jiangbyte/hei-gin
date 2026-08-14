// internal/modules/iam/position/service.go 业务服务。
//
// Author: Charlie

package position

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service èŒä½æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ èŒä½æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º iam.position æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.position",
		Models: []any{&Position{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create åˆ›å»ºèŒä½ã€‚
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Position{
		ID: idgen.Next(), Name: req.Name, Category: req.Category, OwnerDeptID: req.OwnerDeptID,
		Sort: orSort(req.Sort), IsVirtual: req.IsVirtual, Status: orStatus(req.Status),
		Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update æ›´æ–°èŒä½ã€‚
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"name": req.Name, "category": req.Category, "owner_dept_id": req.OwnerDeptID,
		"sort": orSort(req.Sort), "is_virtual": req.IsVirtual, "status": orStatus(req.Status),
		"description": req.Description,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail èŒä½è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Position, error) {
	return s.repo.GetByID(ctx, id)
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Position, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p)
	return rows, total, current, size, err
}

func orStatus(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}

func orSort(n int) int {
	if n == 0 {
		return 99
	}
	return n
}
