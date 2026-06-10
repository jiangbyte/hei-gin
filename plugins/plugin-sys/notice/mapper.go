package notice

import (
	"hei-gin/sdk/utils"
)

// SysNoticeToNoticeVO 将 notice.SysNotice 映射到 notice.NoticeVO
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

	// *time.Time → *string / string manual conversion
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

// NoticeVOToSysNotice 将 notice.NoticeVO 映射到 notice.SysNotice
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

	// *string → *time.Time manual conversion
	dst.PublishAt = utils.ParseDateTimePtr(src.PublishAt)
	dst.ExpireAt = utils.ParseDateTimePtr(src.ExpireAt)
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}
