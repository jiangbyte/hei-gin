package cg_test_activity

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/framework/core/bind"
	contextx "hei-gin/framework/core/context"
	"hei-gin/framework/core/response"
	"hei-gin/framework/core/schema"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
	"hei-gin/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) func(*gin.RouterGroup) {
	return func(api *gin.RouterGroup) {
		g := api.Group("/v1/admin/biz/cg-test-activity", middleware.RequireAccountType(security.AccountAdmin))
		g.POST("/create", middleware.RequirePermission(d.Perms, "biz:cgtestactivity:create", "Create activity"), s.create)
		g.POST("/update", middleware.RequirePermission(d.Perms, "biz:cgtestactivity:update", "Update activity"), s.update)
		g.POST("/delete", middleware.RequirePermission(d.Perms, "biz:cgtestactivity:delete", "Delete activity"), s.delete)
		g.GET("/detail", middleware.RequirePermission(d.Perms, "biz:cgtestactivity:detail", "Activity detail"), s.detail)
		g.GET("/page", middleware.RequirePermission(d.Perms, "biz:cgtestactivity:page", "Activity page"), s.page)
	}
}

func (s *Service) create(c *gin.Context) {
	var req AddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Create(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
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
	if err := s.Update(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) delete(c *gin.Context) {
	var req IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Delete(c.Request.Context(), req.IDs); err != nil {
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
	rows, total, cur, size, err := s.Page(c.Request.Context(), q, contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}
