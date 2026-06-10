package banner

import (
	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

// ===== Page =====

func BannerPage(c *gin.Context, p *BannerPageParam) {
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

	q := db.DB.WithContext(ctx).Model(&SysBanner{})

	var total int64
	q.Count(&total)

	var rows []SysBanner
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*BannerVO, len(rows))
	for i, r := range rows {
		vos[i] = SysBannerToBannerVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

// ===== Detail =====

func BannerDetail(c *gin.Context, id string) *BannerVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysBanner
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询横幅详情失败: "+err.Error(), 500))
		return nil
	}
	return SysBannerToBannerVO(&e)
}

// ===== Create =====

func BannerCreate(c *gin.Context, vo *BannerVO) {
	ctx := c.Request.Context()

	e := BannerVOToSysBanner(vo)
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加横幅失败: "+err.Error(), 500))
		return
	}
}

// ===== Modify =====

func BannerModify(c *gin.Context, vo *BannerVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var e SysBanner
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询横幅失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{
		"title": vo.Title, "image": vo.Image,
		"link_type": vo.LinkType, "category": vo.Category,
		"type": vo.Type, "position": vo.Position,
		"sort_code": vo.SortCode, "view_count": vo.ViewCount,
		"click_count": vo.ClickCount,
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
	if err := db.DB.WithContext(ctx).Model(&SysBanner{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑横幅失败: "+err.Error(), 500))
		return
	}
}

// ===== Remove =====

func BannerRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	if err := db.DB.WithContext(c.Request.Context()).Where("id IN ?", ids).Delete(&SysBanner{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除横幅失败: "+err.Error(), 500))
		return
	}
}

// ===== Options =====

func BannerOptions(c *gin.Context) []*BannerVO {
	ctx := c.Request.Context()
	var rows []SysBanner
	db.DB.WithContext(ctx).Model(&SysBanner{}).Order("sort_code ASC").Find(&rows)
	vos := make([]*BannerVO, len(rows))
	for i, r := range rows {
		vos[i] = SysBannerToBannerVO(&r)
	}
	return vos
}
