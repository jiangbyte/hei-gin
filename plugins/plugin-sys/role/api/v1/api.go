package v1

import (
	role "hei-gin/plugins/plugin-sys/role"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *role.Service
}

var defaultHandler = newHandler(role.DefaultModule)

func newHandler(module *role.Module) *handler {
	return &handler{service: module.Service()}
}

// RegisterRoutes registers all role routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/role/page
	r.GET("/api/v1/sys/role/page",
		registry.Perm("sys:role:page", "角色分页"),
		defaultHandler.page,
	)

	// POST /api/v1/sys/role/create
	r.POST("/api/v1/sys/role/create",
		registry.Perm("sys:role:create", "添加角色"),
		log.SysLog("添加角色"),
		middleware.NoRepeat(3000),
		defaultHandler.create,
	)

	// POST /api/v1/sys/role/modify
	r.POST("/api/v1/sys/role/modify",
		registry.Perm("sys:role:modify", "编辑角色"),
		log.SysLog("编辑角色"),
		defaultHandler.modify,
	)

	// POST /api/v1/sys/role/remove
	r.POST("/api/v1/sys/role/remove",
		registry.Perm("sys:role:remove", "删除角色"),
		log.SysLog("删除角色"),
		defaultHandler.remove,
	)

	// GET /api/v1/sys/role/detail
	r.GET("/api/v1/sys/role/detail",
		registry.Perm("sys:role:detail", "角色详情"),
		defaultHandler.detail,
	)

	// POST /api/v1/sys/role/grant-permission
	r.POST("/api/v1/sys/role/grant-permission",
		registry.Perm("sys:role:grant-permission", "分配角色权限"),
		log.SysLog("分配角色权限"),
		middleware.NoRepeat(3000),
		defaultHandler.grantPermission,
	)

	// POST /api/v1/sys/role/grant-resource
	r.POST("/api/v1/sys/role/grant-resource",
		registry.Perm("sys:role:grant-resource", "分配角色资源"),
		log.SysLog("分配角色资源"),
		middleware.NoRepeat(3000),
		defaultHandler.grantResource,
	)

	// POST /api/v1/sys/role/refresh-session-acl
	r.POST("/api/v1/sys/role/refresh-session-acl",
		registry.Perm("sys:role:refresh-session-acl", "刷新角色会话权限"),
		log.SysLog("刷新角色会话权限"),
		middleware.NoRepeat(1000),
		defaultHandler.refreshSessionACL,
	)

	// GET /api/v1/sys/role/own-permission
	r.GET("/api/v1/sys/role/own-permission",
		registry.Perm("sys:role:own-permission", "角色权限列表"),
		defaultHandler.ownPermission,
	)

	// GET /api/v1/sys/role/own-permission-detail
	r.GET("/api/v1/sys/role/own-permission-detail",
		registry.Perm("sys:role:own-permission-detail", "角色权限详情"),
		defaultHandler.ownPermissionDetail,
	)

	// GET /api/v1/sys/role/own-resource
	r.GET("/api/v1/sys/role/own-resource",
		registry.Perm("sys:role:own-resource", "角色资源列表"),
		defaultHandler.ownResource,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/role/page
// @Summary      角色管理分页查询
// @Description  访问 /api/v1/sys/role/page，角色管理分页查询
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        query  query  role.RolePageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/page [get]
func (h *handler) page(c *gin.Context) {
	var param role.RolePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
}

// createHandler handles POST /api/v1/sys/role/create
// @Summary      角色管理创建
// @Description  访问 /api/v1/sys/role/create，角色管理创建
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body  role.RoleVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/create [post]
func (h *handler) create(c *gin.Context) {
	var vo role.RoleVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Create(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/role/modify
// @Summary      角色管理修改
// @Description  访问 /api/v1/sys/role/modify，角色管理修改
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body  role.RoleVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/modify [post]
func (h *handler) modify(c *gin.Context) {
	var vo role.RoleVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Modify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/role/remove
// @Summary      角色管理删除
// @Description  访问 /api/v1/sys/role/remove，角色管理删除
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/role/detail
// @Summary      角色管理详情查询
// @Description  访问 /api/v1/sys/role/detail，角色管理详情查询
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}

// grantPermissionHandler handles POST /api/v1/sys/role/grant-permission
// @Summary      角色管理分配权限
// @Description  访问 /api/v1/sys/role/grant-permission，角色管理分配权限
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body  role.GrantPermissionParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/grant-permission [post]
func (h *handler) grantPermission(c *gin.Context) {
	var param role.GrantPermissionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.GrantPermissions(c, &param)
	result.Success(c, nil)
}

// grantResourceHandler handles POST /api/v1/sys/role/grant-resource
// @Summary      角色管理分配资源
// @Description  访问 /api/v1/sys/role/grant-resource，角色管理分配资源
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body  role.GrantResourceParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/grant-resource [post]
func (h *handler) grantResource(c *gin.Context) {
	var param role.GrantResourceParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.GrantResources(c, &param)
	result.Success(c, nil)
}

// refreshSessionACLHandler handles POST /api/v1/sys/role/refresh-session-acl
// @Summary      角色管理刷新会话权限
// @Description  访问 /api/v1/sys/role/refresh-session-acl，角色管理刷新会话权限
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body  role.RefreshRoleSessionACLParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/refresh-session-acl [post]
func (h *handler) refreshSessionACL(c *gin.Context) {
	var param role.RefreshRoleSessionACLParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.RefreshSessionACL(c, &param)
	result.Success(c, nil)
}

// ownPermissionHandler handles GET /api/v1/sys/role/own-permission
// @Summary      角色管理自身权限
// @Description  访问 /api/v1/sys/role/own-permission，角色管理自身权限
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        role_id  query  string  false  "role_id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/own-permission [get]
func (h *handler) ownPermission(c *gin.Context) {
	codes := h.service.OwnPermissionCodes(c, c.Query("role_id"))
	result.Success(c, codes)
}

// ownPermissionDetailHandler handles GET /api/v1/sys/role/own-permission-detail
// @Summary      角色管理自身权限详情
// @Description  访问 /api/v1/sys/role/own-permission-detail，角色管理自身权限详情
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        role_id  query  string  false  "role_id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/own-permission-detail [get]
func (h *handler) ownPermissionDetail(c *gin.Context) {
	details := h.service.OwnPermissionDetails(c, c.Query("role_id"))
	result.Success(c, details)
}

// ownResourceHandler handles GET /api/v1/sys/role/own-resource
// @Summary      角色管理自身资源
// @Description  访问 /api/v1/sys/role/own-resource，角色管理自身资源
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        role_id  query  string  false  "role_id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/role/own-resource [get]
func (h *handler) ownResource(c *gin.Context) {
	ids := h.service.OwnResourceIDs(c, c.Query("role_id"))
	result.Success(c, ids)
}
