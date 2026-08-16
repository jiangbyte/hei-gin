// internal/modules/iam/client/handler.go HTTP 处理器。
//
// Author: Charlie

package client

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		// 操作审计登记（对齐 hei-boot @OperationAudit：iam_clientmodule / iam_clientresource）
		d.AuditReg.RegisterSpecs(
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/client-modules/create", ResourceType: "iam_clientmodule", Action: "create"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/client-modules/update", ResourceType: "iam_clientmodule", Action: "update"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/client-modules/delete", ResourceType: "iam_clientmodule", Action: "delete"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/client-resources/create", ResourceType: "iam_clientresource", Action: "create"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/client-resources/update", ResourceType: "iam_clientresource", Action: "update"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/client-resources/delete", ResourceType: "iam_clientresource", Action: "delete"},
		)
		api.POST("/v1/admin/sys/client-modules/create", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:create", "客户端模块创建"), s.createModule)
		api.POST("/v1/admin/sys/client-modules/update", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:update", "客户端模块更新"), s.updateModule)
		api.POST("/v1/admin/sys/client-modules/delete", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:delete", "客户端模块删除"), s.deleteModule)
		api.GET("/v1/admin/sys/client-modules/detail", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:detail", "客户端模块详情"), s.detailModule)
		api.GET("/v1/admin/sys/client-modules/page", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:page", "客户端模块分页"), s.pageModule)
		api.GET("/v1/admin/sys/client-modules/selector", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:page", "客户端模块选择"), s.selectorModule)
		api.POST("/v1/admin/sys/client-resources/create", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:create", "客户端资源创建"), s.createResource)
		api.POST("/v1/admin/sys/client-resources/update", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:update", "客户端资源更新"), s.updateResource)
		api.POST("/v1/admin/sys/client-resources/delete", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:delete", "客户端资源删除"), s.deleteResource)
		api.GET("/v1/admin/sys/client-resources/detail", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:detail", "客户端资源详情"), s.detailResource)
		api.GET("/v1/admin/sys/client-resources/page", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:page", "客户端资源分页"), s.pageResource)
		api.GET("/v1/admin/sys/client-resources/tree", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:list", "客户端资源树"), s.treeResource)
	}
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
	out, err := s.ModuleSelector(c.Request.Context(), c.Query("account_type"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) createResource(c *gin.Context) {
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

func (s *Service) updateResource(c *gin.Context) {
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

func (s *Service) deleteResource(c *gin.Context) {
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

func (s *Service) detailResource(c *gin.Context) {
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

func (s *Service) pageResource(c *gin.Context) {
	var q ResourcePageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.ResourcePage(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (s *Service) treeResource(c *gin.Context) {
	nodes, err := s.ResourceTree(c.Request.Context(), c.Query("module_id"), c.Query("account_type"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}
