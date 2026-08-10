package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/framework/core/response"
	"hei-gin/framework/core/schema"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET("/v1/admin/sys/audit/page", admin, middleware.RequirePermission(d.Perms, "sys:audit:page", "审计分页"), s.page)
		api.GET("/v1/admin/sys/audit/detail", admin, middleware.RequirePermission(d.Perms, "sys:audit:detail", "审计详情"), s.detail)
	}
}

func (s *Service) page(c *gin.Context) {
	var q PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.Page(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (s *Service) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}
