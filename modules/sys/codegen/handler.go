package codegen

import (
	"net/http"

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
		api.POST("/v1/admin/sys/codegen/create", admin, middleware.RequirePermission(d.Perms, "sys:codegen:create", "代码生成创建"), s.create)
		api.POST("/v1/admin/sys/codegen/update", admin, middleware.RequirePermission(d.Perms, "sys:codegen:update", "代码生成更新"), s.update)
		api.POST("/v1/admin/sys/codegen/delete", admin, middleware.RequirePermission(d.Perms, "sys:codegen:delete", "代码生成删除"), s.delete)
		api.GET("/v1/admin/sys/codegen/detail", admin, middleware.RequirePermission(d.Perms, "sys:codegen:detail", "代码生成详情"), s.detail)
		api.GET("/v1/admin/sys/codegen/page", admin, middleware.RequirePermission(d.Perms, "sys:codegen:page", "代码生成分页"), s.page)
		api.GET("/v1/admin/sys/codegen/tables", admin, middleware.RequirePermission(d.Perms, "sys:codegen:tables", "代码生成表列表"), s.tables)
		api.GET("/v1/admin/sys/codegen/table-columns", admin, middleware.RequirePermission(d.Perms, "sys:codegen:tables", "代码生成表列元数据"), s.tableColumns)
		api.GET("/v1/admin/sys/codegen/fields", admin, middleware.RequirePermission(d.Perms, "sys:codegen:detail", "代码生成字段"), s.fields)
		api.POST("/v1/admin/sys/codegen/fields/update-batch", admin, middleware.RequirePermission(d.Perms, "sys:codegen:update", "代码生成字段批量更新"), s.updateFieldsBatch)
		api.GET("/v1/admin/sys/codegen/parent-resources", admin, middleware.RequirePermission(d.Perms, "sys:codegen:detail", "代码生成父级资源"), s.parentResources)
		api.GET("/v1/admin/sys/codegen/preview", admin, middleware.RequirePermission(d.Perms, "sys:codegen:preview", "代码生成预览"), s.preview)
		api.GET("/v1/admin/sys/codegen/download", admin, middleware.RequirePermission(d.Perms, "sys:codegen:download", "代码生成下载"), s.download)
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

func (s *Service) tables(c *gin.Context) {
	rows, err := s.Tables(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) tableColumns(c *gin.Context) {
	tableName := c.Query("table_name")
	if tableName == "" {
		response.Fail(c, http.StatusBadRequest, 400, "table_name required")
		return
	}
	rows, err := s.TableColumns(c.Request.Context(), tableName)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) fields(c *gin.Context) {
	var q FieldQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	rows, err := s.Fields(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) updateFieldsBatch(c *gin.Context) {
	var req FieldsUpdateBatchParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateFieldsBatch(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) parentResources(c *gin.Context) {
	moduleID := c.Query("module_id")
	rows, err := s.ParentResources(c.Request.Context(), moduleID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (s *Service) preview(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := s.Preview(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) download(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	raw, name, err := s.DownloadZip(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	c.Header("Content-Disposition", "attachment; filename=\""+name+"\"")
	c.Data(http.StatusOK, "application/zip", raw)
}