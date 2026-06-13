package v1

import (
	user "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/auth"
	middleware "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *user.Service
}

var defaultHandler = newHandler(user.DefaultModule)

func newHandler(module *user.Module) *handler {
	return &handler{service: module.Service()}
}

// RegisterRoutes registers all user routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/user/page
	r.GET("/api/v1/sys/user/page",
		registry.Perm("sys:user:page", "用户分页"),
		defaultHandler.page,
	)

	// POST /api/v1/sys/user/create
	r.POST("/api/v1/sys/user/create",
		registry.Perm("sys:user:create", "添加用户"),
		log.SysLog("添加用户"),
		middleware.NoRepeat(3000),
		defaultHandler.create,
	)

	// POST /api/v1/sys/user/modify
	r.POST("/api/v1/sys/user/modify",
		registry.Perm("sys:user:modify", "编辑用户"),
		log.SysLog("编辑用户"),
		defaultHandler.modify,
	)

	// POST /api/v1/sys/user/remove
	r.POST("/api/v1/sys/user/remove",
		registry.Perm("sys:user:remove", "删除用户"),
		log.SysLog("删除用户"),
		defaultHandler.remove,
	)

	// GET /api/v1/sys/user/detail
	r.GET("/api/v1/sys/user/detail",
		registry.Perm("sys:user:detail", "用户详情"),
		defaultHandler.detail,
	)

	// POST /api/v1/sys/user/grant-role
	r.POST("/api/v1/sys/user/grant-role",
		registry.Perm("sys:user:grant-role", "分配用户角色"),
		log.SysLog("分配用户角色"),
		middleware.NoRepeat(3000),
		defaultHandler.grantRole,
	)

	// POST /api/v1/sys/user/grant-permission
	r.POST("/api/v1/sys/user/grant-permission",
		registry.Perm("sys:user:grant-permission", "分配用户权限"),
		log.SysLog("分配用户权限"),
		middleware.NoRepeat(3000),
		defaultHandler.grantPermission,
	)

	// POST /api/v1/sys/user/refresh-session-acl
	r.POST("/api/v1/sys/user/refresh-session-acl",
		registry.Perm("sys:user:refresh-session-acl", "刷新用户会话权限"),
		log.SysLog("刷新用户会话权限"),
		middleware.NoRepeat(1000),
		defaultHandler.refreshSessionACL,
	)

	// POST /api/v1/sys/user/batch-refresh-session-acl
	r.POST("/api/v1/sys/user/batch-refresh-session-acl",
		registry.Perm("sys:user:batch-refresh-session-acl", "批量刷新用户会话权限"),
		log.SysLog("批量刷新用户会话权限"),
		middleware.NoRepeat(1000),
		defaultHandler.batchRefreshSessionACL,
	)

	// GET /api/v1/sys/user/own-permission-detail
	r.GET("/api/v1/sys/user/own-permission-detail",
		registry.Perm("sys:user:own-permission-detail", "用户权限详情"),
		defaultHandler.ownPermissionDetail,
	)

	// GET /api/v1/sys/user/own-roles
	r.GET("/api/v1/sys/user/own-roles",
		registry.Perm("sys:user:own-roles", "用户角色列表"),
		defaultHandler.ownRoles,
	)

	// GET /api/v1/sys/user/current
	r.GET("/api/v1/sys/user/current",
		middleware.CheckLogin(auth.Business),
		defaultHandler.current,
	)

	// GET /api/v1/sys/user/menus
	r.GET("/api/v1/sys/user/menus",
		middleware.CheckLogin(auth.Business),
		defaultHandler.menus,
	)

	// GET /api/v1/sys/user/permissions
	r.GET("/api/v1/sys/user/permissions",
		middleware.CheckLogin(auth.Business),
		defaultHandler.permissions,
	)

	// POST /api/v1/sys/user/update-profile
	r.POST("/api/v1/sys/user/update-profile",
		middleware.CheckLogin(auth.Business),
		log.SysLog("更新个人信息"),
		middleware.NoRepeat(3000),
		defaultHandler.updateProfile,
	)

	// POST /api/v1/sys/user/update-avatar
	r.POST("/api/v1/sys/user/update-avatar",
		middleware.CheckLogin(auth.Business),
		log.SysLog("更新头像"),
		defaultHandler.updateAvatar,
	)

	// POST /api/v1/sys/user/update-password
	r.POST("/api/v1/sys/user/update-password",
		middleware.CheckLogin(auth.Business),
		log.SysLog("修改密码"),
		middleware.NoRepeat(3000),
		defaultHandler.updatePassword,
	)
}

func Register() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/user/page
// @Summary      用户管理分页查询
// @Description  访问 /api/v1/sys/user/page，用户管理分页查询
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        query  query  user.UserPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/page [get]
func (h *handler) page(c *gin.Context) {
	var param user.UserPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
}

