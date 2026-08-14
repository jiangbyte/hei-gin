// internal/modules/message/feedback/param.go 入参定义。
//
// Author: Charlie

package feedback

import "hei-gin/internal/framework/core/schema"

// CreateParam 提交反馈入参。
//
// Author: Charlie
type CreateParam struct {
	Title             string   `json:"title" binding:"required"`
	Content           string   `json:"content" binding:"required"`
	Category          string   `json:"category" binding:"required"`
	Contact           *string  `json:"contact"`
	AttachObjectNames []string `json:"attach_object_names"`
}

// UpdateParam 回复/更新反馈入参。
//
// Author: Charlie
type UpdateParam struct {
	ID     string  `json:"id" binding:"required"`
	Status string  `json:"status" binding:"required"`
	Reply  *string `json:"reply"`
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// PageParam 管理端分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Title                string `form:"title"`
	Category             string `form:"category"`
	Status               string `form:"status"`
	SubmitterAccountType string `form:"submitter_account_type"`
}

// SubmitMeta 提交者信息。
//
// Author: Charlie
type SubmitMeta struct {
	AccountType string
	AccountID   string
	CreatedBy   string
}

// ReplyMeta 回复者信息。
//
// Author: Charlie
type ReplyMeta struct {
	RepliedBy string
	UpdatedBy string
}
