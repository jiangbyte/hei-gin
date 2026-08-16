// internal/modules/sys/job/handler.go HTTP 处理器。
//
// Author: Charlie

package job

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
		// 操作审计登记（对齐 hei-boot @OperationAudit：sys_job；status→enabled、trigger→run）
		d.AuditReg.RegisterSpecs(
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/jobs/create", ResourceType: "sys_job", Action: "create"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/jobs/update", ResourceType: "sys_job", Action: "update"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/jobs/delete", ResourceType: "sys_job", Action: "delete"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/jobs/status", ResourceType: "sys_job", Action: "enabled"},
			audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/admin/sys/jobs/trigger", ResourceType: "sys_job", Action: "run"},
		)
		api.GET("/v1/admin/sys/jobs/page", admin, middleware.RequirePermission(d.Perms, "sys:job:page", "任务分页"), s.page)
		api.GET("/v1/admin/sys/jobs/handlers", admin, middleware.RequirePermission(d.Perms, "sys:job:page", "任务处理器列表"), s.handlers)
		api.GET("/v1/admin/sys/jobs/detail", admin, middleware.RequirePermission(d.Perms, "sys:job:detail", "任务详情"), s.detail)
		api.POST("/v1/admin/sys/jobs/create", admin, middleware.RequirePermission(d.Perms, "sys:job:create", "任务创建"), s.create)
		api.POST("/v1/admin/sys/jobs/update", admin, middleware.RequirePermission(d.Perms, "sys:job:update", "任务更新"), s.update)
		api.POST("/v1/admin/sys/jobs/delete", admin, middleware.RequirePermission(d.Perms, "sys:job:delete", "任务删除"), s.delete)
		api.POST("/v1/admin/sys/jobs/status", admin, middleware.RequirePermission(d.Perms, "sys:job:update", "任务启停"), s.status)
		api.POST("/v1/admin/sys/jobs/trigger", admin, middleware.RequirePermission(d.Perms, "sys:job:update", "任务触发"), s.trigger)
		api.GET("/v1/admin/sys/jobs/logs", admin, middleware.RequirePermission(d.Perms, "sys:job:page", "任务日志"), s.logs)
	}
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

func (s *Service) handlers(c *gin.Context) {
	response.OK(c, s.Handlers())
}

func (s *Service) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.OK(c, row)
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

func (s *Service) status(c *gin.Context) {
	var req StatusParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.SetStatus(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) trigger(c *gin.Context) {
	var req TriggerParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Trigger(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) logs(c *gin.Context) {
	var q LogParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.Logs(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}
