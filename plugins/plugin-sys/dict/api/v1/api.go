package v1

import (
	dict "hei-gin/plugins/plugin-sys/dict"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *dict.Service
}

var defaultHandler = newHandler(dict.DefaultModule)

func newHandler(module *dict.Module) *handler {
	return &handler{service: module.Service()}
}

// RegisterRoutes registers all dict routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/dict/page
	r.GET("/api/v1/sys/dict/page",
		registry.Perm("sys:dict:page", "字典分页"),
		defaultHandler.page,
	)

	// POST /api/v1/sys/dict/create
	r.POST("/api/v1/sys/dict/create",
		registry.Perm("sys:dict:create", "添加字典"),
		log.SysLog("添加字典"),
		middleware.NoRepeat(3000),
		defaultHandler.create,
	)

	// POST /api/v1/sys/dict/modify
	r.POST("/api/v1/sys/dict/modify",
		registry.Perm("sys:dict:modify", "编辑字典"),
		log.SysLog("编辑字典"),
		defaultHandler.modify,
	)

	// POST /api/v1/sys/dict/remove
	r.POST("/api/v1/sys/dict/remove",
		registry.Perm("sys:dict:remove", "删除字典"),
		log.SysLog("删除字典"),
		defaultHandler.remove,
	)

	// GET /api/v1/sys/dict/detail
	r.GET("/api/v1/sys/dict/detail",
		registry.Perm("sys:dict:detail", "字典详情"),
		defaultHandler.detail,
	)

	// GET /api/v1/sys/dict/list
	r.GET("/api/v1/sys/dict/list",
		registry.Perm("sys:dict:list", "字典列表"),
		defaultHandler.list,
	)

	// GET /api/v1/sys/dict/tree
	r.GET("/api/v1/sys/dict/tree",
		defaultHandler.tree,
	)

	// GET /api/v1/sys/dict/get-label
	r.GET("/api/v1/sys/dict/get-label",
		registry.Perm("sys:dict:get-label", "字典标签"),
		defaultHandler.getLabel,
	)

	// GET /api/v1/sys/dict/get-children
	r.GET("/api/v1/sys/dict/get-children",
		registry.Perm("sys:dict:get-children", "字典子项"),
		defaultHandler.getChildren,
	)
}

func Register() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/dict/page
// @Summary      字典管理分页查询
// @Description  访问 /api/v1/sys/dict/page，字典管理分页查询
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        query  query  dict.DictPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/page [get]
func (h *handler) page(c *gin.Context) {
	var param dict.DictPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
}

// listHandler handles GET /api/v1/sys/dict/list
// @Summary      字典管理列表查询
// @Description  访问 /api/v1/sys/dict/list，字典管理列表查询
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        query  query  dict.DictListParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/list [get]
func (h *handler) list(c *gin.Context) {
	var param dict.DictListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := h.service.List(c, &param)
	result.Success(c, data)
}

// treeHandler handles GET /api/v1/sys/dict/tree
// @Summary      字典管理树形查询
// @Description  访问 /api/v1/sys/dict/tree，字典管理树形查询
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        query  query  dict.DictTreeParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/tree [get]
func (h *handler) tree(c *gin.Context) {
	var param dict.DictTreeParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := h.service.Tree(c, &param)
	result.Success(c, data)
}

// createHandler handles POST /api/v1/sys/dict/create
// @Summary      字典管理创建
// @Description  访问 /api/v1/sys/dict/create，字典管理创建
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        body  body  dict.DictVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/create [post]
func (h *handler) create(c *gin.Context) {
	var vo dict.DictVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Create(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/dict/modify
// @Summary      字典管理修改
// @Description  访问 /api/v1/sys/dict/modify，字典管理修改
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        body  body  dict.DictVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/modify [post]
func (h *handler) modify(c *gin.Context) {
	var vo dict.DictVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Modify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/dict/remove
// @Summary      字典管理删除
// @Description  访问 /api/v1/sys/dict/remove，字典管理删除
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/dict/detail
// @Summary      字典管理详情查询
// @Description  访问 /api/v1/sys/dict/detail，字典管理详情查询
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}

// getLabelHandler handles GET /api/v1/sys/dict/get-label
// @Summary      字典管理获取标签
// @Description  访问 /api/v1/sys/dict/get-label，字典管理获取标签
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        type_code  query  string  false  "type_code"
// @Param        value  query  string  false  "value"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/get-label [get]
func (h *handler) getLabel(c *gin.Context) {
	data := h.service.GetLabel(c, c.Query("type_code"), c.Query("value"))
	result.Success(c, data)
}

// getChildrenHandler handles GET /api/v1/sys/dict/get-children
// @Summary      字典管理获取子项
// @Description  访问 /api/v1/sys/dict/get-children，字典管理获取子项
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Param        type_code  query  string  false  "type_code"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/dict/get-children [get]
func (h *handler) getChildren(c *gin.Context) {
	data := h.service.GetChildren(c, c.Query("type_code"))
	result.Success(c, data)
}
