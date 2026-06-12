package config

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) pageQuery(ctx context.Context, p *ConfigPageParam) *gorm.DB {
	q := r.db.WithContext(ctx).Model(&SysConfig{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("config_key LIKE ? OR remark LIKE ?", like, like)
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	return q
}

func (r *repository) Page(ctx context.Context, p *ConfigPageParam) ([]SysConfig, int64) {
	q := r.pageQuery(ctx, p)
	var total int64
	q.Count(&total)
	var rows []SysConfig
	q.Order("sort_code ASC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysConfig, error) {
	var e SysConfig
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *SysConfig) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysConfig{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysConfig{}).Error
}

func (r *repository) ListAll(ctx context.Context) []SysConfig {
	var rows []SysConfig
	r.db.WithContext(ctx).Model(&SysConfig{}).Order("sort_code ASC").Find(&rows)
	return rows
}

func (r *repository) ListByCategory(ctx context.Context, category string) []SysConfig {
	var rows []SysConfig
	r.db.WithContext(ctx).Where("category = ?", category).Order("sort_code ASC").Find(&rows)
	return rows
}

func (r *repository) EditBatch(ctx context.Context, items []ConfigBatchEditItem) error {
	tx := r.db.WithContext(ctx).Begin()
	for _, item := range items {
		up := map[string]interface{}{}
		if item.ConfigKey != nil {
			up["config_key"] = *item.ConfigKey
		}
		if item.ConfigValue != nil {
			up["config_value"] = *item.ConfigValue
		}
		if item.Remark != nil {
			up["remark"] = *item.Remark
		}
		if item.SortCode != 0 {
			up["sort_code"] = item.SortCode
		}
		if len(up) == 0 {
			continue
		}
		if err := tx.Model(&SysConfig{}).Where("id = ?", item.ID).Updates(up).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *repository) EditByCategory(ctx context.Context, category string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysConfig{}).Where("category = ?", category).Updates(up).Error
}
