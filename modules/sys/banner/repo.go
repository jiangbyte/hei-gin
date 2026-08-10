package banner

import (
	"context"

	"gorm.io/gorm"
)

// Repo Banner 持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建 Banner。
func (r *Repo) Create(ctx context.Context, row *Banner) error {
	return r.with(ctx).Create(row).Error
}

// Update 按 ID 更新字段。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Banner{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&Banner{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Banner, error) {
	var row Banner
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page 分页查询。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []Banner, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&Banner{})
	if q.Title != "" {
		db = db.Where("title ILIKE ?", "%"+q.Title+"%")
	}
	if q.Position != "" {
		db = db.Where("position = ?", q.Position)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort asc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// List 按位置列出启用 Banner。
func (r *Repo) List(ctx context.Context, position, status string) ([]Banner, error) {
	db := r.with(ctx).Model(&Banner{}).Where("status = ?", status)
	if position != "" {
		db = db.Where("position = ?", position)
	}
	var rows []Banner
	err := db.Order("sort asc, id asc").Find(&rows).Error
	return rows, err
}
