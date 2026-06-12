package position

import (
	"gorm.io/gorm"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type service struct {
	repo *repository
}

func (s *service) PositionPage(c *gin.Context, p *PositionPageParam) {
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

	vos := make([]*PositionVO, len(rows))
	for i, r := range rows {
		vos[i] = SysPositionToPositionVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *service) PositionDetail(c *gin.Context, id string) *PositionVO {
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
		result.WriteError(c, exception.NewBusinessError("查询职位失败: "+err.Error(), 500))
		return nil
	}
	return SysPositionToPositionVO(e)
}

func (s *service) PositionCreate(c *gin.Context, vo *PositionVO) {
	ctx := c.Request.Context()

	e := PositionVOToSysPosition(vo)
	e.Status = string(enums.StatusEnabled)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加职位失败: "+err.Error(), 500))
		return
	}
}

func (s *service) PositionModify(c *gin.Context, vo *PositionVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	_, err := s.repo.FindByID(ctx, vo.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询职位失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{
		"code": vo.Code, "name": vo.Name, "category": vo.Category,
		"sort_code": vo.SortCode,
	}
	if vo.OrgID != nil {
		up["org_id"] = *vo.OrgID
	}
	if vo.GroupID != nil {
		up["group_id"] = *vo.GroupID
	}
	if vo.Description != nil {
		up["description"] = *vo.Description
	}
	if vo.Extra != nil {
		up["extra"] = *vo.Extra
	}
	if vo.Status != "" {
		up["status"] = vo.Status
	}
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑职位失败: "+err.Error(), 500))
		return
	}
}

func (s *service) PositionRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := s.repo.DeleteByIDs(ctx, ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除职位失败: "+err.Error(), 500))
		return
	}
}

func (s *service) PositionOptions(c *gin.Context) []any {
	ctx := c.Request.Context()
	rows := s.repo.ListAll(ctx)
	vos := make([]any, len(rows))
	for i, r := range rows {
		vos[i] = SysPositionToPositionVO(&r)
	}
	return vos
}

func PositionPage(c *gin.Context, p *PositionPageParam) {
	defaultModule.service.PositionPage(c, p)
}

func PositionDetail(c *gin.Context, id string) *PositionVO {
	return defaultModule.service.PositionDetail(c, id)
}

func PositionCreate(c *gin.Context, vo *PositionVO) {
	defaultModule.service.PositionCreate(c, vo)
}

func PositionModify(c *gin.Context, vo *PositionVO) {
	defaultModule.service.PositionModify(c, vo)
}

func PositionRemove(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.PositionRemove(c, param)
}

func PositionOptions(c *gin.Context) []any {
	return defaultModule.service.PositionOptions(c)
}
