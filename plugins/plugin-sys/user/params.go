package user

import (
	"hei-gin/sdk/utils"
)



type UserVO struct {
	ID           string   `json:"id"`
	Username     *string  `json:"username"`
	Nickname     *string  `json:"nickname"`
	Avatar       *string  `json:"avatar"`
	Motto        *string  `json:"motto"`
	Gender       *string  `json:"gender"`
	Birthday     string   `json:"birthday"`
	Email        *string  `json:"email"`
	Github       *string  `json:"github"`
	Phone        *string  `json:"phone"`
	OrgID        *string  `json:"org_id"`
	PositionID   *string  `json:"position_id"`
	GroupID      *string  `json:"group_id"`
	OrgNames     []string `json:"org_names"`
	GroupNames   []string `json:"group_names"`
	PositionName *string  `json:"position_name"`
	Status       string   `json:"status"`
	LastLoginAt  string   `json:"last_login_at"`
	LastLoginIP  *string  `json:"last_login_ip"`
	LoginCount   int      `json:"login_count"`
	CreatedAt    string   `json:"created_at"`
	CreatedBy    *string  `json:"created_by"`
	UpdatedAt    string   `json:"updated_at"`
	UpdatedBy    *string  `json:"updated_by"`
	RoleIDs      []string `json:"role_ids"`
	Password     *string  `json:"password"`
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
	UserID      string                `json:"user_id"`
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


func toVO(e *SysUser) *UserVO {
	if e == nil {
		return nil
	}
	return &UserVO{
		ID: e.ID, Username: e.Username, Nickname: e.Nickname, Avatar: e.Avatar,
		Motto: e.Motto, Gender: e.Gender, Birthday: utils.FormatDatePtr(e.Birthday),
		Email: e.Email, Github: e.Github, Phone: e.Phone,
		OrgID: e.OrgID, PositionID: e.PositionID, GroupID: e.GroupID,
		Status: e.Status, LastLoginAt: utils.FormatDateTimePtr(e.LastLoginAt), LastLoginIP: e.LastLoginIP,
		LoginCount: e.LoginCount, CreatedAt: utils.FormatDateTimePtr(e.CreatedAt), CreatedBy: e.CreatedBy,
		UpdatedAt: utils.FormatDateTimePtr(e.UpdatedAt), UpdatedBy: e.UpdatedBy,
	}
}
