// internal/modules/iam/position/param.go 入参定义。
//
// Author: Charlie

package position

import "gorm.io/datatypes"

// AddParam 创建职位入参（category 可选，对齐 hei-boot）。
//
// Author: Charlie
type AddParam struct {
	Name        string         `json:"name" binding:"required"`
	Category    string         `json:"category"`
	OwnerDeptID *string        `json:"owner_dept_id"`
	Sort        int            `json:"sort"`
	IsVirtual   bool           `json:"is_virtual"`
	Status      string         `json:"status"`
	Description *string        `json:"description"`
	Extra       datatypes.JSON `json:"extra"`
}

// EditParam 更新职位入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam 职位分页查询。
//
// Author: Charlie
type PageParam struct {
	Current  int    `form:"current" json:"current"`
	Size     int    `form:"size" json:"size"`
	Name     string `form:"name" json:"name"`
	Category string `form:"category" json:"category"`
	Status   string `form:"status" json:"status"`
}

// Normalize 分页规范化。
func (q PageParam) Normalize() (current, size int) {
	current, size = q.Current, q.Size
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return current, size
}

// IDsParam 批量 ID 入参。
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}
