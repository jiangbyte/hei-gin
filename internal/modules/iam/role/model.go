// internal/modules/iam/role/model.go 数据模型。
//
// Author: Charlie

package role

import (
	"time"

	"gorm.io/datatypes"
)

// Role 映射 sys_role 角色。
//
// Author: Charlie
type Role struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Code        string         `gorm:"column:code;size:64;uniqueIndex;not null" json:"code"`
	Name        string         `gorm:"column:name;size:64;not null" json:"name"`
	Category    string         `gorm:"column:category;size:64;not null" json:"category"`
	ScopeType   string         `gorm:"column:scope_type;size:32;not null" json:"scope_type"`
	OwnerDeptID *string        `gorm:"column:owner_dept_id;size:64" json:"owner_dept_id"`
	Sort        int            `gorm:"column:sort;not null;default:99" json:"sort"`
	Status      string         `gorm:"column:status;size:32;not null" json:"status"`
	IsBuiltin   bool           `gorm:"column:is_builtin;not null;default:false" json:"is_builtin"`
	Description *string        `gorm:"column:description" json:"description"`
	Extra       datatypes.JSON `gorm:"column:extra;type:jsonb" json:"extra"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (Role) TableName() string { return "sys_role" }
