// internal/modules/sys/job/handler.go HTTP 处理器（对齐 hei-boot AdminJobController）。
//
// Author: Charlie

package job

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/gojob"
	"hei-gin/internal/framework/platform/module"
)

func (s *Service) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET("/v1/admin/sys/jobs/page", admin, middleware.RequirePermission(d.Perms, "sys:job:page", "任务分页"), s.page)
		api.GET("/v1/admin/sys/jobs/detail", admin, middleware.RequirePermission(d.Perms, "sys:job:detail", "任务详情"), s.detail)
		api.POST("/v1/admin/sys/jobs/create", admin, middleware.RequirePermission(d.Perms, "sys:job:create", "任务创建"), middleware.OperationAudit(d.Audit, "sys_job", "create"), s.create)
		api.POST("/v1/admin/sys/jobs/update", admin, middleware.RequirePermission(d.Perms, "sys:job:update", "任务更新"), middleware.OperationAudit(d.Audit, "sys_job", "update"), s.update)
		api.POST("/v1/admin/sys/jobs/delete", admin, middleware.RequirePermission(d.Perms, "sys:job:delete", "任务删除"), middleware.OperationAudit(d.Audit, "sys_job", "delete"), s.delete)
		api.POST("/v1/admin/sys/jobs/enabled", admin, middleware.RequirePermission(d.Perms, "sys:job:update", "任务启停"), middleware.OperationAudit(d.Audit, "sys_job", "enabled"), s.enabled)
		api.POST("/v1/admin/sys/jobs/run", admin, middleware.RequirePermission(d.Perms, "sys:job:run", "任务立即执行"), middleware.OperationAudit(d.Audit, "sys_job", "run"), s.run)
		api.GET("/v1/admin/sys/job-logs/page", admin, middleware.RequirePermission(d.Perms, "sys:joblog:page", "任务日志分页"), s.logs)
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

func (s *Service) enabled(c *gin.Context) {
	var req EnabledParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.SetEnabled(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) run(c *gin.Context) {
	var req RunParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	executor := gojob.ExecutorSystem
	if sess := contextx.Session(c.Request.Context()); sess != nil && sess.AccountID != "" {
		executor = sess.AccountID
	}
	if err := s.RunNow(c.Request.Context(), req.ID, executor); err != nil {
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
