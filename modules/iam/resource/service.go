package resource

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/iam/relation"
	"hei-gin/modules/shared"
)

// Service 资源服务（权限绑定经 relation 模块，权限注册表经 Perms）。
//
// Author: Charlie
type Service struct {
	repo  *Repo
	rel   *relation.Service
	perms *security.PermissionRegistry
}

// NewService 构造资源服务。
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: NewRepo(db),
		rel:  relation.NewService(db),
	}
}

// New 构建 iam.resource 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	s.perms = d.Perms
	return module.Module{
		Name:   "iam.resource",
		Models: []any{&Resource{}, &ResourceModule{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// CreateResource 创建资源。
func (s *Service) CreateResource(ctx context.Context, req ResourceAddParam) error {
	vis := true
	if req.IsVisible != nil {
		vis = *req.IsVisible
	}
	row := Resource{
		ID: idgen.Next(), ParentID: req.ParentID, Code: req.Code, Name: req.Name, ResourceType: req.ResourceType,
		ModuleID: req.ModuleID, Path: req.Path, Component: req.Component, Redirect: req.Redirect, Icon: req.Icon,
		Color: req.Color, Href: req.Href, Sort: sortOr(req.Sort), IsVisible: vis, IsCache: req.IsCache, IsAffix: req.IsAffix,
		Status: statusOr(req.Status), Description: req.Description, Layout: req.Layout, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.CreateResource(ctx, &row)
}

// UpdateResource 更新资源。
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

// DeleteResources 批量删除资源。
func (s *Service) DeleteResources(ctx context.Context, ids []string) error {
	return s.repo.DeleteResources(ctx, ids)
}

// ResourceDetail 资源详情。
func (s *Service) ResourceDetail(ctx context.Context, id string) (*Resource, error) {
	return s.repo.GetResourceByID(ctx, id)
}

// ResourcePage 资源分页。
func (s *Service) ResourcePage(ctx context.Context, p ResourcePageParam) (rows []Resource, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageResources(ctx, p)
	return rows, total, current, size, err
}

// CurrentAdmin 管理端当前资源。
func (s *Service) CurrentAdmin(ctx context.Context) ([]Resource, error) {
	return s.repo.ListResourcesByClient(ctx, string(security.AccountAdmin))
}

// CurrentPortal 门户当前资源。
func (s *Service) CurrentPortal(ctx context.Context) ([]Resource, error) {
	return s.repo.ListResourcesByClient(ctx, string(security.AccountPortal))
}

// ListGrantModules 资源授权模块选项（含模块下启用资源，空模块过滤）。
func (s *Service) ListGrantModules(ctx context.Context, accountType string) ([]GrantModule, error) {
	typ := accountType
	if typ == "" {
		typ = string(security.AccountAdmin)
	}
	modules, err := s.repo.ListEnabledModules(ctx, typ)
	if err != nil {
		return nil, err
	}
	resources, err := s.repo.ListGrantResources(ctx, typ)
	if err != nil {
		return nil, err
	}
	out := make([]GrantModule, 0, len(modules))
	for _, m := range modules {
		gm := GrantModule{ModuleID: m.ID, Name: m.Name, Resources: []Resource{}}
		for _, res := range resources {
			if res.ModuleID != nil && *res.ModuleID == m.ID {
				gm.Resources = append(gm.Resources, res)
			}
		}
		if len(gm.Resources) > 0 {
			out = append(out, gm)
		}
	}
	return out, nil
}

// ResourceTree 资源树。
func (s *Service) ResourceTree(ctx context.Context, moduleID string) ([]TreeNode, error) {
	rows, err := s.repo.ListResources(ctx, moduleID)
	if err != nil {
		return nil, err
	}
	return buildTree(rows, nil), nil
}

// CreateModule 创建资源模块。
func (s *Service) CreateModule(ctx context.Context, req ModuleAddParam) error {
	client := req.Client
	if client == "" {
		client = string(security.AccountAdmin)
	}
	row := ResourceModule{
		ID: idgen.Next(), Name: req.Name, Code: req.Code, Client: client, Icon: req.Icon, Color: req.Color,
		Sort: sortOr(req.Sort), Status: statusOr(req.Status), Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.CreateModule(ctx, &row)
}

// UpdateModule 更新资源模块。
func (s *Service) UpdateModule(ctx context.Context, req ModuleEditParam) error {
	client := req.Client
	if client == "" {
		client = string(security.AccountAdmin)
	}
	updates := map[string]any{
		"name": req.Name, "code": req.Code, "client": client, "icon": req.Icon, "color": req.Color,
		"sort": sortOr(req.Sort), "status": statusOr(req.Status), "description": req.Description,
	}
	return s.repo.UpdateModule(ctx, req.ID, updates)
}

// DeleteModules 批量删除资源模块。
func (s *Service) DeleteModules(ctx context.Context, ids []string) error {
	return s.repo.DeleteModules(ctx, ids)
}

// ModuleDetail 资源模块详情。
func (s *Service) ModuleDetail(ctx context.Context, id string) (*ResourceModule, error) {
	return s.repo.GetModuleByID(ctx, id)
}

// ModulePage 资源模块分页。
func (s *Service) ModulePage(ctx context.Context, p ModulePageParam) (rows []ResourceModule, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageModules(ctx, p)
	return rows, total, current, size, err
}

// ModuleSelector 资源模块选择器。
func (s *Service) ModuleSelector(ctx context.Context, client string) ([]ModuleOption, error) {
	rows, err := s.repo.ListEnabledModules(ctx, client)
	if err != nil {
		return nil, err
	}
	out := make([]ModuleOption, 0, len(rows))
	for _, r := range rows {
		out = append(out, ModuleOption{ID: r.ID, Name: r.Name, Code: r.Code})
	}
	return out, nil
}

func buildTree(rows []Resource, parent *string) []TreeNode {
	out := make([]TreeNode, 0)
	for _, r := range rows {
		same := (r.ParentID == nil && parent == nil) || (r.ParentID != nil && parent != nil && *r.ParentID == *parent)
		if same {
			out = append(out, TreeNode{Resource: r, Children: buildTree(rows, &r.ID)})
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
