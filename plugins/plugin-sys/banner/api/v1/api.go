package v1

import (
	"hei-gin/sdk/auth"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	banner "hei-gin/plugins/plugin-sys/banner"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all banner routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/banner/page
	r.GET("/api/v1/sys/banner/page",
		registry.Perm("sys:banner:page", "横幅分页"),
		log.SysLog("获取Banner列表"),
		pageHandler,
	)

	// POST /api/v1/sys/banner/create
	r.POST("/api/v1/sys/banner/create",
		registry.Perm("sys:banner:create", "添加横幅"),
		log.SysLog("添加Banner"),
		middleware.NoRepeat(3000),
		createHandler,
	)

	// POST /api/v1/sys/banner/modify
	r.POST("/api/v1/sys/banner/modify",
		registry.Perm("sys:banner:modify", "编辑横幅"),
		log.SysLog("编辑Banner"),
		modifyHandler,
	)

	// POST /api/v1/sys/banner/remove
	r.POST("/api/v1/sys/banner/remove",
		registry.Perm("sys:banner:remove", "删除横幅"),
		log.SysLog("删除Banner"),
		deleteHandler,
	)

	// GET /api/v1/sys/banner/detail
	r.GET("/api/v1/sys/banner/detail",
		registry.Perm("sys:banner:detail", "横幅详情"),
		detailHandler,
	)
}

// pageHandler handles GET /api/v1/sys/banner/page
func pageHandler(c *gin.Context) {
	var param banner.BannerPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	banner.Page(c, &param)
}

// createHandler handles POST /api/v1/sys/banner/create
func createHandler(c *gin.Context) {
	var vo banner.BannerVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	banner.Create(c, &vo, userID)
	result.Success(c, nil)
}

// modifyHandler handles POST /api/v1/sys/banner/modify
func modifyHandler(c *gin.Context) {
	var vo banner.BannerVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	banner.Modify(c, &vo, userID)
	result.Success(c, nil)
}

// deleteHandler handles POST /api/v1/sys/banner/remove
func deleteHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	banner.Remove(c, param.IDs)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/banner/detail
func detailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := banner.Detail(c, id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
