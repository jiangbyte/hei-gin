// Package result 提供跨 IAM 模块共用的只读结果视图（避免模块间循环依赖）。
//
// Author: Charlie
package result

import "time"

// AccountView 账号结果行（与 account.AccountResult 字段一致，供角色/用户组成员回显）。
//
// Author: Charlie
type AccountView struct {
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

// OwnUserResult 主体（角色/用户组）已关联账号结果。
//
// Author: Charlie
type OwnUserResult struct {
	ID         string        `json:"id"`
	Users      []AccountView `json:"users"`
	AccountIDs []string      `json:"account_ids"`
}
