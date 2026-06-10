package v1

import (
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
		removeHandler,
	)

	// GET /api/v1/sys/notice/detail
	r.GET("/api/v1/sys/notice/detail",
		registry.Perm("sys:notice:detail", "通知详情"),
		detailHandler,
	)
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

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterPublicRoutes)
}

// pageHandler handles GET /api/v1/sys/notice/page
func pageHandler(c *gin.Context) {
	var param notice.NoticePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	notice.NoticePage(c, &param)
}

// createHandler handles POST /api/v1/sys/notice/create
func createHandler(c *gin.Context) {
	var vo notice.NoticeVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	notice.NoticeCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/notice/modify
func modifyHandler(c *gin.Context) {
	var vo notice.NoticeVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	notice.NoticeModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/notice/remove
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	notice.NoticeRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/notice/detail
func detailHandler(c *gin.Context) {
	vo := notice.NoticeDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// pagePublicHandler handles GET /api/v1/public/c/notice/page
func pagePublicHandler(c *gin.Context) {
	var param notice.NoticePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	notice.NoticePublicPage(c, &param)
}

// detailPublicHandler handles GET /api/v1/public/c/notice/detail
func detailPublicHandler(c *gin.Context) {
	vo := notice.NoticePublicDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// latestHandler handles GET /api/v1/public/c/notice/latest
func latestHandler(c *gin.Context) {
	var param notice.NoticeLatestParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	result.Success(c, notice.NoticeLatest(c, &param))
}
