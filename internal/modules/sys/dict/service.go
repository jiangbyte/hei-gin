// internal/modules/sys/dict/service.go 业务服务。
//
// Author: Charlie

package dict

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service æ•°æ®å­—å…¸ä¸šåŠ¡æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ å­—å…¸æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º sys.dict æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.dict",
		Models: []any{&Dict{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create åˆ›å»ºå­—å…¸ã€‚
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Dict{
		ID: idgen.Next(), Code: req.Code, Label: req.Label, Value: req.Value, Color: req.Color,
		Category: req.Category, ParentID: req.ParentID, Status: statusOr(req.Status), Sort: req.Sort,
	}
	return s.repo.Create(ctx, &row)
}

// Update æ›´æ–°å­—å…¸ã€‚
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"code": req.Code, "label": req.Label, "value": req.Value, "color": req.Color,
		"category": req.Category, "parent_id": req.ParentID, "status": statusOr(req.Status), "sort": req.Sort,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Dict, error) {
	return s.repo.GetByID(ctx, id)
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Dict, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// Tree å­—å…¸æ ‘ã€‚
func (s *Service) Tree(ctx context.Context, q TreeParam) ([]TreeNode, error) {
	rows, err := s.repo.ListForTree(ctx, q, security.StatusEnabled)
	if err != nil {
		return nil, err
	}
	return buildTree(rows, nil), nil
}

func buildTree(rows []Dict, parent *string) []TreeNode {
	out := make([]TreeNode, 0)
	for _, r := range rows {
		same := (r.ParentID == nil && parent == nil) || (r.ParentID != nil && parent != nil && *r.ParentID == *parent)
		if same {
			out = append(out, TreeNode{Dict: r, Children: buildTree(rows, &r.ID)})
		}
	}
	return out
}

func statusOr(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}
