package user

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
