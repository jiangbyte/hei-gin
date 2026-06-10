package captcha_api

import (
	"hei-gin/sdk/captcha"
	"hei-gin/sdk/registry"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers consumer captcha-related routes.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/public/c/captcha", getCaptchaHandler)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// getCaptchaHandler generates a captcha image and returns it as a base64-encoded string.
func getCaptchaHandler(c *gin.Context) {
	captchaResult, err := captcha.CCaptcha.GetCaptcha()
	if err != nil {
		result.Failure(c, "验证码生成失败", 500)
		return
	}
	result.Success(c, captchaResult)
}
