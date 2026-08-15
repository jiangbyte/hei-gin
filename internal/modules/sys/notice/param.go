// internal/modules/sys/notice/param.go 入参定义。
//
// Author: Charlie

package notice

import (
	"time"

	"hei-gin/internal/framework/core/schema"
)

// CreateParam 创建通知入参。
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

// UpdateParam 更新通知入参。
//
// Author: Charlie
type UpdateParam struct {
	ID string `json:"id" binding:"required"`
	CreateParam
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// ReadParam 标记已读入参。
//
// Author: Charlie
type ReadParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// PinParam 置顶入参。
//
// Author: Charlie
type PinParam struct {
	ID          string     `json:"id" binding:"required"`
	IsPinned    bool       `json:"is_pinned"`
	PinnedUntil *time.Time `json:"pinned_until"`
}

// PageParam 通知分页查询。
//
// Author: Charlie
type PageParam struct {
	schema.PageQuery
	Title  string `form:"title"`
	Status string `form:"status"`
	Kind   string `form:"kind"`
}

// PublishParam 发布通知更新字段。
//
// Author: Charlie
type PublishParam struct {
	Status            string
	PublishAt         time.Time
	SenderAccountID   string
	SenderAccountType string
	UpdatedBy         string
}

// RevokeParam 撤回通知更新字段。
//
// Author: Charlie
type RevokeParam struct {
	Status    string
	RevokedAt time.Time
}

// ReadRecord 已读记录键。
//
// Author: Charlie
type ReadRecord struct {
	NoticeID    string
	AccountType string
	AccountID   string
	ReadAt      time.Time
}
