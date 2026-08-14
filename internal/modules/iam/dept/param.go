// internal/modules/iam/dept/param.go 入参定义。
//
// Author: Charlie

package dept

// AddParam 创建部门入参。
//
// Author: Charlie
type AddParam struct {
	ParentID       *string `json:"parent_id"`
	MasterID       *string `json:"master_id"`
	DeputyMasterID *string `json:"deputy_master_id"`
	Name           string  `json:"name" binding:"required"`
	Category       string  `json:"category" binding:"required"`
	Sort           int     `json:"sort"`
	IsVirtual      bool    `json:"is_virtual"`
	Status         string  `json:"status"`
}

// EditParam 更新部门入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam 部门分页查询。
//
// Author: Charlie
type PageParam struct {
	Current int    `form:"current" json:"current"`
	Size    int    `form:"size" json:"size"`
	Name    string `form:"name" json:"name"`
	Status  string `form:"status" json:"status"`
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
