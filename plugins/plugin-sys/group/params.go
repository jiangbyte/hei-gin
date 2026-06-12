package group

import "hei-gin/sdk/utils"

// GroupVO 用户组视图对象
type GroupVO struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	ParentID    *string  `json:"parent_id"`
	OrgID       string   `json:"org_id"`
	Description *string  `json:"description"`
	Status      string   `json:"status"`
	SortCode    int      `json:"sort_code"`
	OrgNames    []string `json:"org_names"`
	Extra       *string  `json:"extra"`
	CreatedAt   string   `json:"created_at"`
	CreatedBy   *string  `json:"created_by"`
	UpdatedAt   string   `json:"updated_at"`
	UpdatedBy   *string  `json:"updated_by"`
}

// GroupPageParam 用户组分页参数
type GroupPageParam struct {
	Current  int    `json:"current" form:"current"`
	Size     int    `json:"size" form:"size"`
	Keyword  string `json:"keyword" form:"keyword"`
	Category string `json:"category" form:"category"`
	OrgID    string `json:"org_id" form:"org_id"`
}

// GroupTreeParam 用户组树查询参数
type GroupTreeParam struct {
	Category string `json:"category" form:"category"`
	OrgID    string `json:"org_id" form:"org_id"`
}

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
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

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
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
	return dst
}
