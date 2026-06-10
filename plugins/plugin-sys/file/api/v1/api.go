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

func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/sys/file/upload",
		registry.Perm("sys:file:upload", "上传文件"),
		log.SysLog("上传文件"),
		fileUpload,
	)

	r.POST("/api/v1/sys/file/upload/init",
		registry.Perm("sys:file:upload", "分片上传-初始化"),
		log.SysLog("分片上传-初始化"),
		fileUploadInit,
	)
	r.POST("/api/v1/sys/file/upload/chunk",
		registry.Perm("sys:file:upload", "分片上传-上传分片"),
		log.SysLog("分片上传-上传分片"),
		fileUploadChunk,
	)
	r.POST("/api/v1/sys/file/upload/complete",
		registry.Perm("sys:file:upload", "分片上传-完成"),
		log.SysLog("分片上传-完成"),
		fileUploadComplete,
	)
	r.POST("/api/v1/sys/file/upload/abort",
		registry.Perm("sys:file:upload", "分片上传-取消"),
		log.SysLog("分片上传-取消"),
		fileUploadAbort,
	)

	r.GET("/api/v1/sys/file/download",
		registry.Perm("sys:file:download", "下载文件"),
		fileDownload,
	)
	r.GET("/api/v1/sys/file/page",
		registry.Perm("sys:file:page", "文件分页"),
		filePage,
	)
	r.GET("/api/v1/sys/file/detail",
		registry.Perm("sys:file:detail", "文件详情"),
		fileDetail,
	)
	r.POST("/api/v1/sys/file/remove",
		registry.Perm("sys:file:remove", "删除文件"),
		log.SysLog("删除文件"),
		fileRemove,
	)
	r.POST("/api/v1/sys/file/remove-absolute",
		registry.Perm("sys:file:remove", "物理删除文件"),
		log.SysLog("物理删除文件"),
		fileRemoveAbsolute,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	r.POST("/api/v1/c/file/upload",
		middleware.HeiClientCheckLogin(),
		log.SysLog("C端上传文件"),
		clientFileUpload,
	)
}

func fileUpload(c *gin.Context) {
	data, err := file.Upload(c)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

func clientFileUpload(c *gin.Context) {
	data, err := file.Upload(c)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

func fileDownload(c *gin.Context) {
	id := c.Query("id")
	if err := file.Download(c, id); err != nil {
		result.Failure(c, err.Error(), 400)
	}
}

func filePage(c *gin.Context) {
	var param file.FilePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	file.Page(c, &param)
}

func fileDetail(c *gin.Context) {
	id := c.Query("id")
	vo := file.Detail(c, id)
	result.Success(c, vo)
}

func fileRemove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if err := file.Remove(c, param.IDs); err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, nil)
}

func fileRemoveAbsolute(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if err := file.RemoveAbsolute(c, param.IDs); err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, nil)
}

func fileUploadInit(c *gin.Context) {
	var param file.ChunkUploadInitParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	data, err := file.InitChunkUpload(c, &param)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

func fileUploadChunk(c *gin.Context) {
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
	if err := file.UploadChunk(c, &param); err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, nil)
}

func fileUploadComplete(c *gin.Context) {
	var param file.ChunkCompleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	data, err := file.CompleteChunkUpload(c, &param)
	if err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, data)
}

func fileUploadAbort(c *gin.Context) {
	var param file.ChunkAbortParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if err := file.AbortChunkUpload(c, &param); err != nil {
		result.Failure(c, err.Error(), 400)
		return
	}
	result.Success(c, nil)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}
