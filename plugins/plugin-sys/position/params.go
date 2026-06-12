package position

import "hei-gin/sdk/utils"

type PositionVO struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	OrgID       *string  `json:"org_id"`
	GroupID     *string  `json:"group_id"`
	Description *string  `json:"description"`
	Status      string   `json:"status"`
	SortCode    int      `json:"sort_code"`
	OrgNames    []string `json:"org_names"`
	GroupNames  []string `json:"group_names"`
	Extra       *string  `json:"extra"`
	CreatedAt   string   `json:"created_at"`
	CreatedBy   *string  `json:"created_by"`
	UpdatedAt   string   `json:"updated_at"`
	UpdatedBy   *string  `json:"updated_by"`
}

type PositionPageParam struct {
	Current  int    `json:"current" form:"current"`
	Size     int    `json:"size" form:"size"`
	Keyword  string `json:"keyword" form:"keyword"`
	Category string `json:"category" form:"category"`
	OrgID    string `json:"org_id" form:"org_id"`
}

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
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

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
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
	return dst
}
