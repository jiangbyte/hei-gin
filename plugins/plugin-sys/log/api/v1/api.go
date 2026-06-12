package v1

import (
	logPackage "hei-gin/plugins/plugin-sys/log"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all sys/log routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/log/page
	r.GET("/api/v1/sys/log/page",
		registry.Perm("sys:log:page", "日志分页"),
		pageHandler,
	)

	// POST /api/v1/sys/log/create
	r.POST("/api/v1/sys/log/create",
		registry.Perm("sys:log:create", "添加日志"),
		createHandler,
	)

	// POST /api/v1/sys/log/modify
	r.POST("/api/v1/sys/log/modify",
		registry.Perm("sys:log:modify", "编辑日志"),
		modifyHandler,
	)

	// POST /api/v1/sys/log/remove
	r.POST("/api/v1/sys/log/remove",
		registry.Perm("sys:log:remove", "删除日志"),
		log.SysLog("删除操作日志"),
		removeHandler,
	)

	// GET /api/v1/sys/log/detail
	r.GET("/api/v1/sys/log/detail",
		registry.Perm("sys:log:detail", "日志详情"),
		detailHandler,
	)

	// POST /api/v1/sys/log/delete-by-category
	r.POST("/api/v1/sys/log/delete-by-category",
		registry.Perm("sys:log:remove", "删除日志"),
		middleware.NoRepeat(5000),
		deleteByCategoryHandler,
	)

	// GET /api/v1/sys/log/vis/line-chart-data
	r.GET("/api/v1/sys/log/vis/line-chart-data",
		registry.Perm("sys:log:page", "日志分页"),
		visLineChartHandler,
	)

	// GET /api/v1/sys/log/vis/pie-chart-data
	r.GET("/api/v1/sys/log/vis/pie-chart-data",
		registry.Perm("sys:log:page", "日志分页"),
		visPieChartHandler,
	)

	// GET /api/v1/sys/log/op/bar-chart-data
	r.GET("/api/v1/sys/log/op/bar-chart-data",
		registry.Perm("sys:log:page", "日志分页"),
		opBarChartHandler,
	)

	// GET /api/v1/sys/log/op/pie-chart-data
	r.GET("/api/v1/sys/log/op/pie-chart-data",
		registry.Perm("sys:log:page", "日志分页"),
		opPieChartHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/log/page
// @Summary      日志管理分页查询
// @Description  访问 /api/v1/sys/log/page，日志管理分页查询
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Param        query  query  logPackage.LogPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/page [get]
func pageHandler(c *gin.Context) {
	var param logPackage.LogPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	logPackage.LogPage(c, &param)
}

// createHandler handles POST /api/v1/sys/log/create
// @Summary      日志管理创建
// @Description  访问 /api/v1/sys/log/create，日志管理创建
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Param        body  body  logPackage.LogVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/create [post]
func createHandler(c *gin.Context) {
	var vo logPackage.LogVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	logPackage.LogCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/log/modify
// @Summary      日志管理修改
// @Description  访问 /api/v1/sys/log/modify，日志管理修改
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Param        body  body  logPackage.LogVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/modify [post]
func modifyHandler(c *gin.Context) {
	var vo logPackage.LogVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	logPackage.LogModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/log/remove
// @Summary      日志管理删除
// @Description  访问 /api/v1/sys/log/remove，日志管理删除
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/remove [post]
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	logPackage.LogRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/log/detail
// @Summary      日志管理详情查询
// @Description  访问 /api/v1/sys/log/detail，日志管理详情查询
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/detail [get]
func detailHandler(c *gin.Context) {
	vo := logPackage.LogDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// deleteByCategoryHandler handles POST /api/v1/sys/log/delete-by-category
// @Summary      日志管理按分类删除
// @Description  访问 /api/v1/sys/log/delete-by-category，日志管理按分类删除
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Param        body  body  logPackage.LogDeleteByCategoryParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/delete-by-category [post]
func deleteByCategoryHandler(c *gin.Context) {
	var param logPackage.LogDeleteByCategoryParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	logPackage.LogDeleteByCategory(c, &param)
	result.Success(c, nil)
}

// visLineChartHandler handles GET /api/v1/sys/log/vis/line-chart-data
// @Summary      日志管理登录折线图数据
// @Description  访问 /api/v1/sys/log/vis/line-chart-data，日志管理登录折线图数据
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/vis/line-chart-data [get]
func visLineChartHandler(c *gin.Context) {
	data := logPackage.LogLoginBarChart(c)
	result.Success(c, data)
}

// visPieChartHandler handles GET /api/v1/sys/log/vis/pie-chart-data
// @Summary      日志管理登录饼图数据
// @Description  访问 /api/v1/sys/log/vis/pie-chart-data，日志管理登录饼图数据
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/vis/pie-chart-data [get]
func visPieChartHandler(c *gin.Context) {
	data := logPackage.LogLoginPieChart(c)
	result.Success(c, data)
}

// opBarChartHandler handles GET /api/v1/sys/log/op/bar-chart-data
// @Summary      日志管理操作柱状图数据
// @Description  访问 /api/v1/sys/log/op/bar-chart-data，日志管理操作柱状图数据
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/op/bar-chart-data [get]
func opBarChartHandler(c *gin.Context) {
	data := logPackage.LogOpBarChart(c)
	result.Success(c, data)
}

// opPieChartHandler handles GET /api/v1/sys/log/op/pie-chart-data
// @Summary      日志管理操作饼图数据
// @Description  访问 /api/v1/sys/log/op/pie-chart-data，日志管理操作饼图数据
// @Tags         日志管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/log/op/pie-chart-data [get]
func opPieChartHandler(c *gin.Context) {
	data := logPackage.LogOpPieChart(c)
	result.Success(c, data)
}
