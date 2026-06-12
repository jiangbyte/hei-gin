package position

import (
	"gorm.io/gorm"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func PositionPage(c *gin.Context, p *PositionPageParam) {
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

	rows, total := Page(ctx, p)

	vos := make([]*PositionVO, len(rows))
	for i, r := range rows {
		vos[i] = SysPositionToPositionVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func PositionDetail(c *gin.Context, id string) *PositionVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := FindByID(ctx, id)
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

func PositionCreate(c *gin.Context, vo *PositionVO) {
	ctx := c.Request.Context()

	e := PositionVOToSysPosition(vo)
	e.Status = string(enums.StatusEnabled)
	if err := Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加职位失败: "+err.Error(), 500))
		return
	}
}

func PositionModify(c *gin.Context, vo *PositionVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	_, err := FindByID(ctx, vo.ID)
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
	if err := UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑职位失败: "+err.Error(), 500))
		return
	}
}

func PositionRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := DeleteByIDs(ctx, ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除职位失败: "+err.Error(), 500))
		return
	}
}

func PositionOptions(c *gin.Context) []any {
	ctx := c.Request.Context()
	rows := ListAll(ctx)
	vos := make([]any, len(rows))
	for i, r := range rows {
		vos[i] = SysPositionToPositionVO(&r)
	}
	return vos
}
