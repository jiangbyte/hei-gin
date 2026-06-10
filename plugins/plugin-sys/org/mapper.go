package org

import "hei-gin/sdk/utils"

// SysOrgToOrgVO 将 org.SysOrg 映射到 org.OrgVO
func SysOrgToOrgVO(src *SysOrg) *OrgVO {
	if src == nil {
		return nil
	}

	dst := &OrgVO{}

	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.ParentID = src.ParentID
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

// OrgVOToSysOrg 将 org.OrgVO 映射到 org.SysOrg
func OrgVOToSysOrg(src *OrgVO) *SysOrg {
	if src == nil {
		return nil
	}

	dst := &SysOrg{}

	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.ParentID = src.ParentID
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
