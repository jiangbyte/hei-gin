// internal/modules/iam/client/service.go 业务服务。
//
// Author: Charlie

package client

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service å®¢æˆ·ç«¯èµ„æºæœåŠ¡ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ å®¢æˆ·ç«¯èµ„æºæœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º iam.client æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.client",
		Models: []any{&ClientModule{}, &ClientResource{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// ListGrantModules å®¢æˆ·ç«¯èµ„æºæŽˆæƒæ¨¡å—é€‰é¡¹ï¼ˆå«æ¨¡å—ä¸‹å¯ç”¨èµ„æºï¼Œç©ºæ¨¡å—è¿‡æ»¤ï¼‰ã€‚
func (s *Service) ListGrantModules(ctx context.Context, accountType string) ([]GrantModule, error) {
	at := accountType
	if at == "" {
		at = string(security.AccountAdmin)
	}
	modules, err := s.repo.ListEnabledModules(ctx, at)
	if err != nil {
		return nil, err
	}
	resources, err := s.repo.ListGrantResources(ctx, at)
	if err != nil {
		return nil, err
	}
	out := make([]GrantModule, 0, len(modules))
	for _, m := range modules {
		gm := GrantModule{ModuleID: m.ID, Name: m.Name, Resources: []ClientResource{}}
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

// CreateModule åˆ›å»ºå®¢æˆ·ç«¯æ¨¡å—ã€‚
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

// UpdateModule æ›´æ–°å®¢æˆ·ç«¯æ¨¡å—ã€‚
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

// DeleteModules æ‰¹é‡åˆ é™¤å®¢æˆ·ç«¯æ¨¡å—ã€‚
func (s *Service) DeleteModules(ctx context.Context, ids []string) error {
	return s.repo.DeleteModules(ctx, ids)
}

// ModuleDetail å®¢æˆ·ç«¯æ¨¡å—è¯¦æƒ…ã€‚
func (s *Service) ModuleDetail(ctx context.Context, id string) (*ClientModule, error) {
	return s.repo.GetModuleByID(ctx, id)
}

// ModulePage å®¢æˆ·ç«¯æ¨¡å—åˆ†é¡µã€‚
func (s *Service) ModulePage(ctx context.Context, p ModulePageParam) (rows []ClientModule, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageModules(ctx, p)
	return rows, total, current, size, err
}

// ModuleSelector å®¢æˆ·ç«¯æ¨¡å—é€‰æ‹©å™¨ã€‚
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

// CreateResource åˆ›å»ºå®¢æˆ·ç«¯èµ„æºã€‚
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

// UpdateResource æ›´æ–°å®¢æˆ·ç«¯èµ„æºã€‚
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

// DeleteResources æ‰¹é‡åˆ é™¤å®¢æˆ·ç«¯èµ„æºã€‚
func (s *Service) DeleteResources(ctx context.Context, ids []string) error {
	return s.repo.DeleteResources(ctx, ids)
}

// ResourceDetail å®¢æˆ·ç«¯èµ„æºè¯¦æƒ…ã€‚
func (s *Service) ResourceDetail(ctx context.Context, id string) (*ClientResource, error) {
	return s.repo.GetResourceByID(ctx, id)
}

// ResourcePage å®¢æˆ·ç«¯èµ„æºåˆ†é¡µã€‚
func (s *Service) ResourcePage(ctx context.Context, p ResourcePageParam) (rows []ClientResource, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageResources(ctx, p)
	return rows, total, current, size, err
}

// ResourceTree å®¢æˆ·ç«¯èµ„æºæ ‘ã€‚
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
