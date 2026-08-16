// internal/modules/sys/banner/handler.go HTTP 处理器。
//
// Author: Charlie

package banner

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/module"
)

func (s *Service) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET("/v1/admin/sys/banners/list", admin, s.list)
		api.POST("/v1/admin/sys/banners/create", admin, middleware.RequirePermission(d.Perms, "sys:banner:create", "Banner创建"), middleware.OperationAudit(d.Audit, "sys_banner", "create"), s.create)
		api.POST("/v1/admin/sys/banners/update", admin, middleware.RequirePermission(d.Perms, "sys:banner:update", "Banner更新"), middleware.OperationAudit(d.Audit, "sys_banner", "update"), s.update)
		api.POST("/v1/admin/sys/banners/delete", admin, middleware.RequirePermission(d.Perms, "sys:banner:delete", "Banner删除"), middleware.OperationAudit(d.Audit, "sys_banner", "delete"), s.delete)
		api.GET("/v1/admin/sys/banners/detail", admin, middleware.RequirePermission(d.Perms, "sys:banner:detail", "Banner详情"), s.detail)
		api.GET("/v1/admin/sys/banners/page", admin, middleware.RequirePermission(d.Perms, "sys:banner:page", "Banner分页"), s.page)

		// 门户 banner 列表与互动为公开接口（对齐 hei-boot PortalBannerController，web public:true）
		api.GET("/v1/portal/sys/banners/list", s.portalList)
		api.POST("/v1/portal/sys/banners/interaction", middleware.OperationAudit(d.Audit, "sys_banner", "interaction"), s.interaction)
	}
}

func (s *Service) create(c *gin.Context) {
	var req AddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Create(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) update(c *gin.Context) {
	var req EditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Update(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) delete(c *gin.Context) {
	var body IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Delete(c.Request.Context(), body.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
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

func (s *Service) list(c *gin.Context) {
	var q ListParam
	q.Position = c.Query("position")
	q.Category = c.Query("category")
	q.Type = c.Query("type")
	rows, err := s.List(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) portalList(c *gin.Context) {
	var q PortalListParam
	q.Position = c.Query("position")
	q.Category = c.Query("category")
	q.Type = c.Query("type")
	rows, err := s.PortalList(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) interaction(c *gin.Context) {
	var req InteractionParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Interaction(c.Request.Context(), req.ID); err != nil {
		response.Fail(c, http.StatusNotFound, 404, "banner not found")
		return
	}
	response.OK(c, nil)
}
