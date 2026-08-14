// internal/modules/iam/resource/handler.go HTTP 处理器。
//
// Author: Charlie

package resource

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/resources/create", admin, middleware.RequirePermission(d.Perms, "iam:resource:create", "资源创建"), s.create)
		api.POST("/v1/admin/sys/resources/update", admin, middleware.RequirePermission(d.Perms, "iam:resource:update", "资源更新"), s.update)
		api.POST("/v1/admin/sys/resources/delete", admin, middleware.RequirePermission(d.Perms, "iam:resource:delete", "资源删除"), s.delete)
		api.GET("/v1/admin/sys/resources/detail", admin, middleware.RequirePermission(d.Perms, "iam:resource:detail", "资源详情"), s.detail)
		api.GET("/v1/admin/sys/resources/page", admin, middleware.RequirePermission(d.Perms, "iam:resource:page", "资源分页"), s.page)
		api.GET("/v1/admin/sys/resources/current", admin, s.currentAdmin)
		api.GET("/v1/admin/sys/resources/tree", admin, middleware.RequirePermission(d.Perms, "iam:resource:list", "资源树"), s.tree)
		api.POST("/v1/admin/sys/resource-modules/create", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:create", "资源模块创建"), s.createModule)
		api.POST("/v1/admin/sys/resource-modules/update", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:update", "资源模块更新"), s.updateModule)
		api.POST("/v1/admin/sys/resource-modules/delete", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:delete", "资源模块删除"), s.deleteModule)
		api.GET("/v1/admin/sys/resource-modules/detail", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:detail", "资源模块详情"), s.detailModule)
		api.GET("/v1/admin/sys/resource-modules/page", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:page", "资源模块分页"), s.pageModule)
		api.GET("/v1/admin/sys/resource-modules/selector", admin, s.selectorModule)
		api.GET("/v1/portal/sys/resources/current", s.currentPortal)
		api.GET("/v1/admin/permission-registry", admin, middleware.RequirePermission(d.Perms, "iam:resource:grant", "资源授权"), s.permissionRegistry)
		api.POST("/v1/admin/resource-permissions", admin, middleware.RequirePermission(d.Perms, "iam:resource:grant", "资源授权"), s.bindResourcePermissions)
		api.POST("/v1/admin/client-resource-permissions", admin, middleware.RequirePermission(d.Perms, "iam:resource:grant", "资源授权"), s.bindClientResourcePermissions)
		api.POST("/v1/admin/sys/resource-buttons/create", admin, middleware.RequirePermission(d.Perms, "iam:resource:create", "资源创建"), s.createButton)
		api.POST("/v1/admin/sys/resource-buttons/update", admin, middleware.RequirePermission(d.Perms, "iam:resource:update", "资源更新"), s.updateButton)
		api.POST("/v1/admin/sys/resource-buttons/delete", admin, middleware.RequirePermission(d.Perms, "iam:resource:delete", "资源删除"), s.deleteButton)
		api.GET("/v1/admin/sys/resource-buttons/page", admin, middleware.RequirePermission(d.Perms, "iam:resource:list", "资源分页"), s.pageButton)
	}
}

func (s *Service) create(c *gin.Context) {
	var req ResourceAddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.CreateResource(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) update(c *gin.Context) {
	var req ResourceEditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateResource(c.Request.Context(), req); err != nil {
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
	if err := s.DeleteResources(c.Request.Context(), body.IDs); err != nil {
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
	row, err := s.ResourceDetail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}

func (s *Service) page(c *gin.Context) {
	var q ResourcePageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.ResourcePage(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (s *Service) currentAdmin(c *gin.Context) {
	rows, err := s.CurrentAdmin(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) currentPortal(c *gin.Context) {
	rows, err := s.CurrentPortal(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) tree(c *gin.Context) {
	nodes, err := s.ResourceTree(c.Request.Context(), c.Query("module_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}

func (s *Service) createModule(c *gin.Context) {
	var req ModuleAddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.CreateModule(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) updateModule(c *gin.Context) {
	var req ModuleEditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateModule(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) deleteModule(c *gin.Context) {
	var body IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.DeleteModules(c.Request.Context(), body.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) detailModule(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.ModuleDetail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}

func (s *Service) pageModule(c *gin.Context) {
	var q ModulePageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.ModulePage(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (s *Service) selectorModule(c *gin.Context) {
	out, err := s.ModuleSelector(c.Request.Context(), c.Query("client"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}
func (s *Service) permissionRegistry(c *gin.Context) {
	response.OK(c, s.ListPermissions())
}

func (s *Service) bindResourcePermissions(c *gin.Context) {
	var req ResourcePermissionBindParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.BindResourcePermissions(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) bindClientResourcePermissions(c *gin.Context) {
	var req ResourcePermissionBindParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.BindClientResourcePermissions(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) createButton(c *gin.Context) {
	var req ButtonAddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.CreateButton(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) updateButton(c *gin.Context) {
	var req ButtonEditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateButton(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) deleteButton(c *gin.Context) {
	var body IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.DeleteButtons(c.Request.Context(), body.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) pageButton(c *gin.Context) {
	var q ButtonPageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.PageButtons(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}
