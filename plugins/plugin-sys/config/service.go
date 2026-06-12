package config

import (
	"gorm.io/gorm"

	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

// ===== Page =====

func (s *Service) Page(c *gin.Context, p *ConfigPageParam) {
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

	rows, total := s.repo.Page(ctx, p)

	vos := make([]*ConfigVO, len(rows))
	for i, r := range rows {
		vos[i] = SysConfigToConfigVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

// ===== Detail =====

func (s *Service) Detail(c *gin.Context, id string) *ConfigVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询配置详情失败: "+err.Error(), 500))
		return nil
	}
	return SysConfigToConfigVO(e)
}

// ===== Create =====

func (s *Service) Create(c *gin.Context, vo *ConfigVO) {
	ctx := c.Request.Context()

	e := ConfigVOToSysConfig(vo)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加配置失败: "+err.Error(), 500))
		return
	}
}

// ===== Modify =====

func (s *Service) Modify(c *gin.Context, vo *ConfigVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	if _, err := s.repo.FindByID(ctx, vo.ID); err != nil {
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
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑配置失败: "+err.Error(), 500))
		return
	}
}

// ===== Remove =====

func (s *Service) Remove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	if err := s.repo.DeleteByIDs(c.Request.Context(), ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除配置失败: "+err.Error(), 500))
		return
	}
}

// ===== Options =====

func (s *Service) Options(c *gin.Context) []*ConfigVO {
	ctx := c.Request.Context()
	rows := s.repo.ListAll(ctx)
	vos := make([]*ConfigVO, len(rows))
	for i, r := range rows {
		vos[i] = SysConfigToConfigVO(&r)
	}
	return vos
}

// ===== ListByCategory =====

func (s *Service) ListByCategory(c *gin.Context, category string) []*ConfigVO {
	ctx := c.Request.Context()
	rows := s.repo.ListByCategory(ctx, category)
	vos := make([]*ConfigVO, len(rows))
	for i, r := range rows {
		vos[i] = SysConfigToConfigVO(&r)
	}
	return vos
}

// ===== EditBatch =====

func (s *Service) EditBatch(c *gin.Context, param *ConfigBatchEditParam) {
	ctx := c.Request.Context()
	if err := s.repo.EditBatch(ctx, param.Configs); err != nil {
		result.WriteError(c, exception.NewBusinessError("批量编辑配置失败: "+err.Error(), 500))
		return
	}
}

// ===== EditByCategory =====

func (s *Service) EditByCategory(c *gin.Context, param *ConfigCategoryEditParam) {
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
	if err := s.repo.EditByCategory(ctx, param.Category, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("按分类编辑配置失败: "+err.Error(), 500))
		return
	}
}
