package captcha_api

import (
	"github.com/gin-gonic/gin"

	"hei-gin/sdk/captcha"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
)

// RegisterRoutes registers captcha-related routes.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/public/b/captcha", GetCaptcha)
}

// GetCaptcha generates a captcha image and returns it as a base64-encoded string.
func GetCaptcha(c *gin.Context) {
	captchaResult, err := captcha.BCaptcha.GetCaptcha()
	if err != nil {
		result.Failure(c, "验证码生成失败", 500)
		return
	}
	result.Success(c, captchaResult)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
