package role

import (
	"context"

	"gorm.io/gorm"
)

// Repo 角色持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建角色。
func (r *Repo) Create(ctx context.Context, row *Role) error {
	return r.with(ctx).Create(row).Error
}

// Update 更新角色。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Role{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Role{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Role, error) {
	var row Role
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page 分页查询。
func (r *Repo) Page(ctx context.Context, p PageParam) (rows []Role, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Role{})
	if p.Code != "" {
		db = db.Where("code ILIKE ?", "%"+p.Code+"%")
	}
	if p.Name != "" {
		db = db.Where("name ILIKE ?", "%"+p.Name+"%")
	}
	if p.Status != "" {
		db = db.Where("status = ?", p.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort asc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
