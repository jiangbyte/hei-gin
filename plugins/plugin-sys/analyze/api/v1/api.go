package v1

import (
	"hei-gin/sdk/registry"
	"hei-gin/sdk/result"
	analyze "hei-gin/plugins/plugin-sys/analyze"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/analyze/dashboard
	r.GET("/api/v1/sys/analyze/dashboard", dashboardHandler)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// dashboardHandler handles GET /api/v1/sys/analyze/dashboard
func dashboardHandler(c *gin.Context) {
	data := analyze.AnalyzeDashboard(c)
	result.Success(c, data)
}
