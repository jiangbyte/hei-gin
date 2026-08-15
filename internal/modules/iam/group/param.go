// internal/modules/iam/group/param.go 入参定义。
//
// Author: Charlie

package group

import (
	"gorm.io/datatypes"

	"hei-gin/internal/modules/iam/relation"
)

// AddParam 创建用户组入参。
//
// Author: Charlie
type AddParam struct {
	Name        string         `json:"name" binding:"required"`
	OwnerDeptID *string        `json:"owner_dept_id"`
	Description *string        `json:"description"`
	Status      string         `json:"status"`
	Extra       datatypes.JSON `json:"extra"`
}

// EditParam 更新用户组入参。
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam 用户组分页查询。
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

// OwnResourceQuery 用户组资源授权查询入参。
//
// Author: Charlie
type OwnResourceQuery struct {
	ID          string `form:"id" json:"id" binding:"required"`
	AccountType string `form:"account_type" json:"account_type"`
}

// GrantUserParam 用户组成员授权入参。
//
// Author: Charlie
type GrantUserParam struct {
	ID         string   `json:"id" binding:"required"`
	AccountIDs []string `json:"account_ids"`
}

// GrantRoleParam 用户组授权角色入参。
//
// Author: Charlie
type GrantRoleParam struct {
	ID          string   `json:"id" binding:"required"`
	AccountType string   `json:"account_type"`
	RoleIDs     []string `json:"role_ids"`
}

// GrantResourceParam 用户组授权资源（管理端/客户端）入参。
//
// Author: Charlie
type GrantResourceParam struct {
	ID            string                       `json:"id" binding:"required"`
	AccountType   string                       `json:"account_type"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
