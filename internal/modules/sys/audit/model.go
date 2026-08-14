// Package audit 提供操作审计日志查询。
//
// Author: Charlie
package audit

import (
	"time"

	"gorm.io/datatypes"
)

// OperationLog 操作审计日志实体，对应表 sys_operation_audit_log。
//
// Author: Charlie
type OperationLog struct {
	ID           string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Module       string         `gorm:"column:module;size:64;not null" json:"module"`
	ResourceType *string        `gorm:"column:resource_type;size:128" json:"resource_type"`
	ResourceID   *string        `gorm:"column:resource_id;size:128" json:"resource_id"`
	Action       string         `gorm:"column:action;size:64;not null" json:"action"`
	Summary      *string        `gorm:"column:summary;size:255" json:"summary"`
	BeforeData   datatypes.JSON `gorm:"column:before_data;type:jsonb" json:"before_data"`
	AfterData    datatypes.JSON `gorm:"column:after_data;type:jsonb" json:"after_data"`
	AccountID    *string        `gorm:"column:account_id;size:64" json:"account_id"`
	AccountType  *string        `gorm:"column:account_type;size:32" json:"account_type"`
	RequestID    *string        `gorm:"column:request_id;size:64" json:"request_id"`
	IP           *string        `gorm:"column:ip;size:64" json:"ip"`
	UserAgent    *string        `gorm:"column:user_agent;size:512" json:"user_agent"`
	Success      bool           `gorm:"column:success;not null" json:"success"`
	ErrorMessage *string        `gorm:"column:error_message" json:"error_message"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 返回 OperationLog 对应的数据库表名。
func (OperationLog) TableName() string { return "sys_operation_audit_log" }
