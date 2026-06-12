package v1

import (
	analyze "hei-gin/plugins/plugin-sys/analyze"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *analyze.Service
}

var defaultHandler = newHandler(analyze.DefaultModule)

func newHandler(module *analyze.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/analyze/dashboard
	r.GET("/api/v1/sys/analyze/dashboard", defaultHandler.dashboard)
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
func (h *handler) dashboard(c *gin.Context) {
	data := h.service.Dashboard(c)
	result.Success(c, data)
}
