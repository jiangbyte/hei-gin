package codegen

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 代码生成业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造代码生成服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.codegen 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.codegen",
		Models: []any{&Plan{}, &Field{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建方案。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	pk := req.MainPK
	if pk == "" {
		pk = "id"
	}
	sort := req.Sort
	if sort == 0 {
		sort = 99
	}
	row := Plan{
		ID: idgen.Next(), Name: req.Name, GenType: req.GenType, Author: req.Author, Description: req.Description,
		MainTable: req.MainTable, MainPK: pk, MainEntityName: req.MainEntityName, MainModulePath: req.MainModulePath,
		MainBusinessName: req.MainBusinessName, APIPrefix: req.APIPrefix, PermissionPrefix: req.PermissionPrefix,
		ResourceModuleID: req.ResourceModuleID, ParentResourceID: req.ParentResourceID, MenuName: req.MenuName,
		MenuPath: req.MenuPath, ComponentPath: req.ComponentPath, Icon: req.Icon, Sort: sort,
	}
	return s.repo.Create(ctx, &row)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Plan, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Plan, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}
