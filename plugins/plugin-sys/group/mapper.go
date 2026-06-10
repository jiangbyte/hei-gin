package group

import "hei-gin/sdk/utils"

// SysGroupToGroupVO 将 group.SysGroup 映射到 group.GroupVO
func SysGroupToGroupVO(src *SysGroup) *GroupVO {
	if src == nil {
		return nil
	}

	dst := &GroupVO{}

	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.ParentID = src.ParentID
	dst.OrgID = src.OrgID
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

// GroupVOToSysGroup 将 group.GroupVO 映射到 group.SysGroup
func GroupVOToSysGroup(src *GroupVO) *SysGroup {
	if src == nil {
		return nil
	}

	dst := &SysGroup{}

	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.ParentID = src.ParentID
	dst.OrgID = src.OrgID
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
