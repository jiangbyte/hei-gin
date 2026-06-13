package v1

import (
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers SM2-related routes.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/public/b/sm2/public-key", getPublicKeyHandler)
}

func Register() {
	registry.RegisterRoute(RegisterRoutes)
}

// getPublicKeyHandler returns the SM2 public key for frontend encryption.
// @Summary      后台认证获取公钥
// @Description  访问 /api/v1/public/b/sm2/public-key，后台认证获取公钥
// @Tags         后台认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/b/sm2/public-key [get]
func getPublicKeyHandler(c *gin.Context) {
	publicKey := utils.GetPublicKey()
	result.Success(c, publicKey)
}
