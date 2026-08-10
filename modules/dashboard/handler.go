package dashboard

import (
	"github.com/gin-gonic/gin"

	"hei-gin/framework/core/response"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
	"hei-gin/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) func(*gin.RouterGroup) {
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
