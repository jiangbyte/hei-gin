package banner

import (
	"gorm.io/gorm"

	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

func (s *Service) Page(c *gin.Context, p *BannerPageParam) {
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

	vos := make([]*BannerVO, len(rows))
	for i, r := range rows {
		vos[i] = SysBannerToBannerVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *Service) Detail(c *gin.Context, id string) *BannerVO {
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
		result.WriteError(c, exception.NewBusinessError("查询横幅详情失败: "+err.Error(), 500))
		return nil
	}
	return SysBannerToBannerVO(e)
}

func (s *Service) Create(c *gin.Context, vo *BannerVO) {
	ctx := c.Request.Context()

	e := BannerVOToSysBanner(vo)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加横幅失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Modify(c *gin.Context, vo *BannerVO) {
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
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑横幅失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Remove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	if err := s.repo.DeleteByIDs(c.Request.Context(), ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除横幅失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Options(c *gin.Context) []*BannerVO {
	ctx := c.Request.Context()
	rows := s.repo.ListAll(ctx)
	vos := make([]*BannerVO, len(rows))
	for i, r := range rows {
		vos[i] = SysBannerToBannerVO(&r)
	}
	return vos
}
