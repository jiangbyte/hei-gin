package account

import "hei-gin/internal/modules/iam/relation"

// AddParam ç®¡ç†ç«¯åˆ›å»ºè´¦å·å…¥å‚ã€‚
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

// EditParam ç®¡ç†ç«¯æ›´æ–°è´¦å·å…¥å‚ã€‚
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

// PageParam è´¦å·åˆ†é¡µæŸ¥è¯¢ã€‚
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

// OwnResourceQuery è´¦å·èµ„æºæŽˆæƒæŸ¥è¯¢å…¥å‚ï¼ˆaccount_type é»˜è®¤å–è´¦å·è‡ªèº«ç±»åž‹ï¼‰ã€‚
//
// Author: Charlie
type OwnResourceQuery struct {
	ID          string `form:"id" json:"id" binding:"required"`
	AccountType string `form:"account_type" json:"account_type"`
}

// GrantRoleParam è´¦å·æŽˆæƒè§’è‰²å…¥å‚ã€‚
//
// Author: Charlie
type GrantRoleParam struct {
	ID      string   `json:"id" binding:"required"`
	RoleIDs []string `json:"role_ids"`
}

// GrantGroupParam è´¦å·æŽˆæƒç”¨æˆ·ç»„å…¥å‚ã€‚
//
// Author: Charlie
type GrantGroupParam struct {
	ID       string   `json:"id" binding:"required"`
	GroupIDs []string `json:"group_ids"`
}

// GrantDeptParam è´¦å·æŽˆæƒéƒ¨é—¨å…¥å‚ã€‚
//
// Author: Charlie
type GrantDeptParam struct {
	ID            string                   `json:"id" binding:"required"`
	AccountType   string                   `json:"account_type"`
	GrantInfoList []relation.DeptGrantInfo `json:"grant_info_list"`
}

// GrantResourceParam è´¦å·æŽˆæƒèµ„æºï¼ˆç®¡ç†ç«¯/å®¢æˆ·ç«¯ï¼‰å…¥å‚ã€‚
//
// Author: Charlie
type GrantResourceParam struct {
	ID            string                       `json:"id" binding:"required"`
	AccountType   string                       `json:"account_type"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
