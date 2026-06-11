package v1

import (
	"hei-gin/plugins/plugin-client/auth/username"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/registry"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers consumer username-based auth routes (login/register/logout).
func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/public/c/login", username.DoLogin)
	r.POST("/api/v1/public/c/register",
		log.SysLog("注册"),
		middleware.NoRepeat(5000),
		username.DoRegister,
	)
	r.POST("/api/v1/c/logout",
		middleware.HeiCheckLogin("CONSUMER"),
		username.DoLogout,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}
