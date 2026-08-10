package file

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"hei-gin/framework/core/bind"
	"hei-gin/framework/core/response"
	"hei-gin/framework/core/schema"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/file/upload", admin, middleware.RequirePermission(d.Perms, "sys:file:upload", "文件上传"), s.upload)
		api.POST("/v1/admin/sys/file/delete", admin, middleware.RequirePermission(d.Perms, "sys:file:delete", "文件删除"), s.delete)
		api.POST("/v1/admin/sys/file/update", admin, middleware.RequirePermission(d.Perms, "sys:file:update", "文件更新"), s.update)
		api.GET("/v1/admin/sys/file/detail", admin, middleware.RequirePermission(d.Perms, "sys:file:detail", "文件详情"), s.detail)
		api.GET("/v1/admin/sys/file/page", admin, middleware.RequirePermission(d.Perms, "sys:file:page", "文件分页"), s.page)
		api.POST("/v1/admin/sys/file/list_by_ids", admin, middleware.RequirePermission(d.Perms, "sys:file:detail", "文件批量查询"), s.listByIDs)
		api.GET("/v1/admin/sys/file/download", admin, middleware.RequirePermission(d.Perms, "sys:file:url", "文件下载"), s.download)
		api.GET("/v1/files/*object_name", s.publicGet)
	}
}

func (s *Service) upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "file required")
		return
	}
	row, err := s.Upload(c.Request.Context(), fh)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, row)
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

func (s *Service) listByIDs(c *gin.Context) {
	var body IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	rows, err := s.ListByIDs(c.Request.Context(), body.IDs)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) download(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, rc, err := s.OpenByID(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	defer rc.Close()
	c.Header("Content-Disposition", `attachment; filename="`+row.OriginalName+`"`)
	c.DataFromReader(http.StatusOK, row.Size, row.ContentType, rc, nil)
}

func (s *Service) publicGet(c *gin.Context) {
	objectName := strings.TrimPrefix(c.Param("object_name"), "/")
	if objectName == "" {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	ct, rc, err := s.OpenByObjectName(c.Request.Context(), objectName)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	defer rc.Close()
	c.Header("Content-Type", ct)
	_, _ = io.Copy(c.Writer, rc)
}
