// internal/modules/iam/role/handler.go HTTP 处理器。
//
// Author: Charlie

package role

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
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
		// 操作审计登记（对齐 hei-boot @OperationAudit：iam_role）
		d.AuditReg.RegisterSpecs(
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/roles/create", ResourceType: "iam_role", Action: "create"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/roles/update", ResourceType: "iam_role", Action: "update"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/roles/delete", ResourceType: "iam_role", Action: "delete"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/roles/grant-resource", ResourceType: "iam_role", Action: "grant_resource"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/roles/grant-client-resource", ResourceType: "iam_role", Action: "grant_client_resource"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/roles/grant-user", ResourceType: "iam_role", Action: "grant_user"},
		)
		api.POST("/v1/admin/sys/roles/create", admin, middleware.RequirePermission(d.Perms, "iam:role:create", "角色创建"), s.create)
		api.POST("/v1/admin/sys/roles/update", admin, middleware.RequirePermission(d.Perms, "iam:role:update", "角色更新"), s.update)
		api.POST("/v1/admin/sys/roles/delete", admin, middleware.RequirePermission(d.Perms, "iam:role:delete", "角色删除"), s.delete)
		api.GET("/v1/admin/sys/roles/detail", admin, middleware.RequirePermission(d.Perms, "iam:role:detail", "角色详情"), s.detail)
		api.GET("/v1/admin/sys/roles/page", admin, middleware.RequirePermission(d.Perms, "iam:role:page", "角色分页"), s.page)
		api.GET("/v1/admin/sys/roles/own-resource", admin, middleware.RequirePermission(d.Perms, "iam:role:ownresource", "角色已拥有资源"), s.ownResource)
		api.POST("/v1/admin/sys/roles/grant-resource", admin, middleware.RequirePermission(d.Perms, "iam:role:grantresource", "角色资源授权"), s.grantResource)
		api.GET("/v1/admin/sys/roles/own-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:role:ownclientresource", "角色已拥有客户端资源"), s.ownClientResource)
		api.POST("/v1/admin/sys/roles/grant-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:role:grantclientresource", "角色客户端资源授权"), s.grantClientResource)
		api.GET("/v1/admin/sys/roles/own-user", admin, middleware.RequirePermission(d.Perms, "iam:role:ownuser", "角色成员查询"), s.ownUser)
		api.POST("/v1/admin/sys/roles/grant-user", admin, middleware.RequirePermission(d.Perms, "iam:role:grantuser", "角色成员授权"), s.grantUser)
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
	if err := s.Update(c.Request.Context(), req, contextx.Session(c.Request.Context())); err != nil {
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
	if err := s.Delete(c.Request.Context(), body.IDs, contextx.Session(c.Request.Context())); err != nil {
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
	row, err := s.Detail(c.Request.Context(), q.ID, contextx.Session(c.Request.Context()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, 404, "not found")
			return
		}
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
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
func (s *Service) ownResource(c *gin.Context) {
	var q OwnResourceQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := s.OwnResources(c.Request.Context(), q.ID, q.AccountType)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (s *Service) grantResource(c *gin.Context) {
	var req GrantResourceParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.GrantResources(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) ownClientResource(c *gin.Context) {
	var q OwnResourceQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := s.OwnClientResources(c.Request.Context(), q.ID, q.AccountType)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (s *Service) grantClientResource(c *gin.Context) {
	var req GrantResourceParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.GrantClientResources(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) ownUser(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := s.OwnUsers(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (s *Service) grantUser(c *gin.Context) {
	var req GrantUserParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.GrantUsers(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}
