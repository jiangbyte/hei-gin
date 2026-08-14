package dept

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service éƒ¨é—¨æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ éƒ¨é—¨æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New æž„å»º iam.dept æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "iam.dept",
		Models: []any{&Dept{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create åˆ›å»ºéƒ¨é—¨ã€‚
func (s *Service) Create(ctx context.Context, req AddParam) error {
	row := Dept{
		ID: idgen.Next(), ParentID: req.ParentID, MasterID: req.MasterID, DeputyMasterID: req.DeputyMasterID,
		Name: req.Name, Category: req.Category, Sort: req.Sort, IsVirtual: req.IsVirtual,
		Status: orStatus(req.Status), Extra: datatypes.JSON([]byte("{}")),
	}
	if row.Sort == 0 {
		row.Sort = 99
	}
	return s.repo.Create(ctx, &row)
}

// Update æ›´æ–°éƒ¨é—¨ã€‚
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{
		"parent_id": req.ParentID, "master_id": req.MasterID, "deputy_master_id": req.DeputyMasterID,
		"name": req.Name, "category": req.Category, "sort": req.Sort, "is_virtual": req.IsVirtual,
		"status": orStatus(req.Status),
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete æ‰¹é‡åˆ é™¤ã€‚
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail éƒ¨é—¨è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*Dept, error) {
	return s.repo.GetByID(ctx, id)
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Dept, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p)
	return rows, total, current, size, err
}

// Tree éƒ¨é—¨æ ‘ã€‚
func (s *Service) Tree(ctx context.Context) ([]TreeNode, error) {
	rows, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildDeptTree(rows, nil), nil
}

func buildDeptTree(rows []Dept, parent *string) []TreeNode {
	out := make([]TreeNode, 0)
	for _, r := range rows {
		if eqPtr(r.ParentID, parent) {
			n := TreeNode{Dept: r, Children: buildDeptTree(rows, &r.ID)}
			out = append(out, n)
		}
	}
	return out
}

func eqPtr(a *string, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func orStatus(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}
