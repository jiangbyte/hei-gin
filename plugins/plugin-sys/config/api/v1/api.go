package v1

import (
	config "hei-gin/plugins/plugin-sys/config"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *config.Service
}

var defaultHandler = newHandler(config.DefaultModule)

func newHandler(module *config.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/sys/config/page",
		registry.Perm("sys:config:page", "配置分页"),
		defaultHandler.page,
	)

	r.GET("/api/v1/sys/config/list-by-category",
		registry.Perm("sys:config:list", "配置列表"),
		defaultHandler.listByCategory,
	)

	r.POST("/api/v1/sys/config/create",
		registry.Perm("sys:config:create", "添加配置"),
		log.SysLog("添加配置"),
		middleware.NoRepeat(3000),
		defaultHandler.create,
	)

	r.POST("/api/v1/sys/config/modify",
		registry.Perm("sys:config:modify", "编辑配置"),
		log.SysLog("编辑配置"),
		defaultHandler.modify,
	)

	r.POST("/api/v1/sys/config/remove",
		registry.Perm("sys:config:remove", "删除配置"),
		log.SysLog("删除配置"),
		defaultHandler.remove,
	)

	r.GET("/api/v1/sys/config/detail",
		registry.Perm("sys:config:detail", "配置详情"),
		defaultHandler.detail,
	)

	r.POST("/api/v1/sys/config/edit-batch",
		registry.Perm("sys:config:edit", "配置编辑"),
		log.SysLog("批量编辑配置"),
		middleware.NoRepeat(3000),
		defaultHandler.editBatch,
	)

	r.POST("/api/v1/sys/config/edit-by-category",
		registry.Perm("sys:config:edit", "配置编辑"),
		log.SysLog("按分类批量编辑配置"),
		middleware.NoRepeat(3000),
		defaultHandler.editByCategory,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/config/page
// @Summary      系统配置分页查询
// @Description  访问 /api/v1/sys/config/page，系统配置分页查询
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        query  query  config.ConfigPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/config/page [get]
func (h *handler) page(c *gin.Context) {
	var param config.ConfigPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Page(c, &param)
}

// listByCategoryHandler handles GET /api/v1/sys/config/list-by-category
// @Summary      系统配置按分类列表查询
// @Description  访问 /api/v1/sys/config/list-by-category，系统配置按分类列表查询
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        category  query  string  false  "category"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/config/list-by-category [get]
func (h *handler) listByCategory(c *gin.Context) {
	vos := h.service.ListByCategory(c, c.Query("category"))
	result.Success(c, vos)
}

// createHandler handles POST /api/v1/sys/config/create
// @Summary      系统配置创建
// @Description  访问 /api/v1/sys/config/create，系统配置创建
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        body  body  config.ConfigVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/config/create [post]
func (h *handler) create(c *gin.Context) {
	var vo config.ConfigVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Create(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/config/modify
// @Summary      系统配置修改
// @Description  访问 /api/v1/sys/config/modify，系统配置修改
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        body  body  config.ConfigVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/config/modify [post]
func (h *handler) modify(c *gin.Context) {
	var vo config.ConfigVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Modify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/config/remove
// @Summary      系统配置删除
// @Description  访问 /api/v1/sys/config/remove，系统配置删除
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/config/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/config/detail
// @Summary      系统配置详情查询
// @Description  访问 /api/v1/sys/config/detail，系统配置详情查询
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/config/detail [get]
func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}

// editBatchHandler handles POST /api/v1/sys/config/edit-batch
// @Summary      系统配置批量编辑
// @Description  访问 /api/v1/sys/config/edit-batch，系统配置批量编辑
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        body  body  config.ConfigBatchEditParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/config/edit-batch [post]
func (h *handler) editBatch(c *gin.Context) {
	var param config.ConfigBatchEditParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.EditBatch(c, &param)
	result.Success(c, nil)
}

// editByCategoryHandler handles POST /api/v1/sys/config/edit-by-category
// @Summary      系统配置按分类批量编辑
// @Description  访问 /api/v1/sys/config/edit-by-category，系统配置按分类批量编辑
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        body  body  config.ConfigCategoryEditParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/config/edit-by-category [post]
func (h *handler) editByCategory(c *gin.Context) {
	var param config.ConfigCategoryEditParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	h.service.EditByCategory(c, &param)
	result.Success(c, nil)
}
