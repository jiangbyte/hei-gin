// Package config 提供系统运行时配置项管理。
//
// Author: Charlie
package config

import (
	"time"

	"gorm.io/datatypes"
)

// Config 系统配置实体，对应表 sys_config。
//
// Author: Charlie
type Config struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	ConfigKey   string         `gorm:"column:config_key;size:255;uniqueIndex;not null" json:"config_key"`
	ConfigValue *string        `gorm:"column:config_value" json:"config_value"`
	Category    *string        `gorm:"column:category;size:255" json:"category"`
	Remark      *string        `gorm:"column:remark;size:255" json:"remark"`
	SortCode    int            `gorm:"column:sort_code;not null;default:0" json:"sort_code"`
	ValueType   string         `gorm:"column:value_type;size:32;not null;default:STRING" json:"value_type"`
	Label       *string        `gorm:"column:label;size:128" json:"label"`
	Scope       *string        `gorm:"column:scope;size:32" json:"scope"`
	Scene       *string        `gorm:"column:scene;size:64" json:"scene"`
	IsBuiltin   bool           `gorm:"column:is_builtin;not null;default:false" json:"is_builtin"`
	ExtJSON     datatypes.JSON `gorm:"column:ext_json;type:json" json:"ext_json"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Config 对应的数据库表名。
func (Config) TableName() string { return "sys_config" }
