// internal/modules/workspace/handler.go HTTP 处理器。
//
// Author: Charlie

package workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/module"
)

func (s *Service) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET(
			"/v1/admin/workspace/overview",
			admin,
			middleware.RequirePermission(d.Perms, "workspace:overview:view", "工作台总览"),
			s.overview,
		)
		api.GET("/v1/admin/workspace/shortcuts", admin, s.listShortcuts)
		api.POST(
			"/v1/admin/workspace/shortcuts",
			admin,
			middleware.OperationAudit(d.Audit, "workspace_shortcut", "update"),
			s.replaceShortcuts,
		)
	}
}

func (s *Service) overview(c *gin.Context) {
	out, err := s.Overview(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) listShortcuts(c *gin.Context) {
	out, err := s.ListShortcuts(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) replaceShortcuts(c *gin.Context) {
	var req ShortcutSaveParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := s.ReplaceShortcuts(c.Request.Context(), req.ResourceIDs)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}
