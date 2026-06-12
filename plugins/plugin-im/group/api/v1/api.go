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

type handler struct {
	service *group.Service
}

var defaultHandler = newHandler(group.DefaultModule)

func newHandler(module *group.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/sys/im/group/create",
		authMW.HeiCheckLogin(),
		authMW.NoRepeat(3000),
		defaultHandler.create,
	)
	r.GET("/api/v1/sys/im/group/my-groups",
		authMW.HeiCheckLogin(),
		defaultHandler.myGroups,
	)
	r.GET("/api/v1/sys/im/group/detail",
		authMW.HeiCheckLogin(),
		defaultHandler.detail,
	)
	r.POST("/api/v1/sys/im/group/update",
		authMW.HeiCheckLogin(),
		defaultHandler.update,
	)
	r.POST("/api/v1/sys/im/group/dissolve",
		authMW.HeiCheckLogin(),
		log.SysLog("解散群"),
		defaultHandler.dissolve,
	)
	r.POST("/api/v1/sys/im/group/invite",
		authMW.HeiCheckLogin(),
		authMW.NoRepeat(3000),
		defaultHandler.invite,
	)
	r.POST("/api/v1/sys/im/group/join",
		authMW.HeiCheckLogin(),
		defaultHandler.join,
	)
	r.GET("/api/v1/sys/im/group/pending-join-requests",
		authMW.HeiCheckLogin(),
		defaultHandler.pendingJoinRequests,
	)
	r.POST("/api/v1/sys/im/group/handle-join-request",
		authMW.HeiCheckLogin(),
		defaultHandler.handleJoinRequest,
	)
	r.POST("/api/v1/sys/im/group/leave",
		authMW.HeiCheckLogin(),
		defaultHandler.leave,
	)
	r.POST("/api/v1/sys/im/group/kick",
		authMW.HeiCheckLogin(),
		log.SysLog("踢出成员"),
		defaultHandler.kick,
	)
	r.POST("/api/v1/sys/im/group/set-role",
		authMW.HeiCheckLogin(),
		log.SysLog("设置角色"),
		defaultHandler.setRole,
	)
	r.POST("/api/v1/sys/im/group/transfer-owner",
		authMW.HeiCheckLogin(),
		log.SysLog("转让群"),
		defaultHandler.transferOwner,
	)
	r.POST("/api/v1/sys/im/group/set-nickname",
		authMW.HeiCheckLogin(),
		defaultHandler.setNickname,
	)
	r.GET("/api/v1/sys/im/group/messages",
		authMW.HeiCheckLogin(),
		defaultHandler.messages,
	)
	r.GET("/api/v1/sys/im/group/search",
		authMW.HeiCheckLogin(),
		defaultHandler.search,
	)
	r.GET("/api/v1/sys/im/group/search-groups",
		authMW.HeiCheckLogin(),
		defaultHandler.searchGroups,
	)
	r.POST("/api/v1/sys/im/group/send",
		authMW.HeiCheckLogin(),
		middleware.RateLimiter("sys_group_send", 3, 20),
		authMW.NoRepeat(3000),
		defaultHandler.send,
	)
	r.POST("/api/v1/sys/im/group/recall",
		authMW.HeiCheckLogin(),
		defaultHandler.recall,
	)
	r.POST("/api/v1/sys/im/group/mark-read",
		authMW.HeiCheckLogin(),
		defaultHandler.markRead,
	)
	r.POST("/api/v1/sys/im/group/mute",
		authMW.HeiCheckLogin(),
		log.SysLog("禁言"),
		defaultHandler.mute,
	)
	r.POST("/api/v1/sys/im/group/unmute",
		authMW.HeiCheckLogin(),
		log.SysLog("解禁"),
		defaultHandler.unmute,
	)
	r.GET("/api/v1/sys/im/group/members",
		authMW.HeiCheckLogin(),
		defaultHandler.members,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	r.POST("/api/v1/c/im/group/create",
		authMW.HeiClientCheckLogin(),
		defaultHandler.create,
	)
	r.GET("/api/v1/c/im/group/my-groups",
		authMW.HeiClientCheckLogin(),
		defaultHandler.myGroups,
	)
	r.GET("/api/v1/c/im/group/detail",
		authMW.HeiClientCheckLogin(),
		defaultHandler.detail,
	)
	r.POST("/api/v1/c/im/group/join",
		authMW.HeiClientCheckLogin(),
		defaultHandler.join,
	)
	r.POST("/api/v1/c/im/group/leave",
		authMW.HeiClientCheckLogin(),
		defaultHandler.leave,
	)
	r.GET("/api/v1/c/im/group/messages",
		authMW.HeiClientCheckLogin(),
		defaultHandler.messages,
	)
	r.GET("/api/v1/c/im/group/search",
		authMW.HeiClientCheckLogin(),
		defaultHandler.search,
	)
	r.GET("/api/v1/c/im/group/search-groups",
		authMW.HeiClientCheckLogin(),
		defaultHandler.searchGroups,
	)
	r.POST("/api/v1/c/im/group/send",
		authMW.HeiClientCheckLogin(),
		middleware.RateLimiter("c_group_send", 3, 20),
		defaultHandler.send,
	)
	r.POST("/api/v1/c/im/group/recall",
		authMW.HeiClientCheckLogin(),
		defaultHandler.recall,
	)
	r.POST("/api/v1/c/im/group/mark-read",
		authMW.HeiClientCheckLogin(),
		defaultHandler.markRead,
	)
	r.GET("/api/v1/c/im/group/members",
		authMW.HeiClientCheckLogin(),
		defaultHandler.members,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// @Summary      即时通讯群组创建
// @Description  访问 /api/v1/sys/im/group/create，即时通讯群组创建
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/create [post]
// @Router       /api/v1/c/im/group/create [post]
func (h *handler) create(c *gin.Context) {
	h.service.Create(c)
}

// @Summary      即时通讯群组更新
// @Description  访问 /api/v1/sys/im/group/update，即时通讯群组更新
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.UpdateParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/update [post]
func (h *handler) update(c *gin.Context) {
	var param group.UpdateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Update(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组解散群组
// @Description  访问 /api/v1/sys/im/group/dissolve，即时通讯群组解散群组
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.DissolveParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/dissolve [post]
func (h *handler) dissolve(c *gin.Context) {
	var param group.DissolveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Dissolve(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组详情查询
// @Description  访问 /api/v1/sys/im/group/detail，即时通讯群组详情查询
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/detail [get]
// @Router       /api/v1/c/im/group/detail [get]
func (h *handler) detail(c *gin.Context) {
	result.Success(c, h.service.Detail(c))
}

// @Summary      即时通讯群组我的群组列表
// @Description  访问 /api/v1/sys/im/group/my-groups，即时通讯群组我的群组列表
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/my-groups [get]
// @Router       /api/v1/c/im/group/my-groups [get]
func (h *handler) myGroups(c *gin.Context) {
	result.Success(c, h.service.MyGroups(c))
}

// @Summary      即时通讯群组邀请成员
// @Description  访问 /api/v1/sys/im/group/invite，即时通讯群组邀请成员
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.InviteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/invite [post]
func (h *handler) invite(c *gin.Context) {
	var param group.InviteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Invite(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组加入群组
// @Description  访问 /api/v1/sys/im/group/join，即时通讯群组加入群组
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.JoinOrLeaveParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/join [post]
// @Router       /api/v1/c/im/group/join [post]
func (h *handler) join(c *gin.Context) {
	var param group.JoinOrLeaveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Join(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组待处理入群申请列表
// @Description  访问 /api/v1/sys/im/group/pending-join-requests，即时通讯群组待处理入群申请列表
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/pending-join-requests [get]
func (h *handler) pendingJoinRequests(c *gin.Context) {
	result.Success(c, h.service.PendingJoinRequests(c))
}

// @Summary      即时通讯群组处理入群申请
// @Description  访问 /api/v1/sys/im/group/handle-join-request，即时通讯群组处理入群申请
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.HandleJoinRequestParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/handle-join-request [post]
func (h *handler) handleJoinRequest(c *gin.Context) {
	var param group.HandleJoinRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.HandleJoinRequest(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组退出群组
// @Description  访问 /api/v1/sys/im/group/leave，即时通讯群组退出群组
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.JoinOrLeaveParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/leave [post]
// @Router       /api/v1/c/im/group/leave [post]
func (h *handler) leave(c *gin.Context) {
	var param group.JoinOrLeaveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Leave(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组移除成员
// @Description  访问 /api/v1/sys/im/group/kick，即时通讯群组移除成员
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.KickParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/kick [post]
func (h *handler) kick(c *gin.Context) {
	var param group.KickParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Kick(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组设置成员角色
// @Description  访问 /api/v1/sys/im/group/set-role，即时通讯群组设置成员角色
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.SetRoleParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/set-role [post]
func (h *handler) setRole(c *gin.Context) {
	var param group.SetRoleParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.SetRole(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组转让群主
// @Description  访问 /api/v1/sys/im/group/transfer-owner，即时通讯群组转让群主
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.TransferOwnerParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/transfer-owner [post]
func (h *handler) transferOwner(c *gin.Context) {
	var param group.TransferOwnerParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.TransferOwner(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组设置群昵称
// @Description  访问 /api/v1/sys/im/group/set-nickname，即时通讯群组设置群昵称
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.SetNicknameParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/set-nickname [post]
func (h *handler) setNickname(c *gin.Context) {
	var param group.SetNicknameParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.SetMemberNickname(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组成员列表
// @Description  访问 /api/v1/sys/im/group/members，即时通讯群组成员列表
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/members [get]
// @Router       /api/v1/c/im/group/members [get]
func (h *handler) members(c *gin.Context) {
	result.Success(c, h.service.Members(c))
}

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
func (h *handler) messages(c *gin.Context) {
	groupID := c.Query("group_id")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	msgs, hasMore := h.service.Messages(c, groupID, cursor, size)
	result.Success(c, gin.H{"records": msgs, "has_more": hasMore})
}

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
func (h *handler) search(c *gin.Context) {
	groupID := c.Query("group_id")
	keyword := c.Query("keyword")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	msgs, hasMore := h.service.SearchMessages(c, groupID, keyword, cursor, size)
	result.Success(c, gin.H{"records": msgs, "has_more": hasMore})
}

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
func (h *handler) searchGroups(c *gin.Context) {
	keyword := c.Query("keyword")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	result.Success(c, h.service.SearchGroups(c, keyword, size))
}

// @Summary      即时通讯群组发送消息
// @Description  访问 /api/v1/sys/im/group/send，即时通讯群组发送消息
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.SendMessageParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/send [post]
// @Router       /api/v1/c/im/group/send [post]
func (h *handler) send(c *gin.Context) {
	var param group.SendMessageParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.SendMessage(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组撤回消息
// @Description  访问 /api/v1/sys/im/group/recall，即时通讯群组撤回消息
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.RecallMessageParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/recall [post]
// @Router       /api/v1/c/im/group/recall [post]
func (h *handler) recall(c *gin.Context) {
	var param group.RecallMessageParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.RecallMessage(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组标记已读
// @Description  访问 /api/v1/sys/im/group/mark-read，即时通讯群组标记已读
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.MarkReadParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/mark-read [post]
// @Router       /api/v1/c/im/group/mark-read [post]
func (h *handler) markRead(c *gin.Context) {
	var param group.MarkReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.MarkRead(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组禁言
// @Description  访问 /api/v1/sys/im/group/mute，即时通讯群组禁言
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.MuteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/mute [post]
func (h *handler) mute(c *gin.Context) {
	var param group.MuteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.MuteMember(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯群组解除禁言
// @Description  访问 /api/v1/sys/im/group/unmute，即时通讯群组解除禁言
// @Tags         即时通讯群组
// @Accept       json
// @Produce      json
// @Param        body  body  group.UnmuteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/group/unmute [post]
func (h *handler) unmute(c *gin.Context) {
	var param group.UnmuteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.UnmuteMember(c, &param)
	result.Success(c, nil)
}
