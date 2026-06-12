package v1

import (
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
		removeHandler,
	)

	// GET /api/v1/sys/position/detail
	r.GET("/api/v1/sys/position/detail",
		registry.Perm("sys:position:detail", "岗位详情"),
		detailHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/position/page
// @Summary      岗位管理分页查询
// @Description  访问 /api/v1/sys/position/page，岗位管理分页查询
// @Tags         岗位管理
// @Accept       json
// @Produce      json
// @Param        query  query  position.PositionPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/position/page [get]
func pageHandler(c *gin.Context) {
	var param position.PositionPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	position.PositionPage(c, &param)
}

// createHandler handles POST /api/v1/sys/position/create
// @Summary      岗位管理创建
// @Description  访问 /api/v1/sys/position/create，岗位管理创建
// @Tags         岗位管理
// @Accept       json
// @Produce      json
// @Param        body  body  position.PositionVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/position/create [post]
func createHandler(c *gin.Context) {
	var vo position.PositionVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	position.PositionCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/position/modify
// @Summary      岗位管理修改
// @Description  访问 /api/v1/sys/position/modify，岗位管理修改
// @Tags         岗位管理
// @Accept       json
// @Produce      json
// @Param        body  body  position.PositionVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/position/modify [post]
func modifyHandler(c *gin.Context) {
	var vo position.PositionVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	position.PositionModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/position/remove
// @Summary      岗位管理删除
// @Description  访问 /api/v1/sys/position/remove，岗位管理删除
// @Tags         岗位管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/position/remove [post]
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	position.PositionRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/position/detail
// @Summary      岗位管理详情查询
// @Description  访问 /api/v1/sys/position/detail，岗位管理详情查询
// @Tags         岗位管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/position/detail [get]
func detailHandler(c *gin.Context) {
	vo := position.PositionDetail(c, c.Query("id"))
	result.Success(c, vo)
}
