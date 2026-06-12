package banner

import "time"

type SysBanner struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Image       string     `gorm:"size:500;not null" json:"image"`
	URL         *string    `gorm:"size:500" json:"url"`
	LinkType    string     `gorm:"size:16;default:URL" json:"link_type"`
	Summary     *string    `gorm:"size:500" json:"summary"`
	Description *string    `gorm:"type:text" json:"description"`
	Category    string     `gorm:"size:32;not null" json:"category"`
	Type        string     `gorm:"size:32;not null" json:"type"`
	Position    string     `gorm:"size:32;not null" json:"position"`
	SortCode    int        `gorm:"default:0" json:"sort_code"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	ClickCount  int        `gorm:"default:0" json:"click_count"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *string    `gorm:"size:32" json:"updated_by"`
}

func (SysBanner) TableName() string { return "sys_banner" }
