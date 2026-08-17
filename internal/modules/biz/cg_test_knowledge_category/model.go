// Package cg_test_knowledge_category 为代码生成演示的知识分类业务模块。
//
// Author: Charlie
package cg_test_knowledge_category

import (
	"time"

	"gorm.io/datatypes"
)

// Category 演示知识分类实体，对应表 cg_test_knowledge_category。
//
// Author: Charlie
type Category struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	ParentID    *string        `gorm:"column:parent_id;size:64" json:"parent_id"`
	Code        string         `gorm:"column:code;size:64" json:"code"`
	Name        string         `gorm:"column:name;size:120" json:"name"`
	Status      string         `gorm:"column:status;size:32" json:"status"`
	Sort        int            `gorm:"column:sort" json:"sort"`
	IsVisible   bool           `gorm:"column:is_visible" json:"is_visible"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	Extra       datatypes.JSON `gorm:"column:extra;type:json" json:"extra"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
	OwnerDeptID *string        `gorm:"column:owner_dept_id;size:64" json:"owner_dept_id"`
}

// TableName 返回 Category 对应的数据库表名。
func (Category) TableName() string { return "cg_test_knowledge_category" }

// Doc 演示知识文档实体，对应表 cg_test_knowledge_doc。
//
// Author: Charlie
type Doc struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	CategoryID  string         `gorm:"column:category_id;size:64" json:"category_id"`
	Code        string         `gorm:"column:code;size:64" json:"code"`
	Title       string         `gorm:"column:title;size:160" json:"title"`
	Type        string         `gorm:"column:type;size:32" json:"type"`
	Status      string         `gorm:"column:status;size:32" json:"status"`
	Summary     *string        `gorm:"column:summary;size:512" json:"summary"`
	Content     *string        `gorm:"column:content;type:text" json:"content"`
	Author      *string        `gorm:"column:author;size:64" json:"author"`
	PublishedAt *time.Time     `gorm:"column:published_at" json:"published_at"`
	ViewCount   int            `gorm:"column:view_count" json:"view_count"`
	Sort        int            `gorm:"column:sort" json:"sort"`
	IsTop       bool           `gorm:"column:is_top" json:"is_top"`
	Settings    datatypes.JSON `gorm:"column:settings;type:json" json:"settings"`
	Extra       datatypes.JSON `gorm:"column:extra;type:json" json:"extra"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Doc 对应的数据库表名。
func (Doc) TableName() string { return "cg_test_knowledge_doc" }
