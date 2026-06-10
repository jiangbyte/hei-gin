package v1

import (
	"strconv"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/result"
	"hei-gin/plugins/plugin-im/broadcast"

	"github.com/gin-gonic/gin"
	"hei-gin/sdk/registry"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
)

func RegisterRoutes(r *gin.Engine) {
	// POST /api/v1/sys/im/broadcast/send
	r.POST("/api/v1/sys/im/broadcast/send",
		authMW.HeiCheckLogin(),
		registry.Perm("sys:im:broadcast:send", "发送通知"),
		log.SysLog("发送通知"),
		authMW.NoRepeat(5000),
		sendHandler,
	)

	// GET /api/v1/sys/im/broadcast/list
	r.GET("/api/v1/sys/im/broadcast/list",
		authMW.HeiCheckLogin(),
		registry.Perm("sys:im:broadcast:list", "通知列表"),
		listHandler,
	)

	// GET /api/v1/sys/im/broadcast/unread-list
	r.GET("/api/v1/sys/im/broadcast/unread-list",
		authMW.HeiCheckLogin(),
		unreadListHandler,
	)

	// POST /api/v1/sys/im/broadcast/read
	r.POST("/api/v1/sys/im/broadcast/read",
		authMW.HeiCheckLogin(),
		readHandler,
	)

	// GET /api/v1/sys/im/broadcast/detail
	r.GET("/api/v1/sys/im/broadcast/detail",
		authMW.HeiCheckLogin(),
		detailHandler,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	// GET /api/v1/c/im/broadcast/unread-list
	r.GET("/api/v1/c/im/broadcast/unread-list",
		authMW.HeiClientCheckLogin(),
		clientUnreadListHandler,
	)

	// POST /api/v1/c/im/broadcast/read
	r.POST("/api/v1/c/im/broadcast/read",
		authMW.HeiClientCheckLogin(),
		clientReadHandler,
	)

	// GET /api/v1/c/im/broadcast/detail
	r.GET("/api/v1/c/im/broadcast/detail",
		authMW.HeiClientCheckLogin(),
		detailHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// sendHandler handles POST /api/v1/sys/im/broadcast/send
func sendHandler(c *gin.Context) {
	var param broadcast.SendBroadcastParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	broadcast.BroadcastSend(c, &param)
	result.Success(c, nil)
}

// listHandler handles GET /api/v1/sys/im/broadcast/list
func listHandler(c *gin.Context) {
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list, hasMore := broadcast.BroadcastPage(c, cursor, size)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

// unreadListHandler handles GET /api/v1/sys/im/broadcast/unread-list
func unreadListHandler(c *gin.Context) {
	list := broadcast.BroadcastUnreadList(c, string(enums.LoginTypeBusiness))
	result.Success(c, list)
}

// clientUnreadListHandler handles GET /api/v1/c/im/broadcast/unread-list
func clientUnreadListHandler(c *gin.Context) {
	list := broadcast.BroadcastUnreadList(c, string(enums.LoginTypeConsumer))
	result.Success(c, list)
}

// readHandler handles POST /api/v1/sys/im/broadcast/read
func readHandler(c *gin.Context) {
	var param broadcast.ReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	broadcast.BroadcastMarkRead(c, &param)
	result.Success(c, nil)
}

// clientReadHandler handles POST /api/v1/c/im/broadcast/read
func clientReadHandler(c *gin.Context) {
	var param broadcast.ReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	broadcast.BroadcastMarkReadConsumer(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/im/broadcast/detail and GET /api/v1/c/im/broadcast/detail
func detailHandler(c *gin.Context) {
	vo := broadcast.BroadcastDetail(c, c.Query("id"))
	result.Success(c, vo)
}
