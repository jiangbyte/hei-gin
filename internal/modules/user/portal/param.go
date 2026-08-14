package portal

// ProfileUpdateParam 资料更新入参。
//
// Author: Charlie
type ProfileUpdateParam struct {
	Name      *string `json:"name"`
	Nickname  *string `json:"nickname"`
	Signature *string `json:"signature"`
}

// PasswordUpdateParam 密码更新入参。
//
// Author: Charlie
type PasswordUpdateParam struct {
	PasswordKeyID string `json:"password_key_id"`
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"new_password" binding:"required"`
}

// PhoneUpdateParam 手机号更新入参。
//
// Author: Charlie
type PhoneUpdateParam struct {
	PasswordKeyID string `json:"password_key_id"`
	Password      string `json:"password"`
	Phone         string `json:"phone" binding:"required"`
}

// EmailUpdateParam 邮箱更新入参。
//
// Author: Charlie
type EmailUpdateParam struct {
	PasswordKeyID string `json:"password_key_id"`
	Password      string `json:"password"`
	Email         string `json:"email" binding:"required"`
}

// SendCodeParam 绑定验证码发送入参。
//
// Author: Charlie
type SendCodeParam struct {
	Target string `json:"target" binding:"required"`
}
