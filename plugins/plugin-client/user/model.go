package user

import "time"

type ClientUser struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	Username    *string    `gorm:"size:32;index" json:"username"`
	Password    *string    `gorm:"size:255" json:"password"`
	Nickname    *string    `gorm:"size:32" json:"nickname"`
	Avatar      *string    `gorm:"type:longtext" json:"avatar"`
	Motto       *string    `gorm:"size:32" json:"motto"`
	Gender      *string    `gorm:"size:8" json:"gender"`
	Birthday    *time.Time `gorm:"type:date" json:"birthday"`
	Email       *string    `gorm:"size:64" json:"email"`
	Github      *string    `gorm:"size:64" json:"github"`
	Phone       *string    `gorm:"size:32" json:"phone"`
	OrgID       *string    `gorm:"size:32" json:"org_id"`
	PositionID  *string    `gorm:"size:32" json:"position_id"`
	GroupID     *string    `gorm:"size:32" json:"group_id"`
	Status      string     `gorm:"size:16;default:ACTIVE" json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP *string    `gorm:"size:64" json:"last_login_ip"`
	LoginCount  int        `gorm:"default:0" json:"login_count"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *string    `gorm:"size:32" json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *string    `gorm:"size:32" json:"updated_by"`
}

func (ClientUser) TableName() string { return "client_user" }
