package position

import (
	"gorm.io/gorm"

	"hei-gin/sdk/db"
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

	q := db.DB.WithContext(ctx).Model(&SysPosition{})
	if p.Keyword != "" {
		q = q.Where("name LIKE ?", "%"+p.Keyword+"%")
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}

	var total int64
	q.Count(&total)

	var rows []SysPosition
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

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
	var e SysPosition
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询职位失败: "+err.Error(), 500))
		return nil
	}
	return SysPositionToPositionVO(&e)
}

func PositionCreate(c *gin.Context, vo *PositionVO) {
	ctx := c.Request.Context()

	e := PositionVOToSysPosition(vo)
	e.Status = string(enums.StatusEnabled)
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
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

	var e SysPosition
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
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
	if err := db.DB.WithContext(ctx).Model(&SysPosition{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
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
	db.DB.WithContext(ctx).Table("sys_user").Where("position_id IN ?", ids).Update("position_id", nil)
	if err := db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysPosition{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除职位失败: "+err.Error(), 500))
		return
	}
}

func PositionOptions(c *gin.Context) []any {
	ctx := c.Request.Context()
	var rows []SysPosition
	db.DB.WithContext(ctx).Model(&SysPosition{}).Order("sort_code ASC").Find(&rows)
	vos := make([]any, len(rows))
	for i, r := range rows {
		vos[i] = SysPositionToPositionVO(&r)
	}
	return vos
}
