// internal/modules/biz/cg_test_knowledge_category/param.go 入参定义。
//
// Author: Charlie

package cg_test_knowledge_category

import (
	"time"

	"hei-gin/internal/framework/core/schema"
)

// AddParam 创建分类入参。
//
// Author: Charlie
type AddParam struct {
	ParentID    *string        `json:"parent_id"`
	Code        string         `json:"code" binding:"required"`
	Name        string         `json:"name" binding:"required"`
	Status      string         `json:"status" binding:"required"`
	Sort        int            `json:"sort"`
	IsVisible   bool           `json:"is_visible"`
	Description *string        `json:"description"`
	Extra       map[string]any `json:"extra"`
}

// EditParam 更新分类入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// PageParam 分类分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Code     string `form:"code"`
	Name     string `form:"name"`
	Status   string `form:"status"`
	ParentID string `form:"parent_id"`
}

// DocAddParam 创建文档入参。
//
// Author: Charlie
type DocAddParam struct {
	CategoryID  string         `json:"category_id" binding:"required"`
	Code        string         `json:"code" binding:"required"`
	Title       string         `json:"title" binding:"required"`
	Type        string         `json:"type" binding:"required"`
	Status      string         `json:"status" binding:"required"`
	Summary     *string        `json:"summary"`
	Content     *string        `json:"content"`
	Author      *string        `json:"author"`
	PublishedAt *time.Time     `json:"published_at"`
	ViewCount   int            `json:"view_count"`
	Sort        int            `json:"sort"`
	IsTop       bool           `json:"is_top"`
	Settings    map[string]any `json:"settings"`
	Extra       map[string]any `json:"extra"`
}

// DocEditParam 更新文档入参。
//
// Author: Charlie
type DocEditParam struct {
	ID string `json:"id" binding:"required"`
	DocAddParam
}

// DocPageParam 文档分页查询。
//
// Author: Charlie
type DocPageParam struct {
	schema.PageQuery
	CategoryID string `form:"category_id"`
	Code       string `form:"code"`
	Title      string `form:"title"`
	Status     string `form:"status"`
}
