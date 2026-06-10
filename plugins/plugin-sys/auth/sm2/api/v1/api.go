package sm2_api

import (
	"hei-gin/sdk/registry"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers SM2-related routes.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/public/b/sm2/public-key", getPublicKeyHandler)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// getPublicKeyHandler returns the SM2 public key for frontend encryption.
func getPublicKeyHandler(c *gin.Context) {
	publicKey := utils.GetPublicKey()
	result.Success(c, publicKey)
}
