package auth

import (
	"hei-gin/internal/framework/core/security"
)

// LoginResult ç™»å½•ç»“æžœã€‚
//
// Author: Charlie
type LoginResult struct {
	Token           string               `json:"token"`
	AccountID       string               `json:"account_id"`
	AccountType     security.AccountType `json:"account_type"`
	PasswordExpired bool                 `json:"password_expired"`
}

// RegisterResult æ³¨å†Œç»“æžœã€‚
//
// Author: Charlie
type RegisterResult struct {
	AccountID   string               `json:"account_id"`
	Account     string               `json:"account"`
	AccountType security.AccountType `json:"account_type"`
}

// CaptchaResult éªŒè¯ç ç»“æžœã€‚
//
// Author: Charlie
type CaptchaResult struct {
	CaptchaID   string `json:"captcha_id"`
	ImageBase64 string `json:"image_base64"`
	ImageType   string `json:"image_type"`
}

// PasswordKeyResult å¯†ç åŠ å¯†å…¬é’¥ç»“æžœã€‚
//
// Author: Charlie
type PasswordKeyResult struct {
	KeyID     string `json:"key_id"`
	PublicKey string `json:"public_key"`
}

// LogoutResult ç™»å‡ºç»“æžœã€‚
//
// Author: Charlie
type LogoutResult struct {
	Success bool `json:"success"`
}
