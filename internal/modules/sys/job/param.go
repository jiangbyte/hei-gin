// internal/modules/sys/job/param.go 入参定义（对齐 hei-boot / hei-admin）。
//
// Author: Charlie

package job

import "hei-gin/internal/framework/core/schema"

// AddParam 创建任务入参。
//
// Author: Charlie
type AddParam struct {
	Name          string         `json:"name" binding:"required"`
	Handler       string         `json:"handler" binding:"required"`
	TriggerType   string         `json:"trigger_type" binding:"required"`
	TriggerConfig string         `json:"trigger_config" binding:"required"`
	Params        map[string]any `json:"params"`
	Description   *string        `json:"description"`
	Sort          int            `json:"sort"`
	Enabled       *bool          `json:"enabled"`
}

// EditParam 更新任务入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// EnabledParam 启停入参。
//
// Author: Charlie
type EnabledParam struct {
	ID      string `json:"id" binding:"required"`
	Enabled bool   `json:"enabled"`
}

// RunParam 立即执行入参。
//
// Author: Charlie
type RunParam struct {
	ID string `json:"id" binding:"required"`
}

// PageParam 任务分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Name        string `form:"name"`
	TriggerType string `form:"trigger_type"`
	Enabled     *bool  `form:"enabled"`
}

// LogParam 执行日志分页查询。
//
// Author: Charlie
type LogParam struct {
	schema.PageQuery
	JobID   string `form:"job_id"`
	Success *bool  `form:"success"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
