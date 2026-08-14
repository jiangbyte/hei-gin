package oauth

import "time"

// AccountOAuthBinding 对应 sys_account_oauth_binding。
//
// Author: Charlie
type AccountOAuthBinding struct {
	ID         string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	AccountID  string    `gorm:"column:account_id;size:64;not null" json:"account_id"`
	Provider   string    `gorm:"column:provider;size:32;not null" json:"provider"`
	OpenID     string    `gorm:"column:open_id;size:128;not null" json:"open_id"`
	UnionID    *string   `gorm:"column:union_id;size:128" json:"union_id"`
	Nickname   *string   `gorm:"column:nickname;size:128" json:"nickname"`
	Avatar     *string   `gorm:"column:avatar;type:text" json:"avatar"`
	RawProfile string    `gorm:"column:raw_profile;type:jsonb;not null;default:'{}'" json:"raw_profile"`
	BoundAt    time.Time `gorm:"column:bound_at" json:"bound_at"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy  *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy  *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 表名。
func (AccountOAuthBinding) TableName() string { return "sys_account_oauth_binding" }
