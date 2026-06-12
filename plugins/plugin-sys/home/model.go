package home

import "time"

type SysQuickAction struct {
	ID         string     `gorm:"primaryKey;size:32" json:"id"`
	UserID     string     `gorm:"size:32;uniqueIndex:idx_user_resource;not null" json:"user_id"`
	ResourceID string     `gorm:"size:32;uniqueIndex:idx_user_resource;not null" json:"resource_id"`
	SortCode   int        `gorm:"default:0" json:"sort_code"`
	CreatedAt  *time.Time `json:"created_at"`
	CreatedBy  *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt  *time.Time `json:"updated_at"`
	UpdatedBy  *string    `gorm:"size:32" json:"updated_by"`
}

func (SysQuickAction) TableName() string { return "sys_quick_action" }
