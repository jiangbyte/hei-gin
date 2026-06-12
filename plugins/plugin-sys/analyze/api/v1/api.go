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
// @Summary      系统分析仪表盘数据
// @Description  访问 /api/v1/sys/analyze/dashboard，系统分析仪表盘数据
// @Tags         系统分析
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/analyze/dashboard [get]
func dashboardHandler(c *gin.Context) {
	data := analyze.AnalyzeDashboard(c)
	result.Success(c, data)
}
