package v1

import (
	session "hei-gin/plugins/plugin-sys/session"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all session routes.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/session/analysis
	r.GET("/api/v1/sys/session/analysis",
		registry.Perm("sys:session:page", "会话分页"),
		analysisHandler,
	)

	// GET /api/v1/sys/session/page
	r.GET("/api/v1/sys/session/page",
		registry.Perm("sys:session:page", "会话分页"),
		pageHandler,
	)

	// POST /api/v1/sys/session/exit
	r.POST("/api/v1/sys/session/exit",
		registry.Perm("sys:session:exit", "强退会话"),
		exitHandler,
	)

	// GET /api/v1/sys/session/tokens
	r.GET("/api/v1/sys/session/tokens",
		registry.Perm("sys:session:page", "会话分页"),
		tokensHandler,
	)

	// POST /api/v1/sys/session/exit-token
	r.POST("/api/v1/sys/session/exit-token",
		registry.Perm("sys:session:exit", "强退会话"),
		exitTokenHandler,
	)

	// GET /api/v1/sys/session/chart-data
	r.GET("/api/v1/sys/session/chart-data",
		registry.Perm("sys:session:page", "会话分页"),
		chartDataHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// analysisHandler handles GET /api/v1/sys/session/analysis
// @Summary      会话管理分析数据
// @Description  访问 /api/v1/sys/session/analysis，会话管理分析数据
// @Tags         会话管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/session/analysis [get]
func analysisHandler(c *gin.Context) {
	result.Success(c, session.Analysis(c))
}

// pageHandler handles GET /api/v1/sys/session/page
// @Summary      会话管理分页查询
// @Description  访问 /api/v1/sys/session/page，会话管理分页查询
// @Tags         会话管理
// @Accept       json
// @Produce      json
// @Param        query  query  session.SessionPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/session/page [get]
func pageHandler(c *gin.Context) {
	var param session.SessionPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	session.Page(c, &param)
}

// exitHandler handles POST /api/v1/sys/session/exit
// @Summary      会话管理强退会话
// @Description  访问 /api/v1/sys/session/exit，会话管理强退会话
// @Tags         会话管理
// @Accept       json
// @Produce      json
// @Param        body  body  session.SessionExitParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/session/exit [post]
func exitHandler(c *gin.Context) {
	var param session.SessionExitParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	session.Exit(c, &param)
	result.Success(c, nil)
}

// tokensHandler handles GET /api/v1/sys/session/tokens
// @Summary      会话管理令牌列表
// @Description  访问 /api/v1/sys/session/tokens，会话管理令牌列表
// @Tags         会话管理
// @Accept       json
// @Produce      json
// @Param        user_id  query  string  false  "user_id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/session/tokens [get]
func tokensHandler(c *gin.Context) {
	data := session.TokenList(c, c.Query("user_id"))
	result.Success(c, data)
}

// exitTokenHandler handles POST /api/v1/sys/session/exit-token
// @Summary      会话管理强退令牌
// @Description  访问 /api/v1/sys/session/exit-token，会话管理强退令牌
// @Tags         会话管理
// @Accept       json
// @Produce      json
// @Param        body  body  session.SessionExitTokenParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/session/exit-token [post]
func exitTokenHandler(c *gin.Context) {
	var param session.SessionExitTokenParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	session.ExitToken(c, &param)
	result.Success(c, nil)
}

// chartDataHandler handles GET /api/v1/sys/session/chart-data
// @Summary      会话管理图表数据
// @Description  访问 /api/v1/sys/session/chart-data，会话管理图表数据
// @Tags         会话管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/session/chart-data [get]
func chartDataHandler(c *gin.Context) {
	result.Success(c, session.ChartData(c))
}
