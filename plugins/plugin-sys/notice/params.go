package notice

import "hei-gin/sdk/utils"

type NoticeVO struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Category  string  `json:"category"`
	Type      string  `json:"type"`
	Summary   *string `json:"summary"`
	Content   *string `json:"content"`
	Cover     *string `json:"cover"`
	Level     string  `json:"level"`
	IsTop     string  `json:"is_top"`
	Status    string  `json:"status"`
	SortCode  int     `json:"sort_code"`
	Author    *string `json:"author"`
	PublishAt *string `json:"publish_at"`
	ExpireAt  *string `json:"expire_at"`
	CreatedAt string  `json:"created_at"`
	CreatedBy *string `json:"created_by"`
	UpdatedAt string  `json:"updated_at"`
	UpdatedBy *string `json:"updated_by"`
}

type NoticeLatestParam struct {
	Size int `json:"size" form:"size"`
}

type NoticePageParam struct {
	Current  int    `json:"current" form:"current"`
	Size     int    `json:"size" form:"size"`
	Keyword  string `json:"keyword" form:"keyword"`
	Category string `json:"category" form:"category"`
	Status   string `json:"status" form:"status"`
}

func SysNoticeToNoticeVO(src *SysNotice) *NoticeVO {
	if src == nil {
		return nil
	}

	dst := &NoticeVO{}
	dst.ID = src.ID
	dst.Title = src.Title
	dst.Category = src.Category
	dst.Type = src.Type
	dst.Summary = src.Summary
	dst.Content = src.Content
	dst.Cover = src.Cover
	dst.Level = src.Level
	dst.IsTop = src.IsTop
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.Author = src.Author
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	s := utils.FormatDateTimePtr(src.PublishAt)
	if s != "" {
		dst.PublishAt = &s
	}
	s = utils.FormatDateTimePtr(src.ExpireAt)
	if s != "" {
		dst.ExpireAt = &s
	}
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

func NoticeVOToSysNotice(src *NoticeVO) *SysNotice {
	if src == nil {
		return nil
	}

	dst := &SysNotice{}
	dst.ID = src.ID
	dst.Title = src.Title
	dst.Category = src.Category
	dst.Type = src.Type
	dst.Summary = src.Summary
	dst.Content = src.Content
	dst.Cover = src.Cover
	dst.Level = src.Level
	dst.IsTop = src.IsTop
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.Author = src.Author
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.PublishAt = utils.ParseDateTimePtr(src.PublishAt)
	dst.ExpireAt = utils.ParseDateTimePtr(src.ExpireAt)
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
	return dst
}
