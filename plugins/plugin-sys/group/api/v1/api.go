package v1

import (
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	group "hei-gin/plugins/plugin-sys/group"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all group routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/group/page
	r.GET("/api/v1/sys/group/page",
		registry.Perm("sys:group:page", "分组分页"),
		pageHandler,
	)

	// GET /api/v1/sys/group/union-tree
	r.GET("/api/v1/sys/group/union-tree",
		registry.Perm("sys:group:tree", "分组树"),
		unionTreeHandler,
	)

	// GET /api/v1/sys/group/tree
	r.GET("/api/v1/sys/group/tree",
		registry.Perm("sys:group:tree", "分组树"),
		treeHandler,
	)

	// POST /api/v1/sys/group/create
	r.POST("/api/v1/sys/group/create",
		registry.Perm("sys:group:create", "添加分组"),
		log.SysLog("添加用户组"),
		middleware.NoRepeat(3000),
		createHandler,
	)

	// POST /api/v1/sys/group/modify
	r.POST("/api/v1/sys/group/modify",
		registry.Perm("sys:group:modify", "编辑分组"),
		log.SysLog("编辑用户组"),
		modifyHandler,
	)

	// POST /api/v1/sys/group/remove
	r.POST("/api/v1/sys/group/remove",
		registry.Perm("sys:group:remove", "删除分组"),
		log.SysLog("删除用户组"),
		removeHandler,
	)

	// GET /api/v1/sys/group/detail
	r.GET("/api/v1/sys/group/detail",
		registry.Perm("sys:group:detail", "分组详情"),
		detailHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/group/page
func pageHandler(c *gin.Context) {
	var param group.GroupPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	group.GroupPage(c, &param)
}

// unionTreeHandler handles GET /api/v1/sys/group/union-tree
func unionTreeHandler(c *gin.Context) {
	data := group.GroupOptions(c)
	result.Success(c, data)
}

// treeHandler handles GET /api/v1/sys/group/tree
func treeHandler(c *gin.Context) {
	var param group.GroupTreeParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := group.GroupTree(c, &param)
	result.Success(c, data)
}

// createHandler handles POST /api/v1/sys/group/create
func createHandler(c *gin.Context) {
	var vo group.GroupVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	group.GroupCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/group/modify
func modifyHandler(c *gin.Context) {
	var vo group.GroupVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	group.GroupModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/group/remove
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	group.GroupRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/group/detail
func detailHandler(c *gin.Context) {
	vo := group.GroupDetail(c, c.Query("id"))
	result.Success(c, vo)
}
