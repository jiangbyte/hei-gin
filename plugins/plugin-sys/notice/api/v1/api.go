package v1

import (
	notice "hei-gin/plugins/plugin-sys/notice"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *notice.Service
}

var defaultHandler = newHandler(notice.DefaultModule)

func newHandler(module *notice.Module) *handler {
	return &handler{service: module.Service()}
}

// RegisterRoutes registers all notice routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/notice/page
	r.GET("/api/v1/sys/notice/page",
		registry.Perm("sys:notice:page", "通知分页"),
		defaultHandler.page,
	)

	// POST /api/v1/sys/notice/create
	r.POST("/api/v1/sys/notice/create",
		registry.Perm("sys:notice:create", "添加通知"),
		log.SysLog("添加通知"),
		middleware.NoRepeat(3000),
		defaultHandler.create,
	)

	// POST /api/v1/sys/notice/modify
	r.POST("/api/v1/sys/notice/modify",
		registry.Perm("sys:notice:modify", "编辑通知"),
		log.SysLog("编辑通知"),
		defaultHandler.modify,
	)

	// POST /api/v1/sys/notice/remove
	r.POST("/api/v1/sys/notice/remove",
		registry.Perm("sys:notice:remove", "删除通知"),
		log.SysLog("删除通知"),
		defaultHandler.remove,
	)

	// GET /api/v1/sys/notice/detail
	r.GET("/api/v1/sys/notice/detail",
		registry.Perm("sys:notice:detail", "通知详情"),
		defaultHandler.detail,
	)
}

// RegisterPublicRoutes registers public notice routes (no auth required).
func RegisterPublicRoutes(r *gin.Engine) {
	// GET /api/v1/public/c/notice/latest — latest published notices
	r.GET("/api/v1/public/c/notice/latest", defaultHandler.latest)

	// GET /api/v1/public/c/notice/page — paginated published notices
	r.GET("/api/v1/public/c/notice/page", defaultHandler.pagePublic)

	// GET /api/v1/public/c/notice/detail — published notice detail
	r.GET("/api/v1/public/c/notice/detail", defaultHandler.detailPublic)
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
func (h *handler) page(c *gin.Context) {
	var param notice.NoticePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
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
func (h *handler) create(c *gin.Context) {
	var vo notice.NoticeVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Create(c, &vo)
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
func (h *handler) modify(c *gin.Context) {
	var vo notice.NoticeVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Modify(c, &vo)
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
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
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
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
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
func (h *handler) pagePublic(c *gin.Context) {
	var param notice.NoticePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.PublicPage(c, &param)
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
func (h *handler) detailPublic(c *gin.Context) {
	vo := h.service.PublicDetail(c, c.Query("id"))
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
func (h *handler) latest(c *gin.Context) {
	var param notice.NoticeLatestParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	result.Success(c, h.service.Latest(c, &param))
}
