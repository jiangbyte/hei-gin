package group

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 用户组服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造用户组服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 iam.group 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.group",
		Models: []any{&Group{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建用户组。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Group{
		ID: idgen.Next(), Name: req.Name, OwnerDeptID: req.OwnerDeptID,
		Description: req.Description, Status: orStatus(req.Status), Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新用户组。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"name": req.Name, "owner_dept_id": req.OwnerDeptID,
		"description": req.Description, "status": orStatus(req.Status),
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 用户组详情。
func (s *Service) Detail(ctx context.Context, id string) (*Group, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Group, total int64, current, size int, err error) {
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
