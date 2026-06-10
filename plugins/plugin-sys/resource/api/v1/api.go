package v1

import (

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/auth"
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
		modulePage,
	)
	r.POST("/api/v1/sys/module/create",
		registry.Perm("sys:module:create", "添加模块"),
		log.SysLog("添加模块"),
		authmw.NoRepeat(3000),
		moduleCreate,
	)
	r.POST("/api/v1/sys/module/modify",
		registry.Perm("sys:module:modify", "编辑模块"),
		log.SysLog("编辑模块"),
		moduleModify,
	)
	r.POST("/api/v1/sys/module/remove",
		registry.Perm("sys:module:remove", "删除模块"),
		log.SysLog("删除模块"),
		moduleRemove,
	)
	r.GET("/api/v1/sys/module/detail",
		registry.Perm("sys:module:detail", "模块详情"),
		moduleDetail,
	)

	// ---- Resource routes ----
	r.GET("/api/v1/sys/resource/tree",
		registry.Perm("sys:resource:tree", "资源树"),
		resourceTree,
	)
	r.GET("/api/v1/sys/resource/page",
		registry.Perm("sys:resource:page", "资源分页"),
		resourcePage,
	)
	r.POST("/api/v1/sys/resource/create",
		registry.Perm("sys:resource:create", "添加资源"),
		log.SysLog("添加资源"),
		authmw.NoRepeat(3000),
		resourceCreate,
	)
	r.POST("/api/v1/sys/resource/modify",
		registry.Perm("sys:resource:modify", "编辑资源"),
		log.SysLog("编辑资源"),
		resourceModify,
	)
	r.POST("/api/v1/sys/resource/remove",
		registry.Perm("sys:resource:remove", "删除资源"),
		log.SysLog("删除资源"),
		resourceRemove,
	)
	r.GET("/api/v1/sys/resource/detail",
		registry.Perm("sys:resource:detail", "资源详情"),
		resourceDetail,
	)
}

// ---------------------------------------------------------------------------
// Module handlers
// ---------------------------------------------------------------------------

func modulePage(c *gin.Context) {
	param := &resource.ModulePageParam{}
	if err := c.ShouldBindQuery(param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	resource.ModulePage(c, param)
}

func moduleDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Failure(c, "id不能为空", 400)
		return
	}

	data := resource.ModuleDetail(c, id)
	if data == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, data)
}

func moduleCreate(c *gin.Context) {
	vo := &resource.ModuleVO{}
	if err := c.ShouldBindJSON(vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	resource.ModuleCreate(c, vo, userID)
	result.Success(c, nil)
}

func moduleModify(c *gin.Context) {
	vo := &resource.ModuleVO{}
	if err := c.ShouldBindJSON(vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	if vo.ID == "" {
		result.Failure(c, "id不能为空", 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	resource.ModuleModify(c, vo, userID)
	result.Success(c, nil)
}

func moduleRemove(c *gin.Context) {
	param := &utils.IdsParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	if len(param.IDs) == 0 {
		result.Failure(c, "ids不能为空", 400)
		return
	}

	resource.ModuleRemove(c, param.IDs)
	result.Success(c, nil)
}

// ---------------------------------------------------------------------------
// Resource handlers
// ---------------------------------------------------------------------------

func resourceTree(c *gin.Context) {
	data := resource.ResourceTree(c, "")
	result.Success(c, data)
}

func resourcePage(c *gin.Context) {
	param := &resource.ResourcePageParam{}
	if err := c.ShouldBindQuery(param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	resource.ResourcePage(c, param)
}

func resourceDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Failure(c, "id不能为空", 400)
		return
	}

	data := resource.ResourceDetail(c, id)
	if data == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, data)
}

func resourceCreate(c *gin.Context) {
	vo := &resource.ResourceVO{}
	if err := c.ShouldBindJSON(vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	resource.ResourceCreate(c, vo, userID)
	result.Success(c, nil)
}

func resourceModify(c *gin.Context) {
	vo := &resource.ResourceVO{}
	if err := c.ShouldBindJSON(vo); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	if vo.ID == "" {
		result.Failure(c, "id不能为空", 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	resource.ResourceModify(c, vo, userID)
	result.Success(c, nil)
}

func resourceRemove(c *gin.Context) {
	param := &utils.IdsParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		result.Failure(c, "请求参数错误: "+err.Error(), 400)
		return
	}
	if len(param.IDs) == 0 {
		result.Failure(c, "ids不能为空", 400)
		return
	}

	resource.ResourceRemove(c, param.IDs)
	result.Success(c, nil)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}
