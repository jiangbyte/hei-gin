// internal/modules/biz/cg_test_activity/param.go 入参定义。
//
// Author: Charlie

package cg_test_activity

import (
	"time"

	"hei-gin/internal/framework/core/schema"
)

// AddParam åˆ›å»ºæ´»åŠ¨å…¥å‚ã€‚
//
// Author: Charlie
type AddParam struct {
	Code            string         `json:"code" binding:"required"`
	Name            string         `json:"name" binding:"required"`
	Category        *string        `json:"category"`
	Type            string         `json:"type" binding:"required"`
	Status          string         `json:"status" binding:"required"`
	CoverURL        *string        `json:"cover_url"`
	Description     *string        `json:"description"`
	StartAt         time.Time      `json:"start_at" binding:"required"`
	EndAt           *time.Time     `json:"end_at"`
	MaxParticipants int            `json:"max_participants"`
	Price           float64        `json:"price"`
	IsPublic        bool           `json:"is_public"`
	NeedApproval    bool           `json:"need_approval"`
	RuleConfig      map[string]any `json:"rule_config"`
	Extra           map[string]any `json:"extra"`
}

// EditParam æ›´æ–°æ´»åŠ¨å…¥å‚ã€‚
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

// PageParam æ´»åŠ¨åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Code     string `form:"code"`
	Name     string `form:"name"`
	Category string `form:"category"`
	Type     string `form:"type"`
	Status   string `form:"status"`
}
