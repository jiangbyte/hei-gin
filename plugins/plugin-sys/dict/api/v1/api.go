package v1

import (
	"hei-gin/sdk/auth"
	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	dict "hei-gin/plugins/plugin-sys/dict"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all dict routes on the given gin engine.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/dict/page
	r.GET("/api/v1/sys/dict/page",
		registry.Perm("sys:dict:page", "字典分页"),
		dictPage,
	)

	// POST /api/v1/sys/dict/create
	r.POST("/api/v1/sys/dict/create",
		registry.Perm("sys:dict:create", "添加字典"),
		log.SysLog("添加字典"),
		middleware.NoRepeat(3000),
		dictCreate,
	)

	// POST /api/v1/sys/dict/modify
	r.POST("/api/v1/sys/dict/modify",
		registry.Perm("sys:dict:modify", "编辑字典"),
		log.SysLog("编辑字典"),
		dictModify,
	)

	// POST /api/v1/sys/dict/remove
	r.POST("/api/v1/sys/dict/remove",
		registry.Perm("sys:dict:remove", "删除字典"),
		log.SysLog("删除字典"),
		dictRemove,
	)

	// GET /api/v1/sys/dict/detail
	r.GET("/api/v1/sys/dict/detail",
		registry.Perm("sys:dict:detail", "字典详情"),
		dictDetail,
	)

	// GET /api/v1/sys/dict/list
	r.GET("/api/v1/sys/dict/list",
		registry.Perm("sys:dict:list", "字典列表"),
		dictList,
	)

	// GET /api/v1/sys/dict/tree
	r.GET("/api/v1/sys/dict/tree",
		dictTree,
	)

	// GET /api/v1/sys/dict/get-label
	r.GET("/api/v1/sys/dict/get-label",
		registry.Perm("sys:dict:get-label", "字典标签"),
		dictGetLabel,
	)

	// GET /api/v1/sys/dict/get-children
	r.GET("/api/v1/sys/dict/get-children",
		registry.Perm("sys:dict:get-children", "字典子项"),
		dictGetChildren,
	)
}

// dictPage handles GET /api/v1/sys/dict/page
func dictPage(c *gin.Context) {
	var param dict.DictPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	dict.Page(c, &param)
}

// dictList handles GET /api/v1/sys/dict/list
func dictList(c *gin.Context) {
	var param dict.DictListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := dict.DictList(c, &param)
	result.Success(c, data)
}

// dictTree handles GET /api/v1/sys/dict/tree
func dictTree(c *gin.Context) {
	var param dict.DictTreeParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := dict.Tree(c, &param)
	result.Success(c, data)
}

// dictCreate handles POST /api/v1/sys/dict/create
func dictCreate(c *gin.Context) {
	var vo dict.DictVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	dict.Create(c, &vo, userID)
	result.Success(c, nil)
}

// dictModify handles POST /api/v1/sys/dict/modify
func dictModify(c *gin.Context) {
	var vo dict.DictVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if vo.ID == "" {
		result.Failure(c, "ID不能为空", 400)
		return
	}

	userID := auth.GetLoginIDDefaultNull(c)
	dict.Modify(c, &vo, userID)
	result.Success(c, nil)
}

// dictRemove handles POST /api/v1/sys/dict/remove
func dictRemove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if len(param.IDs) == 0 {
		result.Failure(c, "ids不能为空", 400)
		return
	}

	dict.Remove(c, param.IDs)
	result.Success(c, nil)
}

// dictDetail handles GET /api/v1/sys/dict/detail
func dictDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Failure(c, "id不能为空", 400)
		return
	}

	vo := dict.Detail(c, id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}

// dictGetLabel handles GET /api/v1/sys/dict/get-label
func dictGetLabel(c *gin.Context) {
	typeCode := c.Query("type_code")
	value := c.Query("value")

	data := dict.DictGetLabel(c, typeCode, value)
	c.JSON(200, data)
}

// dictGetChildren handles GET /api/v1/sys/dict/get-children
func dictGetChildren(c *gin.Context) {
	typeCode := c.Query("type_code")

	data := dict.DictGetChildren(c, typeCode)
	result.Success(c, data)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}
