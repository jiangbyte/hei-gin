package position

import (
	"hei-gin/sdk/utils"
)

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


func toVO(entity *SysPosition) *PositionVO {
	if entity == nil { return nil }
	return &PositionVO{
		ID: entity.ID, Code: entity.Code, Name: entity.Name, Category: entity.Category,
		Description: entity.Description, Status: entity.Status, SortCode: entity.SortCode,
		Extra: entity.Extra, CreatedAt: utils.FormatDateTimePtr(entity.CreatedAt),
		CreatedBy: entity.CreatedBy, UpdatedAt: utils.FormatDateTimePtr(entity.UpdatedAt),
		UpdatedBy: entity.UpdatedBy,
	}
}

func (p *PositionPageParam) GetCurrent() int { return p.Current }
func (p *PositionPageParam) GetSize() int    { return p.Size }
