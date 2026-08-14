package cg_test_knowledge_category

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
		g := api.Group("/v1/admin/biz/cg-test-knowledge-category", middleware.RequireAccountType(security.AccountAdmin))
		g.POST("/create", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:create", "Create knowledge category"), s.create)
		g.POST("/update", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:update", "Update knowledge category"), s.update)
		g.POST("/delete", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:delete", "Delete knowledge category"), s.delete)
		g.GET("/detail", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:detail", "Knowledge category detail"), s.detail)
		g.GET("/page", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:page", "Knowledge category page"), s.page)
		g.GET("/tree", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:tree", "Knowledge category tree"), s.tree)

		ch := g.Group("/children")
		ch.POST("/create", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:children:create", "Create knowledge doc"), s.createDoc)
		ch.POST("/update", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:children:update", "Update knowledge doc"), s.updateDoc)
		ch.POST("/delete", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:children:delete", "Delete knowledge doc"), s.deleteDoc)
		ch.GET("/detail", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:children:detail", "Knowledge doc detail"), s.detailDoc)
		ch.GET("/page", middleware.RequirePermission(d.Perms, "biz:cgtestknowledgecategory:children:page", "Knowledge doc page"), s.pageDoc)
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
	rows, total, cur, size, err := s.Page(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (s *Service) tree(c *gin.Context) {
	out, err := s.Tree(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) createDoc(c *gin.Context) {
	var req DocAddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.CreateDoc(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) updateDoc(c *gin.Context) {
	var req DocEditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateDoc(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) deleteDoc(c *gin.Context) {
	var req IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.DeleteDocs(c.Request.Context(), req.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) detailDoc(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.DetailDoc(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}

func (s *Service) pageDoc(c *gin.Context) {
	var q DocPageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.PageDocs(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}
