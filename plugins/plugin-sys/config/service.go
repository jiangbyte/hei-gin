package config

import (
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/crud"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func Page(c *gin.Context, param *ConfigPageParam) {
	crud.Page(c, &SysConfig{}, param, func(q *gorm.DB) *gorm.DB {
		if param.Keyword != "" {
			like := "%" + param.Keyword + "%"
			q = q.Where("config_key LIKE ? OR remark LIKE ?", like, like)
		}
		if param.Category != "" {
			q = q.Where("category = ?", param.Category)
		}
		return q
	}, "sort_code ASC", func(e *SysConfig) any { return toVO(e) })
}

func Detail(c *gin.Context, id string) *ConfigVO {
	if id == "" {
		return nil
	}
	ctx := c.Request.Context()
	var entity SysConfig
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(exception.NewBusinessError("查询配置详情失败: "+err.Error(), 500))
	}
	return toVO(&entity)
}

func Create(c *gin.Context, vo *ConfigVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	entity := SysConfig{
		ID:        utils.GenerateID(),
		SortCode:  vo.SortCode,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if vo.ConfigKey != nil {
		entity.ConfigKey = vo.ConfigKey
	}
	if vo.ConfigValue != nil {
		entity.ConfigValue = vo.ConfigValue
	}
	if vo.Remark != nil {
		entity.Remark = vo.Remark
	}
	if vo.Category != nil {
		entity.Category = vo.Category
	}
	if userID != "" {
		entity.CreatedBy = &userID
		entity.UpdatedBy = &userID
	}
	if vo.Extra != nil {
		entity.Extra = vo.Extra
	}
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("添加配置失败: "+err.Error(), 500))
	}
}

func Modify(c *gin.Context, vo *ConfigVO, userID string) {
	ctx := c.Request.Context()
	var entity SysConfig
	if err := db.DB.WithContext(ctx).Where("id = ?", vo.ID).First(&entity).Error; err != nil {
		panic(exception.NewBusinessError("配置不存在: "+err.Error(), 500))
	}
	now := time.Now()
	up := map[string]any{
		"sort_code":  vo.SortCode,
		"updated_at": now,
	}
	if vo.ConfigKey != nil {
		up["config_key"] = *vo.ConfigKey
	}
	if vo.ConfigValue != nil {
		up["config_value"] = *vo.ConfigValue
	}
	if vo.Remark != nil {
		up["remark"] = *vo.Remark
	}
	if vo.Category != nil {
		up["category"] = *vo.Category
	}
	if vo.Extra != nil {
		up["extra"] = *vo.Extra
	}
	if userID != "" {
		up["updated_by"] = userID
	}
	if err := db.DB.WithContext(ctx).Model(&entity).Updates(up).Error; err != nil {
		panic(exception.NewBusinessError("编辑配置失败: "+err.Error(), 500))
	}
}

func Remove(c *gin.Context, ids []string) {
	crud.Remove(c, &SysConfig{}, ids)
}

func Options(c *gin.Context) []any {
	return crud.Options(c, &SysConfig{}, "sort_code ASC", func(e *SysConfig) any { return toVO(e) })
}

func ListByCategory(c *gin.Context, category string) []ConfigVO {
	ctx := c.Request.Context()
	var records []SysConfig
	db.DB.WithContext(ctx).Where("category = ?", category).Order("sort_code ASC").Find(&records)
	return toVOList(records)
}

func EditBatch(c *gin.Context, param *ConfigBatchEditParam, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	tx := db.DB.WithContext(ctx).Begin()
	for _, item := range param.Configs {
		up := map[string]interface{}{"updated_at": now}
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
		if userID != "" {
			up["updated_by"] = userID
		}
		if err := tx.Model(&SysConfig{}).Where("id = ?", item.ID).Updates(up).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBusinessError("批量编辑配置失败: "+err.Error(), 500))
		}
	}
	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
	}
}

func EditByCategory(c *gin.Context, param *ConfigCategoryEditParam, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	up := map[string]interface{}{"updated_at": now}
	if param.ConfigKey != nil {
		up["config_key"] = *param.ConfigKey
	}
	if param.ConfigValue != nil {
		up["config_value"] = *param.ConfigValue
	}
	if param.Remark != nil {
		up["remark"] = *param.Remark
	}
	if userID != "" {
		up["updated_by"] = userID
	}
	if err := db.DB.WithContext(ctx).Model(&SysConfig{}).Where("category = ?", param.Category).Updates(up).Error; err != nil {
		panic(exception.NewBusinessError("按分类编辑配置失败: "+err.Error(), 500))
	}
}
