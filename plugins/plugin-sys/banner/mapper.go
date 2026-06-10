package banner

import "hei-gin/sdk/utils"

// SysBannerToBannerVO 将 banner.SysBanner 映射到 banner.BannerVO
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

	// *time.Time → string manual conversion
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// BannerVOToSysBanner 将 banner.BannerVO 映射到 banner.SysBanner
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
