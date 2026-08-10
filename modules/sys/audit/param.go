package audit

import "hei-gin/framework/core/schema"

// PageParam 审计日志分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Module    string `form:"module"`
	Action    string `form:"action"`
	AccountID string `form:"account_id"`
}
