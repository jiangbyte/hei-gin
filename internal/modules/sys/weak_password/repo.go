package weakpassword

import (
	"context"

	"gorm.io/gorm"
)

// Repo 弱密码持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建弱密码。
func (r *Repo) Create(ctx context.Context, row *WeakPassword) error {
	return r.with(ctx).Create(row).Error
}

// UpdatePassword 更新密码值。
func (r *Repo) UpdatePassword(ctx context.Context, id, password string) error {
	return r.with(ctx).Model(&WeakPassword{}).Where("id = ?", id).Update("password", password).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ?", ids).Delete(&WeakPassword{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*WeakPassword, error) {
	var row WeakPassword
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page 分页查询。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []WeakPassword, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&WeakPassword{})
	if q.Password != "" {
		db = db.Where("password ILIKE ?", "%"+q.Password+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// List 列表查询。
func (r *Repo) List(ctx context.Context, q ListParam) ([]WeakPassword, error) {
	db := r.with(ctx).Model(&WeakPassword{})
	if q.Password != "" {
		db = db.Where("password ILIKE ?", "%"+q.Password+"%")
	}
	var rows []WeakPassword
	err := db.Order("id desc").Find(&rows).Error
	return rows, err
}
