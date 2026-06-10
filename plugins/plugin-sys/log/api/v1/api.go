package v1

import (

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/registry"
	middleware "hei-gin/sdk/auth/middleware"
	sysLog "hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	logPackage "hei-gin/plugins/plugin-sys/log"
)

// RegisterRoutes registers all sys/log routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/log/page
	r.GET("/api/v1/sys/log/page",
		registry.Perm("sys:log:page", "日志分页"),
		logPage,
	)

	// POST /api/v1/sys/log/create
	r.POST("/api/v1/sys/log/create",
		registry.Perm("sys:log:create", "添加日志"),
		logCreate,
	)

	// POST /api/v1/sys/log/modify
	r.POST("/api/v1/sys/log/modify",
		registry.Perm("sys:log:modify", "编辑日志"),
		logModify,
	)

	// POST /api/v1/sys/log/remove
	r.POST("/api/v1/sys/log/remove",
		registry.Perm("sys:log:remove", "删除日志"),
		sysLog.SysLog("删除操作日志"),
		logRemove,
	)

	// GET /api/v1/sys/log/detail
	r.GET("/api/v1/sys/log/detail",
		registry.Perm("sys:log:detail", "日志详情"),
		logDetail,
	)

	// POST /api/v1/sys/log/delete-by-category
	r.POST("/api/v1/sys/log/delete-by-category",
		registry.Perm("sys:log:remove", "删除日志"),
		middleware.NoRepeat(5000),
		logDeleteByCategory,
	)

	// GET /api/v1/sys/log/vis/line-chart-data
	r.GET("/api/v1/sys/log/vis/line-chart-data",
		registry.Perm("sys:log:page", "日志分页"),
		logVisLineChart,
	)

	// GET /api/v1/sys/log/vis/pie-chart-data
	r.GET("/api/v1/sys/log/vis/pie-chart-data",
		registry.Perm("sys:log:page", "日志分页"),
		logVisPieChart,
	)

	// GET /api/v1/sys/log/op/bar-chart-data
	r.GET("/api/v1/sys/log/op/bar-chart-data",
		registry.Perm("sys:log:page", "日志分页"),
		logOpBarChart,
	)

	// GET /api/v1/sys/log/op/pie-chart-data
	r.GET("/api/v1/sys/log/op/pie-chart-data",
		registry.Perm("sys:log:page", "日志分页"),
		logOpPieChart,
	)
}

// logPage handles GET /api/v1/sys/log/page
func logPage(c *gin.Context) {
	param := &logPackage.LogPageParam{}
	if err := c.ShouldBindQuery(param); err != nil {
		param.Current = 1
		param.Size = 10
	}
	logPackage.Page(c, param)
}

// logCreate handles POST /api/v1/sys/log/create
func logCreate(c *gin.Context) {
	vo := &logPackage.LogVO{}
	if err := c.ShouldBindJSON(vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	logPackage.Create(c, vo, userID)
	result.Success(c, nil)
}

// logModify handles POST /api/v1/sys/log/modify
func logModify(c *gin.Context) {
	vo := &logPackage.LogVO{}
	if err := c.ShouldBindJSON(vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	if vo.ID == "" {
		result.Failure(c, "id不能为空", 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	logPackage.Modify(c, vo, userID)
	result.Success(c, nil)
}

// logRemove handles POST /api/v1/sys/log/remove
func logRemove(c *gin.Context) {
	param := &utils.IdsParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	if len(param.IDs) == 0 {
		result.Failure(c, "ids不能为空", 400)
		return
	}

	logPackage.Remove(c, param.IDs)
	result.Success(c, nil)
}

// logDetail handles GET /api/v1/sys/log/detail
func logDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Failure(c, "id不能为空", 400)
		return
	}

	data := logPackage.Detail(c, id)
	if data == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, data)
}

// logDeleteByCategory handles POST /api/v1/sys/log/delete-by-category
func logDeleteByCategory(c *gin.Context) {
	param := &logPackage.LogDeleteByCategoryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	if param.Category == "" {
		result.Failure(c, "category不能为空", 400)
		return
	}

	logPackage.DeleteByCategory(c, param)
	result.Success(c, nil)
}

// logVisLineChart handles GET /api/v1/sys/log/vis/line-chart-data
func logVisLineChart(c *gin.Context) {
	data := logPackage.VisLineChart(c)
	result.Success(c, data)
}

// logVisPieChart handles GET /api/v1/sys/log/vis/pie-chart-data
func logVisPieChart(c *gin.Context) {
	data := logPackage.VisPieChart(c)
	result.Success(c, data)
}

// logOpBarChart handles GET /api/v1/sys/log/op/bar-chart-data
func logOpBarChart(c *gin.Context) {
	data := logPackage.OpBarChart(c)
	result.Success(c, data)
}

// logOpPieChart handles GET /api/v1/sys/log/op/pie-chart-data
func logOpPieChart(c *gin.Context) {
	data := logPackage.OpPieChart(c)
	result.Success(c, data)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
