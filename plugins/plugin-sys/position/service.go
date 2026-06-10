package position

import (
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/crud"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func Page(c *gin.Context, param *PositionPageParam) {
	crud.Page(c, &SysPosition{}, param, func(q *gorm.DB) *gorm.DB {
		if param.Keyword != "" {
			q = q.Where("name LIKE ?", "%"+param.Keyword+"%")
		}
		if param.Category != "" {
			q = q.Where("category = ?", param.Category)
		}
		return q
	}, "created_at DESC", func(e *SysPosition) any { return toVO(e) })
}

func Detail(c *gin.Context, id string) *PositionVO {
	var entity SysPosition
	crud.Detail(c, &entity, id, "职位")
	return toVO(&entity)
}

func Create(c *gin.Context, vo *PositionVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	entity := SysPosition{
		ID:        utils.GenerateID(),
		Code:      vo.Code,
		Name:      vo.Name,
		Category:  vo.Category,
		OrgID:     vo.OrgID,
		GroupID:   vo.GroupID,
		Status:    vo.Status,
		SortCode:  vo.SortCode,
		CreatedAt: &now,
		CreatedBy: &userID,
		UpdatedAt: &now,
		UpdatedBy: &userID,
	}
	if vo.Description != nil {
		entity.Description = vo.Description
	}
	if vo.Extra != nil {
		entity.Extra = vo.Extra
	}
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("添加职位失败: "+err.Error(), 500))
	}
}

func Modify(c *gin.Context, vo *PositionVO, userID string) {
	ctx := c.Request.Context()
	var entity SysPosition
	if err := db.DB.WithContext(ctx).Where("id = ?", vo.ID).First(&entity).Error; err != nil {
		panic(exception.NewBusinessError("职位不存在: "+err.Error(), 500))
	}
	now := time.Now()
	up := map[string]any{
		"code":       vo.Code,
		"name":       vo.Name,
		"category":   vo.Category,
		"org_id":     vo.OrgID,
		"group_id":   vo.GroupID,
		"status":     vo.Status,
		"sort_code":  vo.SortCode,
		"updated_at": now,
		"updated_by": userID,
	}
	if vo.Description != nil {
		up["description"] = *vo.Description
	}
	if vo.Extra != nil {
		up["extra"] = *vo.Extra
	}
	if err := db.DB.WithContext(ctx).Model(&entity).Updates(up).Error; err != nil {
		panic(exception.NewBusinessError("编辑职位失败: "+err.Error(), 500))
	}
}

func Remove(c *gin.Context, ids []string) {
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	db.DB.WithContext(ctx).Table("sys_user").Where("position_id IN ?", ids).Update("position_id", nil)
	db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysPosition{})
}

func Options(c *gin.Context) []any {
	return crud.Options(c, &SysPosition{}, "sort_code ASC", func(e *SysPosition) any { return toVO(e) })
}
