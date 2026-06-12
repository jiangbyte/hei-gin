package config

import "time"

type SysConfig struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	ConfigKey   *string    `gorm:"size:255" json:"config_key"`
	ConfigValue *string    `gorm:"type:text" json:"config_value"`
	Category    *string    `gorm:"size:255;index" json:"category"`
	Remark      *string    `gorm:"size:500" json:"remark"`
	SortCode    int        `gorm:"default:0" json:"sort_code"`
	Extra       *string    `gorm:"type:text" json:"extra"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *string    `gorm:"size:32" json:"updated_by"`
}

func (SysConfig) TableName() string { return "sys_config" }
