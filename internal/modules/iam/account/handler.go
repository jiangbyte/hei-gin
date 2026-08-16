// internal/modules/iam/account/handler.go HTTP 处理器。
//
// Author: Charlie

package account

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		// 操作审计登记（对齐 hei-boot @OperationAudit：iam_account）
		d.AuditReg.RegisterSpecs(
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/accounts/create", ResourceType: "iam_account", Action: "create"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/accounts/update", ResourceType: "iam_account", Action: "update"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/accounts/delete", ResourceType: "iam_account", Action: "delete"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/accounts/grant-role", ResourceType: "iam_account", Action: "grant_role"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/accounts/grant-group", ResourceType: "iam_account", Action: "grant_group"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/accounts/grant-dept", ResourceType: "iam_account", Action: "grant_dept"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/accounts/grant-resource", ResourceType: "iam_account", Action: "grant_resource"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/accounts/grant-client-resource", ResourceType: "iam_account", Action: "grant_client_resource"},
		)
		api.POST("/v1/admin/sys/accounts/create", admin, middleware.RequirePermission(d.Perms, "iam:account:create", "账户创建"), s.create)
		api.POST("/v1/admin/sys/accounts/update", admin, middleware.RequirePermission(d.Perms, "iam:account:update", "账户更新"), s.update)
		api.POST("/v1/admin/sys/accounts/delete", admin, middleware.RequirePermission(d.Perms, "iam:account:delete", "账户删除"), s.delete)
		api.GET("/v1/admin/sys/accounts/detail", admin, middleware.RequirePermission(d.Perms, "iam:account:detail", "账户详情"), s.detail)
		api.GET("/v1/admin/sys/accounts/page", admin, middleware.RequirePermission(d.Perms, "iam:account:page", "账户分页"), s.page)
		api.GET("/v1/admin/sys/accounts/own-role", admin, middleware.RequirePermission(d.Perms, "iam:account:ownrole", "账号已拥有角色"), s.ownRole)
		api.POST("/v1/admin/sys/accounts/grant-role", admin, middleware.RequirePermission(d.Perms, "iam:account:grantrole", "账号角色授权"), s.grantRole)
		api.GET("/v1/admin/sys/accounts/own-group", admin, middleware.RequirePermission(d.Perms, "iam:account:owngroup", "账号已拥有用户组"), s.ownGroup)
		api.POST("/v1/admin/sys/accounts/grant-group", admin, middleware.RequirePermission(d.Perms, "iam:account:grantgroup", "账号用户组授权"), s.grantGroup)
		api.GET("/v1/admin/sys/accounts/own-dept", admin, middleware.RequirePermission(d.Perms, "iam:account:owndept", "账号已拥有部门"), s.ownDept)
		api.POST("/v1/admin/sys/accounts/grant-dept", admin, middleware.RequirePermission(d.Perms, "iam:account:grantdept", "账号部门授权"), s.grantDept)
		api.GET("/v1/admin/sys/accounts/own-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:ownresource", "账号已拥有资源"), s.ownResource)
		api.POST("/v1/admin/sys/accounts/grant-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:grantresource", "账号资源授权"), s.grantResource)
		api.GET("/v1/admin/sys/accounts/own-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:ownclientresource", "账号已拥有客户端资源"), s.ownClientResource)
		api.POST("/v1/admin/sys/accounts/grant-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:grantclientresource", "账号客户端资源授权"), s.grantClientResource)
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
	vo, err := s.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, vo)
}

func (s *Service) page(c *gin.Context) {
	var q PageParam
	_ = c.ShouldBindQuery(&q)
	records, total, cur, size, err := s.Page(c.Request.Context(), q, contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, records)
}
func (s *Service) ownRole(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := s.OwnRoles(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (s *Service) grantRole(c *gin.Context) {
	var req GrantRoleParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.GrantRoles(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) ownGroup(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := s.OwnGroups(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (s *Service) grantGroup(c *gin.Context) {
	var req GrantGroupParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.GrantGroups(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) ownDept(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := s.OwnDepts(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (s *Service) grantDept(c *gin.Context) {
	var req GrantDeptParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.GrantDepts(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
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
