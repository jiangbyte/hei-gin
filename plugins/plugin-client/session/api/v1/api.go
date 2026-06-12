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

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// analysisHandler handles GET /api/v1/client/session/analysis
// @Summary      C端会话分析数据
// @Description  访问 /api/v1/client/session/analysis，C端会话分析数据
// @Tags         C端会话
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client/session/analysis [get]
func analysisHandler(c *gin.Context) {
	data := clientsession.SessionAnalysis(c)
	result.Success(c, data)
}

// pageHandler handles GET /api/v1/client/session/page
// @Summary      C端会话分页查询
// @Description  访问 /api/v1/client/session/page，C端会话分页查询
// @Tags         C端会话
// @Accept       json
// @Produce      json
// @Param        query  query  clientsession.SessionPageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client/session/page [get]
func pageHandler(c *gin.Context) {
	var param clientsession.SessionPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientsession.SessionPage(c, &param)
}

// exitHandler handles POST /api/v1/client/session/exit
// @Summary      C端会话强退会话
// @Description  访问 /api/v1/client/session/exit，C端会话强退会话
// @Tags         C端会话
// @Accept       json
// @Produce      json
// @Param        body  body  clientsession.SessionExitParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client/session/exit [post]
func exitHandler(c *gin.Context) {
	var param clientsession.SessionExitParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientsession.SessionExit(c, param.UserID)
	result.Success(c, nil)
}

// tokensHandler handles GET /api/v1/client/session/tokens
// @Summary      C端会话令牌列表
// @Description  访问 /api/v1/client/session/tokens，C端会话令牌列表
// @Tags         C端会话
// @Accept       json
// @Produce      json
// @Param        user_id  query  string  false  "user_id"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client/session/tokens [get]
func tokensHandler(c *gin.Context) {
	data := clientsession.SessionTokenList(c, c.Query("user_id"))
	result.Success(c, data)
}

// exitTokenHandler handles POST /api/v1/client/session/exit-token
// @Summary      C端会话强退令牌
// @Description  访问 /api/v1/client/session/exit-token，C端会话强退令牌
// @Tags         C端会话
// @Accept       json
// @Produce      json
// @Param        body  body  clientsession.SessionExitTokenParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client/session/exit-token [post]
func exitTokenHandler(c *gin.Context) {
	var param clientsession.SessionExitTokenParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	clientsession.SessionExitToken(c, param.UserID, param.Token)
	result.Success(c, nil)
}

// chartDataHandler handles GET /api/v1/client/session/chart-data
// @Summary      C端会话图表数据
// @Description  访问 /api/v1/client/session/chart-data，C端会话图表数据
// @Tags         C端会话
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/client/session/chart-data [get]
func chartDataHandler(c *gin.Context) {
	data := clientsession.SessionChart(c)
	result.Success(c, data)
}
