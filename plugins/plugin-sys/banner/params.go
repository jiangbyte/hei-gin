package banner

import "hei-gin/sdk/utils"

// BannerVO 横幅视图对象
type BannerVO struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Image       string  `json:"image"`
	URL         *string `json:"url"`
	LinkType    string  `json:"link_type"`
	Summary     *string `json:"summary"`
	Description *string `json:"description"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Position    string  `json:"position"`
	SortCode    int     `json:"sort_code"`
	ViewCount   int     `json:"view_count"`
	ClickCount  int     `json:"click_count"`
	CreatedAt   string  `json:"created_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedAt   string  `json:"updated_at"`
	UpdatedBy   *string `json:"updated_by"`
}

// BannerPageParam 横幅分页参数
type BannerPageParam struct {
	Current int `json:"current" form:"current"`
	Size    int `json:"size" form:"size"`
}

func SysBannerToBannerVO(src *SysBanner) *BannerVO {
	if src == nil {
		return nil
	}

	dst := &BannerVO{}
	dst.ID = src.ID
	dst.Title = src.Title
	dst.Image = src.Image
	dst.URL = src.URL
	dst.LinkType = src.LinkType
	dst.Summary = src.Summary
	dst.Description = src.Description
	dst.Category = src.Category
	dst.Type = src.Type
	dst.Position = src.Position
	dst.SortCode = src.SortCode
	dst.ViewCount = src.ViewCount
	dst.ClickCount = src.ClickCount
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

func BannerVOToSysBanner(src *BannerVO) *SysBanner {
	if src == nil {
		return nil
	}

	dst := &SysBanner{}
	dst.ID = src.ID
	dst.Title = src.Title
	dst.Image = src.Image
	dst.URL = src.URL
	dst.LinkType = src.LinkType
	dst.Summary = src.Summary
	dst.Description = src.Description
	dst.Category = src.Category
	dst.Type = src.Type
	dst.Position = src.Position
	dst.SortCode = src.SortCode
	dst.ViewCount = src.ViewCount
	dst.ClickCount = src.ClickCount
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	return dst
}
