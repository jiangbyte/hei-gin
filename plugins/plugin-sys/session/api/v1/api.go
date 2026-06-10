package v1

import (
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	session "hei-gin/plugins/plugin-sys/session"

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
func analysisHandler(c *gin.Context) {
	result.Success(c, session.Analysis(c))
}

// pageHandler handles GET /api/v1/sys/session/page
func pageHandler(c *gin.Context) {
	var param session.SessionPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	session.Page(c, &param)
}

// exitHandler handles POST /api/v1/sys/session/exit
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
func tokensHandler(c *gin.Context) {
	data := session.TokenList(c, c.Query("user_id"))
	result.Success(c, data)
}

// exitTokenHandler handles POST /api/v1/sys/session/exit-token
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
func chartDataHandler(c *gin.Context) {
	result.Success(c, session.ChartData(c))
}
