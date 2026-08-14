// internal/modules/iam/resource/model.go 数据模型。
//
// Author: Charlie

package resource

import (
	"time"

	"gorm.io/datatypes"
)

// Resource 映射 sys_resource 菜单/按钮等资源。
//
// Author: Charlie
type Resource struct {
	ID           string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	ParentID     *string        `gorm:"column:parent_id;size:64" json:"parent_id"`
	Code         string         `gorm:"column:code;size:64;not null" json:"code"`
	Name         string         `gorm:"column:name;size:64;not null" json:"name"`
	ResourceType string         `gorm:"column:resource_type;size:32;not null" json:"resource_type"`
	ModuleID     *string        `gorm:"column:module_id;size:64" json:"module_id"`
	Path         *string        `gorm:"column:path;size:255" json:"path"`
	Component    *string        `gorm:"column:component;size:255" json:"component"`
	Redirect     *string        `gorm:"column:redirect;size:255" json:"redirect"`
	Icon         *string        `gorm:"column:icon;size:255" json:"icon"`
	Color        *string        `gorm:"column:color;size:32" json:"color"`
	Href         *string        `gorm:"column:href;size:255" json:"href"`
	Sort         int            `gorm:"column:sort;not null;default:99" json:"sort"`
	IsVisible    bool           `gorm:"column:is_visible;not null;default:true" json:"is_visible"`
	IsCache      bool           `gorm:"column:is_cache;not null;default:false" json:"is_cache"`
	IsAffix      bool           `gorm:"column:is_affix;not null;default:false" json:"is_affix"`
	Status       string         `gorm:"column:status;size:32;not null" json:"status"`
	Description  *string        `gorm:"column:description" json:"description"`
	Layout       *string        `gorm:"column:layout;size:255" json:"layout"`
	Extra        datatypes.JSON `gorm:"column:extra;type:jsonb" json:"extra"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy    *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy    *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (Resource) TableName() string { return "sys_resource" }

// ResourceModule 映射 sys_resource_module 资源模块。
//
// Author: Charlie
type ResourceModule struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Name        string         `gorm:"column:name;size:64;not null" json:"name"`
	Code        string         `gorm:"column:code;size:64;uniqueIndex;not null" json:"code"`
	Client      string         `gorm:"column:client;size:32;not null" json:"client"`
	Icon        *string        `gorm:"column:icon;size:255" json:"icon"`
	Color       *string        `gorm:"column:color;size:32" json:"color"`
	Sort        int            `gorm:"column:sort;not null;default:99" json:"sort"`
	Status      string         `gorm:"column:status;size:32;not null" json:"status"`
	Description *string        `gorm:"column:description" json:"description"`
	Extra       datatypes.JSON `gorm:"column:extra;type:jsonb" json:"extra"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (ResourceModule) TableName() string { return "sys_resource_module" }

// 资源类型常量。
const ResourceTypeButton = "BUTTON"
