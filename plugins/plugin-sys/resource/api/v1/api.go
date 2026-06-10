package v1

import (
	"github.com/gin-gonic/gin"

	"hei-gin/sdk/registry"
	authmw "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	resource "hei-gin/plugins/plugin-sys/resource"
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
func modulePageHandler(c *gin.Context) {
	var param resource.ModulePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	resource.ModulePage(c, &param)
}

// moduleDetailHandler handles GET /api/v1/sys/module/detail
func moduleDetailHandler(c *gin.Context) {
	vo := resource.ModuleDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// moduleCreateHandler handles POST /api/v1/sys/module/create
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
func resourceTreeHandler(c *gin.Context) {
	data := resource.ResourceTree(c, "")
	result.Success(c, data)
}

// resourcePageHandler handles GET /api/v1/sys/resource/page
func resourcePageHandler(c *gin.Context) {
	var param resource.ResourcePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	resource.ResourcePage(c, &param)
}

// resourceDetailHandler handles GET /api/v1/sys/resource/detail
func resourceDetailHandler(c *gin.Context) {
	vo := resource.ResourceDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// resourceCreateHandler handles POST /api/v1/sys/resource/create
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
func resourceRemoveHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	resource.ResourceRemove(c, &param)
	result.Success(c, nil)
}
