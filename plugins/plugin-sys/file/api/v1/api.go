package v1

import (
	"strconv"

	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	file "hei-gin/plugins/plugin-sys/file"

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
func uploadHandler(c *gin.Context) {
	data, err := file.FileUpload(c)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

// clientUploadHandler handles POST /api/v1/c/file/upload
func clientUploadHandler(c *gin.Context) {
	data, err := file.FileUpload(c)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

// downloadHandler handles GET /api/v1/sys/file/download
func downloadHandler(c *gin.Context) {
	if err := file.FileDownload(c, c.Query("id")); err != nil {
		result.Failure(c, err.Error(), 400)
	}
}

// pageHandler handles GET /api/v1/sys/file/page
func pageHandler(c *gin.Context) {
	var param file.FilePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	file.FilePage(c, &param)
}

// detailHandler handles GET /api/v1/sys/file/detail
func detailHandler(c *gin.Context) {
	vo := file.FileDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// removeHandler handles POST /api/v1/sys/file/remove
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
func uploadChunkHandler(c *gin.Context) {
	chunkIndex, _ := strconv.Atoi(c.PostForm("chunk_index"))
	totalChunks, _ := strconv.Atoi(c.PostForm("total_chunks"))
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
