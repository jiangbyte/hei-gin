package weakpassword

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
		api.POST("/v1/admin/sys/weak-password/create", admin, middleware.RequirePermission(d.Perms, "sys:weak-password:create", "弱密码创建"), s.create)
		api.POST("/v1/admin/sys/weak-password/update", admin, middleware.RequirePermission(d.Perms, "sys:weak-password:update", "弱密码更新"), s.update)
		api.POST("/v1/admin/sys/weak-password/delete", admin, middleware.RequirePermission(d.Perms, "sys:weak-password:delete", "弱密码删除"), s.delete)
		api.GET("/v1/admin/sys/weak-password/detail", admin, middleware.RequirePermission(d.Perms, "sys:weak-password:detail", "弱密码详情"), s.detail)
		api.GET("/v1/admin/sys/weak-password/page", admin, middleware.RequirePermission(d.Perms, "sys:weak-password:page", "弱密码分页"), s.page)
		api.GET("/v1/admin/sys/weak-password/list", admin, middleware.RequirePermission(d.Perms, "sys:weak-password:list", "弱密码列表"), s.list)
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
	q.Password = c.Query("password")
	rows, err := s.List(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}
