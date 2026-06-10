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
		sessionAnalysis,
	)

	// GET /api/v1/sys/session/page
	r.GET("/api/v1/sys/session/page",
		registry.Perm("sys:session:page", "会话分页"),
		sessionPage,
	)

	// POST /api/v1/sys/session/exit
	r.POST("/api/v1/sys/session/exit",
		registry.Perm("sys:session:exit", "强退会话"),
		sessionExit,
	)

	// GET /api/v1/sys/session/tokens
	r.GET("/api/v1/sys/session/tokens",
		registry.Perm("sys:session:page", "会话分页"),
		sessionTokens,
	)

	// POST /api/v1/sys/session/exit-token
	r.POST("/api/v1/sys/session/exit-token",
		registry.Perm("sys:session:exit", "强退会话"),
		sessionExitToken,
	)

	// GET /api/v1/sys/session/chart-data
	r.GET("/api/v1/sys/session/chart-data",
		registry.Perm("sys:session:page", "会话分页"),
		sessionChartData,
	)
}

// sessionAnalysis handles GET /api/v1/sys/session/analysis
func sessionAnalysis(c *gin.Context) {
	data := session.Analysis(c)
	result.Success(c, data)
}

// sessionPage handles GET /api/v1/sys/session/page
func sessionPage(c *gin.Context) {
	var param session.SessionPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	session.Page(c, &param)
}

// sessionExit handles POST /api/v1/sys/session/exit
func sessionExit(c *gin.Context) {
	var param session.SessionExitParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	session.Exit(c, param.UserID)
	result.Success(c, nil)
}

// sessionTokens handles GET /api/v1/sys/session/tokens
func sessionTokens(c *gin.Context) {
	var param struct {
		UserID string `form:"user_id"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	data := session.TokenList(c, param.UserID)
	result.Success(c, data)
}

// sessionExitToken handles POST /api/v1/sys/session/exit-token
func sessionExitToken(c *gin.Context) {
	var param session.SessionExitTokenParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}

	session.ExitToken(c, param.UserID, param.Token)
	result.Success(c, nil)
}

// sessionChartData handles GET /api/v1/sys/session/chart-data
func sessionChartData(c *gin.Context) {
	data := session.ChartData(c)
	result.Success(c, data)
}
func init() {
	registry.RegisterRoute(RegisterRoutes)
}
