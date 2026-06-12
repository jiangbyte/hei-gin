package v1

import (
	"strconv"

	"hei-gin/plugins/plugin-im/group"
	"hei-gin/plugins/plugin-im/message"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/web/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/im/message/page
	r.GET("/api/v1/sys/im/message/page",
		authMW.HeiCheckLogin(),
		pageHandler,
	)

	// GET /api/v1/sys/im/message/detail
	r.GET("/api/v1/sys/im/message/detail",
		authMW.HeiCheckLogin(),
		detailHandler,
	)

	// GET /api/v1/sys/im/message/unread-count
	r.GET("/api/v1/sys/im/message/unread-count",
		authMW.HeiCheckLogin(),
		unreadCountHandler,
	)

	// POST /api/v1/sys/im/message/send
	r.POST("/api/v1/sys/im/message/send",
		authMW.HeiCheckLogin(),
		middleware.RateLimiter("sys_send", 5, 20),
		authMW.NoRepeat(3000),
		sendHandler,
	)

	// POST /api/v1/sys/im/message/recall
	r.POST("/api/v1/sys/im/message/recall",
		authMW.HeiCheckLogin(),
		recallHandler,
	)

	// POST /api/v1/sys/im/message/forward
	r.POST("/api/v1/sys/im/message/forward",
		authMW.HeiCheckLogin(),
		forwardHandler,
	)

	// POST /api/v1/sys/im/message/delete
	r.POST("/api/v1/sys/im/message/delete",
		authMW.HeiCheckLogin(),
		deleteHandler,
	)

	// GET /api/v1/sys/im/message/search
	r.GET("/api/v1/sys/im/message/search",
		authMW.HeiCheckLogin(),
		searchHandler,
	)

	// POST /api/v1/sys/im/message/mark-read
	r.POST("/api/v1/sys/im/message/mark-read",
		authMW.HeiCheckLogin(),
		markReadHandler,
	)

	// POST /api/v1/sys/im/message/mark-all-read
	r.POST("/api/v1/sys/im/message/mark-all-read",
		authMW.HeiCheckLogin(),
		markAllReadHandler,
	)

	// POST /api/v1/sys/im/message/remove
	r.POST("/api/v1/sys/im/message/remove",
		authMW.HeiCheckLogin(),
		removeHandler,
	)

	// GET /api/v1/sys/im/conversation/list
	r.GET("/api/v1/sys/im/conversation/list",
		authMW.HeiCheckLogin(),
		conversationsHandler,
	)

	// GET /api/v1/sys/im/conversation/messages
	r.GET("/api/v1/sys/im/conversation/messages",
		authMW.HeiCheckLogin(),
		conversationMessagesHandler,
	)

	// POST /api/v1/sys/im/conversation/read
	r.POST("/api/v1/sys/im/conversation/read",
		authMW.HeiCheckLogin(),
		conversationReadHandler,
	)

	// POST /api/v1/sys/im/conversation/get-or-create
	r.POST("/api/v1/sys/im/conversation/get-or-create",
		authMW.HeiCheckLogin(),
		getOrCreateConversationHandler,
	)

	// POST /api/v1/sys/im/file/upload
	r.POST("/api/v1/sys/im/file/upload",
		authMW.HeiCheckLogin(),
		uploadFileHandler,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	// POST /api/v1/c/im/message/send
	r.POST("/api/v1/c/im/message/send",
		authMW.HeiClientCheckLogin(),
		middleware.RateLimiter("c_send", 5, 20),
		clientSendHandler,
	)

	// POST /api/v1/c/im/message/recall
	r.POST("/api/v1/c/im/message/recall",
		authMW.HeiClientCheckLogin(),
		clientRecallHandler,
	)

	// POST /api/v1/c/im/message/forward
	r.POST("/api/v1/c/im/message/forward",
		authMW.HeiClientCheckLogin(),
		clientForwardHandler,
	)

	// POST /api/v1/c/im/message/delete
	r.POST("/api/v1/c/im/message/delete",
		authMW.HeiClientCheckLogin(),
		clientDeleteHandler,
	)

	// GET /api/v1/c/im/message/search
	r.GET("/api/v1/c/im/message/search",
		authMW.HeiClientCheckLogin(),
		clientSearchHandler,
	)

	// POST /api/v1/c/im/message/mark-read
	r.POST("/api/v1/c/im/message/mark-read",
		authMW.HeiClientCheckLogin(),
		clientMarkReadHandler,
	)

	// POST /api/v1/c/im/message/mark-all-read
	r.POST("/api/v1/c/im/message/mark-all-read",
		authMW.HeiClientCheckLogin(),
		clientMarkAllReadHandler,
	)

	// POST /api/v1/c/im/message/remove
	r.POST("/api/v1/c/im/message/remove",
		authMW.HeiClientCheckLogin(),
		clientRemoveHandler,
	)

	// GET /api/v1/c/im/conversation/list
	r.GET("/api/v1/c/im/conversation/list",
		authMW.HeiClientCheckLogin(),
		clientConversationsHandler,
	)

	// GET /api/v1/c/im/conversation/messages
	r.GET("/api/v1/c/im/conversation/messages",
		authMW.HeiClientCheckLogin(),
		clientConversationMessagesHandler,
	)

	// POST /api/v1/c/im/conversation/read
	r.POST("/api/v1/c/im/conversation/read",
		authMW.HeiClientCheckLogin(),
		clientConversationReadHandler,
	)

	// POST /api/v1/c/im/conversation/get-or-create
	r.POST("/api/v1/c/im/conversation/get-or-create",
		authMW.HeiClientCheckLogin(),
		clientGetOrCreateConversationHandler,
	)

	// POST /api/v1/c/im/file/upload
	r.POST("/api/v1/c/im/file/upload",
		authMW.HeiClientCheckLogin(),
		clientUploadFileHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// pageHandler handles GET /api/v1/sys/im/message/page
// @Summary      即时通讯消息分页查询
// @Description  访问 /api/v1/sys/im/message/page，即时通讯消息分页查询
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        query  query  message.MessagePageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/page [get]
func pageHandler(c *gin.Context) {
	var param message.MessagePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessagePage(c, &param)
}

// detailHandler handles GET /api/v1/sys/im/message/detail
// @Summary      即时通讯消息详情查询
// @Description  访问 /api/v1/sys/im/message/detail，即时通讯消息详情查询
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/detail [get]
func detailHandler(c *gin.Context) {
	message.MessageDetail(c)
}

// unreadCountHandler handles GET /api/v1/sys/im/message/unread-count
// @Summary      即时通讯消息未读数
// @Description  访问 /api/v1/sys/im/message/unread-count，即时通讯消息未读数
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/unread-count [get]
func unreadCountHandler(c *gin.Context) {
	message.MessageUnreadCount(c)
}

// sendHandler handles POST /api/v1/sys/im/message/send
// @Summary      即时通讯消息发送消息
// @Description  访问 /api/v1/sys/im/message/send，即时通讯消息发送消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.MessageSendParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/send [post]
func sendHandler(c *gin.Context) {
	var param message.MessageSendParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageSend(c, &param)
	result.Success(c, nil)
}

// clientSendHandler handles POST /api/v1/c/im/message/send
// @Summary      即时通讯消息发送消息
// @Description  访问 /api/v1/c/im/message/send，即时通讯消息发送消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.MessageSendParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/message/send [post]
func clientSendHandler(c *gin.Context) {
	var param message.MessageSendParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageSend(c, &param)
	result.Success(c, nil)
}

// recallHandler handles POST /api/v1/sys/im/message/recall
// @Summary      即时通讯消息撤回消息
// @Description  访问 /api/v1/sys/im/message/recall，即时通讯消息撤回消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.RecallParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/recall [post]
func recallHandler(c *gin.Context) {
	var param message.RecallParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageRecall(c, &param)
	result.Success(c, nil)
}

// clientRecallHandler handles POST /api/v1/c/im/message/recall
// @Summary      即时通讯消息撤回消息
// @Description  访问 /api/v1/c/im/message/recall，即时通讯消息撤回消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.RecallParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/message/recall [post]
func clientRecallHandler(c *gin.Context) {
	var param message.RecallParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageRecall(c, &param)
	result.Success(c, nil)
}

// forwardHandler handles POST /api/v1/sys/im/message/forward
// @Summary      即时通讯消息转发消息
// @Description  访问 /api/v1/sys/im/message/forward，即时通讯消息转发消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.ForwardParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/forward [post]
func forwardHandler(c *gin.Context) {
	var param message.ForwardParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageForward(c, &param)
	result.Success(c, nil)
}

// clientForwardHandler handles POST /api/v1/c/im/message/forward
// @Summary      即时通讯消息转发消息
// @Description  访问 /api/v1/c/im/message/forward，即时通讯消息转发消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.ForwardParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/message/forward [post]
func clientForwardHandler(c *gin.Context) {
	var param message.ForwardParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageForward(c, &param)
	result.Success(c, nil)
}

// deleteHandler handles POST /api/v1/sys/im/message/delete
// @Summary      即时通讯消息删除消息
// @Description  访问 /api/v1/sys/im/message/delete，即时通讯消息删除消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.DeleteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/delete [post]
func deleteHandler(c *gin.Context) {
	var param message.DeleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageRemove(c, &param)
	result.Success(c, nil)
}

// clientDeleteHandler handles POST /api/v1/c/im/message/delete
// @Summary      即时通讯消息删除消息
// @Description  访问 /api/v1/c/im/message/delete，即时通讯消息删除消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.DeleteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/message/delete [post]
func clientDeleteHandler(c *gin.Context) {
	var param message.DeleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageRemove(c, &param)
	result.Success(c, nil)
}

// searchHandler handles GET /api/v1/sys/im/message/search
// @Summary      即时通讯消息搜索
// @Description  访问 /api/v1/sys/im/message/search，即时通讯消息搜索
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        query  query  message.SearchParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/search [get]
func searchHandler(c *gin.Context) {
	var param message.SearchParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	list, hasMore := message.MessageSearch(c, &param)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

// clientSearchHandler handles GET /api/v1/c/im/message/search
// @Summary      即时通讯消息搜索
// @Description  访问 /api/v1/c/im/message/search，即时通讯消息搜索
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        query  query  message.SearchParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/message/search [get]
func clientSearchHandler(c *gin.Context) {
	var param message.SearchParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	list, hasMore := message.MessageSearch(c, &param)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

// markReadHandler handles POST /api/v1/sys/im/message/mark-read
// @Summary      即时通讯消息标记已读
// @Description  访问 /api/v1/sys/im/message/mark-read，即时通讯消息标记已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/mark-read [post]
func markReadHandler(c *gin.Context) {
	var param utils.IdParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageMarkRead(c, &param)
	result.Success(c, nil)
}

// clientMarkReadHandler handles POST /api/v1/c/im/message/mark-read
// @Summary      即时通讯消息标记已读
// @Description  访问 /api/v1/c/im/message/mark-read，即时通讯消息标记已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/message/mark-read [post]
func clientMarkReadHandler(c *gin.Context) {
	var param utils.IdParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageMarkRead(c, &param)
	result.Success(c, nil)
}

// markAllReadHandler handles POST /api/v1/sys/im/message/mark-all-read
// @Summary      即时通讯消息全部标记已读
// @Description  访问 /api/v1/sys/im/message/mark-all-read，即时通讯消息全部标记已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/mark-all-read [post]
func markAllReadHandler(c *gin.Context) {
	message.MessageMarkAllRead(c)
	result.Success(c, nil)
}

// clientMarkAllReadHandler handles POST /api/v1/c/im/message/mark-all-read
// @Summary      即时通讯消息全部标记已读
// @Description  访问 /api/v1/c/im/message/mark-all-read，即时通讯消息全部标记已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/message/mark-all-read [post]
func clientMarkAllReadHandler(c *gin.Context) {
	message.MessageMarkAllRead(c)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/im/message/remove
// @Summary      即时通讯消息删除
// @Description  访问 /api/v1/sys/im/message/remove，即时通讯消息删除
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.DeleteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/remove [post]
func removeHandler(c *gin.Context) {
	var param message.DeleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageRemove(c, &param)
	result.Success(c, nil)
}

// clientRemoveHandler handles POST /api/v1/c/im/message/remove
// @Summary      即时通讯消息删除
// @Description  访问 /api/v1/c/im/message/remove，即时通讯消息删除
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.DeleteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/message/remove [post]
func clientRemoveHandler(c *gin.Context) {
	var param message.DeleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageRemove(c, &param)
	result.Success(c, nil)
}

// conversationsHandler handles GET /api/v1/sys/im/conversation/list
// @Summary      即时通讯消息会话列表
// @Description  访问 /api/v1/sys/im/conversation/list，即时通讯消息会话列表
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        cursor  query  string  false  "cursor"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/conversation/list [get]
func conversationsHandler(c *gin.Context) {
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list, hasMore := message.MessageConversations(c, cursor, size)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

// clientConversationsHandler handles GET /api/v1/c/im/conversation/list
// @Summary      即时通讯消息会话列表
// @Description  访问 /api/v1/c/im/conversation/list，即时通讯消息会话列表
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        cursor  query  string  false  "cursor"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/conversation/list [get]
func clientConversationsHandler(c *gin.Context) {
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list, hasMore := message.MessageConversations(c, cursor, size)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

// conversationMessagesHandler handles GET /api/v1/sys/im/conversation/messages
// @Summary      即时通讯消息会话消息列表
// @Description  访问 /api/v1/sys/im/conversation/messages，即时通讯消息会话消息列表
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        conversation_id  query  string  false  "conversation_id"
// @Param        cursor  query  string  false  "cursor"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/conversation/messages [get]
func conversationMessagesHandler(c *gin.Context) {
	cid := c.Query("conversation_id")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}

	var messages []message.ConversationMessageVO
	var hasMore bool
	if len(cid) > 6 && cid[:6] == "group:" {
		gid := cid[6:]
		msgs, more := group.Messages(c.Request.Context(), gid, cursor, size)
		messages = make([]message.ConversationMessageVO, len(msgs))
		for i, m := range msgs {
			messages[i] = message.ConversationMessageVO{
				ID: m.ID, SenderID: m.SenderID, SenderType: m.SenderType,
				Content: m.Content, MsgType: m.MsgType, Extra: m.Extra,
				CreatedAt: m.CreatedAt,
			}
		}
		hasMore = more
	} else {
		messages, hasMore = message.MessageConversationMessages(c, cid, cursor, size)
	}
	result.Success(c, gin.H{
		"records":  messages,
		"has_more": hasMore,
	})
}

// clientConversationMessagesHandler handles GET /api/v1/c/im/conversation/messages
// @Summary      即时通讯消息会话消息列表
// @Description  访问 /api/v1/c/im/conversation/messages，即时通讯消息会话消息列表
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        conversation_id  query  string  false  "conversation_id"
// @Param        cursor  query  string  false  "cursor"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/conversation/messages [get]
func clientConversationMessagesHandler(c *gin.Context) {
	cid := c.Query("conversation_id")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}

	var messages []message.ConversationMessageVO
	var hasMore bool
	if len(cid) > 6 && cid[:6] == "group:" {
		gid := cid[6:]
		msgs, more := group.Messages(c.Request.Context(), gid, cursor, size)
		messages = make([]message.ConversationMessageVO, len(msgs))
		for i, m := range msgs {
			messages[i] = message.ConversationMessageVO{
				ID: m.ID, SenderID: m.SenderID, SenderType: m.SenderType,
				Content: m.Content, MsgType: m.MsgType, Extra: m.Extra,
				CreatedAt: m.CreatedAt,
			}
		}
		hasMore = more
	} else {
		messages, hasMore = message.MessageConversationMessages(c, cid, cursor, size)
	}
	result.Success(c, gin.H{
		"records":  messages,
		"has_more": hasMore,
	})
}

// conversationReadHandler handles POST /api/v1/sys/im/conversation/read
// @Summary      即时通讯消息会话已读
// @Description  访问 /api/v1/sys/im/conversation/read，即时通讯消息会话已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.ConversationReadParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/conversation/read [post]
func conversationReadHandler(c *gin.Context) {
	var param message.ConversationReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if len(param.ConversationID) > 6 && param.ConversationID[:6] == "group:" {
		group.MarkConversationRead(c.Request.Context(), param.ConversationID[6:], auth.GetLoginID(c), string(enums.LoginTypeBusiness))
	} else {
		message.MessageMarkConversationRead(c, &param)
	}
	result.Success(c, nil)
}

// clientConversationReadHandler handles POST /api/v1/c/im/conversation/read
// @Summary      即时通讯消息会话已读
// @Description  访问 /api/v1/c/im/conversation/read，即时通讯消息会话已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.ConversationReadParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/conversation/read [post]
func clientConversationReadHandler(c *gin.Context) {
	var param message.ConversationReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if len(param.ConversationID) > 6 && param.ConversationID[:6] == "group:" {
		group.MarkConversationRead(c.Request.Context(), param.ConversationID[6:], auth.Consumer.GetLoginID(c), string(enums.LoginTypeConsumer))
	} else {
		message.MessageMarkConversationRead(c, &param)
	}
	result.Success(c, nil)
}

// getOrCreateConversationHandler handles POST /api/v1/sys/im/conversation/get-or-create
// @Summary      即时通讯消息获取或创建会话
// @Description  访问 /api/v1/sys/im/conversation/get-or-create，即时通讯消息获取或创建会话
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.GetOrCreateConversationParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/conversation/get-or-create [post]
func getOrCreateConversationHandler(c *gin.Context) {
	var param message.GetOrCreateConversationParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageGetOrCreateConversation(c, &param)
}

// clientGetOrCreateConversationHandler handles POST /api/v1/c/im/conversation/get-or-create
// @Summary      即时通讯消息获取或创建会话
// @Description  访问 /api/v1/c/im/conversation/get-or-create，即时通讯消息获取或创建会话
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.GetOrCreateConversationParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/conversation/get-or-create [post]
func clientGetOrCreateConversationHandler(c *gin.Context) {
	var param message.GetOrCreateConversationParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MessageGetOrCreateConversation(c, &param)
}

// uploadFileHandler handles POST /api/v1/sys/im/file/upload
// @Summary      即时通讯消息上传文件
// @Description  访问 /api/v1/sys/im/file/upload，即时通讯消息上传文件
// @Tags         即时通讯消息
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "上传文件"
// @Param        engine  formData  string  false  "存储引擎"
// @Param        bucket  formData  string  false  "存储桶"
// @Param        msg_type  formData  string  false  "消息类型"
// @Param        conversation_id  formData  string  false  "会话 ID"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/file/upload [post]
func uploadFileHandler(c *gin.Context) {
	message.UploadFile(c, auth.GetLoginID(c), string(enums.LoginTypeBusiness))
}

// clientUploadFileHandler handles POST /api/v1/c/im/file/upload
// @Summary      即时通讯消息上传文件
// @Description  访问 /api/v1/c/im/file/upload，即时通讯消息上传文件
// @Tags         即时通讯消息
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "上传文件"
// @Param        engine  formData  string  false  "存储引擎"
// @Param        bucket  formData  string  false  "存储桶"
// @Param        msg_type  formData  string  false  "消息类型"
// @Param        conversation_id  formData  string  false  "会话 ID"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/im/file/upload [post]
func clientUploadFileHandler(c *gin.Context) {
	message.UploadFile(c, auth.Consumer.GetLoginID(c), string(enums.LoginTypeConsumer))
}
