package client

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 客户端资源服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造客户端资源服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 iam.client 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.client",
		Models: []any{&ClientModule{}, &ClientResource{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// CreateModule 创建客户端模块。
func (s *Service) CreateModule(ctx context.Context, req ModuleAddParam) error {
	at := req.AccountType
	if at == "" {
		at = string(security.AccountAdmin)
	}
	row := ClientModule{
		ID: idgen.Next(), Name: req.Name, Code: req.Code, AccountType: at, Icon: req.Icon, Color: req.Color,
		Sort: sortOr(req.Sort), Status: statusOr(req.Status), Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.CreateModule(ctx, &row)
}

// UpdateModule 更新客户端模块。
func (s *Service) UpdateModule(ctx context.Context, req ModuleEditParam) error {
	at := req.AccountType
	if at == "" {
		at = string(security.AccountAdmin)
	}
	updates := map[string]any{
		"name": req.Name, "code": req.Code, "account_type": at, "icon": req.Icon, "color": req.Color,
		"sort": sortOr(req.Sort), "status": statusOr(req.Status), "description": req.Description,
	}
	return s.repo.UpdateModule(ctx, req.ID, updates)
}

// DeleteModules 批量删除客户端模块。
func (s *Service) DeleteModules(ctx context.Context, ids []string) error {
	return s.repo.DeleteModules(ctx, ids)
}

// ModuleDetail 客户端模块详情。
func (s *Service) ModuleDetail(ctx context.Context, id string) (*ClientModule, error) {
	return s.repo.GetModuleByID(ctx, id)
}

// ModulePage 客户端模块分页。
func (s *Service) ModulePage(ctx context.Context, p ModulePageParam) (rows []ClientModule, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageModules(ctx, p)
	return rows, total, current, size, err
}

// ModuleSelector 客户端模块选择器。
func (s *Service) ModuleSelector(ctx context.Context, accountType string) ([]ModuleOption, error) {
	rows, err := s.repo.ListEnabledModules(ctx, accountType)
	if err != nil {
		return nil, err
	}
	out := make([]ModuleOption, 0, len(rows))
	for _, r := range rows {
		out = append(out, ModuleOption{ID: r.ID, Name: r.Name, Code: r.Code})
	}
	return out, nil
}

// CreateResource 创建客户端资源。
func (s *Service) CreateResource(ctx context.Context, req ResourceAddParam) error {
	vis := true
	if req.IsVisible != nil {
		vis = *req.IsVisible
	}
	row := ClientResource{
		ID: idgen.Next(), ParentID: req.ParentID, Code: req.Code, Name: req.Name, ResourceType: req.ResourceType,
		ModuleID: req.ModuleID, Path: req.Path, Component: req.Component, Redirect: req.Redirect, Icon: req.Icon,
		Color: req.Color, Href: req.Href, Sort: sortOr(req.Sort), IsVisible: vis, IsCache: req.IsCache, IsAffix: req.IsAffix,
		Status: statusOr(req.Status), Description: req.Description, Layout: req.Layout, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.CreateResource(ctx, &row)
}

// UpdateResource 更新客户端资源。
func (s *Service) UpdateResource(ctx context.Context, req ResourceEditParam) error {
	vis := true
	if req.IsVisible != nil {
		vis = *req.IsVisible
	}
	updates := map[string]any{
		"parent_id": req.ParentID, "code": req.Code, "name": req.Name, "resource_type": req.ResourceType,
		"module_id": req.ModuleID, "path": req.Path, "component": req.Component, "redirect": req.Redirect,
		"icon": req.Icon, "color": req.Color, "href": req.Href, "sort": sortOr(req.Sort),
		"is_visible": vis, "is_cache": req.IsCache, "is_affix": req.IsAffix, "status": statusOr(req.Status),
		"description": req.Description, "layout": req.Layout,
	}
	return s.repo.UpdateResource(ctx, req.ID, updates)
}

// DeleteResources 批量删除客户端资源。
func (s *Service) DeleteResources(ctx context.Context, ids []string) error {
	return s.repo.DeleteResources(ctx, ids)
}

// ResourceDetail 客户端资源详情。
func (s *Service) ResourceDetail(ctx context.Context, id string) (*ClientResource, error) {
	return s.repo.GetResourceByID(ctx, id)
}

// ResourcePage 客户端资源分页。
func (s *Service) ResourcePage(ctx context.Context, p ResourcePageParam) (rows []ClientResource, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageResources(ctx, p)
	return rows, total, current, size, err
}

// ResourceTree 客户端资源树。
func (s *Service) ResourceTree(ctx context.Context, moduleID string) ([]TreeNode, error) {
	rows, err := s.repo.ListResources(ctx, moduleID)
	if err != nil {
		return nil, err
	}
	return buildTree(rows, nil), nil
}

func buildTree(rows []ClientResource, parent *string) []TreeNode {
	out := make([]TreeNode, 0)
	for _, r := range rows {
		same := (r.ParentID == nil && parent == nil) || (r.ParentID != nil && parent != nil && *r.ParentID == *parent)
		if same {
			out = append(out, TreeNode{ClientResource: r, Children: buildTree(rows, &r.ID)})
		}
	}
	return out
}

func sortOr(n int) int {
	if n == 0 {
		return 99
	}
	return n
}

func statusOr(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}
