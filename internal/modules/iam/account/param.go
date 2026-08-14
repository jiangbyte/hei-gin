// internal/modules/iam/account/param.go 入参定义。
//
// Author: Charlie

package account

import "hei-gin/internal/modules/iam/relation"

// AddParam 管理端创建账号入参。
//
// Author: Charlie
type AddParam struct {
	Account       string  `json:"account" binding:"required"`
	Password      string  `json:"password" binding:"required"`
	AccountType   string  `json:"account_type" binding:"required"`
	AccountStatus string  `json:"account_status"`
	Name          *string `json:"name"`
	Nickname      *string `json:"nickname"`
	Avatar        *string `json:"avatar"`
	Signature     *string `json:"signature"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Remark        *string `json:"remark"`
}

// EditParam 管理端更新账号入参。
//
// Author: Charlie
type EditParam struct {
	ID            string  `json:"id" binding:"required"`
	Account       string  `json:"account" binding:"required"`
	Password      *string `json:"password"`
	AccountType   string  `json:"account_type" binding:"required"`
	AccountStatus string  `json:"account_status"`
	Name          *string `json:"name"`
	Nickname      *string `json:"nickname"`
	Avatar        *string `json:"avatar"`
	Signature     *string `json:"signature"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Remark        *string `json:"remark"`
}

// PageParam 账号分页查询。
//
// Author: Charlie
type PageParam struct {
	Current       int    `form:"current" json:"current"`
	Size          int    `form:"size" json:"size"`
	Account       string `form:"account" json:"account"`
	Name          string `form:"name" json:"name"`
	AccountType   string `form:"account_type" json:"account_type"`
	AccountStatus string `form:"account_status" json:"account_status"`
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

// OwnResourceQuery 账号资源授权查询入参（account_type 默认取账号自身类型）。
//
// Author: Charlie
type OwnResourceQuery struct {
	ID          string `form:"id" json:"id" binding:"required"`
	AccountType string `form:"account_type" json:"account_type"`
}

// GrantRoleParam 账号授权角色入参。
//
// Author: Charlie
type GrantRoleParam struct {
	ID      string   `json:"id" binding:"required"`
	RoleIDs []string `json:"role_ids"`
}

// GrantGroupParam 账号授权用户组入参。
//
// Author: Charlie
type GrantGroupParam struct {
	ID       string   `json:"id" binding:"required"`
	GroupIDs []string `json:"group_ids"`
}

// GrantDeptParam 账号授权部门入参。
//
// Author: Charlie
type GrantDeptParam struct {
	ID            string                   `json:"id" binding:"required"`
	AccountType   string                   `json:"account_type"`
	GrantInfoList []relation.DeptGrantInfo `json:"grant_info_list"`
}

// GrantResourceParam 账号授权资源（管理端/客户端）入参。
//
// Author: Charlie
type GrantResourceParam struct {
	ID            string                       `json:"id" binding:"required"`
	AccountType   string                       `json:"account_type"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
