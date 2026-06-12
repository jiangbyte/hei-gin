package v1

import (
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers consumer SM2-related routes.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/public/c/sm2/public-key", getPublicKeyHandler)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// getPublicKeyHandler returns the SM2 public key for frontend encryption.
// @Summary      C端认证获取公钥
// @Description  访问 /api/v1/public/c/sm2/public-key，C端认证获取公钥
// @Tags         C端认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/c/sm2/public-key [get]
func getPublicKeyHandler(c *gin.Context) {
	publicKey := utils.GetPublicKey()
	result.Success(c, publicKey)
}
