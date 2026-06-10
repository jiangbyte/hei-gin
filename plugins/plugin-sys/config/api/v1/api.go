package v1

import (
	"hei-gin/sdk/auth"
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
		deleteHandler,
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

func pageHandler(c *gin.Context) {
	var param config.ConfigPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	config.Page(c, &param)
}

func listByCategoryHandler(c *gin.Context) {
	var param config.ConfigListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	vos := config.ListByCategory(c, param.Category)
	result.Success(c, vos)
}

func createHandler(c *gin.Context) {
	var vo config.ConfigVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginIDDefaultNull(c)
	config.Create(c, &vo, userID)
	result.Success(c, nil)
}

func modifyHandler(c *gin.Context) {
	var vo config.ConfigVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginIDDefaultNull(c)
	config.Modify(c, &vo, userID)
	result.Success(c, nil)
}

func deleteHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	config.Remove(c, param.IDs)
	result.Success(c, nil)
}

func detailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := config.Detail(c, id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}

func editBatchHandler(c *gin.Context) {
	var param config.ConfigBatchEditParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginIDDefaultNull(c)
	config.EditBatch(c, &param, userID)
	result.Success(c, nil)
}

func editByCategoryHandler(c *gin.Context) {
	var param config.ConfigCategoryEditParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginIDDefaultNull(c)
	config.EditByCategory(c, &param, userID)
	result.Success(c, nil)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
