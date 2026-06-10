package notice

import (
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/crud"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func parseTime(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", *s)
	if err != nil {
		return nil
	}
	return &t
}

func Page(c *gin.Context, param *NoticePageParam) {
	crud.Page(c, &SysNotice{}, param, func(q *gorm.DB) *gorm.DB {
		if param.Keyword != "" {
			q = q.Where("title LIKE ?", "%"+param.Keyword+"%")
		}
		if param.Category != "" {
			q = q.Where("category = ?", param.Category)
		}
		if param.Status != "" {
			q = q.Where("status = ?", param.Status)
		}
		return q
	}, "created_at DESC", func(e *SysNotice) any { return entToVO(e) })
}

func Detail(c *gin.Context, id string) *NoticeVO {
	if id == "" {
		return nil
	}
	ctx := c.Request.Context()
	var entity SysNotice
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
	}
	return entToVO(&entity)
}

func Create(c *gin.Context, vo *NoticeVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	entity := SysNotice{
		ID:        utils.GenerateID(),
		Title:     vo.Title,
		Category:  vo.Category,
		Type:      vo.Type,
		SortCode:  vo.SortCode,
		CreatedAt: &now,
		CreatedBy: &userID,
		UpdatedAt: &now,
		UpdatedBy: &userID,
	}
	if vo.Summary != nil {
		entity.Summary = vo.Summary
	}
	if vo.Content != nil {
		entity.Content = vo.Content
	}
	if vo.Cover != nil {
		entity.Cover = vo.Cover
	}
	if vo.Level != "" {
		entity.Level = vo.Level
	}
	if vo.Status != "" {
		entity.Status = vo.Status
	}
	if vo.IsTop != "" {
		entity.IsTop = vo.IsTop
	}
	if vo.Author != nil {
		entity.Author = vo.Author
	}
	if vo.PublishAt != nil {
		entity.PublishAt = parseTime(vo.PublishAt)
	}
	if vo.ExpireAt != nil {
		entity.ExpireAt = parseTime(vo.ExpireAt)
	}
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("添加通知失败: "+err.Error(), 500))
	}
}

func Modify(c *gin.Context, vo *NoticeVO, userID string) {
	ctx := c.Request.Context()
	var entity SysNotice
	if err := db.DB.WithContext(ctx).Where("id = ?", vo.ID).First(&entity).Error; err != nil {
		panic(exception.NewBusinessError("通知不存在: "+err.Error(), 500))
	}
	now := time.Now()
	up := map[string]any{
		"title":      vo.Title,
		"category":   vo.Category,
		"type":       vo.Type,
		"sort_code":  vo.SortCode,
		"updated_at": now,
		"updated_by": userID,
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
		up["publish_at"] = parseTime(vo.PublishAt)
	}
	if vo.ExpireAt != nil {
		up["expire_at"] = parseTime(vo.ExpireAt)
	}
	if err := db.DB.WithContext(ctx).Model(&entity).Updates(up).Error; err != nil {
		panic(exception.NewBusinessError("编辑通知失败: "+err.Error(), 500))
	}
}

func Remove(c *gin.Context, ids []string) {
	crud.Remove(c, &SysNotice{}, ids)
}

func Options(c *gin.Context) []any {
	return crud.Options(c, &SysNotice{}, "sort_code ASC", func(e *SysNotice) any { return entToVO(e) })
}

func DetailByID(c *gin.Context, id string) *NoticeVO {
	if id == "" {
		return nil
	}
	ctx := c.Request.Context()
	var entity SysNotice
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		panic(exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
	}
	return entToVO(&entity)
}

func Latest(c *gin.Context, param *NoticeLatestParam) []*NoticeVO {
	ctx := c.Request.Context()
	var records []SysNotice
	db.DB.WithContext(ctx).
		Where("status = ?", "ENABLED").
		Order("is_top DESC, sort_code DESC, created_at DESC").
		Limit(param.Size).
		Find(&records)
	vos := make([]*NoticeVO, len(records))
	for i, r := range records {
		vos[i] = entToVO(&r)
	}
	return vos
}

func entToVO(entity *SysNotice) *NoticeVO {
	vo := &NoticeVO{
		ID: entity.ID, Title: entity.Title, Category: entity.Category,
		Type: entity.Type, SortCode: entity.SortCode,
	}
	if entity.Summary != nil {
		vo.Summary = entity.Summary
	}
	if entity.Content != nil {
		vo.Content = entity.Content
	}
	if entity.Cover != nil {
		vo.Cover = entity.Cover
	}
	if entity.Level != "" {
		vo.Level = entity.Level
	}
	if entity.Status != "" {
		vo.Status = entity.Status
	}
	if entity.IsTop != "" {
		vo.IsTop = entity.IsTop
	}
	if entity.Author != nil {
		vo.Author = entity.Author
	}
	if entity.PublishAt != nil {
		s := entity.PublishAt.Format("2006-01-02 15:04:05")
		vo.PublishAt = &s
	}
	if entity.ExpireAt != nil {
		s := entity.ExpireAt.Format("2006-01-02 15:04:05")
		vo.ExpireAt = &s
	}
	if entity.CreatedAt != nil {
		vo.CreatedAt = entity.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if entity.CreatedBy != nil {
		vo.CreatedBy = entity.CreatedBy
	}
	if entity.UpdatedAt != nil {
		vo.UpdatedAt = entity.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	if entity.UpdatedBy != nil {
		vo.UpdatedBy = entity.UpdatedBy
	}
	return vo
}

func PublicPage(c *gin.Context, param *NoticePageParam) {
	crud.Page(c, &SysNotice{}, param, func(q *gorm.DB) *gorm.DB {
		q = q.Where("status = ?", "ENABLED")
		if param.Keyword != "" {
			q = q.Where("title LIKE ?", "%"+param.Keyword+"%")
		}
		if param.Category != "" {
			q = q.Where("category = ?", param.Category)
		}
		return q
	}, "is_top DESC, sort_code DESC, created_at DESC", func(e *SysNotice) any { return entToVO(e) })
}

func PublicDetail(c *gin.Context, id string) *NoticeVO {
	if id == "" {
		return nil
	}
	ctx := c.Request.Context()
	var entity SysNotice
	if err := db.DB.WithContext(ctx).Where("id = ? AND status = ?", id, "ENABLED").First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(exception.NewBusinessError("查询通知详情失败: "+err.Error(), 500))
	}
	return entToVO(&entity)
}
