package v1

import (
	"strconv"

	"hei-gin/sdk/result"
	"hei-gin/plugins/plugin-im/group"

	"github.com/gin-gonic/gin"
	"hei-gin/sdk/registry"
	"hei-gin/sdk/middleware"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
)

func RegisterRoutes(r *gin.Engine) {
	// POST /api/v1/sys/im/group/create
	r.POST("/api/v1/sys/im/group/create",
		authMW.HeiCheckLogin(),
		authMW.NoRepeat(3000),
		createHandler,
	)

	// GET /api/v1/sys/im/group/my-groups
	r.GET("/api/v1/sys/im/group/my-groups",
		authMW.HeiCheckLogin(),
		myGroupsHandler,
	)

	// GET /api/v1/sys/im/group/detail
	r.GET("/api/v1/sys/im/group/detail",
		authMW.HeiCheckLogin(),
		detailHandler,
	)

	// POST /api/v1/sys/im/group/update
	r.POST("/api/v1/sys/im/group/update",
		authMW.HeiCheckLogin(),
		updateHandler,
	)

	// POST /api/v1/sys/im/group/dissolve
	r.POST("/api/v1/sys/im/group/dissolve",
		authMW.HeiCheckLogin(),
		log.SysLog("解散群"),
		dissolveHandler,
	)

	// POST /api/v1/sys/im/group/invite
	r.POST("/api/v1/sys/im/group/invite",
		authMW.HeiCheckLogin(),
		authMW.NoRepeat(3000),
		inviteHandler,
	)

	// POST /api/v1/sys/im/group/join
	r.POST("/api/v1/sys/im/group/join",
		authMW.HeiCheckLogin(),
		joinHandler,
	)

	// GET /api/v1/sys/im/group/pending-join-requests
	r.GET("/api/v1/sys/im/group/pending-join-requests",
		authMW.HeiCheckLogin(),
		pendingJoinRequestsHandler,
	)

	// POST /api/v1/sys/im/group/handle-join-request
	r.POST("/api/v1/sys/im/group/handle-join-request",
		authMW.HeiCheckLogin(),
		handleJoinRequestHandler,
	)

	// POST /api/v1/sys/im/group/leave
	r.POST("/api/v1/sys/im/group/leave",
		authMW.HeiCheckLogin(),
		leaveHandler,
	)

	// POST /api/v1/sys/im/group/kick
	r.POST("/api/v1/sys/im/group/kick",
		authMW.HeiCheckLogin(),
		log.SysLog("踢出成员"),
		kickHandler,
	)

	// POST /api/v1/sys/im/group/set-role
	r.POST("/api/v1/sys/im/group/set-role",
		authMW.HeiCheckLogin(),
		log.SysLog("设置角色"),
		setRoleHandler,
	)

	// POST /api/v1/sys/im/group/transfer-owner
	r.POST("/api/v1/sys/im/group/transfer-owner",
		authMW.HeiCheckLogin(),
		log.SysLog("转让群"),
		transferOwnerHandler,
	)

	// POST /api/v1/sys/im/group/set-nickname
	r.POST("/api/v1/sys/im/group/set-nickname",
		authMW.HeiCheckLogin(),
		setNicknameHandler,
	)

	// GET /api/v1/sys/im/group/messages
	r.GET("/api/v1/sys/im/group/messages",
		authMW.HeiCheckLogin(),
		messagesHandler,
	)

	// GET /api/v1/sys/im/group/search
	r.GET("/api/v1/sys/im/group/search",
		authMW.HeiCheckLogin(),
		searchHandler,
	)

	// GET /api/v1/sys/im/group/search-groups
	r.GET("/api/v1/sys/im/group/search-groups",
		authMW.HeiCheckLogin(),
		searchGroupsHandler,
	)

	// POST /api/v1/sys/im/group/send
	r.POST("/api/v1/sys/im/group/send",
		authMW.HeiCheckLogin(),
		middleware.RateLimiter("sys_group_send", 3, 20),
		authMW.NoRepeat(3000),
		sendHandler,
	)

	// POST /api/v1/sys/im/group/recall
	r.POST("/api/v1/sys/im/group/recall",
		authMW.HeiCheckLogin(),
		recallHandler,
	)

	// POST /api/v1/sys/im/group/mark-read
	r.POST("/api/v1/sys/im/group/mark-read",
		authMW.HeiCheckLogin(),
		markReadHandler,
	)

	// POST /api/v1/sys/im/group/mute
	r.POST("/api/v1/sys/im/group/mute",
		authMW.HeiCheckLogin(),
		log.SysLog("禁言"),
		muteHandler,
	)

	// POST /api/v1/sys/im/group/unmute
	r.POST("/api/v1/sys/im/group/unmute",
		authMW.HeiCheckLogin(),
		log.SysLog("解禁"),
		unmuteHandler,
	)

	// GET /api/v1/sys/im/group/members
	r.GET("/api/v1/sys/im/group/members",
		authMW.HeiCheckLogin(),
		membersHandler,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	// POST /api/v1/c/im/group/create
	r.POST("/api/v1/c/im/group/create",
		authMW.HeiClientCheckLogin(),
		clientCreateHandler,
	)

	// GET /api/v1/c/im/group/my-groups
	r.GET("/api/v1/c/im/group/my-groups",
		authMW.HeiClientCheckLogin(),
		clientMyGroupsHandler,
	)

	// GET /api/v1/c/im/group/detail
	r.GET("/api/v1/c/im/group/detail",
		authMW.HeiClientCheckLogin(),
		detailHandler,
	)

	// POST /api/v1/c/im/group/join
	r.POST("/api/v1/c/im/group/join",
		authMW.HeiClientCheckLogin(),
		clientJoinHandler,
	)

	// POST /api/v1/c/im/group/leave
	r.POST("/api/v1/c/im/group/leave",
		authMW.HeiClientCheckLogin(),
		clientLeaveHandler,
	)

	// GET /api/v1/c/im/group/messages
	r.GET("/api/v1/c/im/group/messages",
		authMW.HeiClientCheckLogin(),
		messagesHandler,
	)

	// GET /api/v1/c/im/group/search
	r.GET("/api/v1/c/im/group/search",
		authMW.HeiClientCheckLogin(),
		searchHandler,
	)

	// GET /api/v1/c/im/group/search-groups
	r.GET("/api/v1/c/im/group/search-groups",
		authMW.HeiClientCheckLogin(),
		searchGroupsHandler,
	)

	// POST /api/v1/c/im/group/send
	r.POST("/api/v1/c/im/group/send",
		authMW.HeiClientCheckLogin(),
		middleware.RateLimiter("c_group_send", 3, 20),
		sendHandler,
	)

	// POST /api/v1/c/im/group/recall
	r.POST("/api/v1/c/im/group/recall",
		authMW.HeiClientCheckLogin(),
		recallHandler,
	)

	// POST /api/v1/c/im/group/mark-read
	r.POST("/api/v1/c/im/group/mark-read",
		authMW.HeiClientCheckLogin(),
		markReadHandler,
	)

	// GET /api/v1/c/im/group/members
	r.GET("/api/v1/c/im/group/members",
		authMW.HeiClientCheckLogin(),
		membersHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// createHandler handles POST /api/v1/sys/im/group/create
func createHandler(c *gin.Context) {
	group.GroupCreate(c)
}

// clientCreateHandler handles POST /api/v1/c/im/group/create
func clientCreateHandler(c *gin.Context) {
	group.GroupCreate(c)
}

// updateHandler handles POST /api/v1/sys/im/group/update
func updateHandler(c *gin.Context) {
	var param group.UpdateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupUpdate(c, &param)
	result.Success(c, nil)
}

// dissolveHandler handles POST /api/v1/sys/im/group/dissolve
func dissolveHandler(c *gin.Context) {
	var param group.DissolveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupDissolve(c, &param)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/im/group/detail and GET /api/v1/c/im/group/detail
func detailHandler(c *gin.Context) {
	vo := group.GroupDetail(c)
	result.Success(c, vo)
}

// myGroupsHandler handles GET /api/v1/sys/im/group/my-groups
func myGroupsHandler(c *gin.Context) {
	list := group.GroupMyGroups(c)
	result.Success(c, list)
}

// clientMyGroupsHandler handles GET /api/v1/c/im/group/my-groups
func clientMyGroupsHandler(c *gin.Context) {
	list := group.GroupMyGroups(c)
	result.Success(c, list)
}

// inviteHandler handles POST /api/v1/sys/im/group/invite
func inviteHandler(c *gin.Context) {
	var param group.InviteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupInvite(c, &param)
	result.Success(c, nil)
}

// joinHandler handles POST /api/v1/sys/im/group/join
func joinHandler(c *gin.Context) {
	var param group.JoinOrLeaveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupJoin(c, &param)
	result.Success(c, nil)
}

// clientJoinHandler handles POST /api/v1/c/im/group/join
func clientJoinHandler(c *gin.Context) {
	var param group.JoinOrLeaveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupJoin(c, &param)
	result.Success(c, nil)
}

// leaveHandler handles POST /api/v1/sys/im/group/leave
func leaveHandler(c *gin.Context) {
	var param group.JoinOrLeaveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupLeave(c, &param)
	result.Success(c, nil)
}

// clientLeaveHandler handles POST /api/v1/c/im/group/leave
func clientLeaveHandler(c *gin.Context) {
	var param group.JoinOrLeaveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupLeave(c, &param)
	result.Success(c, nil)
}

// kickHandler handles POST /api/v1/sys/im/group/kick
func kickHandler(c *gin.Context) {
	var param group.KickParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupKick(c, &param)
	result.Success(c, nil)
}

// setRoleHandler handles POST /api/v1/sys/im/group/set-role
func setRoleHandler(c *gin.Context) {
	var param group.SetRoleParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupSetRole(c, &param)
	result.Success(c, nil)
}

// transferOwnerHandler handles POST /api/v1/sys/im/group/transfer-owner
func transferOwnerHandler(c *gin.Context) {
	var param group.TransferOwnerParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupTransferOwner(c, &param)
	result.Success(c, nil)
}

// setNicknameHandler handles POST /api/v1/sys/im/group/set-nickname
func setNicknameHandler(c *gin.Context) {
	var param group.SetNicknameParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupSetMemberNickname(c, &param)
	result.Success(c, nil)
}

// pendingJoinRequestsHandler handles GET /api/v1/sys/im/group/pending-join-requests
func pendingJoinRequestsHandler(c *gin.Context) {
	requests := group.GroupPendingJoinRequests(c)
	result.Success(c, requests)
}

// handleJoinRequestHandler handles POST /api/v1/sys/im/group/handle-join-request
func handleJoinRequestHandler(c *gin.Context) {
	var param group.HandleJoinRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupHandleJoinRequest(c, &param)
	result.Success(c, nil)
}

// membersHandler handles GET /api/v1/sys/im/group/members and GET /api/v1/c/im/group/members
func membersHandler(c *gin.Context) {
	list := group.GroupMembers(c)
	result.Success(c, list)
}

// messagesHandler handles GET /api/v1/sys/im/group/messages and GET /api/v1/c/im/group/messages
func messagesHandler(c *gin.Context) {
	groupID := c.Query("group_id")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	msgs, hasMore := group.GroupMessages(c, groupID, cursor, size)
	result.Success(c, gin.H{"records": msgs, "has_more": hasMore})
}

// searchHandler handles GET /api/v1/sys/im/group/search and GET /api/v1/c/im/group/search
func searchHandler(c *gin.Context) {
	groupID := c.Query("group_id")
	keyword := c.Query("keyword")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	msgs, hasMore := group.GroupSearchMessages(c, groupID, keyword, cursor, size)
	result.Success(c, gin.H{"records": msgs, "has_more": hasMore})
}

// searchGroupsHandler handles GET /api/v1/sys/im/group/search-groups and GET /api/v1/c/im/group/search-groups
func searchGroupsHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list := group.GroupSearchGroups(c, keyword, size)
	result.Success(c, list)
}

// sendHandler handles POST /api/v1/sys/im/group/send and POST /api/v1/c/im/group/send
func sendHandler(c *gin.Context) {
	var param group.SendMessageParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupSendMessage(c, &param)
	result.Success(c, nil)
}

// recallHandler handles POST /api/v1/sys/im/group/recall and POST /api/v1/c/im/group/recall
func recallHandler(c *gin.Context) {
	var param group.RecallMessageParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupRecallMessage(c, &param)
	result.Success(c, nil)
}

// markReadHandler handles POST /api/v1/sys/im/group/mark-read and POST /api/v1/c/im/group/mark-read
func markReadHandler(c *gin.Context) {
	var param group.MarkReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupMarkRead(c, &param)
	result.Success(c, nil)
}

// muteHandler handles POST /api/v1/sys/im/group/mute
func muteHandler(c *gin.Context) {
	var param group.MuteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupMuteMember(c, &param)
	result.Success(c, nil)
}

// unmuteHandler handles POST /api/v1/sys/im/group/unmute
func unmuteHandler(c *gin.Context) {
	var param group.UnmuteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupUnmuteMember(c, &param)
	result.Success(c, nil)
}
