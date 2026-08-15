// Package notice 提供站内通知公告管理。
//
// Author: Charlie
package notice

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// Notice 通知公告实体，对应表 sys_notice。
//
// Author: Charlie
type Notice struct {
	ID                 string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Kind               string         `gorm:"column:kind;size:32" json:"kind"`
	Title              string         `gorm:"column:title;size:255" json:"title"`
	Content            string         `gorm:"column:content;type:text" json:"content"`
	ContentType        string         `gorm:"column:content_type;size:32" json:"content_type"`
	Category           *string        `gorm:"column:category;size:32" json:"category"`
	Severity           string         `gorm:"column:severity;size:32" json:"severity"`
	TargetScope        string         `gorm:"column:target_scope;size:32" json:"target_scope"`
	TargetAccountTypes datatypes.JSON `gorm:"column:target_account_types;type:jsonb" json:"target_account_types"`
	TargetAccountIDs   datatypes.JSON `gorm:"column:target_account_ids;type:jsonb" json:"target_account_ids"`
	TargetDeptIDs      datatypes.JSON `gorm:"column:target_dept_ids;type:jsonb" json:"target_dept_ids"`
	TargetRoleIDs      datatypes.JSON `gorm:"column:target_role_ids;type:jsonb" json:"target_role_ids"`
	PublishLocations   datatypes.JSON `gorm:"column:publish_locations;type:jsonb" json:"publish_locations"`
	IsPinned           bool           `gorm:"column:is_pinned" json:"is_pinned"`
	PinnedUntil        *time.Time     `gorm:"column:pinned_until" json:"pinned_until"`
	SenderAccountType  *string        `gorm:"column:sender_account_type;size:32" json:"sender_account_type"`
	SenderAccountID    *string        `gorm:"column:sender_account_id;size:64" json:"sender_account_id"`
	SourceType         *string        `gorm:"column:source_type;size:64" json:"source_type"`
	SourceID           *string        `gorm:"column:source_id;size:64" json:"source_id"`
	Status             string         `gorm:"column:status;size:32" json:"status"`
	PublishAt          *time.Time     `gorm:"column:publish_at" json:"publish_at"`
	RevokedAt          *time.Time     `gorm:"column:revoked_at" json:"revoked_at"`
	ExpireAt           *time.Time     `gorm:"column:expire_at" json:"expire_at"`
	ViewCount          int            `gorm:"column:view_count" json:"view_count"`
	IsRead             bool           `gorm:"-" json:"is_read"`
	Extra              datatypes.JSON `gorm:"column:extra;type:jsonb" json:"extra"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy          *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy          *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Notice 对应的数据库表名。
func (Notice) TableName() string { return "sys_notice" }

// NoticeRead 通知已读记录实体，对应表 sys_notice_read。
//
// Author: Charlie
type NoticeRead struct {
	ID          string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	NoticeID    string    `gorm:"column:notice_id;size:64" json:"notice_id"`
	AccountType string    `gorm:"column:account_type;size:32" json:"account_type"`
	AccountID   string    `gorm:"column:account_id;size:64" json:"account_id"`
	ReadAt      time.Time `gorm:"column:read_at" json:"read_at"`
}

// TableName 返回 NoticeRead 对应的数据库表名。
func (NoticeRead) TableName() string { return "sys_notice_read" }

func jsonList(v any) datatypes.JSON {
	if v == nil {
		b, _ := json.Marshal([]any{})
		return b
	}
	b, err := json.Marshal(v)
	if err != nil {
		b, _ = json.Marshal([]any{})
	}
	return b
}

func jsonObj(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return b
}
