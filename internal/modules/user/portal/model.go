// Package portal 提供门户端用户资料管理。
//
// Author: Charlie
package portal

import (
	"time"
)

// Profile 门户端用户资料实体，对应表 portal_user_profile。
//
// Author: Charlie
type Profile struct {
	AccountID string    `gorm:"column:account_id;primaryKey;size:64" json:"account_id"`
	Name      *string   `gorm:"column:name;size:64" json:"name"`
	Nickname  *string   `gorm:"column:nickname;size:64" json:"nickname"`
	Avatar    *string   `gorm:"column:avatar;type:text" json:"avatar"`
	Signature *string   `gorm:"column:signature;type:text" json:"signature"`
	Phone     *string   `gorm:"column:phone;size:32" json:"phone"`
	Email     *string   `gorm:"column:email;size:128" json:"email"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Profile 对应的数据库表名。
func (Profile) TableName() string { return "portal_user_profile" }

// AccountPassword 账号密码更新用投影，对应表 sys_account。
//
// Author: Charlie
type AccountPassword struct {
	ID           string `gorm:"column:id;primaryKey"`
	PasswordHash string `gorm:"column:password_hash"`
}

// TableName 返回 AccountPassword 对应的数据库表名。
func (AccountPassword) TableName() string { return "sys_account" }
