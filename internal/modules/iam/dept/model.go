// internal/modules/iam/dept/model.go 数据模型。
//
// Author: Charlie

package dept

import (
	"time"

	"gorm.io/datatypes"
)

// Dept 映射 sys_dept 部门。
//
// Author: Charlie
type Dept struct {
	ID             string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	ParentID       *string        `gorm:"column:parent_id;size:64" json:"parent_id"`
	MasterID       *string        `gorm:"column:master_id;size:64" json:"master_id"`
	DeputyMasterID *string        `gorm:"column:deputy_master_id;size:64" json:"deputy_master_id"`
	Name           string         `gorm:"column:name;size:64;not null" json:"name"`
	Category       string         `gorm:"column:category;size:64;not null" json:"category"`
	Sort           int            `gorm:"column:sort;not null;default:99" json:"sort"`
	IsVirtual      bool           `gorm:"column:is_virtual;not null;default:false" json:"is_virtual"`
	Status         string         `gorm:"column:status;size:32;not null" json:"status"`
	Extra          datatypes.JSON `gorm:"column:extra;type:jsonb" json:"extra"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy      *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy      *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (Dept) TableName() string { return "sys_dept" }
