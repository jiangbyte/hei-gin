// internal/modules/iam/role/param.go 入参定义。
//
// Author: Charlie

package role

import (
	"gorm.io/datatypes"

	"hei-gin/internal/modules/iam/relation"
)

// AddParam 创建角色入参。
//
// Author: Charlie
type AddParam struct {
	Code        string         `json:"code" binding:"required"`
	Name        string         `json:"name" binding:"required"`
	Category    string         `json:"category"`
	ScopeType   string         `json:"scope_type"`
	OwnerDeptID *string        `json:"owner_dept_id"`
	Sort        int            `json:"sort"`
	Status      string         `json:"status"`
	IsBuiltin   *bool          `json:"is_builtin"`
	Description *string        `json:"description"`
	Extra       datatypes.JSON `json:"extra"`
}

// EditParam 更新角色入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam 角色分页查询。
//
// Author: Charlie
type PageParam struct {
	Current   int    `form:"current" json:"current"`
	Size      int    `form:"size" json:"size"`
	Code      string `form:"code" json:"code"`
	Name      string `form:"name" json:"name"`
	Category  string `form:"category" json:"category"`
	ScopeType string `form:"scope_type" json:"scope_type"`
	Status    string `form:"status" json:"status"`
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

// OwnResourceQuery 角色资源授权查询入参。
//
// Author: Charlie
type OwnResourceQuery struct {
	ID          string `form:"id" json:"id" binding:"required"`
	AccountType string `form:"account_type" json:"account_type"`
}

// GrantUserParam 角色成员授权入参。
//
// Author: Charlie
type GrantUserParam struct {
	ID         string   `json:"id" binding:"required"`
	AccountIDs []string `json:"account_ids"`
}

// GrantResourceParam 角色授权资源（管理端/客户端）入参。
//
// Author: Charlie
type GrantResourceParam struct {
	ID            string                       `json:"id" binding:"required"`
	AccountType   string                       `json:"account_type"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
