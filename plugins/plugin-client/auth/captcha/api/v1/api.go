package v1

import (
	"hei-gin/sdk/captcha"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/web/result"

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
// @Summary      C端认证获取验证码
// @Description  访问 /api/v1/public/c/captcha，C端认证获取验证码
// @Tags         C端认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/c/captcha [get]
func getCaptchaHandler(c *gin.Context) {
	captchaResult, err := captcha.CCaptcha.GetCaptcha()
	if err != nil {
		result.Failure(c, "验证码生成失败", 500)
		return
	}
	result.Success(c, captchaResult)
}
