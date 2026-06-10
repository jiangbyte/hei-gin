package v1

import (
	"hei-gin/sdk/auth"
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
		orgPage,
	)

	// GET /api/v1/sys/org/tree
	r.GET("/api/v1/sys/org/tree",
		registry.Perm("sys:org:tree", "组织树"),
		orgTree,
	)

	// POST /api/v1/sys/org/create
	r.POST("/api/v1/sys/org/create",
		registry.Perm("sys:org:create", "添加组织"),
		log.SysLog("添加组织"),
		middleware.NoRepeat(3000),
		orgCreate,
	)

	// POST /api/v1/sys/org/modify
	r.POST("/api/v1/sys/org/modify",
		registry.Perm("sys:org:modify", "编辑组织"),
		log.SysLog("编辑组织"),
		orgModify,
	)

	// POST /api/v1/sys/org/remove
	r.POST("/api/v1/sys/org/remove",
		registry.Perm("sys:org:remove", "删除组织"),
		log.SysLog("删除组织"),
		orgRemove,
	)

	// GET /api/v1/sys/org/detail
	r.GET("/api/v1/sys/org/detail",
		registry.Perm("sys:org:detail", "组织详情"),
		orgDetail,
	)
}

// orgPage handles GET /api/v1/sys/org/page
func orgPage(c *gin.Context) {
	var param org.OrgPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	org.Page(c, &param)
}

// orgTree handles GET /api/v1/sys/org/tree
func orgTree(c *gin.Context) {
	var param org.OrgTreeParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := org.Tree(c, &param)
	result.Success(c, data)
}

// orgCreate handles POST /api/v1/sys/org/create
func orgCreate(c *gin.Context) {
	var vo org.OrgVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	if vo.Code == "" || vo.Name == "" || vo.Category == "" {
		result.Failure(c, "组织编码、名称、类别不能为空", 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	org.Create(c, &vo, userID)
	result.Success(c, nil)
}

// orgModify handles POST /api/v1/sys/org/modify
func orgModify(c *gin.Context) {
	var vo org.OrgVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	if vo.ID == "" {
		result.Failure(c, "ID不能为空", 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	org.Modify(c, &vo, userID)
	result.Success(c, nil)
}

// orgRemove handles POST /api/v1/sys/org/remove
func orgRemove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	org.Remove(c, param.IDs)
	result.Success(c, nil)
}

// orgDetail handles GET /api/v1/sys/org/detail
func orgDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Success(c, nil)
		return
	}

	vo := org.Detail(c, id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
