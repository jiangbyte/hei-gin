package v1

import (
	"hei-gin/sdk/captcha"
	"hei-gin/sdk/registry"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers captcha-related routes.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/public/b/captcha", getCaptchaHandler)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// getCaptchaHandler generates a captcha image and returns it as a base64-encoded string.
// @Summary      后台认证获取验证码
// @Description  访问 /api/v1/public/b/captcha，后台认证获取验证码
// @Tags         后台认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/b/captcha [get]
func getCaptchaHandler(c *gin.Context) {
	captchaResult, err := captcha.BCaptcha.GetCaptcha()
	if err != nil {
		result.Failure(c, "验证码生成失败", 500)
		return
	}
	result.Success(c, captchaResult)
}
