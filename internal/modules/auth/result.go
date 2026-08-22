// internal/modules/auth/result.go 出参定义。
//
// Author: Charlie

package auth

import (
	"hei-gin/internal/framework/core/security"
)

// LoginResult 登录结果（对齐 hei-boot 登录响应：含安全墙字段）。
//
// Author: Charlie
type LoginResult struct {
	Token                     string               `json:"token"`
	AccountID                 string               `json:"account_id"`
	AccountType               security.AccountType `json:"account_type"`
	PasswordExpired           bool                 `json:"password_expired"`
	ForceBindEmail            bool                 `json:"force_bind_email"`
	ForceBindPhone            bool                 `json:"force_bind_phone"`
	PasswordExpiryWarningDays int                  `json:"password_expiry_warning_days"`
}

// RegisterResult 注册结果。
//
// Author: Charlie
type RegisterResult struct {
	AccountID   string               `json:"account_id"`
	Account     string               `json:"account"`
	AccountType security.AccountType `json:"account_type"`
}

// CaptchaResult 验证码结果。
//
// Author: Charlie
type CaptchaResult struct {
	CaptchaID   string `json:"captcha_id"`
	ImageBase64 string `json:"image_base64"`
	ImageType   string `json:"image_type"`
}

// PasswordKeyResult 密码加密公钥结果。
//
// Author: Charlie
type PasswordKeyResult struct {
	KeyID     string `json:"key_id"`
	PublicKey string `json:"public_key"`
}

// LogoutResult 登出结果。
//
// Author: Charlie
type LogoutResult struct {
	Success bool `json:"success"`
}

// SiteFooterResult 站点页脚公开信息：版权与备案。
//
// Author: Charlie
type SiteFooterResult struct {
	CopyrightText string `json:"copyright_text"`
	CopyrightURL  string `json:"copyright_url"`
	IcpNumber     string `json:"icp_number"`
	IcpURL        string `json:"icp_url"`
	PsbNumber     string `json:"psb_number"`
	PsbURL        string `json:"psb_url"`
}
