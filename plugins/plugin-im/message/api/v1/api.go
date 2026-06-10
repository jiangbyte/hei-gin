package v1

import (
	"strconv"

	"hei-gin/sdk/middleware"
	"hei-gin/sdk/auth"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/plugins/plugin-im/message"
	"hei-gin/plugins/plugin-im/group"

	"github.com/gin-gonic/gin"
	"hei-gin/sdk/registry"
)

func RegisterRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/sys/im").Use(authMW.HeiCheckLogin())
	{
		// Message
		g.GET("/message/page", pageHandler)
		g.GET("/message/detail", detailHandler)
		g.GET("/message/unread-count", unreadCountHandler)
		g.POST("/message/send", middleware.RateLimiter("sys_send", 5, 20), authMW.NoRepeat(3000), sendHandler)
		g.POST("/message/recall", recallHandler)
		g.POST("/message/forward", forwardHandler)
		g.POST("/message/delete", deleteHandler)
		g.GET("/message/search", searchHandler)
		g.POST("/message/mark-read", markReadHandler)
		g.POST("/message/mark-all-read", markAllReadHandler)
		g.POST("/message/remove", removeHandler)

		// Conversation
		g.GET("/conversation/list", conversationsHandler)
		g.GET("/conversation/messages", conversationMessagesHandler)
		g.POST("/conversation/read", conversationReadHandler)
		g.POST("/conversation/get-or-create", getOrCreateConversationHandler)
		g.POST("/file/upload", uploadFileHandler)
	}
}

func pageHandler(c *gin.Context) {
	var param message.MessagePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	message.Page(c, userID, &param)
}

func detailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := message.Detail(id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}

func unreadCountHandler(c *gin.Context) {
	userID := auth.GetLoginID(c)
	count := message.UnreadCount(userID)
	result.Success(c, message.UnreadCountVO{Count: count})
}
  func sendHandler(c *gin.Context) {
	var param message.MessageSendParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	convIDs := message.Send(c, &param, userID, string(enums.LoginTypeBusiness))
	data := gin.H{}
	if len(convIDs) > 0 {
		data["conversation_id"] = convIDs[0]
	}
	result.Success(c, data)
}

