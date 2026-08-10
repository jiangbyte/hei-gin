// Package feedback 提供用户反馈管理。
package feedback

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// Feedback 用户反馈实体，对应表 msg_feedback。
//
// Author: Charlie
type Feedback struct {
	ID                   string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Title                string         `gorm:"column:title;size:255" json:"title"`
	Content              string         `gorm:"column:content;type:text" json:"content"`
	Category             string         `gorm:"column:category;size:64" json:"category"`
	Contact              *string        `gorm:"column:contact;size:255" json:"contact"`
	AttachObjectNames    datatypes.JSON `gorm:"column:attach_object_names;type:jsonb" json:"attach_object_names"`
	Status               string         `gorm:"column:status;size:32" json:"status"`
	Reply                *string        `gorm:"column:reply;type:text" json:"reply"`
	RepliedBy            *string        `gorm:"column:replied_by;size:64" json:"replied_by"`
	RepliedAt            *time.Time     `gorm:"column:replied_at" json:"replied_at"`
	SubmitterAccountType string         `gorm:"column:submitter_account_type;size:32" json:"submitter_account_type"`
	SubmitterAccountID   string         `gorm:"column:submitter_account_id;size:64" json:"submitter_account_id"`
	CreatedAt            time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy            *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy            *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Feedback 对应的数据库表名。
func (Feedback) TableName() string { return "msg_feedback" }

func jsonList(v []string) datatypes.JSON {
	if v == nil {
		v = []string{}
	}
	b, _ := json.Marshal(v)
	return b
}
