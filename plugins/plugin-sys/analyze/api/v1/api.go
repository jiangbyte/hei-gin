package v1

import (
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	analyze "hei-gin/plugins/plugin-sys/analyze"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/analyze/dashboard — NO permission middleware
	r.GET("/api/v1/sys/analyze/dashboard", dashboard)
}

func dashboard(c *gin.Context) {
	data := analyze.Dashboard(c)
	result.Success(c, data)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
