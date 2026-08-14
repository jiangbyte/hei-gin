package notice

import (
	"time"

	"hei-gin/internal/framework/core/schema"
)

// CreateParam åˆ›å»ºé€šçŸ¥å…¥å‚ã€‚
//
// Author: Charlie
type CreateParam struct {
	Kind               string         `json:"kind" binding:"required"`
	Title              string         `json:"title" binding:"required"`
	Content            string         `json:"content" binding:"required"`
	ContentType        string         `json:"content_type" binding:"required"`
	Category           *string        `json:"category"`
	Severity           string         `json:"severity" binding:"required"`
	TargetScope        string         `json:"target_scope" binding:"required"`
	TargetAccountTypes []string       `json:"target_account_types"`
	TargetAccountIDs   []string       `json:"target_account_ids"`
	TargetDeptIDs      []string       `json:"target_dept_ids"`
	TargetRoleIDs      []string       `json:"target_role_ids"`
	PublishLocations   map[string]any `json:"publish_locations"`
	IsPinned           bool           `json:"is_pinned"`
	PinnedUntil        *time.Time     `json:"pinned_until"`
	Status             string         `json:"status" binding:"required"`
	PublishAt          *time.Time     `json:"publish_at"`
	ExpireAt           *time.Time     `json:"expire_at"`
	Extra              map[string]any `json:"extra"`
}

// UpdateParam æ›´æ–°é€šçŸ¥å…¥å‚ã€‚
//
// Author: Charlie
type UpdateParam struct {
	ID string `json:"id" binding:"required"`
	CreateParam
}

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// ReadParam æ ‡è®°å·²è¯»å…¥å‚ã€‚
//
// Author: Charlie
type ReadParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// PinParam ç½®é¡¶å…¥å‚ã€‚
//
// Author: Charlie
type PinParam struct {
	ID          string     `json:"id" binding:"required"`
	IsPinned    bool       `json:"is_pinned"`
	PinnedUntil *time.Time `json:"pinned_until"`
}

// PageParam é€šçŸ¥åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Title  string `form:"title"`
	Status string `form:"status"`
	Kind   string `form:"kind"`
}

// PublishParam å‘å¸ƒé€šçŸ¥æ›´æ–°å­—æ®µã€‚
//
// Author: Charlie
type PublishParam struct {
	Status            string
	PublishAt         time.Time
	SenderAccountID   string
	SenderAccountType string
	UpdatedBy         string
}

// RevokeParam æ’¤å›žé€šçŸ¥æ›´æ–°å­—æ®µã€‚
//
// Author: Charlie
type RevokeParam struct {
	Status    string
	RevokedAt time.Time
}

// ReadRecord å·²è¯»è®°å½•é”®ã€‚
//
// Author: Charlie
type ReadRecord struct {
	NoticeID    string
	AccountType string
	AccountID   string
	ReadAt      time.Time
}
