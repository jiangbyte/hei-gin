package org

import (
	"hei-gin/sdk/utils"
)

type OrgVO struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	ParentID    *string `json:"parent_id"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
	SortCode    int     `json:"sort_code"`
	Extra       *string `json:"extra"`
	CreatedAt   string  `json:"created_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedAt   string  `json:"updated_at"`
	UpdatedBy   *string `json:"updated_by"`
}

type OrgPageParam struct {
	Current  int    `json:"current" form:"current"`
	Size     int    `json:"size" form:"size"`
	ParentID string `json:"parent_id" form:"parent_id"`
	Keyword  string `json:"keyword" form:"keyword"`
}

type OrgTreeParam struct {
	Category string `json:"category" form:"category"`
}


func toVO(entity *SysOrg) *OrgVO {
	if entity == nil { return nil }
	return &OrgVO{
		ID:          entity.ID,
		Code:        entity.Code,
		Name:        entity.Name,
		Category:    entity.Category,
		ParentID:    entity.ParentID,
		Description: entity.Description,
		Status:      entity.Status,
		SortCode:    entity.SortCode,
		Extra:       entity.Extra,
		CreatedAt:   utils.FormatDateTimePtr(entity.CreatedAt),
		CreatedBy:   entity.CreatedBy,
		UpdatedAt:   utils.FormatDateTimePtr(entity.UpdatedAt),
		UpdatedBy:   entity.UpdatedBy,
	}
}


func (p *OrgPageParam) GetCurrent() int { return p.Current }
func (p *OrgPageParam) GetSize() int    { return p.Size }
