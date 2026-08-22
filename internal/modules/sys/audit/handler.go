// internal/modules/sys/audit/handler.go HTTP 处理器。
//
// Author: Charlie

package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/module"
)

func (s *Service) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET("/v1/admin/sys/audit/page", admin, middleware.RequirePermission(d.Perms, "sys:audit:page", "审计分页"), s.page)
		api.GET("/v1/admin/sys/audit/my-page", admin, s.myPage)
		api.GET("/v1/admin/sys/audit/detail", admin, middleware.RequirePermission(d.Perms, "sys:audit:detail", "审计详情"), s.detail)
		api.GET("/v1/admin/sys/audit/my-detail", admin, s.myDetail)

		portal := middleware.RequireAccountType(security.AccountPortal)
		api.GET("/v1/portal/sys/audit/my-page", portal, s.myPage)
		api.GET("/v1/portal/sys/audit/my-detail", portal, s.myDetail)
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

func (s *Service) myPage(c *gin.Context) {
	accountID := contextx.AccountID(c.Request.Context())
	if accountID == "" {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	var q PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.MyPage(c.Request.Context(), accountID, q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (s *Service) myDetail(c *gin.Context) {
	accountID := contextx.AccountID(c.Request.Context())
	if accountID == "" {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.MyDetail(c.Request.Context(), accountID, q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}
