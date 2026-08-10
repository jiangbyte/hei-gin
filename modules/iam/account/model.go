package account

import (
	"time"
)

// Account 映射 sys_account 账号主表。
//
// Author: Charlie
type Account struct {
	ID                 string     `gorm:"column:id;primaryKey;size:64" json:"id"`
	PasswordHash       string     `gorm:"column:password_hash;size:255;not null" json:"-"`
	AccountType        string     `gorm:"column:account_type;size:32;not null" json:"account_type"`
	AccountStatus      string     `gorm:"column:account_status;size:32;not null" json:"account_status"`
	CancelledAt        *time.Time `gorm:"column:cancelled_at" json:"cancelled_at"`
	CancelledBy        *string    `gorm:"column:cancelled_by;size:64" json:"cancelled_by"`
	CancelReason       *string    `gorm:"column:cancel_reason" json:"cancel_reason"`
	CancelNotifyEmail  *string    `gorm:"column:cancel_notify_email;size:128" json:"cancel_notify_email"`
	CancelNotifyPhone  *string    `gorm:"column:cancel_notify_phone;size:32" json:"cancel_notify_phone"`
	LastLoginIP        *string    `gorm:"column:last_login_ip;size:64" json:"last_login_ip"`
	LastLoginAddress   *string    `gorm:"column:last_login_address;size:255" json:"last_login_address"`
	LastLoginTime      *time.Time `gorm:"column:last_login_time" json:"last_login_time"`
	LastLoginDevice    *string    `gorm:"column:last_login_device" json:"last_login_device"`
	LatestLoginIP      *string    `gorm:"column:latest_login_ip;size:64" json:"latest_login_ip"`
	LatestLoginAddress *string    `gorm:"column:latest_login_address;size:255" json:"latest_login_address"`
	LatestLoginTime    *time.Time `gorm:"column:latest_login_time" json:"latest_login_time"`
	LatestLoginDevice  *string    `gorm:"column:latest_login_device" json:"latest_login_device"`
	CreatedAt          time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy          *string    `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy          *string    `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (Account) TableName() string { return "sys_account" }

// Identity 映射 sys_account_identity 登录身份。
//
// Author: Charlie
type Identity struct {
	ID           string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	AccountID    string    `gorm:"column:account_id;size:64;not null;index" json:"account_id"`
	IdentityType string    `gorm:"column:identity_type;size:32;not null" json:"identity_type"`
	Identifier   string    `gorm:"column:identifier;size:128;not null" json:"identifier"`
	Verified     bool      `gorm:"column:verified;not null;default:false" json:"verified"`
	IsPrimary    bool      `gorm:"column:is_primary;not null;default:false" json:"is_primary"`
	BindStatus   string    `gorm:"column:bind_status;size:32;not null" json:"bind_status"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy    *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy    *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (Identity) TableName() string { return "sys_account_identity" }

// AdminUserProfile 映射 admin_user_profile 管理端资料。
//
// Author: Charlie
type AdminUserProfile struct {
	AccountID string    `gorm:"column:account_id;primaryKey;size:64" json:"account_id"`
	Name      *string   `gorm:"column:name;size:64" json:"name"`
	Nickname  *string   `gorm:"column:nickname;size:64" json:"nickname"`
	Avatar    *string   `gorm:"column:avatar" json:"avatar"`
	Signature *string   `gorm:"column:signature" json:"signature"`
	Phone     *string   `gorm:"column:phone;size:32" json:"phone"`
	Email     *string   `gorm:"column:email;size:128" json:"email"`
	Remark    *string   `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (AdminUserProfile) TableName() string { return "admin_user_profile" }

// PortalUserProfile 映射 portal_user_profile 门户资料。
//
// Author: Charlie
type PortalUserProfile struct {
	AccountID string    `gorm:"column:account_id;primaryKey;size:64" json:"account_id"`
	Name      *string   `gorm:"column:name;size:64" json:"name"`
	Nickname  *string   `gorm:"column:nickname;size:64" json:"nickname"`
	Avatar    *string   `gorm:"column:avatar" json:"avatar"`
	Signature *string   `gorm:"column:signature" json:"signature"`
	Phone     *string   `gorm:"column:phone;size:32" json:"phone"`
	Email     *string   `gorm:"column:email;size:128" json:"email"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (PortalUserProfile) TableName() string { return "portal_user_profile" }

// 身份类型与绑定状态常量。
const (
	IdentityAccount = "ACCOUNT"
	IdentityEmail   = "EMAIL"
	IdentityPhone   = "PHONE"
	BindBound       = "BOUND"
)
