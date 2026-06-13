package v1

import (
	"strconv"
	"strings"

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

type handler struct {
	service *message.Service
}

var defaultHandler = newHandler(message.DefaultModule)

func newHandler(module *message.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/sys/im/message/page",
		authMW.CheckLogin(auth.Business),
		defaultHandler.page,
	)
	r.GET("/api/v1/sys/im/message/detail",
		authMW.CheckLogin(auth.Business),
		defaultHandler.detail,
	)
	r.GET("/api/v1/sys/im/message/unread-count",
		authMW.CheckLogin(auth.Business),
		defaultHandler.unreadCount,
	)
	r.POST("/api/v1/sys/im/message/send",
		authMW.CheckLogin(auth.Business),
		middleware.RateLimiter("sys_send", 5, 20),
		authMW.NoRepeat(3000),
		defaultHandler.send,
	)
	r.POST("/api/v1/sys/im/message/recall",
		authMW.CheckLogin(auth.Business),
		defaultHandler.recall,
	)
	r.POST("/api/v1/sys/im/message/forward",
		authMW.CheckLogin(auth.Business),
		defaultHandler.forward,
	)
	r.POST("/api/v1/sys/im/message/delete",
		authMW.CheckLogin(auth.Business),
		defaultHandler.delete,
	)
	r.GET("/api/v1/sys/im/message/search",
		authMW.CheckLogin(auth.Business),
		defaultHandler.search,
	)
	r.POST("/api/v1/sys/im/message/mark-read",
		authMW.CheckLogin(auth.Business),
		defaultHandler.markRead,
	)
	r.POST("/api/v1/sys/im/message/mark-all-read",
		authMW.CheckLogin(auth.Business),
		defaultHandler.markAllRead,
	)
	r.POST("/api/v1/sys/im/message/remove",
		authMW.CheckLogin(auth.Business),
		defaultHandler.remove,
	)
	r.GET("/api/v1/sys/im/conversation/list",
		authMW.CheckLogin(auth.Business),
		defaultHandler.conversations,
	)
	r.GET("/api/v1/sys/im/conversation/messages",
		authMW.CheckLogin(auth.Business),
		defaultHandler.conversationMessages,
	)
	r.POST("/api/v1/sys/im/conversation/read",
		authMW.CheckLogin(auth.Business),
		defaultHandler.conversationRead,
	)
	r.POST("/api/v1/sys/im/conversation/get-or-create",
		authMW.CheckLogin(auth.Business),
		defaultHandler.getOrCreateConversation,
	)
	r.POST("/api/v1/sys/im/file/upload",
		authMW.CheckLogin(auth.Business),
		defaultHandler.uploadFile,
	)
}

func RegisterClientRoutes(r *gin.Engine) {
	r.POST("/api/v1/c/im/message/send",
		authMW.CheckLogin(auth.Consumer),
		middleware.RateLimiter("c_send", 5, 20),
		defaultHandler.send,
	)
	r.POST("/api/v1/c/im/message/recall",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.recall,
	)
	r.POST("/api/v1/c/im/message/forward",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.forward,
	)
	r.POST("/api/v1/c/im/message/delete",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.delete,
	)
	r.GET("/api/v1/c/im/message/search",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.search,
	)
	r.POST("/api/v1/c/im/message/mark-read",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.markRead,
	)
	r.POST("/api/v1/c/im/message/mark-all-read",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.markAllRead,
	)
	r.POST("/api/v1/c/im/message/remove",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.remove,
	)
	r.GET("/api/v1/c/im/conversation/list",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.conversations,
	)
	r.GET("/api/v1/c/im/conversation/messages",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.conversationMessages,
	)
	r.POST("/api/v1/c/im/conversation/read",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.conversationRead,
	)
	r.POST("/api/v1/c/im/conversation/get-or-create",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.getOrCreateConversation,
	)
	r.POST("/api/v1/c/im/file/upload",
		authMW.CheckLogin(auth.Consumer),
		defaultHandler.clientUploadFile,
	)
}

func Register() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// @Summary      即时通讯消息分页查询
// @Description  访问 /api/v1/sys/im/message/page，即时通讯消息分页查询
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        query  query  message.MessagePageParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/page [get]
func (h *handler) page(c *gin.Context) {
	var param message.MessagePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Page(c, &param)
}

// @Summary      即时通讯消息详情查询
// @Description  访问 /api/v1/sys/im/message/detail，即时通讯消息详情查询
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/detail [get]
func (h *handler) detail(c *gin.Context) {
	h.service.Detail(c)
}

// @Summary      即时通讯消息未读数
// @Description  访问 /api/v1/sys/im/message/unread-count，即时通讯消息未读数
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/unread-count [get]
func (h *handler) unreadCount(c *gin.Context) {
	h.service.UnreadCount(c)
}

