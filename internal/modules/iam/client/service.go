// internal/modules/iam/client/service.go 业务服务。
//
// Author: Charlie

package client

import (
	"context"
	"sort"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/iam/relation"
)

// Service 客户端资源服务。
//
// Author: Charlie
type Service struct {
	repo *Repo
	rel  *relation.Service
}

// NewService 构造客户端资源服务。
func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepo(db), rel: relation.NewService(db)}
}

// New 构建 iam.client 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.client",
		Models: []any{&ClientModule{}, &ClientResource{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// ListGrantModules 客户端资源授权模块选项树（模块 → 菜单 → 按钮权限；对齐 hei-boot ClientResourceServiceImpl.listGrantModules）。
func (s *Service) ListGrantModules(ctx context.Context, accountType string) ([]GrantModule, error) {
	at := accountType
	if at == "" {
		at = string(security.AccountAdmin)
	}
	// 1. 加载启用模块（按账号类型过滤）
	modules, err := s.repo.ListEnabledModules(ctx, at)
	if err != nil {
		return nil, err
	}
	// 2. 加载模块下启用资源
	resources, err := s.repo.ListGrantResources(ctx, at)
	if err != nil {
		return nil, err
	}
	if len(resources) == 0 {
		return []GrantModule{}, nil
	}
	resourceIDs := make([]string, 0, len(resources))
	resourceMap := make(map[string]*ClientResource, len(resources))
	for i := range resources {
		resourceIDs = append(resourceIDs, resources[i].ID)
		resourceMap[resources[i].ID] = &resources[i]
	}
	// 3. 按资源汇总权限选项（CLIENT_RESOURCE_PERMISSION 关系）
	permissionMap, err := s.repo.GrantPermissions(ctx, resourceIDs, at)
	if err != nil {
		return nil, err
	}
	// 4. 按钮/动作权限挂到父菜单（无关系时回退 code/name）
	childPermissionMap := make(map[string][]PermissionOption)
	for i := range resources {
		res := &resources[i]
		if res.ResourceType != ResourceTypeButton && res.ResourceType != ResourceTypeAction {
			continue
		}
		if res.ParentID == nil || *res.ParentID == "" {
			continue
		}
		options := permissionMap[res.ID]
		if len(options) == 0 {
			options = []PermissionOption{{ID: res.ID, PermissionKey: res.Code, Title: res.Name}}
		}
		childPermissionMap[*res.ParentID] = append(childPermissionMap[*res.ParentID], options...)
	}
	// 5. 组装模块 → 菜单 → 按钮授权树
	moduleMap := make(map[string]*GrantModule, len(modules))
	moduleSort := make(map[string]int, len(modules))
	for i := range modules {
		m := &modules[i]
		moduleMap[m.ID] = &GrantModule{ID: m.ID, Title: m.Name, Menu: []GrantMenuOption{}}
		moduleSort[m.ID] = m.Sort
	}
	for i := range resources {
		res := &resources[i]
		if !GrantMenuTypes[res.ResourceType] || res.ModuleID == nil || *res.ModuleID == "" {
			continue
		}
		moduleID := *res.ModuleID
		mod := moduleMap[moduleID]
		if mod == nil {
			mod = &GrantModule{ID: moduleID, Title: moduleID, Menu: []GrantMenuOption{}}
			moduleMap[moduleID] = mod
			moduleSort[moduleID] = 99
		}
		menu := GrantMenuOption{
			ID: res.ID, ModuleID: moduleID, ParentID: res.ParentID, Title: res.Name,
			Button: append([]PermissionOption{}, permissionMap[res.ID]...),
		}
		if res.ParentID != nil && *res.ParentID != "" {
			if p, ok := resourceMap[*res.ParentID]; ok {
				menu.ParentIDName = p.Name
			} else {
				menu.ParentIDName = res.Name
			}
		} else {
			menu.ParentIDName = res.Name
		}
		menu.Button = append(menu.Button, childPermissionMap[res.ID]...)
		mod.Menu = append(mod.Menu, menu)
	}
	// 6. 过滤空模块并按排序返回
	out := make([]GrantModule, 0, len(moduleMap))
	for _, mod := range moduleMap {
		if len(mod.Menu) == 0 {
			continue
		}
		out = append(out, *mod)
	}
	sort.Slice(out, func(i, j int) bool {
		if moduleSort[out[i].ID] != moduleSort[out[j].ID] {
			return moduleSort[out[i].ID] < moduleSort[out[j].ID]
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
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

// DeleteResources 批量删除客户端资源（先清关联，再删；对齐 hei-boot ClientResourceServiceImpl.delete）。
func (s *Service) DeleteResources(ctx context.Context, ids []string) error {
	_ = s.rel.DeleteBySubjectIDs(ctx, relation.SubjectClientResource, ids, "")
	_ = s.rel.DeleteByTargetIDs(ctx, relation.TargetClientResource, ids, "")
	return s.repo.DeleteResources(ctx, ids)
}

// ResourceDetail 客户端资源详情。
func (s *Service) ResourceDetail(ctx context.Context, id string) (*ClientResource, error) {
	row, err := s.repo.GetResourceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.withNames(ctx, []*ClientResource{row})
	return row, nil
}

// ResourcePage 客户端资源分页。
func (s *Service) ResourcePage(ctx context.Context, p ResourcePageParam) (rows []ClientResource, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageResources(ctx, p)
	if err != nil {
		return nil, 0, current, size, err
	}
	s.withNames(ctx, toPtrs2(rows))
	return rows, total, current, size, nil
}

// ResourceTree 客户端资源树（module_id 或 account_type 过滤；对齐 hei-boot ClientResourceServiceImpl.tree）。
func (s *Service) ResourceTree(ctx context.Context, moduleID, accountType string) ([]TreeNode, error) {
	rows, err := s.repo.ListResources(ctx, moduleID, accountType)
	if err != nil {
		return nil, err
	}
	s.withNames(ctx, toPtrs2(rows))
	return buildTree(rows, nil), nil
}

// withNames 批量回填 parent_id_name / module_id_name / account_type（对齐 hei-boot tree 的 transBatch + loadModuleAccountTypes）。
func (s *Service) withNames(ctx context.Context, rows []*ClientResource) {
	moduleIDs := make([]string, 0, len(rows))
	parentIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.ModuleID != nil && *row.ModuleID != "" {
			moduleIDs = append(moduleIDs, *row.ModuleID)
		}
		if row.ParentID != nil && *row.ParentID != "" {
			parentIDs = append(parentIDs, *row.ParentID)
		}
	}
	modules := map[string]ClientModule{}
	if len(moduleIDs) > 0 {
		ms, err := s.repo.GetModulesByIDs(ctx, moduleIDs)
		if err == nil {
			for i := range ms {
				modules[ms[i].ID] = ms[i]
			}
		}
	}
	parents := map[string]string{}
	if len(parentIDs) > 0 {
		ps, err := s.repo.GetResourcesByIDs(ctx, parentIDs)
		if err == nil {
			for i := range ps {
				parents[ps[i].ID] = ps[i].Name
			}
		}
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.ModuleID != nil {
			if m, ok := modules[*row.ModuleID]; ok {
				name := m.Name
				row.ModuleIDName = &name
				at := m.AccountType
				row.AccountType = &at
			}
		}
		if row.ParentID != nil {
			if n, ok := parents[*row.ParentID]; ok {
				row.ParentIDName = &n
			}
		}
	}
}

func toPtrs2(rows []ClientResource) []*ClientResource {
	out := make([]*ClientResource, len(rows))
	for i := range rows {
		out[i] = &rows[i]
	}
	return out
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
