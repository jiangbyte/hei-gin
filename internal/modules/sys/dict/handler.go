// internal/modules/sys/dict/handler.go HTTP 处理器。
//
// Author: Charlie

package dict

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
		// 操作审计登记（对齐 hei-boot @OperationAudit：sys_dict）
		d.AuditReg.RegisterSpecs(
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/dicts/create", ResourceType: "sys_dict", Action: "create"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/dicts/update", ResourceType: "sys_dict", Action: "update"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/dicts/delete", ResourceType: "sys_dict", Action: "delete"},
		)
		api.POST("/v1/admin/sys/dicts/create", admin, middleware.RequirePermission(d.Perms, "sys:dict:create", "字典创建"), s.create)
		api.POST("/v1/admin/sys/dicts/update", admin, middleware.RequirePermission(d.Perms, "sys:dict:update", "字典更新"), s.update)
		api.POST("/v1/admin/sys/dicts/delete", admin, middleware.RequirePermission(d.Perms, "sys:dict:delete", "字典删除"), s.delete)
		api.GET("/v1/admin/sys/dicts/detail", admin, middleware.RequirePermission(d.Perms, "sys:dict:detail", "字典详情"), s.detail)
		api.GET("/v1/admin/sys/dicts/page", admin, middleware.RequirePermission(d.Perms, "sys:dict:page", "字典分页"), s.page)
		api.GET("/v1/admin/sys/dicts/tree", admin, s.tree)

		// 门户字典树为公开接口（对齐 hei-boot PortalDictController，web public:true）
		api.GET("/v1/portal/sys/dicts/tree", s.portalTree)
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

func (s *Service) tree(c *gin.Context) {
	var q TreeParam
	q.Code = c.Query("code")
	q.Category = c.Query("category")
	nodes, err := s.Tree(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}

func (s *Service) portalTree(c *gin.Context) {
	var q TreeParam
	q.Category = c.Query("category")
	nodes, err := s.Tree(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}
