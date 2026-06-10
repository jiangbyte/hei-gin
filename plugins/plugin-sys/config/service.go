package config

import (
	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

// ===== Page =====

func ConfigPage(c *gin.Context, p *ConfigPageParam) {
	ctx := c.Request.Context()
	if p.Current < 1 {
		p.Current = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}

	q := db.DB.WithContext(ctx).Model(&SysConfig{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("config_key LIKE ? OR remark LIKE ?", like, like)
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}

	var total int64
	q.Count(&total)

	var rows []SysConfig
	q.Order("sort_code ASC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*ConfigVO, len(rows))
	for i, r := range rows {
		vos[i] = SysConfigToConfigVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

// ===== Detail =====

func ConfigDetail(c *gin.Context, id string) *ConfigVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysConfig
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询配置详情失败: "+err.Error(), 500))
		return nil
	}
	return SysConfigToConfigVO(&e)
}

// ===== Create =====

func ConfigCreate(c *gin.Context, vo *ConfigVO) {
	ctx := c.Request.Context()

	e := ConfigVOToSysConfig(vo)
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加配置失败: "+err.Error(), 500))
		return
	}
}

// ===== Modify =====

func ConfigModify(c *gin.Context, vo *ConfigVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var e SysConfig
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询配置失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{
		"sort_code": vo.SortCode,
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
	if err := db.DB.WithContext(ctx).Model(&SysConfig{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑配置失败: "+err.Error(), 500))
		return
	}
}

// ===== Remove =====

func ConfigRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	if err := db.DB.WithContext(c.Request.Context()).Where("id IN ?", ids).Delete(&SysConfig{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除配置失败: "+err.Error(), 500))
		return
	}
}

// ===== Options =====

func ConfigOptions(c *gin.Context) []*ConfigVO {
	ctx := c.Request.Context()
	var rows []SysConfig
	db.DB.WithContext(ctx).Model(&SysConfig{}).Order("sort_code ASC").Find(&rows)
	vos := make([]*ConfigVO, len(rows))
	for i, r := range rows {
		vos[i] = SysConfigToConfigVO(&r)
	}
	return vos
}

// ===== ListByCategory =====

func ConfigListByCategory(c *gin.Context, category string) []*ConfigVO {
	ctx := c.Request.Context()
	var rows []SysConfig
	db.DB.WithContext(ctx).Where("category = ?", category).Order("sort_code ASC").Find(&rows)
	vos := make([]*ConfigVO, len(rows))
	for i, r := range rows {
		vos[i] = SysConfigToConfigVO(&r)
	}
	return vos
}

// ===== EditBatch =====

func ConfigEditBatch(c *gin.Context, param *ConfigBatchEditParam) {
	ctx := c.Request.Context()
	tx := db.DB.WithContext(ctx).Begin()
	for _, item := range param.Configs {
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
			result.WriteError(c, exception.NewBusinessError("批量编辑配置失败: "+err.Error(), 500))
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
		return
	}
}

// ===== EditByCategory =====

func ConfigEditByCategory(c *gin.Context, param *ConfigCategoryEditParam) {
	ctx := c.Request.Context()

	up := map[string]interface{}{}
	if param.ConfigKey != nil {
		up["config_key"] = *param.ConfigKey
	}
	if param.ConfigValue != nil {
		up["config_value"] = *param.ConfigValue
	}
	if param.Remark != nil {
		up["remark"] = *param.Remark
	}
	if len(up) == 0 {
		return
	}
	if err := db.DB.WithContext(ctx).Model(&SysConfig{}).Where("category = ?", param.Category).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("按分类编辑配置失败: "+err.Error(), 500))
		return
	}
}
