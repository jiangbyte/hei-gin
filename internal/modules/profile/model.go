// internal/modules/profile/model.go 用户中心资料模型。
//
// Author: Charlie

package profile

import (
	"time"
)

// ProfileTableAdmin 管理端用户资料表（对齐 hei-boot profile_user_admin）。
const ProfileTableAdmin = "profile_user_admin"

// ProfileTablePortal 门户端用户资料表（对齐 hei-boot profile_user_portal）。
const ProfileTablePortal = "profile_user_portal"

// Profile 用户资料实体（admin/portal 共用结构；portal 表无 remark 列）。
//
// Author: Charlie
type Profile struct {
	AccountID string    `gorm:"column:account_id;primaryKey;size:64" json:"account_id"`
	Nickname  *string   `gorm:"column:nickname;size:64" json:"nickname"`
	Avatar    *string   `gorm:"column:avatar;type:text" json:"avatar"`
	Signature *string   `gorm:"column:signature;type:text" json:"signature"`
	Phone     *string   `gorm:"column:phone;size:32" json:"phone"`
	Email     *string   `gorm:"column:email;size:128" json:"email"`
	Remark    *string   `gorm:"column:remark;type:text" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// AdminProfileModel 管理端资料迁移模型（表 profile_user_admin）。
//
// Author: Charlie
type AdminProfileModel Profile

// TableName 返回表名。
func (AdminProfileModel) TableName() string { return ProfileTableAdmin }

// PortalProfileModel 门户端资料迁移模型（表 profile_user_portal）。
//
// Author: Charlie
type PortalProfileModel Profile

// TableName 返回表名。
func (PortalProfileModel) TableName() string { return ProfileTablePortal }

// AccountPassword 账号密码更新用投影，对应表 sys_account。
//
// Author: Charlie
type AccountPassword struct {
	ID           string `gorm:"column:id;primaryKey"`
	PasswordHash string `gorm:"column:password_hash"`
}

// TableName 返回表名。
func (AccountPassword) TableName() string { return "sys_account" }
