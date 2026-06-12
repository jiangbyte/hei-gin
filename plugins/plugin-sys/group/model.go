package group

import "time"

type SysGroup struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	Code        string     `gorm:"size:32;uniqueIndex;not null" json:"code"`
	Name        string     `gorm:"size:64;not null" json:"name"`
	Category    string     `gorm:"size:32;not null" json:"category"`
	ParentID    *string    `gorm:"size:32;index" json:"parent_id"`
	OrgID       string     `gorm:"size:32;not null" json:"org_id"`
	Description *string    `gorm:"size:500" json:"description"`
	Status      string     `gorm:"size:16;default:ENABLED" json:"status"`
	SortCode    int        `gorm:"default:0" json:"sort_code"`
	Extra       *string    `gorm:"type:text" json:"extra"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *string    `gorm:"size:32" json:"updated_by"`
}

func (SysGroup) TableName() string { return "sys_group" }
