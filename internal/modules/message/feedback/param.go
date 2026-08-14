package feedback

import "hei-gin/internal/framework/core/schema"

// CreateParam æäº¤åé¦ˆå…¥å‚ã€‚
//
// Author: Charlie
type CreateParam struct {
	Title             string   `json:"title" binding:"required"`
	Content           string   `json:"content" binding:"required"`
	Category          string   `json:"category" binding:"required"`
	Contact           *string  `json:"contact"`
	AttachObjectNames []string `json:"attach_object_names"`
}

// UpdateParam å›žå¤/æ›´æ–°åé¦ˆå…¥å‚ã€‚
//
// Author: Charlie
type UpdateParam struct {
	ID     string  `json:"id" binding:"required"`
	Status string  `json:"status" binding:"required"`
	Reply  *string `json:"reply"`
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// PageParam ç®¡ç†ç«¯åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Title                string `form:"title"`
	Category             string `form:"category"`
	Status               string `form:"status"`
	SubmitterAccountType string `form:"submitter_account_type"`
}

// SubmitMeta æäº¤è€…ä¿¡æ¯ã€‚
//
// Author: Charlie
type SubmitMeta struct {
	AccountType string
	AccountID   string
	CreatedBy   string
}

// ReplyMeta å›žå¤è€…ä¿¡æ¯ã€‚
//
// Author: Charlie
type ReplyMeta struct {
	RepliedBy string
	UpdatedBy string
}
