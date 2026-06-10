package v1

import (
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/registry"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
	logPackage "hei-gin/plugins/plugin-sys/log"

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
func pageHandler(c *gin.Context) {
	var param logPackage.LogPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	logPackage.LogPage(c, &param)
}

// createHandler handles POST /api/v1/sys/log/create
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
func detailHandler(c *gin.Context) {
	vo := logPackage.LogDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// deleteByCategoryHandler handles POST /api/v1/sys/log/delete-by-category
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
func visLineChartHandler(c *gin.Context) {
	data := logPackage.LogLoginBarChart(c)
	result.Success(c, data)
}

// visPieChartHandler handles GET /api/v1/sys/log/vis/pie-chart-data
func visPieChartHandler(c *gin.Context) {
	data := logPackage.LogLoginPieChart(c)
	result.Success(c, data)
}

// opBarChartHandler handles GET /api/v1/sys/log/op/bar-chart-data
func opBarChartHandler(c *gin.Context) {
	data := logPackage.LogOpBarChart(c)
	result.Success(c, data)
}

// opPieChartHandler handles GET /api/v1/sys/log/op/pie-chart-data
func opPieChartHandler(c *gin.Context) {
	data := logPackage.LogOpPieChart(c)
	result.Success(c, data)
}
