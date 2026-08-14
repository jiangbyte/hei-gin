// internal/modules/biz/cg_test_catalog/param.go 入参定义。
//
// Author: Charlie

package cg_test_catalog

import (
	"hei-gin/internal/framework/core/schema"
)

// AddParam åˆ›å»ºç›®å½•å…¥å‚ã€‚
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

// EditParam æ›´æ–°ç›®å½•å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// PageParam ç›®å½•åˆ†é¡µæŸ¥è¯¢ã€‚
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
