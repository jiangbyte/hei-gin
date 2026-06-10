package v1

import (
	"hei-gin/sdk/auth"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	notice "hei-gin/plugins/plugin-sys/notice"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all notice routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/notice/page
	r.GET("/api/v1/sys/notice/page",
		registry.Perm("sys:notice:page", "通知分页"),
		pageHandler,
	)

	// POST /api/v1/sys/notice/create
	r.POST("/api/v1/sys/notice/create",
		registry.Perm("sys:notice:create", "添加通知"),
		log.SysLog("添加通知"),
		middleware.NoRepeat(3000),
		createHandler,
	)

	// POST /api/v1/sys/notice/modify
	r.POST("/api/v1/sys/notice/modify",
		registry.Perm("sys:notice:modify", "编辑通知"),
		log.SysLog("编辑通知"),
		modifyHandler,
	)

	// POST /api/v1/sys/notice/remove
	r.POST("/api/v1/sys/notice/remove",
		registry.Perm("sys:notice:remove", "删除通知"),
		log.SysLog("删除通知"),
		deleteHandler,
	)

	// GET /api/v1/sys/notice/detail
	r.GET("/api/v1/sys/notice/detail",
		registry.Perm("sys:notice:detail", "通知详情"),
		detailHandler,
	)
}

// pageHandler handles GET /api/v1/sys/notice/page
func pageHandler(c *gin.Context) {
	var param notice.NoticePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	notice.Page(c, &param)
}

// createHandler handles POST /api/v1/sys/notice/create
func createHandler(c *gin.Context) {
	var vo notice.NoticeVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	notice.Create(c, &vo, userID)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/notice/modify
func modifyHandler(c *gin.Context) {
	var vo notice.NoticeVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	notice.Modify(c, &vo, userID)
	result.Success(c, nil)
}

// deleteHandler handles POST /api/v1/sys/notice/remove
func deleteHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	notice.Remove(c, param.IDs)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/notice/detail
func detailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := notice.Detail(c, id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}

// RegisterPublicRoutes registers public notice routes (no auth required).
func RegisterPublicRoutes(r *gin.Engine) {
	// GET /api/v1/public/c/notice/latest — latest published notices
	r.GET("/api/v1/public/c/notice/latest", latestHandler)

	// GET /api/v1/public/c/notice/page — paginated published notices
	r.GET("/api/v1/public/c/notice/page", pagePublicHandler)

	// GET /api/v1/public/c/notice/detail — published notice detail
	r.GET("/api/v1/public/c/notice/detail", detailPublicHandler)
}

// pagePublicHandler handles GET /api/v1/public/c/notice/page
func pagePublicHandler(c *gin.Context) {
	var param notice.NoticePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	notice.PublicPage(c, &param)
}

// detailPublicHandler handles GET /api/v1/public/c/notice/detail
func detailPublicHandler(c *gin.Context) {
	id := c.Query("id")
	vo := notice.PublicDetail(c, id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}

// latestHandler handles GET /api/v1/public/c/notice/latest
func latestHandler(c *gin.Context) {
	var param notice.NoticeLatestParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if param.Size < 1 {
		param.Size = 5
	}
	if param.Size > 20 {
		param.Size = 20
	}
	data := notice.Latest(c, &param)
	result.Success(c, data)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterPublicRoutes)
}
