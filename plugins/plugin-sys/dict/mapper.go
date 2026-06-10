package dict

import "hei-gin/sdk/utils"

// SysDictToDictVO 将 dict.SysDict 映射到 dict.DictVO
func SysDictToDictVO(src *SysDict) *DictVO {
	if src == nil {
		return nil
	}

	dst := &DictVO{}

	dst.ID = src.ID
	dst.Code = src.Code
	dst.Label = src.Label
	dst.Value = src.Value
	dst.Color = src.Color
	dst.Category = src.Category
	dst.ParentID = src.ParentID
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// *time.Time → string manual conversion
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// DictVOToSysDict 将 dict.DictVO 映射到 dict.SysDict
func DictVOToSysDict(src *DictVO) *SysDict {
	if src == nil {
		return nil
	}

	dst := &SysDict{}

	dst.ID = src.ID
	dst.Code = src.Code
	dst.Label = src.Label
	dst.Value = src.Value
	dst.Color = src.Color
	dst.Category = src.Category
	dst.ParentID = src.ParentID
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// string → *time.Time manual conversion
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}
