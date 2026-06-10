package user

type (
	PermissionItem struct {
		PermissionCode      string  `json:"permission_code"`
		Scope               string  `json:"scope"`
		CustomScopeGroupIds *string `json:"custom_scope_group_ids"`
		CustomScopeOrgIds   *string `json:"custom_scope_org_ids"`
	}

	UpdateStatusParam struct {
		IDs    []string `json:"ids"`
		Status string   `json:"status"`
	}

	BatchImportParam struct {
		Users []BatchImportUser `json:"users"`
	}

	BatchImportUser struct {
		Username *string `json:"username"`
		Nickname *string `json:"nickname"`
		Phone    *string `json:"phone"`
		Email    *string `json:"email"`
		Gender   *string `json:"gender"`
		Password *string `json:"password"`
	}

	SysQuickAction struct {
		ID         string `gorm:"primaryKey;size:32" json:"id"`
		UserID     string `gorm:"size:32;uniqueIndex:idx_user_resource;not null" json:"user_id"`
		ResourceID string `gorm:"size:32;uniqueIndex:idx_user_resource;not null" json:"resource_id"`
		SortCode   int    `gorm:"default:0" json:"sort_code"`
	}
)

// UserVO 用户视图对象
// 与 SysUser 对应字段已添加 mapstruct 标签以供代码生成
// Birthday / LastLoginAt / CreatedAt / UpdatedAt 为 *time.Time → string 手动转换
// OrgNames / GroupNames / PositionName / RoleIDs 为 VO 扩展字段，无 SysUser 对应
type UserVO struct {
	ID           string   `json:"id" mapstruct:"ID"`
	Username     *string  `json:"username" mapstruct:"Username"`
	Nickname     *string  `json:"nickname" mapstruct:"Nickname"`
	Avatar       *string  `json:"avatar" mapstruct:"Avatar"`
	Motto        *string  `json:"motto" mapstruct:"Motto"`
	Gender       *string  `json:"gender" mapstruct:"Gender"`
	Birthday     string   `json:"birthday"`
	Email        *string  `json:"email" mapstruct:"Email"`
	Github       *string  `json:"github" mapstruct:"Github"`
	Phone        *string  `json:"phone" mapstruct:"Phone"`
	OrgID        *string  `json:"org_id" mapstruct:"OrgID"`
	PositionID   *string  `json:"position_id" mapstruct:"PositionID"`
	GroupID      *string  `json:"group_id" mapstruct:"GroupID"`
	Status       string   `json:"status" mapstruct:"Status"`
	LastLoginAt  string   `json:"last_login_at"`
	LastLoginIP  *string  `json:"last_login_ip" mapstruct:"LastLoginIP"`
	LoginCount   int      `json:"login_count" mapstruct:"LoginCount"`
	CreatedAt    string   `json:"created_at"`
	CreatedBy    *string  `json:"created_by" mapstruct:"CreatedBy"`
	UpdatedAt    string   `json:"updated_at"`
	UpdatedBy    *string  `json:"updated_by" mapstruct:"UpdatedBy"`
	Password     *string  `json:"password" mapstruct:"Password"`
	RoleIDs      []string `json:"role_ids"`
	OrgNames     []string `json:"org_names"`
	GroupNames   []string `json:"group_names"`
	PositionName *string  `json:"position_name"`
}

type UserPageParam struct {
	Current int    `json:"current" form:"current"`
	Size    int    `json:"size" form:"size"`
	Keyword string `json:"keyword" form:"keyword"`
	Status  string `json:"status" form:"status"`
}

type GrantRoleParam struct {
	UserID  string   `json:"user_id"`
	RoleIDs []string `json:"role_ids"`
}

type GrantUserPermissionParam struct {
	UserID      string         `json:"user_id"`
	Permissions []PermissionItem `json:"permissions"`
}

type UpdateProfileParam struct {
	Username *string `json:"username"`
	Nickname *string `json:"nickname"`
	Motto    *string `json:"motto"`
	Gender   *string `json:"gender"`
	Birthday string  `json:"birthday"`
	Email    *string `json:"email"`
	Github   *string `json:"github"`
	Phone    *string `json:"phone"`
}

type UpdateAvatarParam struct {
	Avatar string `json:"avatar"`
}

type UpdatePasswordParam struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type PermissionDetail struct {
	PermissionCode      string  `json:"permission_code"`
	Scope               string  `json:"scope"`
	CustomScopeGroupIds *string `json:"custom_scope_group_ids"`
	CustomScopeOrgIds   *string `json:"custom_scope_org_ids"`
}