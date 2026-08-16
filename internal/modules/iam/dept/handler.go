// internal/modules/iam/dept/handler.go HTTP 处理器。
//
// Author: Charlie

package dept

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
		// 操作审计登记（对齐 hei-boot @OperationAudit：iam_dept）
		d.AuditReg.RegisterSpecs(
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/depts/create", ResourceType: "iam_dept", Action: "create"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/depts/update", ResourceType: "iam_dept", Action: "update"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/depts/delete", ResourceType: "iam_dept", Action: "delete"},
		)
		api.POST("/v1/admin/sys/depts/create", admin, middleware.RequirePermission(d.Perms, "iam:dept:create", "部门创建"), s.create)
		api.POST("/v1/admin/sys/depts/update", admin, middleware.RequirePermission(d.Perms, "iam:dept:update", "部门更新"), s.update)
		api.POST("/v1/admin/sys/depts/delete", admin, middleware.RequirePermission(d.Perms, "iam:dept:delete", "部门删除"), s.delete)
		api.GET("/v1/admin/sys/depts/detail", admin, middleware.RequirePermission(d.Perms, "iam:dept:detail", "部门详情"), s.detail)
		api.GET("/v1/admin/sys/depts/page", admin, middleware.RequirePermission(d.Perms, "iam:dept:page", "部门分页"), s.page)
		api.GET("/v1/admin/sys/depts/tree", admin, middleware.RequirePermission(d.Perms, "iam:dept:tree", "部门树"), s.tree)
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

func (s *Service) tree(c *gin.Context) {
	nodes, err := s.Tree(c.Request.Context(), contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}
