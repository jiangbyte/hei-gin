package notice

import "time"

type SysNotice struct {
	ID        string     `gorm:"primaryKey;size:32" json:"id"`
	Title     string     `gorm:"size:255;not null;index" json:"title"`
	Summary   *string    `gorm:"size:500" json:"summary"`
	Content   *string    `gorm:"type:text" json:"content"`
	Cover     *string    `gorm:"size:500" json:"cover"`
	Category  string     `gorm:"size:32;not null;index" json:"category"`
	Type      string     `gorm:"size:32;not null" json:"type"`
	Level     string     `gorm:"size:16;default:NORMAL" json:"level"`
	Status    string     `gorm:"size:16;default:DRAFT" json:"status"`
	SortCode  int        `gorm:"default:0" json:"sort_code"`
	IsTop     string     `gorm:"size:8;default:NO" json:"is_top"`
	Author    *string    `gorm:"size:64" json:"author"`
	PublishAt *time.Time `json:"publish_at"`
	ExpireAt  *time.Time `json:"expire_at"`
	CreatedAt *time.Time `json:"created_at"`
	CreatedBy *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *string    `gorm:"size:32" json:"updated_by"`
}

func (SysNotice) TableName() string { return "sys_notice" }
