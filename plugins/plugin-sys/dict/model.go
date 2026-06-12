package dict

import "time"

type SysDict struct {
	ID        string     `gorm:"primaryKey;size:32" json:"id"`
	Code      string     `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Label     *string    `gorm:"size:255" json:"label"`
	Value     *string    `gorm:"size:255" json:"value"`
	Color     *string    `gorm:"size:32" json:"color"`
	Category  *string    `gorm:"size:64;index" json:"category"`
	ParentID  *string    `gorm:"size:32;index" json:"parent_id"`
	Status    string     `gorm:"size:16;default:ENABLED" json:"status"`
	SortCode  int        `gorm:"default:0" json:"sort_code"`
	CreatedAt *time.Time `json:"created_at"`
	CreatedBy *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *string    `gorm:"size:32" json:"updated_by"`
}

func (SysDict) TableName() string { return "sys_dict" }
