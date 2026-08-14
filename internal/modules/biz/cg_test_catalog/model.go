// Package cg_test_catalog 为代码生成演示的目录业务模块。
package cg_test_catalog

import (
	"time"

	"gorm.io/datatypes"
)

// Catalog 演示目录实体，对应表 cg_test_catalog。
//
// Author: Charlie
type Catalog struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	ParentID    *string        `gorm:"column:parent_id;size:64" json:"parent_id"`
	Code        string         `gorm:"column:code;size:64" json:"code"`
	Name        string         `gorm:"column:name;size:120" json:"name"`
	Category    *string        `gorm:"column:category;size:32" json:"category"`
	Status      string         `gorm:"column:status;size:32" json:"status"`
	Sort        int            `gorm:"column:sort" json:"sort"`
	IsVisible   bool           `gorm:"column:is_visible" json:"is_visible"`
	Icon        *string        `gorm:"column:icon;size:128" json:"icon"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	Extra       datatypes.JSON `gorm:"column:extra;type:jsonb" json:"extra"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
	OwnerDeptID *string        `gorm:"column:owner_dept_id;size:64" json:"owner_dept_id"`
}

// TableName 返回 Catalog 对应的数据库表名。
func (Catalog) TableName() string { return "cg_test_catalog" }
