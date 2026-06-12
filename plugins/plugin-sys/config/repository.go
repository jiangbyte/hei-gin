package config

import (
	"context"

	"hei-gin/sdk/db"

	"gorm.io/gorm"
)

func pageQuery(ctx context.Context, p *ConfigPageParam) *gorm.DB {
	q := db.DB.WithContext(ctx).Model(&SysConfig{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("config_key LIKE ? OR remark LIKE ?", like, like)
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	return q
}

func Page(ctx context.Context, p *ConfigPageParam) ([]SysConfig, int64) {
	q := pageQuery(ctx, p)
	var total int64
	q.Count(&total)
	var rows []SysConfig
	q.Order("sort_code ASC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func FindByID(ctx context.Context, id string) (*SysConfig, error) {
	var e SysConfig
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func Create(ctx context.Context, entity *SysConfig) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysConfig{}).Where("id = ?", id).Updates(up).Error
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysConfig{}).Error
}

func ListAll(ctx context.Context) []SysConfig {
	var rows []SysConfig
	db.DB.WithContext(ctx).Model(&SysConfig{}).Order("sort_code ASC").Find(&rows)
	return rows
}

func ListByCategory(ctx context.Context, category string) []SysConfig {
	var rows []SysConfig
	db.DB.WithContext(ctx).Where("category = ?", category).Order("sort_code ASC").Find(&rows)
	return rows
}

func EditBatch(ctx context.Context, items []ConfigBatchEditItem) error {
	tx := db.DB.WithContext(ctx).Begin()
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

func EditByCategory(ctx context.Context, category string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysConfig{}).Where("category = ?", category).Updates(up).Error
}
