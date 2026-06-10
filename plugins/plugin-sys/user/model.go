package user

import (
	"time"
)

// SysUser entity, table sys_user
type SysUser struct {
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

func (SysUser) TableName() string { return "sys_user" }

type RelUserRole struct {
	ID     string `gorm:"primaryKey;size:32" json:"id"`
	UserID string `gorm:"size:32;uniqueIndex:idx_user_role" json:"user_id"`
	RoleID string `gorm:"size:32;uniqueIndex:idx_user_role;index" json:"role_id"`
}

func (RelUserRole) TableName() string { return "rel_user_role" }

type RelUserPermission struct {
	ID                  string  `gorm:"primaryKey;size:32" json:"id"`
	UserID              string  `gorm:"size:32;uniqueIndex:idx_user_perm" json:"user_id"`
	PermissionCode      string  `gorm:"size:255;uniqueIndex:idx_user_perm;index" json:"permission_code"`
	Scope               string  `gorm:"size:32;default:ALL" json:"scope"`
	CustomScopeGroupIds *string `gorm:"type:text" json:"custom_scope_group_ids"`
	CustomScopeOrgIds   *string `gorm:"type:text" json:"custom_scope_org_ids"`
}

func (RelUserPermission) TableName() string { return "rel_user_permission" }

type RelRolePermission struct {
	ID                  string  `gorm:"primaryKey;size:32" json:"id"`
	RoleID              string  `gorm:"size:32;uniqueIndex:idx_role_perm" json:"role_id"`
	PermissionCode      string  `gorm:"size:255;uniqueIndex:idx_role_perm;index" json:"permission_code"`
	Scope               string  `gorm:"size:32;default:ALL" json:"scope"`
	CustomScopeGroupIds *string `gorm:"type:text" json:"custom_scope_group_ids"`
	CustomScopeOrgIds   *string `gorm:"type:text" json:"custom_scope_org_ids"`
}

func (RelRolePermission) TableName() string { return "rel_role_permission" }

type RelRoleResource struct {
	ID         string `gorm:"primaryKey;size:32" json:"id"`
	RoleID     string `gorm:"size:32;uniqueIndex:idx_role_resource" json:"role_id"`
	ResourceID string `gorm:"size:32;uniqueIndex:idx_role_resource;index" json:"resource_id"`
}

func (RelRoleResource) TableName() string { return "rel_role_resource" }
