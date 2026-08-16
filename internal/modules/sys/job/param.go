// internal/modules/sys/job/param.go 入参定义（对齐 hei-boot / hei-fastapi）。
//
// Author: Charlie

package job

import "hei-gin/internal/framework/core/schema"

// AddParam 创建任务入参。
//
// Author: Charlie
type AddParam struct {
	JobName       string         `json:"job_name" binding:"required"`
	ExecuteClass  string         `json:"execute_class" binding:"required"`
	ExecuteType   string         `json:"execute_type" binding:"required"`
	TriggerConfig string         `json:"trigger_config" binding:"required"`
	ExecuteParam  map[string]any `json:"execute_param"`
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
	JobName     string `form:"job_name"`
	ExecuteType string `form:"execute_type"`
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