func recallHandler(c *gin.Context) {
	var param message.RecallParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	message.Recall(userID, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

func forwardHandler(c *gin.Context) {
	var param message.ForwardParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	message.Forward(c, userID, string(enums.LoginTypeBusiness), &param)
	result.Success(c, nil)
}

func deleteHandler(c *gin.Context) {
	var param struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	message.Remove(userID, param.IDs)
	result.Success(c, nil)
}

func searchHandler(c *gin.Context) {
	var param message.SearchParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	list, hasMore := message.Search(c, userID, &param)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

func markReadHandler(c *gin.Context) {
	var param utils.IdParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MarkRead(param.ID)
	result.Success(c, nil)
}

func markAllReadHandler(c *gin.Context) {
	userID := auth.GetLoginID(c)
	message.MarkAllRead(userID)
	result.Success(c, nil)
}

func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	message.Remove(userID, param.IDs)
	result.Success(c, nil)
}

func conversationsHandler(c *gin.Context) {
	userID := auth.GetLoginID(c)
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list, hasMore := message.Conversations(userID, string(enums.LoginTypeBusiness), cursor, size)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

func conversationMessagesHandler(c *gin.Context) {
	userID := auth.GetLoginID(c)
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
		messages, hasMore = message.BusinessConversationMessages(c.Request.Context(), userID, cid, cursor, size)
	}
	result.Success(c, gin.H{
		"records":  messages,
		"has_more": hasMore,
	})
}

func conversationReadHandler(c *gin.Context) {
	var param struct {
		ConversationID string `json:"conversation_id"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	if len(param.ConversationID) > 6 && param.ConversationID[:6] == "group:" {
		group.MarkConversationRead(c.Request.Context(), param.ConversationID[6:], userID, string(enums.LoginTypeBusiness))
	} else {
		message.MarkConversationRead(userID, param.ConversationID)
	}
	result.Success(c, nil)
}

func getOrCreateConversationHandler(c *gin.Context) {
	var param message.GetOrCreateConversationParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.GetLoginID(c)
	cid, displayName := message.GetOrCreateConversation(userID, string(enums.LoginTypeBusiness), &param)
	result.Success(c, gin.H{"conversation_id": cid, "display_name": displayName})
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}

// ==================== Consumer (C端) Routes ====================

func RegisterClientRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/c/im").Use(authMW.HeiClientCheckLogin())
	{
		g.GET("/message/page", clientPageHandler)
		g.GET("/message/detail", clientDetailHandler)
		g.GET("/message/unread-count", clientUnreadCountHandler)
		g.POST("/message/send", middleware.RateLimiter("c_send", 5, 20), clientSendHandler)
		g.POST("/message/recall", clientRecallHandler)
		g.POST("/message/forward", clientForwardHandler)
		g.POST("/message/delete", clientDeleteHandler)
		g.GET("/message/search", clientSearchHandler)
		g.POST("/message/mark-read", clientMarkReadHandler)
		g.POST("/message/mark-all-read", clientMarkAllReadHandler)
		g.POST("/message/remove", clientRemoveHandler)

		g.GET("/conversation/list", clientConversationsHandler)
		g.GET("/conversation/messages", clientConversationMessagesHandler)
		g.POST("/conversation/read", clientConversationReadHandler)
		g.POST("/conversation/get-or-create", clientGetOrCreateConversationHandler)
		g.POST("/file/upload", clientUploadFileHandler)
	}
}

// Client auth helpers
func clientUserID(c *gin.Context) (string, string) {
	return auth.Consumer.GetLoginID(c), string(enums.LoginTypeConsumer)
}

func clientPageHandler(c *gin.Context) {
	var param message.MessagePageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.Consumer.GetLoginID(c)
	message.Page(c, userID, &param)
}

func clientDetailHandler(c *gin.Context) {
	id := c.Query("id")
	vo := message.Detail(id)
	if vo == nil {
		result.Success(c, nil)
		return
	}
	result.Success(c, vo)
}

func clientUnreadCountHandler(c *gin.Context) {
	userID := auth.Consumer.GetLoginID(c)
	count := message.UnreadCount(userID)
	result.Success(c, message.UnreadCountVO{Count: count})
}

func clientSendHandler(c *gin.Context) {
	var param message.MessageSendParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID, userType := clientUserID(c)
	convIDs := message.Send(c, &param, userID, userType)
	data := gin.H{}
	if len(convIDs) > 0 {
		data["conversation_id"] = convIDs[0]
	}
	result.Success(c, data)
}

func clientRecallHandler(c *gin.Context) {
	var param message.RecallParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID, userType := clientUserID(c)
	message.Recall(userID, userType, &param)
	result.Success(c, nil)
}

func clientForwardHandler(c *gin.Context) {
	var param message.ForwardParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID, userType := clientUserID(c)
	message.Forward(c, userID, userType, &param)
	result.Success(c, nil)
}

func clientDeleteHandler(c *gin.Context) {
	var param struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.Consumer.GetLoginID(c)
	message.Remove(userID, param.IDs)
	result.Success(c, nil)
}

func clientSearchHandler(c *gin.Context) {
	var param message.SearchParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.Consumer.GetLoginID(c)
	list, hasMore := message.Search(c, userID, &param)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

func clientMarkReadHandler(c *gin.Context) {
	var param utils.IdParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	message.MarkRead(param.ID)
	result.Success(c, nil)
}

func clientMarkAllReadHandler(c *gin.Context) {
	userID := auth.Consumer.GetLoginID(c)
	message.MarkAllRead(userID)
	result.Success(c, nil)
}

func clientRemoveHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.Consumer.GetLoginID(c)
	message.Remove(userID, param.IDs)
	result.Success(c, nil)
}

func clientConversationsHandler(c *gin.Context) {
	userID := auth.Consumer.GetLoginID(c)
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list, hasMore := message.Conversations(userID, string(enums.LoginTypeConsumer), cursor, size)
	result.Success(c, gin.H{"records": list, "has_more": hasMore})
}

func clientConversationMessagesHandler(c *gin.Context) {
	userID := auth.Consumer.GetLoginID(c)
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
		messages, hasMore = message.ConsumerConversationMessages(c.Request.Context(), userID, cid, cursor, size)
	}
	result.Success(c, gin.H{
		"records":  messages,
		"has_more": hasMore,
	})
}

func clientConversationReadHandler(c *gin.Context) {
	var param struct {
		ConversationID string `json:"conversation_id"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID := auth.Consumer.GetLoginID(c)
	if len(param.ConversationID) > 6 && param.ConversationID[:6] == "group:" {
		group.MarkConversationRead(c.Request.Context(), param.ConversationID[6:], userID, string(enums.LoginTypeConsumer))
	} else {
		message.MarkConversationRead(userID, param.ConversationID)
	}
	result.Success(c, nil)
}

func clientGetOrCreateConversationHandler(c *gin.Context) {
	var param message.GetOrCreateConversationParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	userID, userType := clientUserID(c)
	cid, displayName := message.GetOrCreateConversation(userID, userType, &param)
	result.Success(c, gin.H{"conversation_id": cid, "display_name": displayName})
}

func uploadFileHandler(c *gin.Context) {
	userID := auth.GetLoginID(c)
	data, err := message.UploadFile(c, userID, string(enums.LoginTypeBusiness))
	if err != nil {
		if appErr, ok := err.(*message.AppError); ok {
			result.Failure(c, appErr.Message, appErr.Code)
		} else {
			result.Failure(c, err.Error(), 400)
		}
		return
	}
	result.Success(c, data)
}

func clientUploadFileHandler(c *gin.Context) {
	userID := auth.Consumer.GetLoginID(c)
	data, err := message.UploadFile(c, userID, string(enums.LoginTypeConsumer))
	if err != nil {
		if appErr, ok := err.(*message.AppError); ok {
			result.Failure(c, appErr.Message, appErr.Code)
		} else {
			result.Failure(c, err.Error(), 400)
		}
		return
	}
	result.Success(c, data)
}
