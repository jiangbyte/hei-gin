// internal/modules/sys/banner/param.go 入参定义。
//
// Author: Charlie

package banner

import (
	"time"

	"gorm.io/datatypes"

	"hei-gin/internal/framework/core/schema"
)

// AddParam 创建 Banner 入参。
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

// EditParam 更新 Banner 入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam Banner 分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Title             string `form:"title"`
	Position          string `form:"position"`
	Status            string `form:"status"`
	TargetAccountType string `form:"target_account_type"`
}

// ListParam Banner 列表查询。
//
// Author: Charlie
type ListParam struct {
	Position string `form:"position"`
	Category string `form:"category"`
	Type     string `form:"type"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// InteractionParam Banner 互动上报入参。
//
// Author: Charlie
type InteractionParam struct {
	ID string `json:"id" binding:"required"`
}

// PortalListParam 门户端 Banner 列表查询。
//
// Author: Charlie
type PortalListParam struct {
	Position string `form:"position"`
	Category string `form:"category"`
	Type     string `form:"type"`
}
