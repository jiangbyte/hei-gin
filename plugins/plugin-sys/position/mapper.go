package position

import "hei-gin/sdk/utils"

// SysPositionToPositionVO 将 position.SysPosition 映射到 position.PositionVO
func SysPositionToPositionVO(src *SysPosition) *PositionVO {
	if src == nil {
		return nil
	}

	dst := &PositionVO{}

	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.OrgID = src.OrgID
	dst.GroupID = src.GroupID
	dst.Description = src.Description
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.Extra = src.Extra
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// *time.Time → string manual conversion
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// PositionVOToSysPosition 将 position.PositionVO 映射到 position.SysPosition
func PositionVOToSysPosition(src *PositionVO) *SysPosition {
	if src == nil {
		return nil
	}

	dst := &SysPosition{}

	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.OrgID = src.OrgID
	dst.GroupID = src.GroupID
	dst.Description = src.Description
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.Extra = src.Extra
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// string → *time.Time manual conversion
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}
