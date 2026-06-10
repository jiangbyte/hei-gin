package notice

import (
	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func NoticePage(c *gin.Context, p *NoticePageParam) {
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

	q := db.DB.WithContext(ctx).Model(&SysNotice{})
	if p.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+p.Keyword+"%")
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}

	var total int64
	q.Count(&total)

	var rows []SysNotice
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*NoticeVO, len(rows))
	for i, r := range rows {
		vos[i] = SysNoticeToNoticeVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func NoticeDetail(c *gin.Context, id string) *NoticeVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysNotice
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
		return nil
	}
	return SysNoticeToNoticeVO(&e)
}

func NoticeCreate(c *gin.Context, vo *NoticeVO) {
	ctx := c.Request.Context()

	e := NoticeVOToSysNotice(vo)
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加通知失败: "+err.Error(), 500))
		return
	}
}

func NoticeModify(c *gin.Context, vo *NoticeVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var e SysNotice
	if err := db.DB.WithContext(ctx).Where("id = ?", vo.ID).First(&e).Error; err != nil {
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
	if err := db.DB.WithContext(ctx).Model(&SysNotice{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑通知失败: "+err.Error(), 500))
		return
	}
}

func NoticeRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysNotice{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除通知失败: "+err.Error(), 500))
		return
	}
}

func NoticeOptions(c *gin.Context) []any {
	ctx := c.Request.Context()
	var rows []SysNotice
	db.DB.WithContext(ctx).Model(&SysNotice{}).Order("sort_code ASC").Find(&rows)
	vos := make([]any, len(rows))
	for i, r := range rows {
		vos[i] = SysNoticeToNoticeVO(&r)
	}
	return vos
}

func NoticeDetailByID(c *gin.Context, id string) *NoticeVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysNotice
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
		return nil
	}
	return SysNoticeToNoticeVO(&e)
}

func NoticeLatest(c *gin.Context, param *NoticeLatestParam) []*NoticeVO {
	if param.Size < 1 {
		param.Size = 5
	}
	if param.Size > 20 {
		param.Size = 20
	}

	ctx := c.Request.Context()
	var rows []SysNotice
	db.DB.WithContext(ctx).
		Where("status = ?", string(enums.StatusEnabled)).
		Order("is_top DESC, sort_code DESC, created_at DESC").
		Limit(param.Size).
		Find(&rows)
	vos := make([]*NoticeVO, len(rows))
	for i, r := range rows {
		vos[i] = SysNoticeToNoticeVO(&r)
	}
	return vos
}

func NoticePublicPage(c *gin.Context, p *NoticePageParam) {
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

	q := db.DB.WithContext(ctx).Model(&SysNotice{}).Where("status = ?", string(enums.StatusEnabled))
	if p.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+p.Keyword+"%")
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}

	var total int64
	q.Count(&total)

	var rows []SysNotice
	q.Order("is_top DESC, sort_code DESC, created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*NoticeVO, len(rows))
	for i, r := range rows {
		vos[i] = SysNoticeToNoticeVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func NoticePublicDetail(c *gin.Context, id string) *NoticeVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysNotice
	if err := db.DB.WithContext(ctx).Where("id = ? AND status = ?", id, string(enums.StatusEnabled)).First(&e).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
		return nil
	}
	return SysNoticeToNoticeVO(&e)
}
