// internal/modules/biz/cg_test_knowledge_category/repo.go 持久化仓储。
//
// Author: Charlie

package cg_test_knowledge_category

import (
	"context"

	"gorm.io/gorm"
)

// Repo 知识分类持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// CreateCategory 创建分类。
func (r *Repo) CreateCategory(ctx context.Context, row *Category) error {
	return r.with(ctx).Create(row).Error
}

// UpdateCategory 更新分类。
func (r *Repo) UpdateCategory(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Category{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCategoriesByIDs 批量删除分类。
func (r *Repo) DeleteCategoriesByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Category{}).Error
}

// GetCategoryByID 按主键查分类。
func (r *Repo) GetCategoryByID(ctx context.Context, id string) (*Category, error) {
	var row Category
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageCategories 分页查分类。
func (r *Repo) PageCategories(ctx context.Context, p PageParam) (rows []Category, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Category{})
	if p.Code != "" {
		db = db.Where("code ILIKE ?", "%"+p.Code+"%")
	}
	if p.Name != "" {
		db = db.Where("name ILIKE ?", "%"+p.Name+"%")
	}
	if p.Status != "" {
		db = db.Where("status = ?", p.Status)
	}
	if p.ParentID != "" {
		db = db.Where("parent_id = ?", p.ParentID)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort ASC, created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// ListAllCategories 列出全部分类。
func (r *Repo) ListAllCategories(ctx context.Context) ([]Category, error) {
	var rows []Category
	err := r.with(ctx).Order("sort ASC").Find(&rows).Error
	return rows, err
}

// CreateDoc 创建文档。
func (r *Repo) CreateDoc(ctx context.Context, row *Doc) error {
	return r.with(ctx).Create(row).Error
}

// UpdateDoc 更新文档。
func (r *Repo) UpdateDoc(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Doc{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteDocsByIDs 批量删除文档。
func (r *Repo) DeleteDocsByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Doc{}).Error
}

// GetDocByID 按主键查文档。
func (r *Repo) GetDocByID(ctx context.Context, id string) (*Doc, error) {
	var row Doc
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// PageDocs 分页查文档。
func (r *Repo) PageDocs(ctx context.Context, p DocPageParam) (rows []Doc, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Doc{})
	if p.CategoryID != "" {
		db = db.Where("category_id = ?", p.CategoryID)
	}
	if p.Code != "" {
		db = db.Where("code ILIKE ?", "%"+p.Code+"%")
	}
	if p.Title != "" {
		db = db.Where("title ILIKE ?", "%"+p.Title+"%")
	}
	if p.Status != "" {
		db = db.Where("status = ?", p.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort ASC, created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
