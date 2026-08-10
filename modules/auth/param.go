package auth

// LoginParam 登录入参。
//
// Author: Charlie
type LoginParam struct {
	Account       string `json:"account" binding:"required,min=3,max=128"`
	Password      string `json:"password"`
	PasswordKeyID string `json:"password_key_id" binding:"required"`
	IdentityType  string `json:"identity_type"`
	LoginMode     string `json:"login_mode"`
	OTPCode       string `json:"otp_code"`
	RememberMe    bool   `json:"remember_me"`
	CaptchaID     string `json:"captcha_id" binding:"required"`
	CaptchaValue  string `json:"captcha_value" binding:"required"`
}

// RegisterParam 门户注册入参。
//
// Author: Charlie
type RegisterParam struct {
	Account       string `json:"account" binding:"required,min=3,max=64"`
	Password      string `json:"password" binding:"required"`
	PasswordKeyID string `json:"password_key_id" binding:"required"`
	Name          string `json:"name"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	CaptchaID     string `json:"captcha_id" binding:"required"`
	CaptchaValue  string `json:"captcha_value" binding:"required"`
}
