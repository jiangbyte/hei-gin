package config

import (
	"hei-gin/sdk/utils"
)

type ConfigVO struct {
	ID          string  `json:"id"`
	ConfigKey   *string `json:"config_key"`
	ConfigValue *string `json:"config_value"`
	Category    *string `json:"category"`
	Remark      *string `json:"remark"`
	SortCode    int     `json:"sort_code"`
	Extra       *string `json:"extra"`
	CreatedAt   string  `json:"created_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedAt   string  `json:"updated_at"`
	UpdatedBy   *string `json:"updated_by"`
}

type ConfigPageParam struct {
	Current  int    `json:"current" form:"current"`
	Size     int    `json:"size" form:"size"`
	Category string `json:"category" form:"category"`
	Keyword  string `json:"keyword" form:"keyword"`
}

type ConfigListParam struct {
	Category string `json:"category" form:"category"`
}
type ConfigBatchEditItem struct {
	ID          string  `json:"id"`
	ConfigKey   *string `json:"config_key"`
	ConfigValue *string `json:"config_value"`
	Remark      *string `json:"remark"`
	SortCode    int     `json:"sort_code"`
}
type ConfigBatchEditParam struct {
	Configs []ConfigBatchEditItem `json:"configs"`
}
type ConfigCategoryEditParam struct {
	Category    string  `json:"category"`
	ConfigKey   *string `json:"config_key"`
	ConfigValue *string `json:"config_value"`
	Remark      *string `json:"remark"`
}


func toVO(e *SysConfig) *ConfigVO {
	v := &ConfigVO{
		ID: e.ID, SortCode: e.SortCode,
		CreatedAt: utils.FormatDateTime(*e.CreatedAt),
		UpdatedAt: utils.FormatDateTime(*e.UpdatedAt),
	}
	if e.ConfigKey != nil { v.ConfigKey = e.ConfigKey }
	if e.ConfigValue != nil { v.ConfigValue = e.ConfigValue }
	if e.Category != nil { v.Category = e.Category }
	if e.Remark != nil { v.Remark = e.Remark }
	if e.Extra != nil { v.Extra = e.Extra }
	if e.CreatedBy != nil { v.CreatedBy = e.CreatedBy }
	if e.UpdatedBy != nil { v.UpdatedBy = e.UpdatedBy }
	return v
}

func toVOList(records []SysConfig) []ConfigVO {
	result := make([]ConfigVO, len(records))
	for i, r := range records {
		v := toVO(&r)
		if v != nil {
			result[i] = *v
		}
	}
	return result
}

func (p *ConfigPageParam) GetCurrent() int { return p.Current }
func (p *ConfigPageParam) GetSize() int    { return p.Size }
