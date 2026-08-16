// internal/modules/biz/cg_test_knowledge_category/service.go 业务服务。
//
// Author: Charlie

package cg_test_knowledge_category

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
)

// Service 知识分类服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 biz.cg_test_knowledge_category 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "biz.cg_test_knowledge_category",
		Order:  93,
		Models: []any{&Category{}, &Doc{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建分类。
func (s *Service) Create(ctx context.Context, accountID string, req AddParam) error {
	row := Category{
		ID: idgen.Next(), ParentID: req.ParentID, Code: req.Code, Name: req.Name, Status: req.Status,
		Sort: req.Sort, IsVisible: req.IsVisible, Description: req.Description, Extra: mustJSON(req.Extra),
		CreatedBy: &accountID, UpdatedBy: &accountID,
	}
	return s.repo.CreateCategory(ctx, &row)
}

// Update 更新分类。
func (s *Service) Update(ctx context.Context, accountID string, req EditParam) error {
	return s.repo.UpdateCategory(ctx, req.ID, map[string]any{
		"parent_id": req.ParentID, "code": req.Code, "name": req.Name, "status": req.Status,
		"sort": req.Sort, "is_visible": req.IsVisible, "description": req.Description,
		"extra": mustJSON(req.Extra), "updated_by": accountID,
	})
}

// Delete 批量删除分类。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteCategoriesByIDs(ctx, ids)
}

// Detail 分类详情。
func (s *Service) Detail(ctx context.Context, id string) (*Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, p PageParam) (rows []Category, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageCategories(ctx, p)
	return rows, total, current, size, err
}

// Tree 分类树。
func (s *Service) Tree(ctx context.Context) ([]TreeNode, error) {
	rows, err := s.repo.ListAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	byParent := map[string][]Category{}
	var roots []Category
	for _, r := range rows {
		if r.ParentID == nil || *r.ParentID == "" {
			roots = append(roots, r)
			continue
		}
		byParent[*r.ParentID] = append(byParent[*r.ParentID], r)
	}
	var build func(Category) TreeNode
	build = func(n Category) TreeNode {
		node := TreeNode{Category: n, Children: []TreeNode{}}
		for _, ch := range byParent[n.ID] {
			node.Children = append(node.Children, build(ch))
		}
		return node
	}
	out := make([]TreeNode, 0, len(roots))
	for _, r := range roots {
		out = append(out, build(r))
	}
	return out, nil
}

// CreateDoc 创建文档。
func (s *Service) CreateDoc(ctx context.Context, accountID string, req DocAddParam) error {
	row := Doc{
		ID: idgen.Next(), CategoryID: req.CategoryID, Code: req.Code, Title: req.Title, Type: req.Type,
		Status: req.Status, Summary: req.Summary, Content: req.Content, Author: req.Author,
		PublishedAt: req.PublishedAt, ViewCount: req.ViewCount, Sort: req.Sort, IsTop: req.IsTop,
		Settings: mustJSON(req.Settings), Extra: mustJSON(req.Extra), CreatedBy: &accountID, UpdatedBy: &accountID,
	}
	return s.repo.CreateDoc(ctx, &row)
}

// UpdateDoc 更新文档。
func (s *Service) UpdateDoc(ctx context.Context, accountID string, req DocEditParam) error {
	return s.repo.UpdateDoc(ctx, req.ID, map[string]any{
		"category_id": req.CategoryID, "code": req.Code, "title": req.Title, "type": req.Type, "status": req.Status,
		"summary": req.Summary, "content": req.Content, "author": req.Author, "published_at": req.PublishedAt,
		"view_count": req.ViewCount, "sort": req.Sort, "is_top": req.IsTop, "settings": mustJSON(req.Settings),
		"extra": mustJSON(req.Extra), "updated_by": accountID,
	})
}

// DeleteDocs 批量删除文档。
func (s *Service) DeleteDocs(ctx context.Context, ids []string) error {
	return s.repo.DeleteDocsByIDs(ctx, ids)
}

// DetailDoc 文档详情。
func (s *Service) DetailDoc(ctx context.Context, id string) (*Doc, error) {
	return s.repo.GetDocByID(ctx, id)
}

// PageDocs 文档分页。
func (s *Service) PageDocs(ctx context.Context, p DocPageParam) (rows []Doc, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageDocs(ctx, p)
	return rows, total, current, size, err
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
