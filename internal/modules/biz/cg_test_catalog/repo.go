// internal/modules/biz/cg_test_catalog/repo.go 持久化仓储。
//
// Author: Charlie

package cg_test_catalog

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
	"hei-gin/internal/framework/platform/db/dialect"
)

// Repo 目录持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建目录。
func (r *Repo) Create(ctx context.Context, row *Catalog) error {
	return r.with(ctx).Create(row).Error
}

// Update 更新目录。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Catalog{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Catalog{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Catalog, error) {
	var row Catalog
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page 分页查询；sess 非空时按 owner_dept_id 数据范围过滤。
func (r *Repo) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Catalog, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Catalog{})
	if sess != nil {
		db = datascope.ApplyKey(db, sess, "biz:cgtestcatalog:page", "owner_dept_id", "created_by")
	}
	if p.Code != "" {
		db = db.Where(dialect.ILike(db, "code"), dialect.Contains(p.Code))
	}
	if p.Name != "" {
		db = db.Where(dialect.ILike(db, "name"), dialect.Contains(p.Name))
	}
	if p.Category != "" {
		db = db.Where("category = ?", p.Category)
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

// ListAll 列出全部（树形）；sess 非空时按数据范围过滤。
func (r *Repo) ListAll(ctx context.Context, sess *security.SessionPayload) ([]Catalog, error) {
	var rows []Catalog
	db := r.with(ctx)
	if sess != nil {
		db = datascope.ApplyKey(db, sess, "biz:cgtestcatalog:list", "owner_dept_id", "created_by")
	}
	err := db.Order("sort ASC").Find(&rows).Error
	return rows, err
}
