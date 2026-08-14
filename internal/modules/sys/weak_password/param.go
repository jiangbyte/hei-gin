package weakpassword

import "hei-gin/internal/framework/core/schema"

// AddParam åˆ›å»ºå¼±å¯†ç å…¥å‚ã€‚
//
// Author: Charlie
type AddParam struct {
	Password string `json:"password" binding:"required"`
}

// EditParam æ›´æ–°å¼±å¯†ç å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PageParam å¼±å¯†ç åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Password string `form:"password"`
}

// ListParam å¼±å¯†ç åˆ—è¡¨æŸ¥è¯¢ã€‚
//
// Author: Charlie
type ListParam struct {
	Password string `form:"password"`
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
