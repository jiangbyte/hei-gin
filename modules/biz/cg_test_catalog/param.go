package cg_test_catalog

import (
	"hei-gin/framework/core/schema"
)

// AddParam 创建目录入参。
//
// Author: Charlie
type AddParam struct {
	ParentID    *string        `json:"parent_id"`
	Code        string         `json:"code" binding:"required"`
	Name        string         `json:"name" binding:"required"`
	Category    *string        `json:"category"`
	Status      string         `json:"status" binding:"required"`
	Sort        int            `json:"sort"`
	IsVisible   bool           `json:"is_visible"`
	Icon        *string        `json:"icon"`
	Description *string        `json:"description"`
	Extra       map[string]any `json:"extra"`
}

// EditParam 更新目录入参。
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

// PageParam 目录分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Code     string `form:"code"`
	Name     string `form:"name"`
	Category string `form:"category"`
	Status   string `form:"status"`
	ParentID string `form:"parent_id"`
}
