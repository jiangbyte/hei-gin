// internal/modules/iam/group/param.go 入参定义。
//
// Author: Charlie

package group

import "hei-gin/internal/modules/iam/relation"

// AddParam åˆ›å»ºç”¨æˆ·ç»„å…¥å‚ã€‚
//
// Author: Charlie
type AddParam struct {
	Name        string  `json:"name" binding:"required"`
	OwnerDeptID *string `json:"owner_dept_id"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}

// EditParam æ›´æ–°ç”¨æˆ·ç»„å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam ç”¨æˆ·ç»„åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	Current int    `form:"current" json:"current"`
	Size    int    `form:"size" json:"size"`
	Name    string `form:"name" json:"name"`
	Status  string `form:"status" json:"status"`
}

// Normalize åˆ†é¡µè§„èŒƒåŒ–ã€‚
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

// IDsParam æ‰¹é‡ ID å…¥å‚ã€‚
//
// Author: Charlie
type IDsParam struct {
	IDs []string `json:"ids" binding:"required"`
}

// OwnResourceQuery ç”¨æˆ·ç»„èµ„æºæŽˆæƒæŸ¥è¯¢å…¥å‚ã€‚
//
// Author: Charlie
type OwnResourceQuery struct {
	ID          string `form:"id" json:"id" binding:"required"`
	AccountType string `form:"account_type" json:"account_type"`
}

// GrantUserParam ç”¨æˆ·ç»„æˆå‘˜æŽˆæƒå…¥å‚ã€‚
//
// Author: Charlie
type GrantUserParam struct {
	ID         string   `json:"id" binding:"required"`
	AccountIDs []string `json:"account_ids"`
}

// GrantRoleParam ç”¨æˆ·ç»„æŽˆæƒè§’è‰²å…¥å‚ã€‚
//
// Author: Charlie
type GrantRoleParam struct {
	ID          string   `json:"id" binding:"required"`
	AccountType string   `json:"account_type"`
	RoleIDs     []string `json:"role_ids"`
}

// GrantResourceParam ç”¨æˆ·ç»„æŽˆæƒèµ„æºï¼ˆç®¡ç†ç«¯/å®¢æˆ·ç«¯ï¼‰å…¥å‚ã€‚
//
// Author: Charlie
type GrantResourceParam struct {
	ID            string                       `json:"id" binding:"required"`
	AccountType   string                       `json:"account_type"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
