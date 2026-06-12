package v1

import (
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	"hei-gin/sdk/utils"
	home "hei-gin/plugins/plugin-sys/home"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all home routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/home
	r.GET("/api/v1/sys/home",
		middleware.HeiCheckLogin(),
		getHandler,
	)

	// POST /api/v1/sys/home/quick-actions/add
	r.POST("/api/v1/sys/home/quick-actions/add",
		middleware.HeiCheckLogin(),
		log.SysLog("添加快捷方式"),
		addQuickActionHandler,
	)

	// POST /api/v1/sys/home/quick-actions/remove
	r.POST("/api/v1/sys/home/quick-actions/remove",
		middleware.HeiCheckLogin(),
		log.SysLog("移除快捷方式"),
		removeQuickActionHandler,
	)

	// POST /api/v1/sys/home/quick-actions/sort
	r.POST("/api/v1/sys/home/quick-actions/sort",
		middleware.HeiCheckLogin(),
		log.SysLog("排序快捷方式"),
		sortQuickActionsHandler,
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
func getHandler(c *gin.Context) {
	data := home.HomeGet(c)
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
func addQuickActionHandler(c *gin.Context) {
	var param home.AddQuickActionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	home.HomeAddQuickAction(c, &param)
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
func removeQuickActionHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	home.HomeRemoveQuickAction(c, &param)
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
func sortQuickActionsHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	home.HomeSortQuickActions(c, &param)
	result.Success(c, nil)
}
