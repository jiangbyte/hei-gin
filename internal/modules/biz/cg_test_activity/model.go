// Package cg_test_activity 为代码生成演示的活动业务模块。
//
// Author: Charlie
package cg_test_activity

import (
	"time"

	"gorm.io/datatypes"
)

// Activity 演示活动实体，对应表 cg_test_activity。
//
// Author: Charlie
type Activity struct {
	ID              string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Code            string         `gorm:"column:code;size:64" json:"code"`
	Name            string         `gorm:"column:name;size:120" json:"name"`
	Category        *string        `gorm:"column:category;size:32" json:"category"`
	Type            string         `gorm:"column:type;size:32" json:"type"`
	Status          string         `gorm:"column:status;size:32" json:"status"`
	CoverURL        *string        `gorm:"column:cover_url;size:512" json:"cover_url"`
	Description     *string        `gorm:"column:description;type:text" json:"description"`
	StartAt         time.Time      `gorm:"column:start_at" json:"start_at"`
	EndAt           *time.Time     `gorm:"column:end_at" json:"end_at"`
	MaxParticipants int            `gorm:"column:max_participants" json:"max_participants"`
	Price           float64        `gorm:"column:price" json:"price"`
	IsPublic        bool           `gorm:"column:is_public" json:"is_public"`
	NeedApproval    bool           `gorm:"column:need_approval" json:"need_approval"`
	RuleConfig      datatypes.JSON `gorm:"column:rule_config;type:json" json:"rule_config"`
	Extra           datatypes.JSON `gorm:"column:extra;type:json" json:"extra"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy       *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy       *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
	OwnerDeptID     *string        `gorm:"column:owner_dept_id;size:64" json:"owner_dept_id"`
}

// TableName 返回 Activity 对应的数据库表名。
func (Activity) TableName() string { return "cg_test_activity" }
