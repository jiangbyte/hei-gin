package file

import "hei-gin/internal/framework/core/schema"

// EditParam æ›´æ–°æ–‡ä»¶å…ƒæ•°æ®å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID           string  `json:"id" binding:"required"`
	OriginalName *string `json:"original_name"`
}

// PageParam æ–‡ä»¶åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	OriginalName string `form:"original_name"`
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
