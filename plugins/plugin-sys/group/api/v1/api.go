package v1

import (
	group "hei-gin/plugins/plugin-sys/group"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *group.Service
}

var defaultHandler = newHandler(group.DefaultModule)

func newHandler(module *group.Module) *handler {
	return &handler{service: module.Service()}
}

// RegisterRoutes registers all group routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/group/page
	r.GET("/api/v1/sys/group/page",
		registry.Perm("sys:group:page", "分组分页"),
		defaultHandler.page,
	)

	// GET /api/v1/sys/group/union-tree
	r.GET("/api/v1/sys/group/union-tree",
		registry.Perm("sys:group:tree", "分组树"),
		defaultHandler.unionTree,
	)

	// GET /api/v1/sys/group/tree
	r.GET("/api/v1/sys/group/tree",
		registry.Perm("sys:group:tree", "分组树"),
		defaultHandler.tree,
	)

	// POST /api/v1/sys/group/create
	r.POST("/api/v1/sys/group/create",
		registry.Perm("sys:group:create", "添加分组"),
		log.SysLog("添加用户组"),
		middleware.NoRepeat(3000),
		defaultHandler.create,
	)

	// POST /api/v1/sys/group/modify
	r.POST("/api/v1/sys/group/modify",
		registry.Perm("sys:group:modify", "编辑分组"),
		log.SysLog("编辑用户组"),
		defaultHandler.modify,
	)

	// POST /api/v1/sys/group/remove
	r.POST("/api/v1/sys/group/remove",
		registry.Perm("sys:group:remove", "删除分组"),
		log.SysLog("删除用户组"),
		defaultHandler.remove,
	)

	// GET /api/v1/sys/group/detail
	r.GET("/api/v1/sys/group/detail",
		registry.Perm("sys:group:detail", "分组详情"),
		defaultHandler.detail,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/group/page
// @Summary      部门分组分页查询
// @Description  访问 /api/v1/sys/group/page，部门分组分页查询
// @Tags         部门分组
// @Accept       json
// @Produce      json
// @Param        query  query  group.GroupPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/group/page [get]
func (h *handler) page(c *gin.Context) {
	var param group.GroupPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
}

// unionTreeHandler handles GET /api/v1/sys/group/union-tree
// @Summary      部门分组联合树形查询
// @Description  访问 /api/v1/sys/group/union-tree，部门分组联合树形查询
// @Tags         部门分组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/group/union-tree [get]
func (h *handler) unionTree(c *gin.Context) {
	data := h.service.Options(c)
	result.Success(c, data)
}

// treeHandler handles GET /api/v1/sys/group/tree
// @Summary      部门分组树形查询
// @Description  访问 /api/v1/sys/group/tree，部门分组树形查询
// @Tags         部门分组
// @Accept       json
// @Produce      json
// @Param        query  query  group.GroupTreeParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/group/tree [get]
func (h *handler) tree(c *gin.Context) {
	var param group.GroupTreeParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := h.service.Tree(c, &param)
	result.Success(c, data)
}

// createHandler handles POST /api/v1/sys/group/create
// @Summary      部门分组创建
// @Description  访问 /api/v1/sys/group/create，部门分组创建
// @Tags         部门分组
// @Accept       json
// @Produce      json
// @Param        body  body  group.GroupVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/group/create [post]
func (h *handler) create(c *gin.Context) {
	var vo group.GroupVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Create(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/group/modify
// @Summary      部门分组修改
// @Description  访问 /api/v1/sys/group/modify，部门分组修改
// @Tags         部门分组
// @Accept       json
// @Produce      json
// @Param        body  body  group.GroupVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/group/modify [post]
func (h *handler) modify(c *gin.Context) {
	var vo group.GroupVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Modify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/group/remove
// @Summary      部门分组删除
// @Description  访问 /api/v1/sys/group/remove，部门分组删除
// @Tags         部门分组
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/group/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/group/detail
// @Summary      部门分组详情查询
// @Description  访问 /api/v1/sys/group/detail，部门分组详情查询
// @Tags         部门分组
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/group/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}
