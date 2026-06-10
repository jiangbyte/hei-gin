package v1

import (
	"hei-gin/sdk/registry"
	"hei-gin/sdk/result"
	clientsession "hei-gin/plugins/plugin-client/session"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all client session routes.
func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/client/session/analysis
	r.GET("/api/v1/client/session/analysis",
		registry.Perm("sys:session:page", "会话分页"),
		analysisHandler,
	)

	// GET /api/v1/client/session/page
	r.GET("/api/v1/client/session/page",
		registry.Perm("sys:session:page", "会话分页"),
		pageHandler,
	)

	// POST /api/v1/client/session/exit
	r.POST("/api/v1/client/session/exit",
		registry.Perm("sys:session:exit", "强退会话"),
		exitHandler,
	)

	// GET /api/v1/client/session/tokens
	r.GET("/api/v1/client/session/tokens",
		registry.Perm("sys:session:page", "会话分页"),
		tokensHandler,
	)

	// POST /api/v1/client/session/exit-token
	r.POST("/api/v1/client/session/exit-token",
		registry.Perm("sys:session:exit", "强退会话"),
		exitTokenHandler,
	)

	// GET /api/v1/client/session/chart-data
	r.GET("/api/v1/client/session/chart-data",
		registry.Perm("sys:session:page", "会话分页"),
		chartDataHandler,
	)
}

// analysisHandler handles GET /api/v1/client/session/analysis
func analysisHandler(c *gin.Context) {
	data := clientsession.Analysis(c)
	result.Success(c, data)
}

// pageHandler handles GET /api/v1/client/session/page
func pageHandler(c *gin.Context) {
	var param clientsession.SessionPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientsession.Page(c, &param)
}

// exitHandler handles POST /api/v1/client/session/exit
func exitHandler(c *gin.Context) {
	var param clientsession.SessionExitParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientsession.Exit(c, param.UserID)
	result.Success(c, nil)
}

// tokensHandler handles GET /api/v1/client/session/tokens
func tokensHandler(c *gin.Context) {
	var param struct {
		UserID string `form:"user_id"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := clientsession.TokenList(c, param.UserID)
	result.Success(c, data)
}

// exitTokenHandler handles POST /api/v1/client/session/exit-token
func exitTokenHandler(c *gin.Context) {
	var param clientsession.SessionExitTokenParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientsession.ExitToken(c, param.UserID, param.Token)
	result.Success(c, nil)
}

// chartDataHandler handles GET /api/v1/client/session/chart-data
func chartDataHandler(c *gin.Context) {
	data := clientsession.ChartData(c)
	result.Success(c, data)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
