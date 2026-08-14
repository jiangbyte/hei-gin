package dict

import "hei-gin/internal/framework/core/schema"

// AddParam åˆ›å»ºå­—å…¸å…¥å‚ã€‚
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

// EditParam æ›´æ–°å­—å…¸å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam å­—å…¸åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Code     string `form:"code"`
	Category string `form:"category"`
	Status   string `form:"status"`
}

// TreeParam å­—å…¸æ ‘æŸ¥è¯¢ã€‚
//
// Author: Charlie
type TreeParam struct {
	Code     string `form:"code"`
	Category string `form:"category"`
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
