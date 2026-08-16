// internal/modules/sys/weak_password/handler.go HTTP 处理器。
//
// Author: Charlie

package weakpassword

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
		// 操作审计登记（对齐 hei-boot @OperationAudit：sys_weakpassword）
		d.AuditReg.RegisterSpecs(
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/weak-password/create", ResourceType: "sys_weakpassword", Action: "create"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/weak-password/update", ResourceType: "sys_weakpassword", Action: "update"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/weak-password/delete", ResourceType: "sys_weakpassword", Action: "delete"},
		)
		api.POST("/v1/admin/sys/weak-password/create", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:create", "弱密码创建"), s.create)
		api.POST("/v1/admin/sys/weak-password/update", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:update", "弱密码更新"), s.update)
		api.POST("/v1/admin/sys/weak-password/delete", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:delete", "弱密码删除"), s.delete)
		api.GET("/v1/admin/sys/weak-password/detail", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:detail", "弱密码详情"), s.detail)
		api.GET("/v1/admin/sys/weak-password/page", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:page", "弱密码分页"), s.page)
		api.GET("/v1/admin/sys/weak-password/list", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:list", "弱密码列表"), s.list)
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
