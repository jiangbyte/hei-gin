// internal/modules/iam/group/handler.go HTTP 处理器。
//
// Author: Charlie

package group

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
		api.POST("/v1/admin/sys/groups/create", admin, middleware.RequirePermission(d.Perms, "iam:group:create", "ç”¨æˆ·ç»„åˆ›å»º"), s.create)
		api.POST("/v1/admin/sys/groups/update", admin, middleware.RequirePermission(d.Perms, "iam:group:update", "ç”¨æˆ·ç»„æ›´æ–°"), s.update)
		api.POST("/v1/admin/sys/groups/delete", admin, middleware.RequirePermission(d.Perms, "iam:group:delete", "ç”¨æˆ·ç»„åˆ é™¤"), s.delete)
		api.GET("/v1/admin/sys/groups/detail", admin, middleware.RequirePermission(d.Perms, "iam:group:detail", "ç”¨æˆ·ç»„è¯¦æƒ…"), s.detail)
		api.GET("/v1/admin/sys/groups/page", admin, middleware.RequirePermission(d.Perms, "iam:group:page", "ç”¨æˆ·ç»„åˆ†é¡µ"), s.page)
		api.GET("/v1/admin/sys/groups/own-user", admin, middleware.RequirePermission(d.Perms, "iam:group:ownuser", "ç”¨æˆ·ç»„æˆå‘˜æŸ¥è¯¢"), s.ownUser)
		api.POST("/v1/admin/sys/groups/grant-user", admin, middleware.RequirePermission(d.Perms, "iam:group:grantuser", "ç”¨æˆ·ç»„æˆå‘˜æŽˆæƒ"), s.grantUser)
		api.GET("/v1/admin/sys/groups/own-role", admin, middleware.RequirePermission(d.Perms, "iam:group:ownrole", "ç”¨æˆ·ç»„å·²æ‹¥æœ‰è§’è‰²"), s.ownRole)
		api.POST("/v1/admin/sys/groups/grant-role", admin, middleware.RequirePermission(d.Perms, "iam:group:grantrole", "ç”¨æˆ·ç»„è§’è‰²æŽˆæƒ"), s.grantRole)
		api.GET("/v1/admin/sys/groups/own-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:ownresource", "ç”¨æˆ·ç»„å·²æ‹¥æœ‰èµ„æº"), s.ownResource)
		api.POST("/v1/admin/sys/groups/grant-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:grantresource", "ç”¨æˆ·ç»„èµ„æºæŽˆæƒ"), s.grantResource)
		api.GET("/v1/admin/sys/groups/own-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:ownclientresource", "ç”¨æˆ·ç»„å·²æ‹¥æœ‰å®¢æˆ·ç«¯èµ„æº"), s.ownClientResource)
		api.POST("/v1/admin/sys/groups/grant-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:grantclientresource", "ç”¨æˆ·ç»„å®¢æˆ·ç«¯èµ„æºæŽˆæƒ"), s.grantClientResource)
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
