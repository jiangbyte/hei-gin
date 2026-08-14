package role

import "hei-gin/internal/modules/iam/relation"

// AddParam åˆ›å»ºè§’è‰²å…¥å‚ã€‚
//
// Author: Charlie
type AddParam struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category"`
	ScopeType   string  `json:"scope_type"`
	OwnerDeptID *string `json:"owner_dept_id"`
	Sort        int     `json:"sort"`
	Status      string  `json:"status"`
	Description *string `json:"description"`
}

// EditParam æ›´æ–°è§’è‰²å…¥å‚ã€‚
//
// Author: Charlie
type EditParam struct {
	ID string `json:"id" binding:"required"`
	AddParam
}

// PageParam è§’è‰²åˆ†é¡µæŸ¥è¯¢ã€‚
//
// Author: Charlie
type PageParam struct {
	Current int    `form:"current" json:"current"`
	Size    int    `form:"size" json:"size"`
	Code    string `form:"code" json:"code"`
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

// OwnResourceQuery è§’è‰²èµ„æºæŽˆæƒæŸ¥è¯¢å…¥å‚ã€‚
//
// Author: Charlie
type OwnResourceQuery struct {
	ID          string `form:"id" json:"id" binding:"required"`
	AccountType string `form:"account_type" json:"account_type"`
}

// GrantUserParam è§’è‰²æˆå‘˜æŽˆæƒå…¥å‚ã€‚
//
// Author: Charlie
type GrantUserParam struct {
	ID         string   `json:"id" binding:"required"`
	AccountIDs []string `json:"account_ids"`
}

// GrantResourceParam è§’è‰²æŽˆæƒèµ„æºï¼ˆç®¡ç†ç«¯/å®¢æˆ·ç«¯ï¼‰å…¥å‚ã€‚
//
// Author: Charlie
type GrantResourceParam struct {
	ID            string                       `json:"id" binding:"required"`
	AccountType   string                       `json:"account_type"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