// @Summary      即时通讯消息发送消息
// @Description  访问 /api/v1/sys/im/message/send，即时通讯消息发送消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.MessageSendParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/send [post]
// @Router       /api/v1/c/im/message/send [post]
func (h *handler) send(c *gin.Context) {
	var param message.MessageSendParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Send(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯消息撤回消息
// @Description  访问 /api/v1/sys/im/message/recall，即时通讯消息撤回消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.RecallParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/recall [post]
// @Router       /api/v1/c/im/message/recall [post]
func (h *handler) recall(c *gin.Context) {
	var param message.RecallParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Recall(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯消息转发消息
// @Description  访问 /api/v1/sys/im/message/forward，即时通讯消息转发消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.ForwardParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/forward [post]
// @Router       /api/v1/c/im/message/forward [post]
func (h *handler) forward(c *gin.Context) {
	var param message.ForwardParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Forward(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯消息删除消息
// @Description  访问 /api/v1/sys/im/message/delete，即时通讯消息删除消息
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.DeleteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/delete [post]
// @Router       /api/v1/c/im/message/delete [post]
func (h *handler) delete(c *gin.Context) {
	var param message.DeleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯消息搜索
// @Description  访问 /api/v1/sys/im/message/search，即时通讯消息搜索
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        query  query  message.SearchParam  false  "查询参数"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/search [get]
// @Router       /api/v1/c/im/message/search [get]
func (h *handler) search(c *gin.Context) {
	var param message.SearchParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	list, hasMore := h.service.Search(c, &param)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

// @Summary      即时通讯消息标记已读
// @Description  访问 /api/v1/sys/im/message/mark-read，即时通讯消息标记已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  utils.IdParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/mark-read [post]
// @Router       /api/v1/c/im/message/mark-read [post]
func (h *handler) markRead(c *gin.Context) {
	var param utils.IdParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.MarkRead(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯消息全部标记已读
// @Description  访问 /api/v1/sys/im/message/mark-all-read，即时通讯消息全部标记已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/mark-all-read [post]
// @Router       /api/v1/c/im/message/mark-all-read [post]
func (h *handler) markAllRead(c *gin.Context) {
	h.service.MarkAllRead(c)
	result.Success(c, nil)
}

// @Summary      即时通讯消息删除
// @Description  访问 /api/v1/sys/im/message/remove，即时通讯消息删除
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.DeleteParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/message/remove [post]
// @Router       /api/v1/c/im/message/remove [post]
func (h *handler) remove(c *gin.Context) {
	var param message.DeleteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Remove(c, &param)
	result.Success(c, nil)
}

// @Summary      即时通讯消息会话列表
// @Description  访问 /api/v1/sys/im/conversation/list，即时通讯消息会话列表
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        cursor  query  string  false  "cursor"
// @Param        size  query  string  false  "size"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/conversation/list [get]
// @Router       /api/v1/c/im/conversation/list [get]
func (h *handler) conversations(c *gin.Context) {
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list, hasMore := h.service.Conversations(c, cursor, size)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

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
// @Router       /api/v1/c/im/conversation/messages [get]
func (h *handler) conversationMessages(c *gin.Context) {
	cid := c.Query("conversation_id")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}

	var messagesVO []message.ConversationMessageVO
	var hasMore bool
	if strings.HasPrefix(cid, "group:") {
		gid := strings.TrimPrefix(cid, "group:")
		msgs, more := group.DefaultModule.Service().Messages(c, gid, cursor, size)
		messagesVO = make([]message.ConversationMessageVO, len(msgs))
		for i, m := range msgs {
			messagesVO[i] = message.ConversationMessageVO{
				ID:         m.ID,
				SenderID:   m.SenderID,
				SenderType: m.SenderType,
				Content:    m.Content,
				MsgType:    m.MsgType,
				Extra:      m.Extra,
				CreatedAt:  m.CreatedAt,
			}
		}
		hasMore = more
	} else {
		messagesVO, hasMore = h.service.ConversationMessages(c, cid, cursor, size)
	}
	result.Success(c, gin.H{"records": messagesVO, "has_more": hasMore})
}

// @Summary      即时通讯消息会话已读
// @Description  访问 /api/v1/sys/im/conversation/read，即时通讯消息会话已读
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.ConversationReadParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/conversation/read [post]
// @Router       /api/v1/c/im/conversation/read [post]
func (h *handler) conversationRead(c *gin.Context) {
	var param message.ConversationReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	if strings.HasPrefix(param.ConversationID, "group:") {
		userID := auth.Business.GetLoginID(c)
		userType := string(enums.LoginTypeBusiness)
		if strings.HasPrefix(c.Request.URL.Path, "/api/v") && strings.Contains(c.Request.URL.Path, "/c/") {
			userID = auth.Consumer.GetLoginID(c)
			userType = string(enums.LoginTypeConsumer)
		}
		group.DefaultModule.Service().MarkConversationReadWithContext(
			c.Request.Context(),
			strings.TrimPrefix(param.ConversationID, "group:"),
			userID,
			userType,
		)
	} else {
		h.service.MarkConversationRead(c, &param)
	}
	result.Success(c, nil)
}

// @Summary      即时通讯消息获取或创建会话
// @Description  访问 /api/v1/sys/im/conversation/get-or-create，即时通讯消息获取或创建会话
// @Tags         即时通讯消息
// @Accept       json
// @Produce      json
// @Param        body  body  message.GetOrCreateConversationParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/sys/im/conversation/get-or-create [post]
// @Router       /api/v1/c/im/conversation/get-or-create [post]
func (h *handler) getOrCreateConversation(c *gin.Context) {
	var param message.GetOrCreateConversationParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.GetOrCreateConversation(c, &param)
}

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
func (h *handler) uploadFile(c *gin.Context) {
	h.service.UploadFile(c, auth.Business.GetLoginID(c), string(enums.LoginTypeBusiness))
}

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
func (h *handler) clientUploadFile(c *gin.Context) {
	h.service.UploadFile(c, auth.Consumer.GetLoginID(c), string(enums.LoginTypeConsumer))
}
