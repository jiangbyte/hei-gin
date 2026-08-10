package position

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 职位服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造职位服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 iam.position 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.position",
		Models: []any{&Position{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建职位。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Position{
		ID: idgen.Next(), Name: req.Name, Category: req.Category, OwnerDeptID: req.OwnerDeptID,
		Sort: orSort(req.Sort), IsVirtual: req.IsVirtual, Status: orStatus(req.Status),
		Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新职位。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"name": req.Name, "category": req.Category, "owner_dept_id": req.OwnerDeptID,
		"sort": orSort(req.Sort), "is_virtual": req.IsVirtual, "status": orStatus(req.Status),
		"description": req.Description,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 职位详情。
func (s *Service) Detail(ctx context.Context, id string) (*Position, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
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
