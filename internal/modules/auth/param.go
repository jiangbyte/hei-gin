// internal/modules/auth/param.go 入参定义。
//
// Author: Charlie

package auth

// LoginParam 登录入参。
//
// Author: Charlie
type LoginParam struct {
	Account       string `json:"account" binding:"required,min=3,max=128"`
	Password      string `json:"password"`
	PasswordKeyID string `json:"password_key_id"`
	IdentityType  string `json:"identity_type"`
	LoginMode     string `json:"login_mode"`
	OTPCode       string `json:"otp_code"`
	RememberMe    *bool  `json:"remember_me"`
	CaptchaID     string `json:"captcha_id" binding:"required"`
	CaptchaValue  string `json:"captcha_value" binding:"required"`
}

// RegisterParam 门户注册入参（对齐 hei-boot RegisterParam）。
type RegisterParam struct {
	RegisterChannel string `json:"register_channel"`
	Account         string `json:"account"`
	Password        string `json:"password" binding:"required"`
	PasswordKeyID   string `json:"password_key_id" binding:"required"`
	Name            string `json:"name"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	OTPCode         string `json:"otp_code"`
	CaptchaID       string `json:"captcha_id" binding:"required"`
	CaptchaValue    string `json:"captcha_value" binding:"required"`
}

// SendLoginCodeParam 发送登录 OTP。
//
// Author: Charlie
type SendLoginCodeParam struct {
	Account      string `json:"account"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Target       string `json:"target"`
	Channel      string `json:"channel"`
	CaptchaID    string `json:"captcha_id" binding:"required"`
	CaptchaValue string `json:"captcha_value" binding:"required"`
}

// ForgotPasswordParam 忘记密码。
//
// Author: Charlie
type ForgotPasswordParam struct {
	Email        string `json:"email" binding:"required,email"`
	CaptchaID    string `json:"captcha_id" binding:"required"`
	CaptchaValue string `json:"captcha_value" binding:"required"`
}

// ResetPasswordParam 重置密码。
//
// Author: Charlie
type ResetPasswordParam struct {
	Email         string `json:"email"`
	Token         string `json:"token" binding:"required"`
	Password      string `json:"password" binding:"required"`
	PasswordKeyID string `json:"password_key_id" binding:"required"`
	CaptchaID     string `json:"captcha_id" binding:"required"`
	CaptchaValue  string `json:"captcha_value" binding:"required"`
}

// ForgotPasswordByPhoneParam 通过手机找回密码。
//
// Author: Charlie
type ForgotPasswordByPhoneParam struct {
	Phone        string `json:"phone" binding:"required"`
	CaptchaID    string `json:"captcha_id" binding:"required"`
	CaptchaValue string `json:"captcha_value" binding:"required"`
}

// ResetPasswordByPhoneParam 通过手机 OTP 重置密码。
//
// Author: Charlie
type ResetPasswordByPhoneParam struct {
	Phone         string `json:"phone" binding:"required"`
	OTPCode       string `json:"otp_code" binding:"required"`
	Password      string `json:"password" binding:"required"`
	PasswordKeyID string `json:"password_key_id" binding:"required"`
	CaptchaID     string `json:"captcha_id" binding:"required"`
	CaptchaValue  string `json:"captcha_value" binding:"required"`
}