// createHandler handles POST /api/v1/sys/user/create
// @Summary      用户管理创建
// @Description  访问 /api/v1/sys/user/create，用户管理创建
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.UserVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/create [post]
func (h *handler) create(c *gin.Context) {
	var vo user.UserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Create(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/user/modify
// @Summary      用户管理修改
// @Description  访问 /api/v1/sys/user/modify，用户管理修改
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.UserVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/modify [post]
func (h *handler) modify(c *gin.Context) {
	var vo user.UserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Modify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/user/remove
// @Summary      用户管理删除
// @Description  访问 /api/v1/sys/user/remove，用户管理删除
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/user/detail
// @Summary      用户管理详情查询
// @Description  访问 /api/v1/sys/user/detail，用户管理详情查询
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}

// grantRoleHandler handles POST /api/v1/sys/user/grant-role
// @Summary      用户管理分配角色
// @Description  访问 /api/v1/sys/user/grant-role，用户管理分配角色
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.GrantRoleParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/grant-role [post]
func (h *handler) grantRole(c *gin.Context) {
	var param user.GrantRoleParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.GrantRole(c, &param)
	result.Success(c, nil)
}

// grantPermissionHandler handles POST /api/v1/sys/user/grant-permission
// @Summary      用户管理分配权限
// @Description  访问 /api/v1/sys/user/grant-permission，用户管理分配权限
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.GrantUserPermissionParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/grant-permission [post]
func (h *handler) grantPermission(c *gin.Context) {
	var param user.GrantUserPermissionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.GrantPermission(c, &param)
	result.Success(c, nil)
}

// refreshSessionACLHandler handles POST /api/v1/sys/user/refresh-session-acl
// @Summary      用户管理刷新会话权限
// @Description  访问 /api/v1/sys/user/refresh-session-acl，用户管理刷新会话权限
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.RefreshSessionACLParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/refresh-session-acl [post]
func (h *handler) refreshSessionACL(c *gin.Context) {
	var param user.RefreshSessionACLParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.RefreshSessionACL(c, &param)
	result.Success(c, nil)
}

// batchRefreshSessionACLHandler handles POST /api/v1/sys/user/batch-refresh-session-acl
// @Summary      用户管理批量刷新会话权限
// @Description  访问 /api/v1/sys/user/batch-refresh-session-acl，用户管理批量刷新会话权限
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.BatchRefreshSessionACLParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/batch-refresh-session-acl [post]
func (h *handler) batchRefreshSessionACL(c *gin.Context) {
	var param user.BatchRefreshSessionACLParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.BatchRefreshSessionACL(c, &param)
	result.Success(c, nil)
}

// ownPermissionDetailHandler handles GET /api/v1/sys/user/own-permission-detail
// @Summary      用户管理自身权限详情
// @Description  访问 /api/v1/sys/user/own-permission-detail，用户管理自身权限详情
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        user_id  query  string  false  "user_id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/own-permission-detail [get]
func (h *handler) ownPermissionDetail(c *gin.Context) {
	data := h.service.OwnPermissionDetails(c, c.Query("user_id"))
	result.Success(c, data)
}

// ownRolesHandler handles GET /api/v1/sys/user/own-roles
// @Summary      用户管理自身角色列表
// @Description  访问 /api/v1/sys/user/own-roles，用户管理自身角色列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        user_id  query  string  false  "user_id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/own-roles [get]
func (h *handler) ownRoles(c *gin.Context) {
	data := h.service.OwnRoles(c, c.Query("user_id"))
	result.Success(c, data)
}

// currentHandler handles GET /api/v1/sys/user/current
// @Summary      用户管理当前信息
// @Description  访问 /api/v1/sys/user/current，用户管理当前信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/current [get]
func (h *handler) current(c *gin.Context) {
	vo := h.service.Current(c)
	result.Success(c, vo)
}

// menusHandler handles GET /api/v1/sys/user/menus
// @Summary      用户管理菜单列表
// @Description  访问 /api/v1/sys/user/menus，用户管理菜单列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/menus [get]
func (h *handler) menus(c *gin.Context) {
	data := h.service.Menus(c)
	result.Success(c, data)
}

// permissionsHandler handles GET /api/v1/sys/user/permissions
// @Summary      用户管理权限列表
// @Description  访问 /api/v1/sys/user/permissions，用户管理权限列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/permissions [get]
func (h *handler) permissions(c *gin.Context) {
	data := h.service.Permissions(c)
	result.Success(c, data)
}

// updateProfileHandler handles POST /api/v1/sys/user/update-profile
// @Summary      用户管理更新个人信息
// @Description  访问 /api/v1/sys/user/update-profile，用户管理更新个人信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.UpdateProfileParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/update-profile [post]
func (h *handler) updateProfile(c *gin.Context) {
	var param user.UpdateProfileParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.UpdateProfile(c, &param)
	result.Success(c, nil)
}

// updateAvatarHandler handles POST /api/v1/sys/user/update-avatar
// @Summary      用户管理更新头像
// @Description  访问 /api/v1/sys/user/update-avatar，用户管理更新头像
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.UpdateAvatarParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/update-avatar [post]
func (h *handler) updateAvatar(c *gin.Context) {
	var param user.UpdateAvatarParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.UpdateAvatar(c, &param)
	result.Success(c, nil)
}

// updatePasswordHandler handles POST /api/v1/sys/user/update-password
// @Summary      用户管理修改密码
// @Description  访问 /api/v1/sys/user/update-password，用户管理修改密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  user.UpdatePasswordParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/user/update-password [post]
func (h *handler) updatePassword(c *gin.Context) {
	var param user.UpdatePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.UpdatePassword(c, &param)
	result.Success(c, nil)
}
