// internal/modules/dashboard/handler.go HTTP 处理器。
//
// Author: Charlie

package dashboard

import (
	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/module"
)

func (s *Service) registerRoutes(d *module.Deps) func(*gin.RouterGroup) {
	return func(api *gin.RouterGroup) {
		api.GET(
			"/v1/admin/dashboard/overview",
			middleware.RequireAccountType(security.AccountAdmin),
			middleware.RequirePermission(d.Perms, "dashboard:overview:view", "Dashboard overview"),
			s.overview,
		)
	}
}

func (s *Service) overview(c *gin.Context) {
	response.OK(c, s.Overview(c.Request.Context()))
}
