package v1

import (
	"strconv"

	"hei-gin/sdk/auth"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	"hei-gin/plugins/plugin-im/friend"

	"github.com/gin-gonic/gin"
)

func RegisterSysRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/sys/im/friend").Use(authMW.HeiCheckLogin())
	{
		g.POST("/send-request", authMW.NoRepeat(3000), sendRequestHandler)
		g.POST("/accept", acceptHandler)
		g.POST("/reject", rejectHandler)
		g.GET("/list", listHandler)
		g.GET("/pending-requests", pendingRequestsHandler)
		g.POST("/remove", removeHandler)
		g.POST("/block", blockHandler)
		g.POST("/unblock", unblockHandler)
		g.GET("/block-list", blockListHandler)
		g.POST("/remark", remarkHandler)
		g.GET("/search", searchHandler)
	}
}

func RegisterClientRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/c/im/friend").Use(authMW.HeiClientCheckLogin())
	{
		g.POST("/send-request", clientSendRequestHandler)
		g.POST("/accept", clientAcceptHandler)
		g.POST("/reject", clientRejectHandler)
		g.GET("/list", clientListHandler)
		g.GET("/pending-requests", clientPendingRequestsHandler)
		g.POST("/remove", clientRemoveHandler)
		g.POST("/block", clientBlockHandler)
		g.POST("/unblock", clientUnblockHandler)
		g.GET("/block-list", clientBlockListHandler)
		g.POST("/remark", clientRemarkHandler)
		g.GET("/search", clientSearchHandler)
	}
}

func getLoginID(c *gin.Context) (string, string) {
	path := c.Request.URL.Path
	if len(path) > 8 && path[:8] == "/api/v1/c" {
		return auth.Consumer.GetLoginID(c), "CONSUMER"
	}
	return auth.GetLoginID(c), "BUSINESS"
}


// ==================== Sys Handlers ====================

func sysUserID(c *gin.Context) (string, string) {
	return auth.GetLoginID(c), "BUSINESS"
}

func sendRequestHandler(c *gin.Context) {
	var p friend.SendRequestParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := sysUserID(c)
	friend.SendRequest(c.Request.Context(), uid, ut, &p)
	result.Success(c, nil)
}

func acceptHandler(c *gin.Context) {
	var p friend.HandleRequestParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := sysUserID(c)
	friend.AcceptRequest(c.Request.Context(), uid, ut, &p)
	result.Success(c, nil)
}

func rejectHandler(c *gin.Context) {
	var p friend.HandleRequestParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := sysUserID(c)
	friend.RejectRequest(c.Request.Context(), uid, ut, &p)
	result.Success(c, nil)
}

func listHandler(c *gin.Context) {
	uid, ut := sysUserID(c)
	list := friend.FriendList(c.Request.Context(), uid, ut)
	result.Success(c, list)
}

func pendingRequestsHandler(c *gin.Context) {
	uid, ut := sysUserID(c)
	incoming, outgoing := friend.PendingRequests(c.Request.Context(), uid, ut)
	result.Success(c, gin.H{"incoming": incoming, "outgoing": outgoing})
}

func removeHandler(c *gin.Context) {
	var p struct {
		FriendID   string `json:"friend_id"`
		FriendType string `json:"friend_type"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := sysUserID(c)
	friend.RemoveFriend(c.Request.Context(), uid, ut, p.FriendID, p.FriendType)
	result.Success(c, nil)
}

func blockHandler(c *gin.Context) {
	var p friend.BlockParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := sysUserID(c)
	friend.BlockUser(c.Request.Context(), uid, ut, p.BlockedID, p.BlockedType)
	result.Success(c, nil)
}

func unblockHandler(c *gin.Context) {
	var p friend.BlockParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := sysUserID(c)
	friend.UnblockUser(c.Request.Context(), uid, ut, p.BlockedID, p.BlockedType)
	result.Success(c, nil)
}

func blockListHandler(c *gin.Context) {
	uid, ut := sysUserID(c)
	list := friend.BlockList(c.Request.Context(), uid, ut)
	result.Success(c, list)
}

func remarkHandler(c *gin.Context) {
	var p friend.RemarkParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := sysUserID(c)
	friend.UpdateFriendRemark(c.Request.Context(), uid, ut, p.FriendID, p.FriendType, p.Remark)
	result.Success(c, nil)
}

func searchHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	results := friend.SearchUsers(c.Request.Context(), keyword, size)
	result.Success(c, results)
}

// ==================== Client Handlers ====================

func clientUserID(c *gin.Context) (string, string) {
	return auth.Consumer.GetLoginID(c), "CONSUMER"
}

func clientSendRequestHandler(c *gin.Context) {
	var p friend.SendRequestParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := clientUserID(c)
	friend.SendRequest(c.Request.Context(), uid, ut, &p)
	result.Success(c, nil)
}

func clientAcceptHandler(c *gin.Context) {
	var p friend.HandleRequestParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := clientUserID(c)
	friend.AcceptRequest(c.Request.Context(), uid, ut, &p)
	result.Success(c, nil)
}

func clientRejectHandler(c *gin.Context) {
	var p friend.HandleRequestParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := clientUserID(c)
	friend.RejectRequest(c.Request.Context(), uid, ut, &p)
	result.Success(c, nil)
}

func clientListHandler(c *gin.Context) {
	uid, ut := clientUserID(c)
	list := friend.FriendList(c.Request.Context(), uid, ut)
	result.Success(c, list)
}

func clientPendingRequestsHandler(c *gin.Context) {
	uid, ut := clientUserID(c)
	incoming, outgoing := friend.PendingRequests(c.Request.Context(), uid, ut)
	result.Success(c, gin.H{"incoming": incoming, "outgoing": outgoing})
}

func clientRemoveHandler(c *gin.Context) {
	var p struct {
		FriendID   string `json:"friend_id"`
		FriendType string `json:"friend_type"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := clientUserID(c)
	friend.RemoveFriend(c.Request.Context(), uid, ut, p.FriendID, p.FriendType)
	result.Success(c, nil)
}

func clientBlockHandler(c *gin.Context) {
	var p friend.BlockParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := clientUserID(c)
	friend.BlockUser(c.Request.Context(), uid, ut, p.BlockedID, p.BlockedType)
	result.Success(c, nil)
}

func clientUnblockHandler(c *gin.Context) {
	var p friend.BlockParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := clientUserID(c)
	friend.UnblockUser(c.Request.Context(), uid, ut, p.BlockedID, p.BlockedType)
	result.Success(c, nil)
}

func clientBlockListHandler(c *gin.Context) {
	uid, ut := clientUserID(c)
	list := friend.BlockList(c.Request.Context(), uid, ut)
	result.Success(c, list)
}

func clientRemarkHandler(c *gin.Context) {
	var p friend.RemarkParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	uid, ut := clientUserID(c)
	friend.UpdateFriendRemark(c.Request.Context(), uid, ut, p.FriendID, p.FriendType, p.Remark)
	result.Success(c, nil)
}

func clientSearchHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	results := friend.SearchUsers(c.Request.Context(), keyword, size)
	result.Success(c, results)
}

func init() {
	registry.RegisterRoute(RegisterSysRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}
