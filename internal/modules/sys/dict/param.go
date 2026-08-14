// internal/modules/sys/dict/param.go 入参定义。
//
// Author: Charlie

package dict

import "hei-gin/internal/framework/core/schema"

// AddParam 创建字典入参。
//
// Author: Charlie
type AddParam struct {
	Code     string  `json:"code" binding:"required"`
	Label    *string `json:"label"`
	Value    *string `json:"value"`
	Color    *string `json:"color"`
	Category *string `json:"category"`
	ParentID *string `json:"parent_id"`
	Status   string  `json:"status"`
	Sort     int     `json:"sort"`
}

// EditParam 更新字典入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam 字典分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Code     string `form:"code"`
	Category string `form:"category"`
	Status   string `form:"status"`
}

// TreeParam 字典树查询。
//
// Author: Charlie
type TreeParam struct {
	Code     string `form:"code"`
	Category string `form:"category"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
