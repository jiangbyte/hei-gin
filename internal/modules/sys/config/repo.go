// internal/modules/sys/config/repo.go 持久化仓储。
//
// Author: Charlie

package config

import (
	"context"

	"gorm.io/gorm"
)

// Repo 系统配置持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建配置。
func (r *Repo) Create(ctx context.Context, row *Config) error {
	return r.with(ctx).Create(row).Error
}

// Update 按 ID 更新。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&Config{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 删除非内置配置。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("id IN ? AND is_builtin = ?", ids, false).Delete(&Config{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*Config, error) {
	var row Config
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// GetByKey 按配置键查询。
func (r *Repo) GetByKey(ctx context.Context, key string) (*Config, error) {
	var row Config
	if err := r.with(ctx).Where("config_key = ?", key).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ListByIDs 按主键集合查询。
func (r *Repo) ListByIDs(ctx context.Context, ids []string) ([]Config, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []Config
	if err := r.with(ctx).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// Page 分页查询。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []Config, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&Config{})
	if q.ConfigKey != "" {
		db = db.Where("config_key ILIKE ?", "%"+q.ConfigKey+"%")
	}
	if q.Category != "" {
		db = db.Where("category ILIKE ?", "%"+q.Category+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("sort_code asc, id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// ListByKeys 按配置键批量查询。
func (r *Repo) ListByKeys(ctx context.Context, keys []string) ([]Config, error) {
	var rows []Config
	err := r.with(ctx).Where("config_key IN ?", keys).Find(&rows).Error
	return rows, err
}

// List 按分类/范围列出配置。
func (r *Repo) List(ctx context.Context, q ListParam) ([]Config, error) {
	db := r.with(ctx).Model(&Config{})
	if q.Category != "" {
		db = db.Where("category = ?", q.Category)
	}
	if q.Scope != "" {
		db = db.Where("scope = ?", q.Scope)
	}
	var rows []Config
	err := db.Order("sort_code asc, id desc").Find(&rows).Error
	return rows, err
}
