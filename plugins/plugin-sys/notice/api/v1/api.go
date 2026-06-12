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
// @Summary      通知公告分页查询
// @Description  访问 /api/v1/sys/notice/page，通知公告分页查询
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Param        query  query  notice.NoticePageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/notice/page [get]
func pageHandler(c *gin.Context) {
	var param notice.NoticePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	notice.NoticePage(c, &param)
}

// createHandler handles POST /api/v1/sys/notice/create
// @Summary      通知公告创建
// @Description  访问 /api/v1/sys/notice/create，通知公告创建
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Param        body  body  notice.NoticeVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/notice/create [post]
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
// @Summary      通知公告修改
// @Description  访问 /api/v1/sys/notice/modify，通知公告修改
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Param        body  body  notice.NoticeVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/notice/modify [post]
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
// @Summary      通知公告删除
// @Description  访问 /api/v1/sys/notice/remove，通知公告删除
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/notice/remove [post]
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
// @Summary      通知公告详情查询
// @Description  访问 /api/v1/sys/notice/detail，通知公告详情查询
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/notice/detail [get]
func detailHandler(c *gin.Context) {
	vo := notice.NoticeDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// pagePublicHandler handles GET /api/v1/public/c/notice/page
// @Summary      通知公告分页查询
// @Description  访问 /api/v1/public/c/notice/page，通知公告分页查询
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Param        query  query  notice.NoticePageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/c/notice/page [get]
func pagePublicHandler(c *gin.Context) {
	var param notice.NoticePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	notice.NoticePublicPage(c, &param)
}

// detailPublicHandler handles GET /api/v1/public/c/notice/detail
// @Summary      通知公告详情查询
// @Description  访问 /api/v1/public/c/notice/detail，通知公告详情查询
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/c/notice/detail [get]
func detailPublicHandler(c *gin.Context) {
	vo := notice.NoticePublicDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// latestHandler handles GET /api/v1/public/c/notice/latest
// @Summary      通知公告最新公告
// @Description  访问 /api/v1/public/c/notice/latest，通知公告最新公告
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Param        query  query  notice.NoticeLatestParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/c/notice/latest [get]
func latestHandler(c *gin.Context) {
	var param notice.NoticeLatestParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	result.Success(c, notice.NoticeLatest(c, &param))
}
