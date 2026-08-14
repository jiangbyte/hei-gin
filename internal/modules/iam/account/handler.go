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
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/accounts/create", admin, middleware.RequirePermission(d.Perms, "iam:account:create", "è´¦æˆ·åˆ›å»º"), s.create)
		api.POST("/v1/admin/sys/accounts/update", admin, middleware.RequirePermission(d.Perms, "iam:account:update", "è´¦æˆ·æ›´æ–°"), s.update)
		api.POST("/v1/admin/sys/accounts/delete", admin, middleware.RequirePermission(d.Perms, "iam:account:delete", "è´¦æˆ·åˆ é™¤"), s.delete)
		api.GET("/v1/admin/sys/accounts/detail", admin, middleware.RequirePermission(d.Perms, "iam:account:detail", "è´¦æˆ·è¯¦æƒ…"), s.detail)
		api.GET("/v1/admin/sys/accounts/page", admin, middleware.RequirePermission(d.Perms, "iam:account:page", "è´¦æˆ·åˆ†é¡µ"), s.page)
		api.GET("/v1/admin/sys/accounts/own-role", admin, middleware.RequirePermission(d.Perms, "iam:account:ownrole", "è´¦å·å·²æ‹¥æœ‰è§’è‰²"), s.ownRole)
		api.POST("/v1/admin/sys/accounts/grant-role", admin, middleware.RequirePermission(d.Perms, "iam:account:grantrole", "è´¦å·è§’è‰²æŽˆæƒ"), s.grantRole)
		api.GET("/v1/admin/sys/accounts/own-group", admin, middleware.RequirePermission(d.Perms, "iam:account:owngroup", "è´¦å·å·²æ‹¥æœ‰ç”¨æˆ·ç»„"), s.ownGroup)
		api.POST("/v1/admin/sys/accounts/grant-group", admin, middleware.RequirePermission(d.Perms, "iam:account:grantgroup", "è´¦å·ç”¨æˆ·ç»„æŽˆæƒ"), s.grantGroup)
		api.GET("/v1/admin/sys/accounts/own-dept", admin, middleware.RequirePermission(d.Perms, "iam:account:owndept", "è´¦å·å·²æ‹¥æœ‰éƒ¨é—¨"), s.ownDept)
		api.POST("/v1/admin/sys/accounts/grant-dept", admin, middleware.RequirePermission(d.Perms, "iam:account:grantdept", "è´¦å·éƒ¨é—¨æŽˆæƒ"), s.grantDept)
		api.GET("/v1/admin/sys/accounts/own-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:ownresource", "è´¦å·å·²æ‹¥æœ‰èµ„æº"), s.ownResource)
		api.POST("/v1/admin/sys/accounts/grant-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:grantresource", "è´¦å·èµ„æºæŽˆæƒ"), s.grantResource)
		api.GET("/v1/admin/sys/accounts/own-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:ownclientresource", "è´¦å·å·²æ‹¥æœ‰å®¢æˆ·ç«¯èµ„æº"), s.ownClientResource)
		api.POST("/v1/admin/sys/accounts/grant-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:grantclientresource", "è´¦å·å®¢æˆ·ç«¯èµ„æºæŽˆæƒ"), s.grantClientResource)
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
