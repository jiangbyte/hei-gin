package username_api

import (
	"hei-gin/plugins/plugin-sys/auth/username"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/registry"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers username-based auth routes (login/register/logout).
func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/public/b/login", username.DoLogin)
	r.POST("/api/v1/public/b/register",
		log.SysLog("注册"),
		middleware.NoRepeat(5000),
		username.DoRegister,
	)
	r.POST("/api/v1/b/logout",
		middleware.HeiCheckLogin(),
		username.DoLogout,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}
