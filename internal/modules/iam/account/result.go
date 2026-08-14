// internal/modules/iam/account/result.go 出参定义。
//
// Author: Charlie

package account

import (
	"time"

	"hei-gin/internal/modules/iam/client"
	"hei-gin/internal/modules/iam/group"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/resource"
	"hei-gin/internal/modules/iam/role"
)

// AccountResult è´¦å·è¯¦æƒ…/åˆ†é¡µè¡Œã€‚
//
// Author: Charlie
type AccountResult struct {
	ID                 string     `json:"id"`
	Account            string     `json:"account"`
	AccountType        string     `json:"account_type"`
	AccountStatus      string     `json:"account_status"`
	Name               *string    `json:"name"`
	Nickname           *string    `json:"nickname"`
	Avatar             *string    `json:"avatar"`
	Signature          *string    `json:"signature"`
	Phone              *string    `json:"phone"`
	Email              *string    `json:"email"`
	Remark             *string    `json:"remark"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CancelledBy        *string    `json:"cancelled_by"`
	CancelReason       *string    `json:"cancel_reason"`
	LastLoginIP        *string    `json:"last_login_ip"`
	LastLoginAddress   *string    `json:"last_login_address"`
	LastLoginTime      *time.Time `json:"last_login_time"`
	LastLoginDevice    *string    `json:"last_login_device"`
	LatestLoginIP      *string    `json:"latest_login_ip"`
	LatestLoginAddress *string    `json:"latest_login_address"`
	LatestLoginTime    *time.Time `json:"latest_login_time"`
	LatestLoginDevice  *string    `json:"latest_login_device"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          *string    `json:"created_by"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UpdatedBy          *string    `json:"updated_by"`
}

// OwnRoleResult è´¦å·å·²æ‹¥æœ‰è§’è‰²ç»“æžœã€‚
//
// Author: Charlie
type OwnRoleResult struct {
	ID      string      `json:"id"`
	Roles   []role.Role `json:"roles"`
	RoleIDs []string    `json:"role_ids"`
}

// OwnGroupResult è´¦å·å·²æ‹¥æœ‰ç”¨æˆ·ç»„ç»“æžœã€‚
//
// Author: Charlie
type OwnGroupResult struct {
	ID       string        `json:"id"`
	Groups   []group.Group `json:"groups"`
	GroupIDs []string      `json:"group_ids"`
}

// OwnDeptResult è´¦å·å·²æ‹¥æœ‰éƒ¨é—¨æŽˆæƒç»“æžœã€‚
//
// Author: Charlie
type OwnDeptResult struct {
	ID            string                   `json:"id"`
	GrantInfoList []relation.DeptGrantInfo `json:"grant_info_list"`
}

// OwnResourceResult è´¦å·å·²æ‹¥æœ‰ç®¡ç†ç«¯èµ„æºæŽˆæƒç»“æžœã€‚
//
// Author: Charlie
type OwnResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []resource.GrantModule       `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

// OwnClientResourceResult è´¦å·å·²æ‹¥æœ‰å®¢æˆ·ç«¯èµ„æºæŽˆæƒç»“æžœã€‚
//
// Author: Charlie
type OwnClientResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []client.GrantModule         `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
