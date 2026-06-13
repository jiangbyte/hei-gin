package v1

import (
	home "hei-gin/plugins/plugin-sys/home"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *home.Service
}

var defaultHandler = newHandler(home.DefaultModule)

func newHandler(module *home.Module) *handler {
	return &handler{service: module.Service()}
}

// RegisterRoutes registers all home routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/home
	r.GET("/api/v1/sys/home",
		middleware.CheckLogin(auth.Business),
		defaultHandler.get,
	)

	// POST /api/v1/sys/home/quick-actions/add
	r.POST("/api/v1/sys/home/quick-actions/add",
		middleware.CheckLogin(auth.Business),
		log.SysLog("添加快捷方式"),
		defaultHandler.addQuickAction,
	)

	// POST /api/v1/sys/home/quick-actions/remove
	r.POST("/api/v1/sys/home/quick-actions/remove",
		middleware.CheckLogin(auth.Business),
		log.SysLog("移除快捷方式"),
		defaultHandler.removeQuickAction,
	)

	// POST /api/v1/sys/home/quick-actions/sort
	r.POST("/api/v1/sys/home/quick-actions/sort",
		middleware.CheckLogin(auth.Business),
		log.SysLog("排序快捷方式"),
		defaultHandler.sortQuickActions,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// getHandler handles GET /api/v1/sys/home
// @Summary      首页配置接口调用
// @Description  访问 /api/v1/sys/home，首页配置接口调用
// @Tags         首页配置
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/home [get]
func (h *handler) get(c *gin.Context) {
	data := h.service.Get(c)
	result.Success(c, data)
}

// addQuickActionHandler handles POST /api/v1/sys/home/quick-actions/add
// @Summary      首页配置新增
// @Description  访问 /api/v1/sys/home/quick-actions/add，首页配置新增
// @Tags         首页配置
// @Accept       json
// @Produce      json
// @Param        body  body  home.AddQuickActionParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/home/quick-actions/add [post]
func (h *handler) addQuickAction(c *gin.Context) {
	var param home.AddQuickActionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.AddQuickAction(c, &param)
	result.Success(c, nil)
}

// removeQuickActionHandler handles POST /api/v1/sys/home/quick-actions/remove
// @Summary      首页配置删除
// @Description  访问 /api/v1/sys/home/quick-actions/remove，首页配置删除
// @Tags         首页配置
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/home/quick-actions/remove [post]
func (h *handler) removeQuickAction(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.RemoveQuickAction(c, &param)
	result.Success(c, nil)
}

// sortQuickActionsHandler handles POST /api/v1/sys/home/quick-actions/sort
// @Summary      首页配置排序
// @Description  访问 /api/v1/sys/home/quick-actions/sort，首页配置排序
// @Tags         首页配置
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/home/quick-actions/sort [post]
func (h *handler) sortQuickActions(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.SortQuickActions(c, &param)
	result.Success(c, nil)
}
