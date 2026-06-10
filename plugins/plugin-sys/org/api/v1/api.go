package v1

import (
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	org "hei-gin/plugins/plugin-sys/org"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all org routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/org/page
	r.GET("/api/v1/sys/org/page",
		registry.Perm("sys:org:page", "组织分页"),
		pageHandler,
	)

	// GET /api/v1/sys/org/tree
	r.GET("/api/v1/sys/org/tree",
		registry.Perm("sys:org:tree", "组织树"),
		treeHandler,
	)

	// POST /api/v1/sys/org/create
	r.POST("/api/v1/sys/org/create",
		registry.Perm("sys:org:create", "添加组织"),
		log.SysLog("添加组织"),
		middleware.NoRepeat(3000),
		createHandler,
	)

	// POST /api/v1/sys/org/modify
	r.POST("/api/v1/sys/org/modify",
		registry.Perm("sys:org:modify", "编辑组织"),
		log.SysLog("编辑组织"),
		modifyHandler,
	)

	// POST /api/v1/sys/org/remove
	r.POST("/api/v1/sys/org/remove",
		registry.Perm("sys:org:remove", "删除组织"),
		log.SysLog("删除组织"),
		removeHandler,
	)

	// GET /api/v1/sys/org/detail
	r.GET("/api/v1/sys/org/detail",
		registry.Perm("sys:org:detail", "组织详情"),
		detailHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/org/page
func pageHandler(c *gin.Context) {
	var param org.OrgPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	org.OrgPage(c, &param)
}

// treeHandler handles GET /api/v1/sys/org/tree
func treeHandler(c *gin.Context) {
	var param org.OrgTreeParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	result.Success(c, org.OrgTree(c, &param))
}

// createHandler handles POST /api/v1/sys/org/create
func createHandler(c *gin.Context) {
	var vo org.OrgVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	org.OrgCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/org/modify
func modifyHandler(c *gin.Context) {
	var vo org.OrgVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	org.OrgModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/org/remove
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	org.OrgRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/org/detail
func detailHandler(c *gin.Context) {
	vo := org.OrgDetail(c, c.Query("id"))
	result.Success(c, vo)
}
