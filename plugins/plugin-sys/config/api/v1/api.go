package v1

import (
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	config "hei-gin/plugins/plugin-sys/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/sys/config/page",
		registry.Perm("sys:config:page", "配置分页"),
		pageHandler,
	)

	r.GET("/api/v1/sys/config/list-by-category",
		registry.Perm("sys:config:list", "配置列表"),
		listByCategoryHandler,
	)

	r.POST("/api/v1/sys/config/create",
		registry.Perm("sys:config:create", "添加配置"),
		log.SysLog("添加配置"),
		middleware.NoRepeat(3000),
		createHandler,
	)

	r.POST("/api/v1/sys/config/modify",
		registry.Perm("sys:config:modify", "编辑配置"),
		log.SysLog("编辑配置"),
		modifyHandler,
	)

	r.POST("/api/v1/sys/config/remove",
		registry.Perm("sys:config:remove", "删除配置"),
		log.SysLog("删除配置"),
		removeHandler,
	)

	r.GET("/api/v1/sys/config/detail",
		registry.Perm("sys:config:detail", "配置详情"),
		detailHandler,
	)

	r.POST("/api/v1/sys/config/edit-batch",
		registry.Perm("sys:config:edit", "配置编辑"),
		log.SysLog("批量编辑配置"),
		middleware.NoRepeat(3000),
		editBatchHandler,
	)

	r.POST("/api/v1/sys/config/edit-by-category",
		registry.Perm("sys:config:edit", "配置编辑"),
		log.SysLog("按分类批量编辑配置"),
		middleware.NoRepeat(3000),
		editByCategoryHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/config/page
func pageHandler(c *gin.Context) {
	var param config.ConfigPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	config.ConfigPage(c, &param)
}

// listByCategoryHandler handles GET /api/v1/sys/config/list-by-category
func listByCategoryHandler(c *gin.Context) {
	vos := config.ConfigListByCategory(c, c.Query("category"))
	result.Success(c, vos)
}

// createHandler handles POST /api/v1/sys/config/create
func createHandler(c *gin.Context) {
	var vo config.ConfigVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	config.ConfigCreate(c, &vo)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/config/modify
func modifyHandler(c *gin.Context) {
	var vo config.ConfigVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	config.ConfigModify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/config/remove
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	config.ConfigRemove(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/config/detail
func detailHandler(c *gin.Context) {
	vo := config.ConfigDetail(c, c.Query("id"))
	result.Success(c, vo)
}

// editBatchHandler handles POST /api/v1/sys/config/edit-batch
func editBatchHandler(c *gin.Context) {
	var param config.ConfigBatchEditParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	config.ConfigEditBatch(c, &param)
	result.Success(c, nil)
}

// editByCategoryHandler handles POST /api/v1/sys/config/edit-by-category
func editByCategoryHandler(c *gin.Context) {
	var param config.ConfigCategoryEditParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	config.ConfigEditByCategory(c, &param)
	result.Success(c, nil)
}
