// Package workspace 提供管理端工作台（快捷应用与本人近期活动）。
//
// Author: Charlie
package workspace

import (
	"time"
)

// WorkspaceShortcut 工作台个人快捷应用，对应表 sys_workspace_shortcut。
//
// Author: Charlie
type WorkspaceShortcut struct {
	ID         string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	AccountID  string    `gorm:"column:account_id;size:64;not null" json:"account_id"`
	ResourceID string    `gorm:"column:resource_id;size:64;not null" json:"resource_id"`
	Sort       int       `gorm:"column:sort;not null" json:"sort"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy  *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy  *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 WorkspaceShortcut 对应的数据库表名。
func (WorkspaceShortcut) TableName() string { return "sys_workspace_shortcut" }

// MenuResource 工作台菜单资源投影，对应表 sys_resource。
//
// Author: Charlie
type MenuResource struct {
	ID           string  `gorm:"column:id;primaryKey;size:64" json:"id"`
	Code         string  `gorm:"column:code;size:64;not null" json:"code"`
	Name         string  `gorm:"column:name;size:64;not null" json:"name"`
	ResourceType string  `gorm:"column:resource_type;size:32;not null" json:"resource_type"`
	Path         *string `gorm:"column:path;size:255" json:"path"`
	Icon         *string `gorm:"column:icon;size:255" json:"icon"`
	Status       string  `gorm:"column:status;size:32;not null" json:"status"`
}

// TableName 返回 MenuResource 对应的数据库表名。
func (MenuResource) TableName() string { return "sys_resource" }

// AuditActivity 工作台只读审计日志投影，对应表 sys_operation_audit_log。
//
// Author: Charlie
type AuditActivity struct {
	ID           string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	Module       string    `gorm:"column:module;size:64;not null" json:"module"`
	ResourceType *string   `gorm:"column:resource_type;size:128" json:"resource_type"`
	ModuleLabel  *string   `gorm:"column:module_label;size:128" json:"module_label"`
	Action       string    `gorm:"column:action;size:64;not null" json:"action"`
	ActionName   *string   `gorm:"column:action_name;size:128" json:"action_name"`
	ActionType   *string   `gorm:"column:action_type;size:32" json:"action_type"`
	Summary      *string   `gorm:"column:summary;size:2000" json:"summary"`
	Success      bool      `gorm:"column:success;not null" json:"success"`
	IP           *string   `gorm:"column:ip;size:64" json:"ip"`
	UserAgent    *string   `gorm:"column:user_agent;size:512" json:"user_agent"`
	OperatorName *string   `gorm:"column:operator_name;size:128" json:"operator_name"`
	DurationMs   *int      `gorm:"column:duration_ms" json:"duration_ms"`
	ResourceID   *string   `gorm:"column:resource_id;size:128" json:"resource_id"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 返回 AuditActivity 对应的数据库表名。
func (AuditActivity) TableName() string { return "sys_operation_audit_log" }
