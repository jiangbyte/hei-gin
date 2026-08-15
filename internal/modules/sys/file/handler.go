// internal/modules/sys/file/handler.go HTTP 处理器（对齐 hei-boot Admin/Portal/PublicFileController）。
//
// Author: Charlie

package file

import (
	"io"
	"net/http"
	"strings"

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
		api.POST("/v1/admin/sys/file/upload", admin, middleware.RequirePermission(d.Perms, "sys:file:upload", "文件上传"), s.upload)
		api.POST("/v1/admin/sys/file/delete", admin, middleware.RequirePermission(d.Perms, "sys:file:delete", "文件删除"), s.delete)
		api.POST("/v1/admin/sys/file/update", admin, middleware.RequirePermission(d.Perms, "sys:file:update", "文件更新"), s.update)
		api.GET("/v1/admin/sys/file/detail", admin, middleware.RequirePermission(d.Perms, "sys:file:detail", "文件详情"), s.detail)
		api.GET("/v1/admin/sys/file/page", admin, middleware.RequirePermission(d.Perms, "sys:file:page", "文件分页"), s.page)
		api.POST("/v1/admin/sys/file/list_by_ids", admin, middleware.RequirePermission(d.Perms, "sys:file:detail", "文件批量查询"), s.listByIDs)
		api.GET("/v1/admin/sys/file/download", admin, middleware.RequirePermission(d.Perms, "sys:file:url", "文件下载"), s.download)
		api.POST("/v1/admin/sys/file/url", admin, middleware.RequirePermission(d.Perms, "sys:file:url", "文件URL"), s.url)
		api.POST("/v1/admin/sys/file/presigned_url", admin, middleware.RequirePermission(d.Perms, "sys:file:presignedurl", "文件预签名URL"), s.presignedURL)
		api.GET("/v1/files", s.publicGet)
		api.GET("/v1/files/*object_name", s.publicGet)

		portal := middleware.RequireAccountType(security.AccountPortal)
		api.POST("/v1/portal/sys/file/upload", portal, s.portalUpload)
		api.GET("/v1/portal/sys/file/detail", portal, s.portalDetail)
		api.POST("/v1/portal/sys/file/list_by_ids", portal, s.portalListByIDs)
		api.GET("/v1/portal/sys/file/download", portal, s.portalDownload)
		api.POST("/v1/portal/sys/file/url", portal, s.portalURL)
		api.POST("/v1/portal/sys/file/presigned_url", portal, s.portalPresignedURL)
	}
}

func (s *Service) upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "file required")
		return
	}
	row, err := s.Upload(c.Request.Context(), fh, c.PostForm("storage_provider"), accountID(c))
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

// publicGet 公开文件访问：/v1/files?object_name= 或 /v1/files/**（校验元数据存在，防越权读存储）。
func (s *Service) publicGet(c *gin.Context) {
	objectName := strings.TrimPrefix(c.Param("object_name"), "/")
	if objectName == "" {
		objectName = c.Query("object_name")
	}
	if objectName == "" {
		response.Fail(c, http.StatusBadRequest, 400, "object_name required")
		return
	}
	ct, rc, err := s.OpenByObjectName(c.Request.Context(), objectName)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	defer rc.Close()
	filename := objectName
	if idx := strings.LastIndex(filename, "/"); idx >= 0 {
		filename = filename[idx+1:]
	}
	c.Header("Content-Disposition", `inline; filename="`+filename+`"`)
	c.Header("Content-Type", ct)
	_, _ = io.Copy(c.Writer, rc)
}

func (s *Service) url(c *gin.Context) {
	var req ObjectNameParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := s.URL(c.Request.Context(), req.ObjectName)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) presignedURL(c *gin.Context) {
	var req ObjectNameParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := s.PresignedURL(c.Request.Context(), req.ObjectName)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.OK(c, out)
}

// ---- Portal（仅本人文件，对齐 hei-boot PortalFileController.assertOwnedByCurrent）----

func (s *Service) portalUpload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "file required")
		return
	}
	row, err := s.Upload(c.Request.Context(), fh, c.PostForm("storage_provider"), accountID(c))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, row)
}

func (s *Service) portalDetail(c *gin.Context) {
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
	if err := s.AssertOwnedByCurrent(row, accountID(c)); err != nil {
		response.Fail(c, http.StatusForbidden, 403, err.Error())
		return
	}
	response.OK(c, row)
}

func (s *Service) portalListByIDs(c *gin.Context) {
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
	me := accountID(c)
	filtered := make([]File, 0, len(rows))
	for i := range rows {
		if s.AssertOwnedByCurrent(&rows[i], me) == nil {
			filtered = append(filtered, rows[i])
		}
	}
	response.OK(c, filtered)
}

func (s *Service) portalDownload(c *gin.Context) {
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
	if err := s.AssertOwnedByCurrent(row, accountID(c)); err != nil {
		response.Fail(c, http.StatusForbidden, 403, err.Error())
		return
	}
	rc, err := s.providerFor(c.Request.Context(), row).Get(c.Request.Context(), toObjectKey(row.ObjectName, s.publicPath()))
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	defer rc.Close()
	c.Header("Content-Disposition", `attachment; filename="`+row.OriginalName+`"`)
	c.DataFromReader(http.StatusOK, row.Size, row.ContentType, rc, nil)
}

func (s *Service) portalURL(c *gin.Context) {
	s.portalObjectName(c, func(ctx *gin.Context, objectName string) (*URLResult, error) {
		return s.URL(ctx, objectName)
	})
}

func (s *Service) portalPresignedURL(c *gin.Context) {
	s.portalObjectName(c, func(ctx *gin.Context, objectName string) (*URLResult, error) {
		return s.PresignedURL(ctx, objectName)
	})
}

// portalObjectName 门户端按 object_name 获取 URL：先校验本人归属。
func (s *Service) portalObjectName(c *gin.Context, fn func(ctx *gin.Context, objectName string) (*URLResult, error)) {
	var req ObjectNameParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	key := toObjectKey(req.ObjectName, s.publicPath())
	if key == "" {
		response.Fail(c, http.StatusNotFound, 404, "file not found")
		return
	}
	row, err := s.repo.FindByObjectName(c.Request.Context(), key)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "file not found")
		return
	}
	if err := s.AssertOwnedByCurrent(row, accountID(c)); err != nil {
		response.Fail(c, http.StatusForbidden, 403, err.Error())
		return
	}
	out, err := fn(c, key)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.OK(c, out)
}

func accountID(c *gin.Context) string {
	if sess := contextx.Session(c.Request.Context()); sess != nil {
		return sess.AccountID
	}
	return ""
}
