package config

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/framework/core/bind"
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
		api.POST("/v1/admin/sys/config/create", admin, middleware.RequirePermission(d.Perms, "sys:config:create", "配置创建"), s.create)
		api.POST("/v1/admin/sys/config/update", admin, middleware.RequirePermission(d.Perms, "sys:config:update", "配置更新"), s.update)
		api.POST("/v1/admin/sys/config/delete", admin, middleware.RequirePermission(d.Perms, "sys:config:delete", "配置删除"), s.delete)
		api.GET("/v1/admin/sys/config/detail", admin, middleware.RequirePermission(d.Perms, "sys:config:detail", "配置详情"), s.detail)
		api.GET("/v1/admin/sys/config/page", admin, middleware.RequirePermission(d.Perms, "sys:config:page", "配置分页"), s.page)
		api.GET("/v1/admin/sys/config/list", admin, middleware.RequirePermission(d.Perms, "sys:config:page", "配置列表"), s.list)
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
	q.Category = c.Query("category")
	q.Scope = c.Query("scope")
	rows, err := s.List(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}
