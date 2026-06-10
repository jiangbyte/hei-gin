package user

import (
	"hei-gin/sdk/utils"
)

type ClientUserVO struct {
	ID          string  `json:"id"`
	Username    *string `json:"username"`
	Password    *string `json:"password"`
	Nickname    *string `json:"nickname"`
	Avatar      *string `json:"avatar"`
	Motto       *string `json:"motto"`
	Gender      *string `json:"gender"`
	Email       *string `json:"email"`
	Github      *string `json:"github"`
	Phone       *string `json:"phone"`
	Status      string  `json:"status"`
	LastLoginIP *string `json:"last_login_ip"`
	LastLoginAt string  `json:"last_login_at"`
	LoginCount  int     `json:"login_count"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ClientUserPageParam struct {
	Current int    `json:"current" form:"current"`
	Size    int    `json:"size" form:"size"`
	Keyword string `json:"keyword" form:"keyword"`
	Status  string `json:"status" form:"status"`
}

type ClientUserCreateParam struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}
type ClientUserModifyParam struct {
	ID       string  `json:"id"`
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Status   string  `json:"status"`
}
type UpdateProfileParam struct {
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Username *string `json:"username"`
}
type UpdateAvatarParam struct {
	Avatar string `json:"avatar"`
}
type UpdatePasswordParam struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}


func toVO(e *ClientUser) ClientUserVO {
	return ClientUserVO{
		ID:          e.ID,
		Username:    e.Username,
		Nickname:    e.Nickname,
		Avatar:      e.Avatar,
		Motto:       e.Motto,
		Gender:      e.Gender,
		Email:       e.Email,
		Github:      e.Github,
		Phone:       e.Phone,
		Status:      e.Status,
		LastLoginIP: e.LastLoginIP,
		LoginCount:  e.LoginCount,
		LastLoginAt: utils.FormatDateTimePtr(e.LastLoginAt),
		CreatedAt:   utils.FormatDateTimePtr(e.CreatedAt),
		UpdatedAt:   utils.FormatDateTimePtr(e.UpdatedAt),
	}
}
