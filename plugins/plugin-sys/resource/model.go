package resource

import "time"

type SysResource struct {
	ID            string     `gorm:"primaryKey;size:32" json:"id"`
	Code          string     `gorm:"size:32;uniqueIndex;not null" json:"code"`
	Name          string     `gorm:"size:64;not null" json:"name"`
	Category      string     `gorm:"size:16;not null" json:"category"`
	Type          string     `gorm:"size:16;not null" json:"type"`
	Description   *string    `gorm:"size:500" json:"description"`
	ParentID      *string    `gorm:"size:32;index" json:"parent_id"`
	RoutePath     *string    `gorm:"size:255" json:"route_path"`
	ComponentPath *string    `gorm:"size:255" json:"component_path"`
	RedirectPath  *string    `gorm:"size:255" json:"redirect_path"`
	Icon          *string    `gorm:"size:64" json:"icon"`
	Color         *string    `gorm:"size:32" json:"color"`
	IsVisible     string     `gorm:"size:8;default:YES" json:"is_visible"`
	IsCache       string     `gorm:"size:8;default:NO" json:"is_cache"`
	IsAffix       string     `gorm:"size:8;default:NO" json:"is_affix"`
	IsBreadcrumb  string     `gorm:"size:8;default:YES" json:"is_breadcrumb"`
	ExternalURL   *string    `gorm:"size:500" json:"external_url"`
	Extra         *string    `gorm:"type:text" json:"extra"`
	Status        string     `gorm:"size:16;default:ENABLED" json:"status"`
	SortCode      int        `gorm:"default:0" json:"sort_code"`
	CreatedAt     *time.Time `json:"created_at"`
	CreatedBy     *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedBy     *string    `gorm:"size:32" json:"updated_by"`
}

func (SysResource) TableName() string { return "sys_resource" }

type SysModule struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	Code        string     `gorm:"size:32;uniqueIndex;not null" json:"code"`
	Name        string     `gorm:"size:64;not null" json:"name"`
	Category    string     `gorm:"size:32;not null" json:"category"`
	Icon        *string    `gorm:"size:64" json:"icon"`
	Color       *string    `gorm:"size:32" json:"color"`
	Description *string    `gorm:"size:500" json:"description"`
	IsVisible   string     `gorm:"size:8;default:YES" json:"is_visible"`
	Status      string     `gorm:"size:16;default:ENABLED" json:"status"`
	SortCode    int        `gorm:"default:0" json:"sort_code"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *string    `gorm:"size:32" json:"updated_by"`
}

func (SysModule) TableName() string { return "sys_module" }
