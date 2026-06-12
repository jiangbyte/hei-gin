package v1

import (
	"github.com/gin-gonic/gin"

	resource "hei-gin/plugins/plugin-sys/resource"
	authmw "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"
)

// RegisterRoutes registers all module and resource routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// ---- Module routes ----
	r.GET("/api/v1/sys/module/page",
		registry.Perm("sys:module:page", "模块分页"),
		modulePageHandler,
	)
	r.POST("/api/v1/sys/module/create",
		registry.Perm("sys:module:create", "添加模块"),
		log.SysLog("添加模块"),
		authmw.NoRepeat(3000),
		moduleCreateHandler,
	)
	r.POST("/api/v1/sys/module/modify",
		registry.Perm("sys:module:modify", "编辑模块"),
		log.SysLog("编辑模块"),
		moduleModifyHandler,
	)
	r.POST("/api/v1/sys/module/remove",
		registry.Perm("sys:module:remove", "删除模块"),
		log.SysLog("删除模块"),
		moduleRemoveHandler,
	)
	r.GET("/api/v1/sys/module/detail",
		registry.Perm("sys:module:detail", "模块详情"),
		moduleDetailHandler,
	)

	// ---- Resource routes ----
	r.GET("/api/v1/sys/resource/tree",
		registry.Perm("sys:resource:tree", "资源树"),
		resourceTreeHandler,
	)
	r.GET("/api/v1/sys/resource/page",
		registry.Perm("sys:resource:page", "资源分页"),
		resourcePageHandler,
	)
	r.POST("/api/v1/sys/resource/create",
		registry.Perm("sys:resource:create", "添加资源"),
		log.SysLog("添加资源"),
		authmw.NoRepeat(3000),
		resourceCreateHandler,
	)
	r.POST("/api/v1/sys/resource/modify",
		registry.Perm("sys:resource:modify", "编辑资源"),
		log.SysLog("编辑资源"),
		resourceModifyHandler,
	)
	r.POST("/api/v1/sys/resource/remove",
		registry.Perm("sys:resource:remove", "删除资源"),
		log.SysLog("删除资源"),
		resourceRemoveHandler,
	)
	r.GET("/api/v1/sys/resource/detail",
		registry.Perm("sys:resource:detail", "资源详情"),
		resourceDetailHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// ---------------------------------------------------------------------------
// Module handlers
// ---------------------------------------------------------------------------

// modulePageHandler handles GET /api/v1/sys/module/page
// @Summary      资源管理分页查询
// @Description  访问 /api/v1/sys/module/page，资源管理分页查询
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        query  query  resource.ModulePageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/module/page [get]
func modulePageHandler(c *gin.Context) {
	var param resource.ModulePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	resource.ModulePage(c, &param)
}

// moduleDetailHandler handles GET /api/v1/sys/module/detail
// @Summary      资源管理详情查询
// @Description  访问 /api/v1/sys/module/detail，资源管理详情查询
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/module/detail [get]
func moduleDetailHandler(c *gin.Context) {
	vo := resource.ModuleDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// moduleCreateHandler handles POST /api/v1/sys/module/create
// @Summary      资源管理创建
// @Description  访问 /api/v1/sys/module/create，资源管理创建
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        body  body  resource.ModuleVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/module/create [post]
func moduleCreateHandler(c *gin.Context) {
	var vo resource.ModuleVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	resource.ModuleCreate(c, &vo)
	result.Success(c, nil)
}

// moduleModifyHandler handles POST /api/v1/sys/module/modify
// @Summary      资源管理修改
// @Description  访问 /api/v1/sys/module/modify，资源管理修改
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        body  body  resource.ModuleVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/module/modify [post]
func moduleModifyHandler(c *gin.Context) {
	var vo resource.ModuleVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	resource.ModuleModify(c, &vo)
	result.Success(c, nil)
}

// moduleRemoveHandler handles POST /api/v1/sys/module/remove
// @Summary      资源管理删除
// @Description  访问 /api/v1/sys/module/remove，资源管理删除
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/module/remove [post]
func moduleRemoveHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	resource.ModuleRemove(c, &param)
	result.Success(c, nil)
}

// ---------------------------------------------------------------------------
// Resource handlers
// ---------------------------------------------------------------------------

// resourceTreeHandler handles GET /api/v1/sys/resource/tree
// @Summary      资源管理树形查询
// @Description  访问 /api/v1/sys/resource/tree，资源管理树形查询
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/resource/tree [get]
func resourceTreeHandler(c *gin.Context) {
	data := resource.ResourceTree(c, "")
	result.Success(c, data)
}

// resourcePageHandler handles GET /api/v1/sys/resource/page
// @Summary      资源管理分页查询
// @Description  访问 /api/v1/sys/resource/page，资源管理分页查询
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        query  query  resource.ResourcePageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/resource/page [get]
func resourcePageHandler(c *gin.Context) {
	var param resource.ResourcePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	resource.ResourcePage(c, &param)
}

// resourceDetailHandler handles GET /api/v1/sys/resource/detail
// @Summary      资源管理详情查询
// @Description  访问 /api/v1/sys/resource/detail，资源管理详情查询
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/resource/detail [get]
func resourceDetailHandler(c *gin.Context) {
	vo := resource.ResourceDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// resourceCreateHandler handles POST /api/v1/sys/resource/create
// @Summary      资源管理创建
// @Description  访问 /api/v1/sys/resource/create，资源管理创建
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        body  body  resource.ResourceVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/resource/create [post]
func resourceCreateHandler(c *gin.Context) {
	var vo resource.ResourceVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	resource.ResourceCreate(c, &vo)
	result.Success(c, nil)
}

// resourceModifyHandler handles POST /api/v1/sys/resource/modify
// @Summary      资源管理修改
// @Description  访问 /api/v1/sys/resource/modify，资源管理修改
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        body  body  resource.ResourceVO  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/resource/modify [post]
func resourceModifyHandler(c *gin.Context) {
	var vo resource.ResourceVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	resource.ResourceModify(c, &vo)
	result.Success(c, nil)
}

// resourceRemoveHandler handles POST /api/v1/sys/resource/remove
// @Summary      资源管理删除
// @Description  访问 /api/v1/sys/resource/remove，资源管理删除
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdsParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/resource/remove [post]
func resourceRemoveHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	resource.ResourceRemove(c, &param)
	result.Success(c, nil)
}
