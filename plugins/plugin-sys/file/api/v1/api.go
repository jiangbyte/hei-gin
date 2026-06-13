package v1

import (
	"hei-gin/sdk/auth"
	"strconv"

	file "hei-gin/plugins/plugin-sys/file"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *file.Service
}

var defaultHandler = newHandler(file.DefaultModule)

func newHandler(module *file.Module) *handler {
	return &handler{service: module.Service()}
}

// RegisterRoutes registers all admin file routes.
func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/sys/file/upload",
		registry.Perm("sys:file:upload", "上传文件"),
		log.SysLog("上传文件"),
		defaultHandler.upload,
	)

	r.POST("/api/v1/sys/file/upload/init",
		registry.Perm("sys:file:upload", "分片上传-初始化"),
		log.SysLog("分片上传-初始化"),
		defaultHandler.uploadInit,
	)
	r.POST("/api/v1/sys/file/upload/chunk",
		registry.Perm("sys:file:upload", "分片上传-上传分片"),
		log.SysLog("分片上传-上传分片"),
		defaultHandler.uploadChunk,
	)
	r.POST("/api/v1/sys/file/upload/complete",
		registry.Perm("sys:file:upload", "分片上传-完成"),
		log.SysLog("分片上传-完成"),
		defaultHandler.uploadComplete,
	)
	r.POST("/api/v1/sys/file/upload/abort",
		registry.Perm("sys:file:upload", "分片上传-取消"),
		log.SysLog("分片上传-取消"),
		defaultHandler.uploadAbort,
	)

	r.GET("/api/v1/sys/file/download",
		registry.Perm("sys:file:download", "下载文件"),
		defaultHandler.download,
	)
	r.GET("/api/v1/sys/file/page",
		registry.Perm("sys:file:page", "文件分页"),
		defaultHandler.page,
	)
	r.GET("/api/v1/sys/file/detail",
		registry.Perm("sys:file:detail", "文件详情"),
		defaultHandler.detail,
	)
	r.POST("/api/v1/sys/file/remove",
		registry.Perm("sys:file:remove", "删除文件"),
		log.SysLog("删除文件"),
		defaultHandler.remove,
	)
	r.POST("/api/v1/sys/file/remove-absolute",
		registry.Perm("sys:file:remove", "物理删除文件"),
		log.SysLog("物理删除文件"),
		defaultHandler.removeAbsolute,
	)
}

// RegisterClientRoutes registers client file routes.
func RegisterClientRoutes(r *gin.Engine) {
	r.POST("/api/v1/c/file/upload",
		middleware.CheckLogin(auth.Consumer),
		log.SysLog("C端上传文件"),
		defaultHandler.clientUpload,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// uploadHandler handles POST /api/v1/sys/file/upload
// @Summary      文件管理上传文件
// @Description  访问 /api/v1/sys/file/upload，文件管理上传文件
// @Tags         文件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "上传文件"
// @Param        engine  formData  string  false  "存储引擎"
// @Param        bucket  formData  string  false  "存储桶"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/upload [post]
func (h *handler) upload(c *gin.Context) {
	data, err := h.service.Upload(c)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

// clientUploadHandler handles POST /api/v1/c/file/upload
// @Summary      文件管理上传文件
// @Description  访问 /api/v1/c/file/upload，文件管理上传文件
// @Tags         文件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "上传文件"
// @Param        engine  formData  string  false  "存储引擎"
// @Param        bucket  formData  string  false  "存储桶"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/file/upload [post]
func (h *handler) clientUpload(c *gin.Context) {
	data, err := h.service.Upload(c)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

// downloadHandler handles GET /api/v1/sys/file/download
// @Summary      文件管理下载
// @Description  访问 /api/v1/sys/file/download，文件管理下载
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/download [get]
func (h *handler) download(c *gin.Context) {
	if err := h.service.Download(c, c.Query("id")); err != nil {
		result.Failure(c, err.Error(), 400)
	}
}

// pageHandler handles GET /api/v1/sys/file/page
// @Summary      文件管理分页查询
// @Description  访问 /api/v1/sys/file/page，文件管理分页查询
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        query  query  file.FilePageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/page [get]
func (h *handler) page(c *gin.Context) {
	var param file.FilePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
}

// detailHandler handles GET /api/v1/sys/file/detail
// @Summary      文件管理详情查询
// @Description  访问 /api/v1/sys/file/detail，文件管理详情查询
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}

// removeHandler handles POST /api/v1/sys/file/remove
// @Summary      文件管理删除
// @Description  访问 /api/v1/sys/file/remove，文件管理删除
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// removeAbsoluteHandler handles POST /api/v1/sys/file/remove-absolute
// @Summary      文件管理接口调用
// @Description  访问 /api/v1/sys/file/remove-absolute，文件管理接口调用
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/remove-absolute [post]
func (h *handler) removeAbsolute(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.RemoveAbsolute(c, &param)
	result.Success(c, nil)
}

// uploadInitHandler handles POST /api/v1/sys/file/upload/init
// @Summary      文件管理初始化分片上传
// @Description  访问 /api/v1/sys/file/upload/init，文件管理初始化分片上传
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        body  body  file.ChunkUploadInitParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/upload/init [post]
func (h *handler) uploadInit(c *gin.Context) {
	var param file.ChunkUploadInitParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data, err := h.service.InitChunkUpload(c, &param)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

// uploadChunkHandler handles POST /api/v1/sys/file/upload/chunk
// @Summary      文件管理上传分片
// @Description  访问 /api/v1/sys/file/upload/chunk，文件管理上传分片
// @Tags         文件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        chunk_index  formData  string  false  "chunk_index"
// @Param        total_chunks  formData  string  false  "total_chunks"
// @Param        upload_id  formData  string  false  "upload_id"
// @Param        checksum  formData  string  false  "checksum"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/upload/chunk [post]
func (h *handler) uploadChunk(c *gin.Context) {
	chunkIndex, err := strconv.Atoi(c.PostForm("chunk_index"))
	if err != nil {
		result.Failure(c, "chunk_index 参数错误", 400)
		return
	}
	totalChunks := 0
	if raw := c.PostForm("total_chunks"); raw != "" {
		totalChunks, err = strconv.Atoi(raw)
		if err != nil {
			result.Failure(c, "total_chunks 参数错误", 400)
			return
		}
	}
	param := file.ChunkUploadParam{
		UploadID:    c.PostForm("upload_id"),
		ChunkIndex:  chunkIndex,
		TotalChunks: totalChunks,
		Checksum:    c.PostForm("checksum"),
	}
	if param.UploadID == "" {
		result.Failure(c, "upload_id 不能为空", 400)
		return
	}

	if err := h.service.UploadChunk(c, &param); err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, nil)
}

// uploadCompleteHandler handles POST /api/v1/sys/file/upload/complete
// @Summary      文件管理完成分片上传
// @Description  访问 /api/v1/sys/file/upload/complete，文件管理完成分片上传
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        body  body  file.ChunkCompleteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/upload/complete [post]
func (h *handler) uploadComplete(c *gin.Context) {
	var param file.ChunkCompleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data, err := h.service.CompleteChunkUpload(c, &param)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

// uploadAbortHandler handles POST /api/v1/sys/file/upload/abort
// @Summary      文件管理取消分片上传
// @Description  访问 /api/v1/sys/file/upload/abort，文件管理取消分片上传
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        body  body  file.ChunkAbortParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/file/upload/abort [post]
func (h *handler) uploadAbort(c *gin.Context) {
	var param file.ChunkAbortParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	if err := h.service.AbortChunkUpload(c, &param); err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, nil)
}
