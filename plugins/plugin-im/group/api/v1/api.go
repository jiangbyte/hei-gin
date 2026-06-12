package v1

import (
	"strconv"

	"hei-gin/plugins/plugin-im/group"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	"hei-gin/sdk/web/middleware"
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
// @Summary      即时通讯群组创建
// @Description  访问 /api/v1/sys/im/group/create，即时通讯群组创建
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/create [post]
func createHandler(c *gin.Context) {
	group.GroupCreate(c)
}

// clientCreateHandler handles POST /api/v1/c/im/group/create
// @Summary      即时通讯群组创建
// @Description  访问 /api/v1/c/im/group/create，即时通讯群组创建
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/group/create [post]
func clientCreateHandler(c *gin.Context) {
	group.GroupCreate(c)
}

// updateHandler handles POST /api/v1/sys/im/group/update
// @Summary      即时通讯群组更新
// @Description  访问 /api/v1/sys/im/group/update，即时通讯群组更新
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.UpdateParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/update [post]
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
// @Summary      即时通讯群组解散群组
// @Description  访问 /api/v1/sys/im/group/dissolve，即时通讯群组解散群组
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.DissolveParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/dissolve [post]
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
// @Summary      即时通讯群组详情查询
// @Description  访问 /api/v1/sys/im/group/detail，即时通讯群组详情查询
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/detail [get]
// @Router       /api/v1/c/im/group/detail [get]
func detailHandler(c *gin.Context) {
	vo := group.GroupDetail(c)
	result.Success(c, vo)
}

// myGroupsHandler handles GET /api/v1/sys/im/group/my-groups
// @Summary      即时通讯群组我的群组列表
// @Description  访问 /api/v1/sys/im/group/my-groups，即时通讯群组我的群组列表
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/my-groups [get]
func myGroupsHandler(c *gin.Context) {
	list := group.GroupMyGroups(c)
	result.Success(c, list)
}

// clientMyGroupsHandler handles GET /api/v1/c/im/group/my-groups
// @Summary      即时通讯群组我的群组列表
// @Description  访问 /api/v1/c/im/group/my-groups，即时通讯群组我的群组列表
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/group/my-groups [get]
func clientMyGroupsHandler(c *gin.Context) {
	list := group.GroupMyGroups(c)
	result.Success(c, list)
}

// inviteHandler handles POST /api/v1/sys/im/group/invite
// @Summary      即时通讯群组邀请成员
// @Description  访问 /api/v1/sys/im/group/invite，即时通讯群组邀请成员
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.InviteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/invite [post]
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
// @Summary      即时通讯群组加入群组
// @Description  访问 /api/v1/sys/im/group/join，即时通讯群组加入群组
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.JoinOrLeaveParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/join [post]
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
// @Summary      即时通讯群组加入群组
// @Description  访问 /api/v1/c/im/group/join，即时通讯群组加入群组
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.JoinOrLeaveParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/group/join [post]
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
// @Summary      即时通讯群组退出群组
// @Description  访问 /api/v1/sys/im/group/leave，即时通讯群组退出群组
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.JoinOrLeaveParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/leave [post]
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
// @Summary      即时通讯群组退出群组
// @Description  访问 /api/v1/c/im/group/leave，即时通讯群组退出群组
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.JoinOrLeaveParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/group/leave [post]
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
// @Summary      即时通讯群组移除成员
// @Description  访问 /api/v1/sys/im/group/kick，即时通讯群组移除成员
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.KickParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/kick [post]
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
// @Summary      即时通讯群组设置成员角色
// @Description  访问 /api/v1/sys/im/group/set-role，即时通讯群组设置成员角色
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.SetRoleParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/set-role [post]
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
// @Summary      即时通讯群组转让群主
// @Description  访问 /api/v1/sys/im/group/transfer-owner，即时通讯群组转让群主
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.TransferOwnerParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/transfer-owner [post]
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
// @Summary      即时通讯群组设置群昵称
// @Description  访问 /api/v1/sys/im/group/set-nickname，即时通讯群组设置群昵称
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.SetNicknameParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/set-nickname [post]
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
// @Summary      即时通讯群组待处理入群申请列表
// @Description  访问 /api/v1/sys/im/group/pending-join-requests，即时通讯群组待处理入群申请列表
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/pending-join-requests [get]
func pendingJoinRequestsHandler(c *gin.Context) {
	requests := group.GroupPendingJoinRequests(c)
	result.Success(c, requests)
}

// handleJoinRequestHandler handles POST /api/v1/sys/im/group/handle-join-request
// @Summary      即时通讯群组处理入群申请
// @Description  访问 /api/v1/sys/im/group/handle-join-request，即时通讯群组处理入群申请
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.HandleJoinRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/handle-join-request [post]
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
// @Summary      即时通讯群组成员列表
// @Description  访问 /api/v1/sys/im/group/members，即时通讯群组成员列表
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/members [get]
// @Router       /api/v1/c/im/group/members [get]
func membersHandler(c *gin.Context) {
	list := group.GroupMembers(c)
	result.Success(c, list)
}

// messagesHandler handles GET /api/v1/sys/im/group/messages and GET /api/v1/c/im/group/messages
// @Summary      即时通讯群组消息列表
// @Description  访问 /api/v1/sys/im/group/messages，即时通讯群组消息列表
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        group_id  query  string  false  "group_id"
// @Param        cursor  query  string  false  "cursor"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/messages [get]
// @Router       /api/v1/c/im/group/messages [get]
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
// @Summary      即时通讯群组搜索
// @Description  访问 /api/v1/sys/im/group/search，即时通讯群组搜索
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        group_id  query  string  false  "group_id"
// @Param        keyword  query  string  false  "keyword"
// @Param        cursor  query  string  false  "cursor"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/search [get]
// @Router       /api/v1/c/im/group/search [get]
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
// @Summary      即时通讯群组群组搜索
// @Description  访问 /api/v1/sys/im/group/search-groups，即时通讯群组群组搜索
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        keyword  query  string  false  "keyword"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/search-groups [get]
// @Router       /api/v1/c/im/group/search-groups [get]
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
// @Summary      即时通讯群组发送消息
// @Description  访问 /api/v1/sys/im/group/send，即时通讯群组发送消息
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.SendMessageParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/send [post]
// @Router       /api/v1/c/im/group/send [post]
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
// @Summary      即时通讯群组撤回消息
// @Description  访问 /api/v1/sys/im/group/recall，即时通讯群组撤回消息
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.RecallMessageParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/recall [post]
// @Router       /api/v1/c/im/group/recall [post]
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
// @Summary      即时通讯群组标记已读
// @Description  访问 /api/v1/sys/im/group/mark-read，即时通讯群组标记已读
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.MarkReadParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/mark-read [post]
// @Router       /api/v1/c/im/group/mark-read [post]
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
// @Summary      即时通讯群组禁言
// @Description  访问 /api/v1/sys/im/group/mute，即时通讯群组禁言
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.MuteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/mute [post]
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
// @Summary      即时通讯群组解除禁言
// @Description  访问 /api/v1/sys/im/group/unmute，即时通讯群组解除禁言
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.UnmuteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/unmute [post]
func unmuteHandler(c *gin.Context) {
	var param group.UnmuteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	group.GroupUnmuteMember(c, &param)
	result.Success(c, nil)
}
