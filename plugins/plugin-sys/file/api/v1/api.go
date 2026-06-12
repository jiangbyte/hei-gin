package v1

import (
	"strconv"

	file "hei-gin/plugins/plugin-sys/file"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all admin file routes.
func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/sys/file/upload",
		registry.Perm("sys:file:upload", "上传文件"),
		log.SysLog("上传文件"),
		uploadHandler,
	)

	r.POST("/api/v1/sys/file/upload/init",
		registry.Perm("sys:file:upload", "分片上传-初始化"),
		log.SysLog("分片上传-初始化"),
		uploadInitHandler,
	)
	r.POST("/api/v1/sys/file/upload/chunk",
		registry.Perm("sys:file:upload", "分片上传-上传分片"),
		log.SysLog("分片上传-上传分片"),
		uploadChunkHandler,
	)
	r.POST("/api/v1/sys/file/upload/complete",
		registry.Perm("sys:file:upload", "分片上传-完成"),
		log.SysLog("分片上传-完成"),
		uploadCompleteHandler,
	)
	r.POST("/api/v1/sys/file/upload/abort",
		registry.Perm("sys:file:upload", "分片上传-取消"),
		log.SysLog("分片上传-取消"),
		uploadAbortHandler,
	)

	r.GET("/api/v1/sys/file/download",
		registry.Perm("sys:file:download", "下载文件"),
		downloadHandler,
	)
	r.GET("/api/v1/sys/file/page",
		registry.Perm("sys:file:page", "文件分页"),
		pageHandler,
	)
	r.GET("/api/v1/sys/file/detail",
		registry.Perm("sys:file:detail", "文件详情"),
		detailHandler,
	)
	r.POST("/api/v1/sys/file/remove",
		registry.Perm("sys:file:remove", "删除文件"),
		log.SysLog("删除文件"),
		removeHandler,
	)
	r.POST("/api/v1/sys/file/remove-absolute",
		registry.Perm("sys:file:remove", "物理删除文件"),
		log.SysLog("物理删除文件"),
		removeAbsoluteHandler,
	)
}

// RegisterClientRoutes registers client file routes.
func RegisterClientRoutes(r *gin.Engine) {
	r.POST("/api/v1/c/file/upload",
		middleware.HeiClientCheckLogin(),
		log.SysLog("C端上传文件"),
		clientUploadHandler,
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
func uploadHandler(c *gin.Context) {
	data, err := file.FileUpload(c)
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
func clientUploadHandler(c *gin.Context) {
	data, err := file.FileUpload(c)
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
func downloadHandler(c *gin.Context) {
	if err := file.FileDownload(c, c.Query("id")); err != nil {
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
func pageHandler(c *gin.Context) {
	var param file.FilePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	file.FilePage(c, &param)
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
func detailHandler(c *gin.Context) {
	vo := file.FileDetail(c, c.Query("id"))
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
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	file.FileRemove(c, &param)
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
func removeAbsoluteHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	file.FileRemoveAbsolute(c, &param)
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
func uploadInitHandler(c *gin.Context) {
	var param file.ChunkUploadInitParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data, err := file.FileInitChunkUpload(c, &param)
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
func uploadChunkHandler(c *gin.Context) {
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

	if err := file.FileUploadChunk(c, &param); err != nil {
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
func uploadCompleteHandler(c *gin.Context) {
	var param file.ChunkCompleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data, err := file.FileCompleteChunkUpload(c, &param)
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
func uploadAbortHandler(c *gin.Context) {
	var param file.ChunkAbortParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	if err := file.FileAbortChunkUpload(c, &param); err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, nil)
}
