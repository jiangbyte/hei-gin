package v1

import (
	"hei-gin/sdk/auth"
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
		rolePage,
	)

	// POST /api/v1/sys/role/create
	r.POST("/api/v1/sys/role/create",
		registry.Perm("sys:role:create", "添加角色"),
		log.SysLog("添加角色"),
		middleware.NoRepeat(3000),
		roleCreate,
	)

	// POST /api/v1/sys/role/modify
	r.POST("/api/v1/sys/role/modify",
		registry.Perm("sys:role:modify", "编辑角色"),
		log.SysLog("编辑角色"),
		roleModify,
	)

	// POST /api/v1/sys/role/remove
	r.POST("/api/v1/sys/role/remove",
		registry.Perm("sys:role:remove", "删除角色"),
		log.SysLog("删除角色"),
		roleRemove,
	)

	// GET /api/v1/sys/role/detail
	r.GET("/api/v1/sys/role/detail",
		registry.Perm("sys:role:detail", "角色详情"),
		roleDetail,
	)

	// POST /api/v1/sys/role/grant-permission
	r.POST("/api/v1/sys/role/grant-permission",
		registry.Perm("sys:role:grant-permission", "分配角色权限"),
		log.SysLog("分配角色权限"),
		middleware.NoRepeat(3000),
		roleGrantPermission,
	)

	// POST /api/v1/sys/role/grant-resource
	r.POST("/api/v1/sys/role/grant-resource",
		registry.Perm("sys:role:grant-resource", "分配角色资源"),
		log.SysLog("分配角色资源"),
		middleware.NoRepeat(3000),
		roleGrantResource,
	)

	// GET /api/v1/sys/role/own-permission
	r.GET("/api/v1/sys/role/own-permission",
		registry.Perm("sys:role:own-permission", "角色权限列表"),
		roleOwnPermission,
	)

	// GET /api/v1/sys/role/own-permission-detail
	r.GET("/api/v1/sys/role/own-permission-detail",
		registry.Perm("sys:role:own-permission", "角色权限列表"),
		roleOwnPermissionDetail,
	)

	// GET /api/v1/sys/role/own-resource
	r.GET("/api/v1/sys/role/own-resource",
		registry.Perm("sys:role:own-resource", "角色资源列表"),
		roleOwnResource,
	)
}

// rolePage handles GET /api/v1/sys/role/page
func rolePage(c *gin.Context) {
	var param role.RolePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	role.RolePage(c, &param)
}

// roleCreate handles POST /api/v1/sys/role/create
func roleCreate(c *gin.Context) {
	var vo role.RoleVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	if vo.Code == "" || vo.Name == "" || vo.Category == "" {
		result.Failure(c, "角色编码、名称、类别不能为空", 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	role.RoleCreate(c, &vo, userID)
	result.Success(c, nil)
}

// roleModify handles POST /api/v1/sys/role/modify
func roleModify(c *gin.Context) {
	var vo role.RoleVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	if vo.ID == "" {
		result.Failure(c, "ID不能为空", 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	role.RoleModify(c, &vo, userID)
	result.Success(c, nil)
}

// roleRemove handles POST /api/v1/sys/role/remove
func roleRemove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	role.RoleRemove(c, param.IDs)
	result.Success(c, nil)
}

// roleDetail handles GET /api/v1/sys/role/detail
func roleDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Success(c, nil)
		return
	}

	vo := role.RoleDetail(c, id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}

// roleGrantPermission handles POST /api/v1/sys/role/grant-permission
func roleGrantPermission(c *gin.Context) {
	var param role.GrantPermissionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	role.RoleGrantPermissions(c, param.RoleID, param.Permissions, userID)
	result.Success(c, nil)
}

// roleGrantResource handles POST /api/v1/sys/role/grant-resource
func roleGrantResource(c *gin.Context) {
	var param role.GrantResourceParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	role.RoleGrantResources(c, param.RoleID, param.ResourceIDs, param.Permissions)
	result.Success(c, nil)
}

// roleOwnPermission handles GET /api/v1/sys/role/own-permission
func roleOwnPermission(c *gin.Context) {
	roleID := c.Query("role_id")
	codes := role.RoleOwnPermissionCodes(c, roleID)
	result.Success(c, codes)
}

// roleOwnPermissionDetail handles GET /api/v1/sys/role/own-permission-detail
func roleOwnPermissionDetail(c *gin.Context) {
	roleID := c.Query("role_id")
	details := role.RoleOwnPermissionDetails(c, roleID)
	result.Success(c, details)
}

// roleOwnResource handles GET /api/v1/sys/role/own-resource
func roleOwnResource(c *gin.Context) {
	roleID := c.Query("role_id")
	ids := role.RoleOwnResourceIDs(c, roleID)
	result.Success(c, ids)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
