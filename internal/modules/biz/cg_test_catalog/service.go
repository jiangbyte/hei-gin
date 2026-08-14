package cg_test_catalog

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service ç›®å½•æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º biz.cg_test_catalog æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "biz.cg_test_catalog",
		Order:  91,
		Models: []any{&Catalog{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create åˆ›å»ºç›®å½•ã€‚
func (s *Service) Create(ctx context.Context, accountID string, req AddParam) error {
	row := Catalog{
		ID: idgen.Next(), ParentID: req.ParentID, Code: req.Code, Name: req.Name, Category: req.Category,
		Status: req.Status, Sort: req.Sort, IsVisible: req.IsVisible, Icon: req.Icon, Description: req.Description,
		Extra: mustJSON(req.Extra), CreatedBy: &accountID, UpdatedBy: &accountID,
	}
	return s.repo.Create(ctx, &row)
}

// Update æ›´æ–°ç›®å½•ã€‚
func (s *Service) Update(ctx context.Context, accountID string, req EditParam) error {
	return s.repo.Update(ctx, req.ID, map[string]any{
		"parent_id": req.ParentID, "code": req.Code, "name": req.Name, "category": req.Category,
		"status": req.Status, "sort": req.Sort, "is_visible": req.IsVisible, "icon": req.Icon,
		"description": req.Description, "extra": mustJSON(req.Extra), "updated_by": accountID,
	})
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail ç›®å½•è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Catalog, error) {
	return s.repo.GetByID(ctx, id)
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Catalog, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p)
	return rows, total, current, size, err
}

// Tree ç›®å½•æ ‘ã€‚
func (s *Service) Tree(ctx context.Context) ([]TreeNode, error) {
	rows, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	byParent := map[string][]Catalog{}
	var roots []Catalog
	for _, r := range rows {
		if r.ParentID == nil || *r.ParentID == "" {
			roots = append(roots, r)
			continue
		}
		byParent[*r.ParentID] = append(byParent[*r.ParentID], r)
	}
	var build func(Catalog) TreeNode
	build = func(n Catalog) TreeNode {
		node := TreeNode{Catalog: n}
		for _, ch := range byParent[n.ID] {
			node.Children = append(node.Children, build(ch))
		}
		if node.Children == nil {
			node.Children = []TreeNode{}
		}
		return node
	}
	out := make([]TreeNode, 0, len(roots))
	for _, r := range roots {
		out = append(out, build(r))
	}
	return out, nil
}

func mustJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return b
}
