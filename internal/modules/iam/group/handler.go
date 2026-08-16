// internal/modules/iam/group/handler.go HTTP 处理器。
//
// Author: Charlie

package group

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
	"hei-gin/internal/framework/platform/module"
)

func (s *Service) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/groups/create", admin, middleware.RequirePermission(d.Perms, "iam:group:create", "用户组创建"), middleware.OperationAudit(d.Audit, "iam_group", "create"), s.create)
		api.POST("/v1/admin/sys/groups/update", admin, middleware.RequirePermission(d.Perms, "iam:group:update", "用户组更新"), middleware.OperationAudit(d.Audit, "iam_group", "update"), s.update)
		api.POST("/v1/admin/sys/groups/delete", admin, middleware.RequirePermission(d.Perms, "iam:group:delete", "用户组删除"), middleware.OperationAudit(d.Audit, "iam_group", "delete"), s.delete)
		api.GET("/v1/admin/sys/groups/detail", admin, middleware.RequirePermission(d.Perms, "iam:group:detail", "用户组详情"), s.detail)
		api.GET("/v1/admin/sys/groups/page", admin, middleware.RequirePermission(d.Perms, "iam:group:page", "用户组分页"), s.page)
		api.GET("/v1/admin/sys/groups/own-user", admin, middleware.RequirePermission(d.Perms, "iam:group:ownuser", "用户组成员查询"), s.ownUser)
		api.POST("/v1/admin/sys/groups/grant-user", admin, middleware.RequirePermission(d.Perms, "iam:group:grantuser", "用户组成员授权"), middleware.OperationAudit(d.Audit, "iam_group", "grant_user"), s.grantUser)
		api.GET("/v1/admin/sys/groups/own-role", admin, middleware.RequirePermission(d.Perms, "iam:group:ownrole", "用户组已拥有角色"), s.ownRole)
		api.POST("/v1/admin/sys/groups/grant-role", admin, middleware.RequirePermission(d.Perms, "iam:group:grantrole", "用户组角色授权"), middleware.OperationAudit(d.Audit, "iam_group", "grant_role"), s.grantRole)
		api.GET("/v1/admin/sys/groups/own-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:ownresource", "用户组已拥有资源"), s.ownResource)
		api.POST("/v1/admin/sys/groups/grant-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:grantresource", "用户组资源授权"), middleware.OperationAudit(d.Audit, "iam_group", "grant_resource"), s.grantResource)
		api.GET("/v1/admin/sys/groups/own-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:ownclientresource", "用户组已拥有客户端资源"), s.ownClientResource)
		api.POST("/v1/admin/sys/groups/grant-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:grantclientresource", "用户组客户端资源授权"), middleware.OperationAudit(d.Audit, "iam_group", "grant_client_resource"), s.grantClientResource)
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

func (s *Service) ownRole(c *gin.Context) {
	var q OwnResourceQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := s.OwnRoles(c.Request.Context(), q.ID, q.AccountType)
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
