package account

import "time"

// AccountResult 账号详情/分页行。
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
