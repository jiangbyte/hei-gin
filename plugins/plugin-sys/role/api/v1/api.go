package v1

import (
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	role "hei-gin/plugins/plugin-sys/role"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all role routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/role/page
	r.GET("/api/v1/sys/role/page",
		registry.Perm("sys:role:page", "角色分页"),
		pageHandler,
	)

	// POST /api/v1/sys/role/create
	r.POST("/api/v1/sys/role/create",
		registry.Perm("sys:role:create", "添加角色"),
		log.SysLog("添加角色"),
		middleware.NoRepeat(3000),
		createHandler,
	)

	// POST /api/v1/sys/role/modify
	r.POST("/api/v1/sys/role/modify",
		registry.Perm("sys:role:modify", "编辑角色"),
		log.SysLog("编辑角色"),
		modifyHandler,
	)

	// POST /api/v1/sys/role/remove
	r.POST("/api/v1/sys/role/remove",
		registry.Perm("sys:role:remove", "删除角色"),
		log.SysLog("删除角色"),
		removeHandler,
	)

	// GET /api/v1/sys/role/detail
	r.GET("/api/v1/sys/role/detail",
		registry.Perm("sys:role:detail", "角色详情"),
		detailHandler,
	)

	// POST /api/v1/sys/role/grant-permission
	r.POST("/api/v1/sys/role/grant-permission",
		registry.Perm("sys:role:grant-permission", "分配角色权限"),
		log.SysLog("分配角色权限"),
		middleware.NoRepeat(3000),
		grantPermissionHandler,
	)

	// POST /api/v1/sys/role/grant-resource
	r.POST("/api/v1/sys/role/grant-resource",
		registry.Perm("sys:role:grant-resource", "分配角色资源"),
		log.SysLog("分配角色资源"),
		middleware.NoRepeat(3000),
		grantResourceHandler,
	)

	// GET /api/v1/sys/role/own-permission
	r.GET("/api/v1/sys/role/own-permission",
		registry.Perm("sys:role:own-permission", "角色权限列表"),
		ownPermissionHandler,
	)

	// GET /api/v1/sys/role/own-permission-detail
	r.GET("/api/v1/sys/role/own-permission-detail",
		registry.Perm("sys:role:own-permission-detail", "角色权限详情"),
		ownPermissionDetailHandler,
	)

	// GET /api/v1/sys/role/own-resource
	r.GET("/api/v1/sys/role/own-resource",
		registry.Perm("sys:role:own-resource", "角色资源列表"),
		ownResourceHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/role/page
func pageHandler(c *gin.Context) {
	var param role.RolePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	role.RolePage(c, &param)
}

// createHandler handles POST /api/v1/sys/role/create
func createHandler(c *gin.Context) {
	var vo role.RoleVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	role.RoleCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/role/modify
func modifyHandler(c *gin.Context) {
	var vo role.RoleVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	role.RoleModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/role/remove
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	role.RoleRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/role/detail
func detailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := role.RoleDetail(c, id)
	result.Success(c, vo)
}

// grantPermissionHandler handles POST /api/v1/sys/role/grant-permission
func grantPermissionHandler(c *gin.Context) {
	var param role.GrantPermissionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	role.RoleGrantPermissions(c, &param)
	result.Success(c, nil)
}

// grantResourceHandler handles POST /api/v1/sys/role/grant-resource
func grantResourceHandler(c *gin.Context) {
	var param role.GrantResourceParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	role.RoleGrantResources(c, &param)
	result.Success(c, nil)
}

// ownPermissionHandler handles GET /api/v1/sys/role/own-permission
func ownPermissionHandler(c *gin.Context) {
	codes := role.RoleOwnPermissionCodes(c, c.Query("role_id"))
	result.Success(c, codes)
}

// ownPermissionDetailHandler handles GET /api/v1/sys/role/own-permission-detail
func ownPermissionDetailHandler(c *gin.Context) {
	details := role.RoleOwnPermissionDetails(c, c.Query("role_id"))
	result.Success(c, details)
}

// ownResourceHandler handles GET /api/v1/sys/role/own-resource
func ownResourceHandler(c *gin.Context) {
	ids := role.RoleOwnResourceIDs(c, c.Query("role_id"))
	result.Success(c, ids)
}
