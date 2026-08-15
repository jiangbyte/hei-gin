// internal/modules/sys/dict/repo.go 持久化仓储。
//
// Author: Charlie

package dict

import (
	"context"

	"gorm.io/gorm"
)

// Repo 数据字典持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// FindByCode 按编码查询字典。
func (r *Repo) FindByCode(ctx context.Context, code string) (*Dict, error) {
	var row Dict
	if err := r.with(ctx).Where("code = ?", code).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Create 创建字典。
func (r *Repo) Create(ctx context.Context, row *Dict) error {
	return r.with(ctx).Create(row).Error
}

// Update 按 ID 更新。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Dict{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Dict{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Dict, error) {
	var row Dict
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page 分页查询。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []Dict, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&Dict{})
	if q.Code != "" {
		db = db.Where("code ILIKE ?", "%"+q.Code+"%")
	}
	if q.Category != "" {
		db = db.Where("category = ?", q.Category)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if q.ParentID != "" {
		db = db.Where("parent_id = ?", q.ParentID)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort asc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// ListForTree 查询树形字典节点。
func (r *Repo) ListForTree(ctx context.Context, q TreeParam, status string) ([]Dict, error) {
	db := r.with(ctx).Model(&Dict{}).Where("status = ?", status)
	if q.Code != "" {
		db = db.Where("code = ? OR parent_id IN (SELECT id FROM sys_dict WHERE code = ?)", q.Code, q.Code)
	}
	if q.Category != "" {
		db = db.Where("category = ?", q.Category)
	}
	var rows []Dict
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}
