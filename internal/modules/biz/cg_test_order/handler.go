// internal/modules/biz/cg_test_order/handler.go HTTP 处理器。
//
// Author: Charlie

package cg_test_order

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/module"
)

func (s *Service) registerRoutes(d *module.Deps) func(*gin.RouterGroup) {
	return func(api *gin.RouterGroup) {
		g := api.Group("/v1/admin/biz/cg-test-order", middleware.RequireAccountType(security.AccountAdmin))
		g.POST("/create", middleware.RequirePermission(d.Perms, "biz:cgtestorder:create", "Create order"), s.create)
		g.POST("/update", middleware.RequirePermission(d.Perms, "biz:cgtestorder:update", "Update order"), s.update)
		g.POST("/delete", middleware.RequirePermission(d.Perms, "biz:cgtestorder:delete", "Delete order"), s.delete)
		g.GET("/detail", middleware.RequirePermission(d.Perms, "biz:cgtestorder:detail", "Order detail"), s.detail)
		g.GET("/page", middleware.RequirePermission(d.Perms, "biz:cgtestorder:page", "Order page"), s.page)

		ch := g.Group("/children")
		ch.POST("/create", middleware.RequirePermission(d.Perms, "biz:cgtestorder:create", "Create order item"), s.createItem)
		ch.POST("/update", middleware.RequirePermission(d.Perms, "biz:cgtestorder:update", "Update order item"), s.updateItem)
		ch.POST("/delete", middleware.RequirePermission(d.Perms, "biz:cgtestorder:delete", "Delete order item"), s.deleteItem)
		ch.GET("/detail", middleware.RequirePermission(d.Perms, "biz:cgtestorder:detail", "Order item detail"), s.detailItem)
		ch.GET("/page", middleware.RequirePermission(d.Perms, "biz:cgtestorder:page", "Order item page"), s.pageItem)
	}
}

func (s *Service) create(c *gin.Context) {
	var req AddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Create(c.Request.Context(), contextx.AccountID(c.Request.Context()), req, contextx.Session(c.Request.Context())); err != nil {
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

func (s *Service) createItem(c *gin.Context) {
	var req ItemAddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.CreateItem(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) updateItem(c *gin.Context) {
	var req ItemEditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateItem(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) deleteItem(c *gin.Context) {
	var req IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.DeleteItems(c.Request.Context(), req.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) detailItem(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.DetailItem(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}

func (s *Service) pageItem(c *gin.Context) {
	var q ItemPageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.PageItems(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}
