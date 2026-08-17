// internal/modules/iam/group/model.go 数据模型。
//
// Author: Charlie

package group

import (
	"time"

	"gorm.io/datatypes"
)

// Group 映射 sys_group 用户组。
//
// Author: Charlie
type Group struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Name        string         `gorm:"column:name;size:64;uniqueIndex;not null" json:"name"`
	OwnerDeptID *string        `gorm:"column:owner_dept_id;size:64" json:"owner_dept_id"`
	Description *string        `gorm:"column:description" json:"description"`
	Status      string         `gorm:"column:status;size:32;not null" json:"status"`
	Extra       datatypes.JSON `gorm:"column:extra;type:json" json:"extra"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (Group) TableName() string { return "sys_group" }
