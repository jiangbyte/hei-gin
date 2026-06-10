package v1

import (
	"strconv"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/result"
	"hei-gin/plugins/plugin-im/friend"

	"github.com/gin-gonic/gin"
	"hei-gin/sdk/registry"
	authMW "hei-gin/sdk/auth/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	// POST /api/v1/sys/im/friend/send-request
	r.POST("/api/v1/sys/im/friend/send-request",
		authMW.HeiCheckLogin(),
		authMW.NoRepeat(3000),
		sendRequestHandler,
	)

	// POST /api/v1/sys/im/friend/accept
	r.POST("/api/v1/sys/im/friend/accept",
		authMW.HeiCheckLogin(),
		acceptHandler,
	)

	// POST /api/v1/sys/im/friend/reject
	r.POST("/api/v1/sys/im/friend/reject",
		authMW.HeiCheckLogin(),
		rejectHandler,
	)

	// GET /api/v1/sys/im/friend/list
	r.GET("/api/v1/sys/im/friend/list",
		authMW.HeiCheckLogin(),
		listHandler,
	)

	// GET /api/v1/sys/im/friend/pending-requests
	r.GET("/api/v1/sys/im/friend/pending-requests",
		authMW.HeiCheckLogin(),
		pendingRequestsHandler,
	)

	// POST /api/v1/sys/im/friend/remove
	r.POST("/api/v1/sys/im/friend/remove",
		authMW.HeiCheckLogin(),
		removeHandler,
	)

	// POST /api/v1/sys/im/friend/block
	r.POST("/api/v1/sys/im/friend/block",
		authMW.HeiCheckLogin(),
		blockHandler,
	)

	// POST /api/v1/sys/im/friend/unblock
	r.POST("/api/v1/sys/im/friend/unblock",
		authMW.HeiCheckLogin(),
		unblockHandler,
	)

	// GET /api/v1/sys/im/friend/block-list
	r.GET("/api/v1/sys/im/friend/block-list",
		authMW.HeiCheckLogin(),
		blockListHandler,
	)

	// POST /api/v1/sys/im/friend/remark
	r.POST("/api/v1/sys/im/friend/remark",
		authMW.HeiCheckLogin(),
		remarkHandler,
	)

	// GET /api/v1/sys/im/friend/search
	r.GET("/api/v1/sys/im/friend/search",
		authMW.HeiCheckLogin(),
		searchHandler,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	// POST /api/v1/c/im/friend/send-request
	r.POST("/api/v1/c/im/friend/send-request",
		authMW.HeiClientCheckLogin(),
		clientSendRequestHandler,
	)

	// POST /api/v1/c/im/friend/accept
	r.POST("/api/v1/c/im/friend/accept",
		authMW.HeiClientCheckLogin(),
		clientAcceptHandler,
	)

	// POST /api/v1/c/im/friend/reject
	r.POST("/api/v1/c/im/friend/reject",
		authMW.HeiClientCheckLogin(),
		clientRejectHandler,
	)

	// GET /api/v1/c/im/friend/list
	r.GET("/api/v1/c/im/friend/list",
		authMW.HeiClientCheckLogin(),
		clientListHandler,
	)

	// GET /api/v1/c/im/friend/pending-requests
	r.GET("/api/v1/c/im/friend/pending-requests",
		authMW.HeiClientCheckLogin(),
		clientPendingRequestsHandler,
	)

	// POST /api/v1/c/im/friend/remove
	r.POST("/api/v1/c/im/friend/remove",
		authMW.HeiClientCheckLogin(),
		clientRemoveHandler,
	)

	// POST /api/v1/c/im/friend/block
	r.POST("/api/v1/c/im/friend/block",
		authMW.HeiClientCheckLogin(),
		clientBlockHandler,
	)

	// POST /api/v1/c/im/friend/unblock
	r.POST("/api/v1/c/im/friend/unblock",
		authMW.HeiClientCheckLogin(),
		clientUnblockHandler,
	)

	// GET /api/v1/c/im/friend/block-list
	r.GET("/api/v1/c/im/friend/block-list",
		authMW.HeiClientCheckLogin(),
		clientBlockListHandler,
	)

	// POST /api/v1/c/im/friend/remark
	r.POST("/api/v1/c/im/friend/remark",
		authMW.HeiClientCheckLogin(),
		clientRemarkHandler,
	)

	// GET /api/v1/c/im/friend/search
	r.GET("/api/v1/c/im/friend/search",
		authMW.HeiClientCheckLogin(),
		searchHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// sendRequestHandler handles POST /api/v1/sys/im/friend/send-request
func sendRequestHandler(c *gin.Context) {
	var param friend.SendRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendSendRequest(c, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

// clientSendRequestHandler handles POST /api/v1/c/im/friend/send-request
func clientSendRequestHandler(c *gin.Context) {
	var param friend.SendRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendSendRequest(c, string(enums.LoginTypeConsumer), &param)
	result.Success(c, nil)
}

// acceptHandler handles POST /api/v1/sys/im/friend/accept
func acceptHandler(c *gin.Context) {
	var param friend.HandleRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendAcceptRequest(c, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

// clientAcceptHandler handles POST /api/v1/c/im/friend/accept
func clientAcceptHandler(c *gin.Context) {
	var param friend.HandleRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendAcceptRequest(c, string(enums.LoginTypeConsumer), &param)
	result.Success(c, nil)
}

// rejectHandler handles POST /api/v1/sys/im/friend/reject
func rejectHandler(c *gin.Context) {
	var param friend.HandleRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendRejectRequest(c, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

// clientRejectHandler handles POST /api/v1/c/im/friend/reject
func clientRejectHandler(c *gin.Context) {
	var param friend.HandleRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendRejectRequest(c, string(enums.LoginTypeConsumer), &param)
	result.Success(c, nil)
}

// listHandler handles GET /api/v1/sys/im/friend/list
func listHandler(c *gin.Context) {
	list := friend.FriendList(c, string(enums.LoginTypeBusiness))
	result.Success(c, list)
}

// clientListHandler handles GET /api/v1/c/im/friend/list
func clientListHandler(c *gin.Context) {
	list := friend.FriendList(c, string(enums.LoginTypeConsumer))
	result.Success(c, list)
}

// pendingRequestsHandler handles GET /api/v1/sys/im/friend/pending-requests
func pendingRequestsHandler(c *gin.Context) {
	incoming, outgoing := friend.FriendPendingRequests(c, string(enums.LoginTypeBusiness))
	result.Success(c, gin.H{"incoming": incoming, "outgoing": outgoing})
}

// clientPendingRequestsHandler handles GET /api/v1/c/im/friend/pending-requests
func clientPendingRequestsHandler(c *gin.Context) {
	incoming, outgoing := friend.FriendPendingRequests(c, string(enums.LoginTypeConsumer))
	result.Success(c, gin.H{"incoming": incoming, "outgoing": outgoing})
}

// removeHandler handles POST /api/v1/sys/im/friend/remove
func removeHandler(c *gin.Context) {
	var param friend.RemoveFriendParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendRemove(c, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

// clientRemoveHandler handles POST /api/v1/c/im/friend/remove
func clientRemoveHandler(c *gin.Context) {
	var param friend.RemoveFriendParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendRemove(c, string(enums.LoginTypeConsumer), &param)
	result.Success(c, nil)
}

// blockHandler handles POST /api/v1/sys/im/friend/block
func blockHandler(c *gin.Context) {
	var param friend.BlockParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendBlock(c, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

// clientBlockHandler handles POST /api/v1/c/im/friend/block
func clientBlockHandler(c *gin.Context) {
	var param friend.BlockParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendBlock(c, string(enums.LoginTypeConsumer), &param)
	result.Success(c, nil)
}

// unblockHandler handles POST /api/v1/sys/im/friend/unblock
func unblockHandler(c *gin.Context) {
	var param friend.BlockParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendUnblock(c, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

// clientUnblockHandler handles POST /api/v1/c/im/friend/unblock
func clientUnblockHandler(c *gin.Context) {
	var param friend.BlockParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendUnblock(c, string(enums.LoginTypeConsumer), &param)
	result.Success(c, nil)
}

// blockListHandler handles GET /api/v1/sys/im/friend/block-list
func blockListHandler(c *gin.Context) {
	list := friend.FriendBlockList(c, string(enums.LoginTypeBusiness))
	result.Success(c, list)
}

// clientBlockListHandler handles GET /api/v1/c/im/friend/block-list
func clientBlockListHandler(c *gin.Context) {
	list := friend.FriendBlockList(c, string(enums.LoginTypeConsumer))
	result.Success(c, list)
}

// remarkHandler handles POST /api/v1/sys/im/friend/remark
func remarkHandler(c *gin.Context) {
	var param friend.RemarkParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendUpdateRemark(c, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

// clientRemarkHandler handles POST /api/v1/c/im/friend/remark
func clientRemarkHandler(c *gin.Context) {
	var param friend.RemarkParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	friend.FriendUpdateRemark(c, string(enums.LoginTypeConsumer), &param)
	result.Success(c, nil)
}

// searchHandler handles GET /api/v1/sys/im/friend/search and GET /api/v1/c/im/friend/search
func searchHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	results := friend.FriendSearch(c, keyword, size)
	result.Success(c, results)
}
