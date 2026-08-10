package role

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 角色服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造角色服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 iam.role 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.role",
		Models: []any{&Role{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建角色。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Role{
		ID: idgen.Next(), Code: req.Code, Name: req.Name,
		Category: orDef(req.Category, "SYS"), ScopeType: orDef(req.ScopeType, "PLATFORM"),
		OwnerDeptID: req.OwnerDeptID, Sort: orSort(req.Sort), Status: orStatus(req.Status),
		Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新角色。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"code": req.Code, "name": req.Name, "category": orDef(req.Category, "SYS"),
		"scope_type": orDef(req.ScopeType, "PLATFORM"), "owner_dept_id": req.OwnerDeptID,
		"sort": orSort(req.Sort), "status": orStatus(req.Status), "description": req.Description,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 角色详情。
func (s *Service) Detail(ctx context.Context, id string) (*Role, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Role, total int64, current, size int, err error) {
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

func orDef(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func orSort(n int) int {
	if n == 0 {
		return 99
	}
	return n
}
