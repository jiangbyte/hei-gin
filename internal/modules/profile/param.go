// internal/modules/profile/param.go 用户中心入参定义（对齐 hei-boot 各 UpdateParam）。
//
// Author: Charlie

package profile

// ProfileUpdateParam 资料更新入参。
//
// Author: Charlie
type ProfileUpdateParam struct {
	Nickname  *string `json:"nickname"`
	Avatar    *string `json:"avatar"`
	Signature *string `json:"signature"`
	Remark    *string `json:"remark"`
}

// PasswordUpdateParam 密码更新入参。
//
// Author: Charlie
type PasswordUpdateParam struct {
	PasswordKeyID string `json:"password_key_id"`
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"new_password" binding:"required"`
	OTPCode       string `json:"otp_code"`
}

// PhoneUpdateParam 手机号更新入参。
//
// Author: Charlie
type PhoneUpdateParam struct {
	PasswordKeyID     string  `json:"password_key_id"`
	Password          string  `json:"password"`
	Phone             *string `json:"phone"`
	PhoneLoginEnabled *bool   `json:"phone_login_enabled"`
	OTPCode           string  `json:"otp_code"`
}

// EmailUpdateParam 邮箱更新入参。
//
// Author: Charlie
type EmailUpdateParam struct {
	PasswordKeyID     string  `json:"password_key_id"`
	Password          string  `json:"password"`
	Email             *string `json:"email"`
	EmailLoginEnabled *bool   `json:"email_login_enabled"`
	OTPCode           string  `json:"otp_code"`
}

// SendCodeParam 绑定验证码发送入参。
//
// Author: Charlie
type SendCodeParam struct {
	Target string `json:"target" binding:"required"`
}
