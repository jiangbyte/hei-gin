// Package weakpassword 提供弱密码黑名单管理。
package weakpassword

import "time"

// WeakPassword 弱密码实体，对应表 sys_weak_password。
//
// Author: Charlie
type WeakPassword struct {
	ID        string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	Password  string    `gorm:"column:password;size:255;uniqueIndex;not null" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 WeakPassword 对应的数据库表名。
func (WeakPassword) TableName() string { return "sys_weak_password" }
