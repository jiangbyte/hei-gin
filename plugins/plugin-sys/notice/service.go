package notice

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

func (s *Service) Page(c *gin.Context, p *NoticePageParam) {
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

	vos := make([]*NoticeVO, len(rows))
	for i, r := range rows {
		vos[i] = SysNoticeToNoticeVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *Service) Detail(c *gin.Context, id string) *NoticeVO {
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
		result.WriteError(c, exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
		return nil
	}
	return SysNoticeToNoticeVO(e)
}

func (s *Service) Create(c *gin.Context, vo *NoticeVO) {
	ctx := c.Request.Context()

	e := NoticeVOToSysNotice(vo)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加通知失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Modify(c *gin.Context, vo *NoticeVO) {
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
		result.WriteError(c, exception.NewBusinessError("查询通知失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{
		"title": vo.Title, "category": vo.Category, "type": vo.Type,
		"sort_code": vo.SortCode,
	}
	if vo.Summary != nil {
		up["summary"] = *vo.Summary
	}
	if vo.Content != nil {
		up["content"] = *vo.Content
	}
	if vo.Cover != nil {
		up["cover"] = *vo.Cover
	}
	if vo.Level != "" {
		up["level"] = vo.Level
	}
	if vo.Status != "" {
		up["status"] = vo.Status
	}
	if vo.IsTop != "" {
		up["is_top"] = vo.IsTop
	}
	if vo.Author != nil {
		up["author"] = *vo.Author
	}
	if vo.PublishAt != nil {
		up["publish_at"] = utils.ParseDateTimePtr(vo.PublishAt)
	}
	if vo.ExpireAt != nil {
		up["expire_at"] = utils.ParseDateTimePtr(vo.ExpireAt)
	}
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑通知失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Remove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := s.repo.DeleteByIDs(ctx, ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除通知失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Options(c *gin.Context) []any {
	ctx := c.Request.Context()
	rows := s.repo.ListAll(ctx)
	vos := make([]any, len(rows))
	for i, r := range rows {
		vos[i] = SysNoticeToNoticeVO(&r)
	}
	return vos
}

func (s *Service) DetailByID(c *gin.Context, id string) *NoticeVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
		return nil
	}
	return SysNoticeToNoticeVO(e)
}

func (s *Service) Latest(c *gin.Context, param *NoticeLatestParam) []*NoticeVO {
	if param.Size < 1 {
		param.Size = 5
	}
	if param.Size > 20 {
		param.Size = 20
	}

	ctx := c.Request.Context()
	rows := s.repo.Latest(ctx, param.Size)
	vos := make([]*NoticeVO, len(rows))
	for i, r := range rows {
		vos[i] = SysNoticeToNoticeVO(&r)
	}
	return vos
}

func (s *Service) PublicPage(c *gin.Context, p *NoticePageParam) {
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

	rows, total := s.repo.PublicPage(ctx, p)

	vos := make([]*NoticeVO, len(rows))
	for i, r := range rows {
		vos[i] = SysNoticeToNoticeVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *Service) PublicDetail(c *gin.Context, id string) *NoticeVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindEnabledByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
		return nil
	}
	return SysNoticeToNoticeVO(e)
}
