package audit

import "hei-gin/internal/framework/core/schema"

// PageParam å®¡è®¡æ—¥å¿—åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Module    string `form:"module"`
	Action    string `form:"action"`
	AccountID string `form:"account_id"`
}
