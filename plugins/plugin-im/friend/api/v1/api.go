package v1

import (
	"strconv"

	"hei-gin/plugins/plugin-im/friend"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
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
// @Summary      即时通讯好友发送申请
// @Description  访问 /api/v1/sys/im/friend/send-request，即时通讯好友发送申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.SendRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/send-request [post]
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
// @Summary      即时通讯好友发送申请
// @Description  访问 /api/v1/c/im/friend/send-request，即时通讯好友发送申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.SendRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/send-request [post]
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
// @Summary      即时通讯好友接受申请
// @Description  访问 /api/v1/sys/im/friend/accept，即时通讯好友接受申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.HandleRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/accept [post]
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
// @Summary      即时通讯好友接受申请
// @Description  访问 /api/v1/c/im/friend/accept，即时通讯好友接受申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.HandleRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/accept [post]
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
// @Summary      即时通讯好友拒绝申请
// @Description  访问 /api/v1/sys/im/friend/reject，即时通讯好友拒绝申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.HandleRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/reject [post]
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
// @Summary      即时通讯好友拒绝申请
// @Description  访问 /api/v1/c/im/friend/reject，即时通讯好友拒绝申请
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.HandleRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/reject [post]
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
// @Summary      即时通讯好友列表查询
// @Description  访问 /api/v1/sys/im/friend/list，即时通讯好友列表查询
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/list [get]
func listHandler(c *gin.Context) {
	list := friend.FriendList(c, string(enums.LoginTypeBusiness))
	result.Success(c, list)
}

// clientListHandler handles GET /api/v1/c/im/friend/list
// @Summary      即时通讯好友列表查询
// @Description  访问 /api/v1/c/im/friend/list，即时通讯好友列表查询
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/list [get]
func clientListHandler(c *gin.Context) {
	list := friend.FriendList(c, string(enums.LoginTypeConsumer))
	result.Success(c, list)
}

// pendingRequestsHandler handles GET /api/v1/sys/im/friend/pending-requests
// @Summary      即时通讯好友待处理申请列表
// @Description  访问 /api/v1/sys/im/friend/pending-requests，即时通讯好友待处理申请列表
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/pending-requests [get]
func pendingRequestsHandler(c *gin.Context) {
	incoming, outgoing := friend.FriendPendingRequests(c, string(enums.LoginTypeBusiness))
	result.Success(c, gin.H{"incoming": incoming, "outgoing": outgoing})
}

// clientPendingRequestsHandler handles GET /api/v1/c/im/friend/pending-requests
// @Summary      即时通讯好友待处理申请列表
// @Description  访问 /api/v1/c/im/friend/pending-requests，即时通讯好友待处理申请列表
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/pending-requests [get]
func clientPendingRequestsHandler(c *gin.Context) {
	incoming, outgoing := friend.FriendPendingRequests(c, string(enums.LoginTypeConsumer))
	result.Success(c, gin.H{"incoming": incoming, "outgoing": outgoing})
}

// removeHandler handles POST /api/v1/sys/im/friend/remove
// @Summary      即时通讯好友删除
// @Description  访问 /api/v1/sys/im/friend/remove，即时通讯好友删除
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.RemoveFriendParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/remove [post]
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
// @Summary      即时通讯好友删除
// @Description  访问 /api/v1/c/im/friend/remove，即时通讯好友删除
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.RemoveFriendParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/remove [post]
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
// @Summary      即时通讯好友拉黑
// @Description  访问 /api/v1/sys/im/friend/block，即时通讯好友拉黑
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.BlockParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/block [post]
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
// @Summary      即时通讯好友拉黑
// @Description  访问 /api/v1/c/im/friend/block，即时通讯好友拉黑
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.BlockParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/block [post]
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
// @Summary      即时通讯好友取消拉黑
// @Description  访问 /api/v1/sys/im/friend/unblock，即时通讯好友取消拉黑
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.BlockParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/unblock [post]
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
// @Summary      即时通讯好友取消拉黑
// @Description  访问 /api/v1/c/im/friend/unblock，即时通讯好友取消拉黑
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.BlockParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/unblock [post]
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
// @Summary      即时通讯好友拉黑列表
// @Description  访问 /api/v1/sys/im/friend/block-list，即时通讯好友拉黑列表
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/block-list [get]
func blockListHandler(c *gin.Context) {
	list := friend.FriendBlockList(c, string(enums.LoginTypeBusiness))
	result.Success(c, list)
}

// clientBlockListHandler handles GET /api/v1/c/im/friend/block-list
// @Summary      即时通讯好友拉黑列表
// @Description  访问 /api/v1/c/im/friend/block-list，即时通讯好友拉黑列表
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/block-list [get]
func clientBlockListHandler(c *gin.Context) {
	list := friend.FriendBlockList(c, string(enums.LoginTypeConsumer))
	result.Success(c, list)
}

// remarkHandler handles POST /api/v1/sys/im/friend/remark
// @Summary      即时通讯好友设置备注
// @Description  访问 /api/v1/sys/im/friend/remark，即时通讯好友设置备注
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.RemarkParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/remark [post]
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
// @Summary      即时通讯好友设置备注
// @Description  访问 /api/v1/c/im/friend/remark，即时通讯好友设置备注
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        body  body  friend.RemarkParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/friend/remark [post]
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
// @Summary      即时通讯好友搜索
// @Description  访问 /api/v1/sys/im/friend/search，即时通讯好友搜索
// @Tags         即时通讯好友
// @Accept       json
// @Produce      json
// @Param        keyword  query  string  false  "keyword"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/friend/search [get]
// @Router       /api/v1/c/im/friend/search [get]
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
