package v1

import (
	"strconv"

	"hei-gin/sdk/auth"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/result"
	"hei-gin/plugins/plugin-im/broadcast"

	"github.com/gin-gonic/gin"
	"hei-gin/sdk/registry"
)

func RegisterRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/sys/im/broadcast").Use(authMW.HeiCheckLogin())
	{
		g.POST("/send", registry.Perm("sys:im:broadcast:send", "发送通知"), log.SysLog("发送通知"), authMW.NoRepeat(5000), sendHandler)
		g.GET("/list", registry.Perm("sys:im:broadcast:list", "通知列表"), listHandler)
		g.GET("/unread-list", unreadListHandler)
		g.POST("/read", readHandler)
		g.GET("/detail", detailHandler)
	}
}

func RegisterClientRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/c/im/broadcast").Use(authMW.HeiClientCheckLogin())
	{
		g.GET("/unread-list", clientUnreadListHandler)
		g.POST("/read", clientReadHandler)
		g.GET("/detail", detailHandler)
	}
}

func sendHandler(c *gin.Context) {
	var p broadcast.SendBroadcastParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID := auth.GetLoginID(c)
	broadcast.Send(userID, &p)
	result.Success(c, nil)
}

func listHandler(c *gin.Context) {
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list, hasMore := broadcast.List(cursor, size)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

func unreadListHandler(c *gin.Context) {
	userID := auth.GetLoginID(c)
	list, _ := broadcast.UnreadList(userID, "BUSINESS")
	result.Success(c, list)
}

func clientUnreadListHandler(c *gin.Context) {
	userID := auth.Consumer.GetLoginID(c)
	list, _ := broadcast.UnreadList(userID, "CONSUMER")
	result.Success(c, list)
}

func readHandler(c *gin.Context) {
	var p struct {
		BroadcastID string `json:"broadcast_id"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID := auth.GetLoginID(c)
	broadcast.MarkRead(userID, "BUSINESS", p.BroadcastID)
	result.Success(c, nil)
}

func clientReadHandler(c *gin.Context) {
	var p struct {
		BroadcastID string `json:"broadcast_id"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID := auth.Consumer.GetLoginID(c)
	broadcast.MarkRead(userID, "CONSUMER", p.BroadcastID)
	result.Success(c, nil)
}

func detailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := broadcast.Detail(id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}
