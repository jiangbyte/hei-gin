// internal/modules/sys/job/param.go 入参定义。
//
// Author: Charlie

package job

import (
	"github.com/robfig/cron/v3"

	"hei-gin/internal/framework/core/schema"
)

// AddParam 创建任务入参。
//
// Author: Charlie
type AddParam struct {
	HandlerKey  string  `json:"handler_key" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	CronExpr    string  `json:"cron_expr" binding:"required"`
	Params      string  `json:"params"`
	Status      string  `json:"status"`
	Description *string `json:"description"`
}

// EditParam 更新任务入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// StatusParam 启停入参。
//
// Author: Charlie
type StatusParam struct {
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// TriggerParam 手动触发入参。
//
// Author: Charlie
type TriggerParam struct {
	ID     string `json:"id" binding:"required"`
	Params string `json:"params"`
}

// PageParam 任务分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Status  string `form:"status"`
	Keyword string `form:"keyword"`
}

// LogParam 执行日志分页查询。
//
// Author: Charlie
type LogParam struct {
	schema.PageQuery
	JobID      string `form:"job_id"`
	HandlerKey string `form:"handler_key"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// cronParser 与调度器一致：6 段 cron（含秒）。
var cronParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

// cronParse 解析 cron 表达式（6 段含秒）。
func cronParse(spec string) (cron.Schedule, error) {
	return cronParser.Parse(spec)
}
