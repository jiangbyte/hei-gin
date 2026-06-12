package user

import "hei-gin/sdk/utils"

// ClientUserVO 客户端用户视图对象
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

func ClientUserToClientUserVO(src *ClientUser) *ClientUserVO {
	if src == nil {
		return nil
	}

	dst := &ClientUserVO{}
	dst.ID = src.ID
	dst.Username = src.Username
	dst.Nickname = src.Nickname
	dst.Avatar = src.Avatar
	dst.Motto = src.Motto
	dst.Gender = src.Gender
	dst.Email = src.Email
	dst.Github = src.Github
	dst.Phone = src.Phone
	dst.Status = src.Status
	dst.LastLoginIP = src.LastLoginIP
	dst.LoginCount = src.LoginCount
	dst.LastLoginAt = utils.FormatDateTimePtr(src.LastLoginAt)
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

func ClientUserVOToClientUser(src *ClientUserVO) *ClientUser {
	if src == nil {
		return nil
	}

	dst := &ClientUser{}
	dst.ID = src.ID
	dst.Username = src.Username
	dst.Nickname = src.Nickname
	dst.Avatar = src.Avatar
	dst.Motto = src.Motto
	dst.Gender = src.Gender
	dst.Email = src.Email
	dst.Github = src.Github
	dst.Phone = src.Phone
	dst.Status = src.Status
	dst.LastLoginIP = src.LastLoginIP
	dst.LoginCount = src.LoginCount
	return dst
}

// ClientUserPageParam 客户端用户分页参数
type ClientUserPageParam struct {
	Current int    `json:"current" form:"current"`
	Size    int    `json:"size" form:"size"`
	Keyword string `json:"keyword" form:"keyword"`
	Status  string `json:"status" form:"status"`
}

// ClientUserCreateParam 添加客户端用户参数
type ClientUserCreateParam struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}

// ClientUserModifyParam 编辑客户端用户参数
type ClientUserModifyParam struct {
	ID       string  `json:"id"`
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Status   string  `json:"status"`
}

// UpdateProfileParam 更新个人信息参数
type UpdateProfileParam struct {
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Username *string `json:"username"`
}

// UpdateAvatarParam 更新头像参数
type UpdateAvatarParam struct {
	Avatar string `json:"avatar"`
}

// UpdatePasswordParam 修改密码参数
type UpdatePasswordParam struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
