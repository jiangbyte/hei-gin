package message

import (
	"strings"
	"time"

	"gorm.io/gorm"

	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/plugins/plugin-im/ws"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

func getLoginID(c *gin.Context) string {
	if v, ok := c.Get("login_id"); ok {
		if uid, ok := v.(string); ok && uid != "" {
			return uid
		}
	}
	if getUserType(c) == string(enums.LoginTypeConsumer) {
		return auth.Consumer.GetLoginID(c)
	}
	return auth.GetLoginID(c)
}

func getUserType(c *gin.Context) string {
	if v, ok := c.Get("login_type"); ok {
		if loginType, ok := v.(string); ok && loginType != "" {
			return loginType
		}
	}
	path := c.Request.URL.Path
	if isConsumerAPIPath(path) {
		return string(enums.LoginTypeConsumer)
	}
	return string(enums.LoginTypeBusiness)
}

func isConsumerAPIPath(path string) bool {
	if !strings.HasPrefix(path, "/api/v") {
		return false
	}
	afterVersionPrefix := path[len("/api/v"):]
	slash := strings.IndexByte(afterVersionPrefix, '/')
	if slash < 0 {
		return false
	}
	return strings.HasPrefix(afterVersionPrefix[slash+1:], "c/")
}

// ==================== MessageSend ====================

func MessageSend(c *gin.Context, param *MessageSendParam) {
	ctx := c.Request.Context()

	senderID := getLoginID(c)
	senderType := getUserType(c)

	if ws.GlobalCrossHub != nil && !ws.GlobalCrossHub.AllowMessage(senderID, enums.LoginTypeEnum(senderType)) {
		result.WriteError(c, exception.NewBusinessError("发送消息过于频繁，请稍后重试", 429))
		return
	}

	msgType := param.MsgType
	if msgType == "" {
		msgType = "TEXT"
	}
	receiverType := param.ReceiverType
	if receiverType == "" {
		receiverType = string(enums.LoginTypeBusiness)
	}
	receiverIDs := normalizeReceiverIDs(param.ReceiverIDs)
	if len(receiverIDs) == 0 {
		result.WriteError(c, exception.NewBusinessError("接收人不能为空", 400))
		return
	}
	if len(receiverIDs) > 200 {
		result.WriteError(c, exception.NewBusinessError("单次最多发送给200个接收人", 400))
		return
	}
	if len(param.Content) > 5000 {
		result.WriteError(c, exception.NewBusinessError("消息内容不能超过5000个字符", 400))
		return
	}

	records := make([]imModel.Message, len(receiverIDs))
	for i, rid := range receiverIDs {
		cid := imModel.GenerateConversationID(senderID, enums.LoginTypeEnum(senderType), rid, enums.LoginTypeEnum(receiverType))
		records[i] = imModel.Message{
			ID:             utils.GenerateID(),
			ConversationID: cid,
			Content:        param.Content,
			Extra:          param.Extra,
			MsgType:        msgType,
			SenderID:       senderID,
			SenderType:     senderType,
			ReceiverID:     rid,
			ReceiverType:   receiverType,
			Status:         "unread",
		}
	}
	if err := CreateMessages(ctx, records); err != nil {
		result.WriteError(c, exception.NewBusinessError("发送消息失败: "+err.Error(), 500))
		return
	}

	conversations := make([]imModel.Conversation, 0, len(records))
	for _, rec := range records {
		lastTime := rec.CreatedAt
		if lastTime == nil {
			now := time.Now()
			lastTime = &now
		}
		conversations = append(conversations, imModel.Conversation{
			ID:       rec.ConversationID,
			FromID:   rec.SenderID,
			FromType: rec.SenderType,
			ToID:     rec.ReceiverID,
			ToType:   rec.ReceiverType,
			LastMsg:  rec.Content,
			LastTime: lastTime,
		})
	}
	_ = UpsertConversations(ctx, conversations)

	pushMessages := make(map[string]ws.Message, len(receiverIDs))
	messageIDs := make(map[string]string, len(receiverIDs))
	for i, rid := range receiverIDs {
		msg := ws.Message{
			Type: ws.MsgNewMessage,
			Payload: ws.NewMessagePayload{
				MessageID:      records[i].ID,
				ConversationID: records[i].ConversationID,
				Content:        param.Content,
				MsgType:        msgType,
				Extra:          param.Extra,
				SenderID:       senderID,
				SenderType:     senderType,
			},
		}
		pushMessages[rid] = msg
		messageIDs[rid] = records[i].ID
	}
	if receiverType == string(enums.LoginTypeConsumer) {
		ws.GlobalCrossHub.SendMessagesToConsumers(pushMessages, messageIDs)
	} else {
		ws.GlobalCrossHub.SendMessagesToUsers(pushMessages, messageIDs)
	}
}

func normalizeReceiverIDs(ids []string) []string {
	result := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

// ==================== MessagePage ====================

func MessagePage(c *gin.Context, param *MessagePageParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if param.Current < 1 {
		param.Current = 1
	}
	if param.Size < 1 {
		param.Size = 10
	}
	if param.Size > 100 {
		param.Size = 100
	}

	records, total := Page(ctx, userID, userType, param.Status, param.Current, param.Size)

	vos := make([]MessageVO, len(records))
	for i, e := range records {
		vos[i] = *ImMessageToMessageVO(&e)
	}
	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

// ==================== MessageUnreadCount ====================

func MessageUnreadCount(c *gin.Context) {
	userID := getLoginID(c)
	userType := getUserType(c)

	count := CountUnread(c.Request.Context(), userID, userType)
	result.Success(c, UnreadCountVO{Count: count})
}

// ==================== MessageDetail ====================

func MessageDetail(c *gin.Context) {
	id := c.Query("id")
	userID := getLoginID(c)
	userType := getUserType(c)

	entity, err := FindOwnedByID(c.Request.Context(), id, userID, userType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Success(c, nil)
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询消息失败: "+err.Error(), 500))
		return
	}
	result.Success(c, ImMessageToMessageVO(entity))
}

// ==================== MessageMarkRead ====================

func MessageMarkRead(c *gin.Context, param *utils.IdParam) {
	userID := getLoginID(c)
	userType := getUserType(c)
	if err := MarkRead(c.Request.Context(), param.ID, userID, userType); err != nil {
		result.WriteError(c, exception.NewBusinessError("标记已读失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageMarkConversationRead ====================

func MessageMarkConversationRead(c *gin.Context, param *ConversationReadParam) {
	receiverID := getLoginID(c)
	receiverType := getUserType(c)

	if param == nil || param.ConversationID == "" {
		return
	}

	if err := MarkConversationRead(c.Request.Context(), param.ConversationID, receiverID, receiverType); err != nil {
		result.WriteError(c, exception.NewBusinessError("标记已读失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageMarkAllRead ====================

func MessageMarkAllRead(c *gin.Context) {
	receiverID := getLoginID(c)
	receiverType := getUserType(c)

	if err := MarkAllRead(c.Request.Context(), receiverID, receiverType); err != nil {
		result.WriteError(c, exception.NewBusinessError("标记全部已读失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageRemove (soft-delete) ====================

func MessageRemove(c *gin.Context, param *DeleteParam) {
	userID := getLoginID(c)
	userType := getUserType(c)

	if len(param.IDs) == 0 {
		return
	}
	if err := SoftDelete(c.Request.Context(), param.IDs, userID, userType); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除消息失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageRecall (within 5 min) ====================

func MessageRecall(c *gin.Context, param *RecallParam) {
	userID := getLoginID(c)
	userType := getUserType(c)

	msg, err := FindByID(c.Request.Context(), param.MessageID)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("消息不存在", 400))
		return
	}
	if msg.SenderID != userID || msg.SenderType != userType {
		result.WriteError(c, exception.NewBusinessError("只能撤回自己的消息", 403))
		return
	}
	if msg.CreatedAt != nil && time.Since(*msg.CreatedAt) > 5*time.Minute {
		result.WriteError(c, exception.NewBusinessError("超过5分钟，无法撤回", 400))
		return
	}
	if err := Recall(c.Request.Context(), param.MessageID); err != nil {
		result.WriteError(c, exception.NewBusinessError("撤回失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageForward ====================

func MessageForward(c *gin.Context, param *ForwardParam) {
	userID := getLoginID(c)
	userType := getUserType(c)
	original, err := FindOwnedByID(c.Request.Context(), param.MessageID, userID, userType)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("消息不存在", 400))
		return
	}

	sendParam := &MessageSendParam{
		Content:      original.Content,
		MsgType:      original.MsgType,
		Extra:        original.Extra,
		ReceiverIDs:  param.TargetIDs,
		ReceiverType: param.TargetType,
	}
	MessageSend(c, sendParam)
}

// ==================== MessageSearch ====================

func MessageSearch(c *gin.Context, param *SearchParam) ([]MessageVO, bool) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if param.Size < 1 {
		param.Size = 20
	}
	if param.Size > 100 {
		param.Size = 100
	}

	records := Search(ctx, userID, userType, param.Keyword, param.Cursor, param.Size)

	hasMore := len(records) > param.Size
	if hasMore {
		records = records[:param.Size]
	}
	vos := make([]MessageVO, len(records))
	for i, e := range records {
		vos[i] = *ImMessageToMessageVO(&e)
	}
	return vos, hasMore
}
