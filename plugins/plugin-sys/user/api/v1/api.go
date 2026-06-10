package v1

import (
	"hei-gin/sdk/registry"
	middleware "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	user "hei-gin/plugins/plugin-sys/user"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all user routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/user/page
	r.GET("/api/v1/sys/user/page",
		registry.Perm("sys:user:page", "用户分页"),
		pageHandler,
	)

	// POST /api/v1/sys/user/create
	r.POST("/api/v1/sys/user/create",
		registry.Perm("sys:user:create", "添加用户"),
		log.SysLog("添加用户"),
		middleware.NoRepeat(3000),
		createHandler,
	)

	// POST /api/v1/sys/user/modify
	r.POST("/api/v1/sys/user/modify",
		registry.Perm("sys:user:modify", "编辑用户"),
		log.SysLog("编辑用户"),
		modifyHandler,
	)

	// POST /api/v1/sys/user/remove
	r.POST("/api/v1/sys/user/remove",
		registry.Perm("sys:user:remove", "删除用户"),
		log.SysLog("删除用户"),
		removeHandler,
	)

	// GET /api/v1/sys/user/detail
	r.GET("/api/v1/sys/user/detail",
		registry.Perm("sys:user:detail", "用户详情"),
		detailHandler,
	)

	// POST /api/v1/sys/user/grant-role
	r.POST("/api/v1/sys/user/grant-role",
		registry.Perm("sys:user:grant-role", "分配用户角色"),
		log.SysLog("分配用户角色"),
		middleware.NoRepeat(3000),
		grantRoleHandler,
	)

	// POST /api/v1/sys/user/grant-permission
	r.POST("/api/v1/sys/user/grant-permission",
		registry.Perm("sys:user:grant-permission", "分配用户权限"),
		log.SysLog("分配用户权限"),
		middleware.NoRepeat(3000),
		grantPermissionHandler,
	)

	// GET /api/v1/sys/user/own-permission-detail
	r.GET("/api/v1/sys/user/own-permission-detail",
		registry.Perm("sys:user:own-permission-detail", "用户权限详情"),
		ownPermissionDetailHandler,
	)

	// GET /api/v1/sys/user/own-roles
	r.GET("/api/v1/sys/user/own-roles",
		registry.Perm("sys:user:own-roles", "用户角色列表"),
		ownRolesHandler,
	)

	// GET /api/v1/sys/user/current
	r.GET("/api/v1/sys/user/current",
		middleware.HeiCheckLogin(),
		currentHandler,
	)

	// GET /api/v1/sys/user/menus
	r.GET("/api/v1/sys/user/menus",
		middleware.HeiCheckLogin(),
		menusHandler,
	)

	// GET /api/v1/sys/user/permissions
	r.GET("/api/v1/sys/user/permissions",
		middleware.HeiCheckLogin(),
		permissionsHandler,
	)

	// POST /api/v1/sys/user/update-profile
	r.POST("/api/v1/sys/user/update-profile",
		middleware.HeiCheckLogin(),
		log.SysLog("更新个人信息"),
		middleware.NoRepeat(3000),
		updateProfileHandler,
	)

	// POST /api/v1/sys/user/update-avatar
	r.POST("/api/v1/sys/user/update-avatar",
		middleware.HeiCheckLogin(),
		log.SysLog("更新头像"),
		updateAvatarHandler,
	)

	// POST /api/v1/sys/user/update-password
	r.POST("/api/v1/sys/user/update-password",
		middleware.HeiCheckLogin(),
		log.SysLog("修改密码"),
		middleware.NoRepeat(3000),
		updatePasswordHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/user/page
func pageHandler(c *gin.Context) {
	var param user.UserPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	user.UserPage(c, &param)
}

// createHandler handles POST /api/v1/sys/user/create
func createHandler(c *gin.Context) {
	var vo user.UserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	user.UserCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/user/modify
func modifyHandler(c *gin.Context) {
	var vo user.UserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	user.UserModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/user/remove
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	user.UserRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/user/detail
func detailHandler(c *gin.Context) {
	vo := user.UserDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// grantRoleHandler handles POST /api/v1/sys/user/grant-role
func grantRoleHandler(c *gin.Context) {
	var param user.GrantRoleParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	user.UserGrantRole(c, &param)
	result.Success(c, nil)
}

// grantPermissionHandler handles POST /api/v1/sys/user/grant-permission
func grantPermissionHandler(c *gin.Context) {
	var param user.GrantUserPermissionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	user.UserGrantPermission(c, &param)
	result.Success(c, nil)
}

// ownPermissionDetailHandler handles GET /api/v1/sys/user/own-permission-detail
func ownPermissionDetailHandler(c *gin.Context) {
	data := user.UserOwnPermissionDetails(c, c.Query("user_id"))
	result.Success(c, data)
}

// ownRolesHandler handles GET /api/v1/sys/user/own-roles
func ownRolesHandler(c *gin.Context) {
	data := user.UserOwnRoles(c, c.Query("user_id"))
	result.Success(c, data)
}

// currentHandler handles GET /api/v1/sys/user/current
func currentHandler(c *gin.Context) {
	vo := user.UserCurrent(c)
	result.Success(c, vo)
}

// menusHandler handles GET /api/v1/sys/user/menus
func menusHandler(c *gin.Context) {
	data := user.UserMenus(c)
	result.Success(c, data)
}

// permissionsHandler handles GET /api/v1/sys/user/permissions
func permissionsHandler(c *gin.Context) {
	data := user.UserPermissions(c)
	result.Success(c, data)
}

// updateProfileHandler handles POST /api/v1/sys/user/update-profile
func updateProfileHandler(c *gin.Context) {
	var param user.UpdateProfileParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	user.UserUpdateProfile(c, &param)
	result.Success(c, nil)
}

// updateAvatarHandler handles POST /api/v1/sys/user/update-avatar
func updateAvatarHandler(c *gin.Context) {
	var param user.UpdateAvatarParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	user.UserUpdateAvatar(c, &param)
	result.Success(c, nil)
}

// updatePasswordHandler handles POST /api/v1/sys/user/update-password
func updatePasswordHandler(c *gin.Context) {
	var param user.UpdatePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	user.UserUpdatePassword(c, &param)
	result.Success(c, nil)
}
