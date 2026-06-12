package v1

import (
	banner "hei-gin/plugins/plugin-sys/banner"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *banner.Service
}

var defaultHandler = newHandler(banner.DefaultModule)

func newHandler(module *banner.Module) *handler {
	return &handler{service: module.Service()}
}

// RegisterRoutes registers all banner routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/banner/page
	r.GET("/api/v1/sys/banner/page",
		registry.Perm("sys:banner:page", "横幅分页"),
		log.SysLog("获取Banner列表"),
		defaultHandler.page,
	)

	// POST /api/v1/sys/banner/create
	r.POST("/api/v1/sys/banner/create",
		registry.Perm("sys:banner:create", "添加横幅"),
		log.SysLog("添加Banner"),
		middleware.NoRepeat(3000),
		defaultHandler.create,
	)

	// POST /api/v1/sys/banner/modify
	r.POST("/api/v1/sys/banner/modify",
		registry.Perm("sys:banner:modify", "编辑横幅"),
		log.SysLog("编辑Banner"),
		defaultHandler.modify,
	)

	// POST /api/v1/sys/banner/remove
	r.POST("/api/v1/sys/banner/remove",
		registry.Perm("sys:banner:remove", "删除横幅"),
		log.SysLog("删除Banner"),
		defaultHandler.remove,
	)

	// GET /api/v1/sys/banner/detail
	r.GET("/api/v1/sys/banner/detail",
		registry.Perm("sys:banner:detail", "横幅详情"),
		defaultHandler.detail,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// @Summary      横幅分页
// @Description  获取横幅分页列表
// @Tags         横幅管理
// @Accept       json
// @Produce      json
// @Param        current  query  int  false  "页码（默认 1）"
// @Param        size     query  int  false  "每页条数（默认 10，最大 100）"
// @Success      200  {object}  map[string]any  "分页数据"
// @Router       /api/v1/sys/banner/page [get]
func (h *handler) page(c *gin.Context) {
	var param banner.BannerPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
}

// @Summary      添加横幅
// @Description  创建一条新的横幅记录
// @Tags         横幅管理
// @Accept       json
// @Produce      json
// @Param        body  body  banner.BannerVO  true  "横幅数据"
// @Success      200  {object}  map[string]any  "创建成功"
// @Router       /api/v1/sys/banner/create [post]
func (h *handler) create(c *gin.Context) {
	var vo banner.BannerVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Create(c, &vo)
	result.Success(c, nil)
}

// @Summary      编辑横幅
// @Description  修改已有的横幅记录
// @Tags         横幅管理
// @Accept       json
// @Produce      json
// @Param        body  body  banner.BannerVO  true  "横幅数据（id 必填）"
// @Success      200  {object}  map[string]any  "修改成功"
// @Router       /api/v1/sys/banner/modify [post]
func (h *handler) modify(c *gin.Context) {
	var vo banner.BannerVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Modify(c, &vo)
	result.Success(c, nil)
}

// @Summary      删除横幅
// @Description  批量删除横幅记录
// @Tags         横幅管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "待删除 ID 列表"
// @Success      200  {object}  map[string]any  "删除成功"
// @Router       /api/v1/sys/banner/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// @Summary      横幅详情
// @Description  根据 ID 获取横幅详情
// @Tags         横幅管理
// @Accept       json
// @Produce      json
// @Param        id   query  string  true  "横幅 ID"
// @Success      200  {object}  map[string]any
// @Router       /api/v1/sys/banner/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}
