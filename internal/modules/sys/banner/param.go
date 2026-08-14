package banner

import (
	"time"

	"gorm.io/datatypes"

	"hei-gin/internal/framework/core/schema"
)

// AddParam åˆ›å»º Banner å…¥å‚ã€‚
//
// Author: Charlie
type AddParam struct {
	Title              string         `json:"title" binding:"required"`
	Image              string         `json:"image" binding:"required"`
	URL                *string        `json:"url"`
	LinkType           string         `json:"link_type"`
	Summary            *string        `json:"summary"`
	Description        *string        `json:"description"`
	Category           string         `json:"category" binding:"required"`
	Type               string         `json:"type" binding:"required"`
	Position           string         `json:"position" binding:"required"`
	TargetAccountTypes datatypes.JSON `json:"target_account_types"`
	Sort               int            `json:"sort"`
	Status             string         `json:"status"`
	StartAt            *time.Time     `json:"start_at"`
	EndAt              *time.Time     `json:"end_at"`
}

// EditParam æ›´æ–° Banner å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam Banner åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Title    string `form:"title"`
	Position string `form:"position"`
	Status   string `form:"status"`
}

// ListParam Banner åˆ—è¡¨æŸ¥è¯¢ã€‚
//
// Author: Charlie
type ListParam struct {
	Position string `form:"position"`
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// InteractionParam Banner äº’åŠ¨ä¸ŠæŠ¥å…¥å‚ã€‚
//
// Author: Charlie
type InteractionParam struct {
	ID string `json:"id" binding:"required"`
}

// PortalListParam é—¨æˆ·ç«¯ Banner åˆ—è¡¨æŸ¥è¯¢ã€‚
//
// Author: Charlie
type PortalListParam struct {
	Position string `form:"position"`
	Category string `form:"category"`
	Type     string `form:"type"`
}
