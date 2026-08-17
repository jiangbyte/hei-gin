// internal/modules/iam/position/model.go 数据模型。
//
// Author: Charlie

package position

import (
	"time"

	"gorm.io/datatypes"
)

// Position 映射 sys_position 职位。
//
// Author: Charlie
type Position struct {
	ID          string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Name        string         `gorm:"column:name;size:64;not null" json:"name"`
	Category    string         `gorm:"column:category;size:32;not null" json:"category"`
	OwnerDeptID *string        `gorm:"column:owner_dept_id;size:64" json:"owner_dept_id"`
	Sort        int            `gorm:"column:sort;not null;default:99" json:"sort"`
	IsVirtual   bool           `gorm:"column:is_virtual;not null;default:false" json:"is_virtual"`
	Status      string         `gorm:"column:status;size:32;not null" json:"status"`
	Description *string        `gorm:"column:description" json:"description"`
	Extra       datatypes.JSON `gorm:"column:extra;type:json" json:"extra"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (Position) TableName() string { return "sys_position" }
