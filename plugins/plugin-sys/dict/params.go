package dict

import "hei-gin/sdk/utils"

// DictVO 字典视图对象
type DictVO struct {
	ID        string  `json:"id"`
	Code      string  `json:"code"`
	Label     *string `json:"label"`
	Value     *string `json:"value"`
	Color     *string `json:"color"`
	Category  *string `json:"category"`
	ParentID  *string `json:"parent_id"`
	Status    string  `json:"status"`
	SortCode  int     `json:"sort_code"`
	CreatedAt string  `json:"created_at"`
	CreatedBy *string `json:"created_by"`
	UpdatedAt string  `json:"updated_at"`
	UpdatedBy *string `json:"updated_by"`
}

// DictPageParam 字典分页参数
type DictPageParam struct {
	Current   int    `json:"current" form:"current"`
	Size      int    `json:"size" form:"size"`
	Keyword   string `json:"keyword" form:"keyword"`
	Category  string `json:"category" form:"category"`
	ParentID  string `json:"parent_id" form:"parent_id"`
	DictGroup string `json:"dict_group" form:"dict_group"`
}

// DictTreeParam 字典树查询参数
type DictTreeParam struct {
	Category  string `json:"category" form:"category"`
	DictGroup string `json:"dict_group" form:"dict_group"`
}

// DictOptionsParam 字典选项查询参数
type DictOptionsParam struct {
	Category string `json:"category" form:"category"`
	ParentID string `json:"parent_id" form:"parent_id"`
}

// DictListParam 字典列表查询参数
type DictListParam struct {
	Category string `json:"category" form:"category"`
	Keyword  string `json:"keyword" form:"keyword"`
}

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
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

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
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
	return dst
}
