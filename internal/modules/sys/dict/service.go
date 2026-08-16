// internal/modules/sys/dict/service.go 业务服务。
//
// Author: Charlie

package dict

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
)

// Service 数据字典业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造字典服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.dict 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.dict",
		Models: []any{&Dict{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建字典（code 唯一 + 分类校验；对齐 hei-boot DictServiceImpl）。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	if err := validateCode(req.Code); err != nil {
		return err
	}
	if err := validateCategory(req.Category); err != nil {
		return err
	}
	if _, err := s.repo.FindByCode(ctx, req.Code); err == nil {
		return fmt.Errorf("字典编码已存在")
	}
	row := Dict{
		ID: idgen.Next(), Code: req.Code, Label: req.Label, Value: req.Value, Color: req.Color,
		Category: req.Category, ParentID: req.ParentID, Status: statusOr(req.Status), Sort: req.Sort,
	}
	return s.repo.Create(ctx, &row)
}

// Update 更新字典（404 + code 唯一排除自身）。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return fmt.Errorf("字典不存在")
	}
	if err := validateCode(req.Code); err != nil {
		return err
	}
	if err := validateCategory(req.Category); err != nil {
		return err
	}
	if existing, err := s.repo.FindByCode(ctx, req.Code); err == nil && existing.ID != req.ID {
		return fmt.Errorf("字典编码已存在")
	}
	updates := map[string]any{
		"code": req.Code, "label": req.Label, "value": req.Value, "color": req.Color,
		"category": req.Category, "parent_id": req.ParentID, "status": statusOr(req.Status), "sort": req.Sort,
	}
	return s.repo.Update(ctx, req.ID, updates)
}

func validateCode(code string) error {
	for _, r := range code {
		if !(r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_') {
			return fmt.Errorf("字典编码仅支持大写字母、数字与下划线")
		}
	}
	return nil
}

func validateCategory(category *string) error {
	if category != nil && *category != "" && *category != "SYS" && *category != "BIZ" {
		return fmt.Errorf("字典分类仅支持 SYS / BIZ")
	}
	return nil
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Dict, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Dict, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// Tree 字典树（含全部状态；断链父节点置为根，对齐 hei-boot DictServiceImpl.tree）。
func (s *Service) Tree(ctx context.Context, q TreeParam) ([]TreeNode, error) {
	rows, err := s.repo.ListForTree(ctx, q, "")
	if err != nil {
		return nil, err
	}
	return buildTree(rows, nil), nil
}

// buildTree 组装字典树：父节点不在集合内（断链/孤儿）时提升为根节点。
func buildTree(rows []Dict, parent *string) []TreeNode {
	ids := make(map[string]struct{}, len(rows))
	for i := range rows {
		ids[rows[i].ID] = struct{}{}
	}
	out := make([]TreeNode, 0)
	for _, r := range rows {
		same := (r.ParentID == nil && parent == nil) || (r.ParentID != nil && parent != nil && *r.ParentID == *parent)
		if !same {
			continue
		}
		out = append(out, TreeNode{Dict: r, Children: buildTree(rows, &r.ID)})
	}
	// 孤儿节点（父 id 指向集合外）作为根节点补回
	if parent == nil {
		for _, r := range rows {
			if r.ParentID == nil {
				continue
			}
			if _, ok := ids[*r.ParentID]; !ok {
				out = append(out, TreeNode{Dict: r, Children: buildTree(rows, &r.ID)})
			}
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
