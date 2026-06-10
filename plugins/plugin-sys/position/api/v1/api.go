package v1

import (
	"hei-gin/sdk/auth"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	position "hei-gin/plugins/plugin-sys/position"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all position routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/position/page
	r.GET("/api/v1/sys/position/page",
		registry.Perm("sys:position:page", "岗位分页"),
		log.SysLog("查看职位列表"),
		pageHandler,
	)

	// POST /api/v1/sys/position/create
	r.POST("/api/v1/sys/position/create",
		registry.Perm("sys:position:create", "添加岗位"),
		log.SysLog("添加职位"),
		middleware.NoRepeat(3000),
		createHandler,
	)

	// POST /api/v1/sys/position/modify
	r.POST("/api/v1/sys/position/modify",
		registry.Perm("sys:position:modify", "编辑岗位"),
		log.SysLog("编辑职位"),
		modifyHandler,
	)

	// POST /api/v1/sys/position/remove
	r.POST("/api/v1/sys/position/remove",
		registry.Perm("sys:position:remove", "删除岗位"),
		log.SysLog("删除职位"),
		deleteHandler,
	)

	// GET /api/v1/sys/position/detail
	r.GET("/api/v1/sys/position/detail",
		registry.Perm("sys:position:detail", "岗位详情"),
		detailHandler,
	)
}

// pageHandler handles GET /api/v1/sys/position/page
func pageHandler(c *gin.Context) {
	var param position.PositionPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	position.Page(c, &param)
}

// createHandler handles POST /api/v1/sys/position/create
func createHandler(c *gin.Context) {
	var vo position.PositionVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	position.Create(c, &vo, userID)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/position/modify
func modifyHandler(c *gin.Context) {
	var vo position.PositionVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	position.Modify(c, &vo, userID)
	result.Success(c, nil)
}

// deleteHandler handles POST /api/v1/sys/position/remove
func deleteHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	position.Remove(c, param.IDs)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/position/detail
func detailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := position.Detail(c, id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
