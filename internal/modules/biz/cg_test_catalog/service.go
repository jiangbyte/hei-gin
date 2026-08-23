// internal/modules/biz/cg_test_catalog/service.go 业务服务。
//
// Author: Charlie

package cg_test_catalog

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/biz/scope"
)

// Service 目录服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 biz.cg_test_catalog 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "biz.cg_test_catalog",
		Order:  91,
		Models: []any{&Catalog{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建目录。
func (s *Service) Create(ctx context.Context, accountID string, req AddParam, sess *security.SessionPayload) error {
	row := Catalog{
		ID: idgen.Next(), ParentID: req.ParentID, Code: req.Code, Name: req.Name, Category: req.Category,
		Status: req.Status, Sort: req.Sort, IsVisible: req.IsVisible, Icon: req.Icon, Description: req.Description,
		Extra: mustJSON(req.Extra), CreatedBy: &accountID, UpdatedBy: &accountID,
		OwnerDeptID: scope.DefaultOwnerDeptID(sess),
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新目录。
func (s *Service) Update(ctx context.Context, accountID string, req EditParam) error {
	return s.repo.Update(ctx, req.ID, map[string]any{
		"parent_id": req.ParentID, "code": req.Code, "name": req.Name, "category": req.Category,
		"status": req.Status, "sort": req.Sort, "is_visible": req.IsVisible, "icon": req.Icon,
		"description": req.Description, "extra": mustJSON(req.Extra), "updated_by": accountID,
	})
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 目录详情。
func (s *Service) Detail(ctx context.Context, id string) (*Catalog, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Catalog, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p, sess)
	return rows, total, current, size, err
}

// Tree 目录树。
func (s *Service) Tree(ctx context.Context, sess *security.SessionPayload) ([]TreeNode, error) {
	rows, err := s.repo.ListAll(ctx, sess)
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
