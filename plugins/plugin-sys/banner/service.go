package banner

import (
	"time"

	"hei-gin/sdk/crud"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func Page(c *gin.Context, param *BannerPageParam) {
	crud.Page(c, &SysBanner{}, param, nil, "created_at DESC", func(e *SysBanner) any { return toVO(e) })
}

func Detail(c *gin.Context, id string) *BannerVO {
	var entity SysBanner
	crud.Detail(c, &entity, id, "横幅")
	return toVO(&entity)
}

func Create(c *gin.Context, vo *BannerVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	entity := SysBanner{
		ID:         utils.GenerateID(),
		Title:      vo.Title,
		Image:      vo.Image,
		LinkType:   vo.LinkType,
		Category:   vo.Category,
		Type:       vo.Type,
		Position:   vo.Position,
		SortCode:   vo.SortCode,
		ViewCount:  vo.ViewCount,
		ClickCount: vo.ClickCount,
		CreatedAt:  &now,
		CreatedBy:  &userID,
		UpdatedAt:  &now,
		UpdatedBy:  &userID,
	}
	if vo.URL != nil {
		entity.URL = vo.URL
	}
	if vo.Summary != nil {
		entity.Summary = vo.Summary
	}
	if vo.Description != nil {
		entity.Description = vo.Description
	}
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("添加横幅失败: "+err.Error(), 500))
	}
}

func Modify(c *gin.Context, vo *BannerVO, userID string) {
	ctx := c.Request.Context()
	var entity SysBanner
	if err := db.DB.WithContext(ctx).Where("id = ?", vo.ID).First(&entity).Error; err != nil {
		panic(exception.NewBusinessError("横幅不存在: "+err.Error(), 500))
	}
	now := time.Now()
	up := map[string]any{
		"title":       vo.Title,
		"image":       vo.Image,
		"link_type":   vo.LinkType,
		"category":    vo.Category,
		"type":        vo.Type,
		"position":    vo.Position,
		"sort_code":   vo.SortCode,
		"view_count":  vo.ViewCount,
		"click_count": vo.ClickCount,
		"updated_at":  now,
		"updated_by":  userID,
	}
	if vo.URL != nil {
		up["url"] = *vo.URL
	}
	if vo.Summary != nil {
		up["summary"] = *vo.Summary
	}
	if vo.Description != nil {
		up["description"] = *vo.Description
	}
	if err := db.DB.WithContext(ctx).Model(&entity).Updates(up).Error; err != nil {
		panic(exception.NewBusinessError("编辑横幅失败: "+err.Error(), 500))
	}
}

func Remove(c *gin.Context, ids []string) {
	crud.Remove(c, &SysBanner{}, ids)
}

func Options(c *gin.Context) []any {
	return crud.Options(c, &SysBanner{}, "sort_code ASC", func(e *SysBanner) any { return toVO(e) })
}

func toVO(entity *SysBanner) *BannerVO {
	vo := &BannerVO{
		ID:         entity.ID,
		Title:      entity.Title,
		Image:      entity.Image,
		LinkType:   entity.LinkType,
		Category:   entity.Category,
		Type:       entity.Type,
		Position:   entity.Position,
		SortCode:   entity.SortCode,
		ViewCount:  entity.ViewCount,
		ClickCount: entity.ClickCount,
	}
	if entity.URL != nil {
		vo.URL = entity.URL
	}
	if entity.Summary != nil {
		vo.Summary = entity.Summary
	}
	if entity.Description != nil {
		vo.Description = entity.Description
	}
	if entity.CreatedAt != nil {
		s := entity.CreatedAt.Format("2006-01-02 15:04:05")
		vo.CreatedAt = &s
	}
	if entity.CreatedBy != nil {
		vo.CreatedBy = entity.CreatedBy
	}
	if entity.UpdatedAt != nil {
		s := entity.UpdatedAt.Format("2006-01-02 15:04:05")
		vo.UpdatedAt = &s
	}
	if entity.UpdatedBy != nil {
		vo.UpdatedBy = entity.UpdatedBy
	}
	return vo
}
