// Package banner 提供系统 Banner 广告位管理。
//
// Author: Charlie
package banner

import (
	"time"

	"gorm.io/datatypes"
)

// Banner 横幅实体，对应表 sys_banner。
//
// Author: Charlie
type Banner struct {
	ID                 string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Title              string         `gorm:"column:title;size:255;not null" json:"title"`
	Image              string         `gorm:"column:image;size:500;not null" json:"image"`
	ImageURL           *string        `gorm:"-" json:"image_url"`
	URL                *string        `gorm:"column:url;size:500" json:"url"`
	LinkType           string         `gorm:"column:link_type;size:16;not null" json:"link_type"`
	Summary            *string        `gorm:"column:summary;size:500" json:"summary"`
	Description        *string        `gorm:"column:description" json:"description"`
	Category           string         `gorm:"column:category;size:32;not null" json:"category"`
	Type               string         `gorm:"column:type;size:32;not null" json:"type"`
	Position           string         `gorm:"column:position;size:32;not null" json:"position"`
	TargetAccountTypes datatypes.JSON `gorm:"column:target_account_types;type:json" json:"target_account_types"`
	Sort               int            `gorm:"column:sort;not null;default:0" json:"sort"`
	InteractionCount   int64          `gorm:"column:interaction_count;not null;default:0" json:"interaction_count"`
	Status             string         `gorm:"column:status;size:32;not null" json:"status"`
	StartAt            *time.Time     `gorm:"column:start_at" json:"start_at"`
	EndAt              *time.Time     `gorm:"column:end_at" json:"end_at"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy          *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy          *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Banner 对应的数据库表名。
func (Banner) TableName() string { return "sys_banner" }
