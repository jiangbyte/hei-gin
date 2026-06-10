package message

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
	"hei-gin/plugins/plugin-im/ws"
	imModel "hei-gin/plugins/plugin-im/model"

	"github.com/gin-gonic/gin"
)


func getLoginID(c *gin.Context) string {
	path := c.Request.URL.Path
	if len(path) > 8 && path[:8] == "/api/v1/c" {
		return auth.Consumer.GetLoginID(c)
	}
	return auth.GetLoginID(c)
}

func getUserType(c *gin.Context) string {
	path := c.Request.URL.Path
	if len(path) > 8 && path[:8] == "/api/v1/c" {
		return string(enums.LoginTypeConsumer)
	}
	return string(enums.LoginTypeBusiness)
}

// ==================== MessageSend ====================

func MessageSend(c *gin.Context, param *MessageSendParam) {
	ctx := c.Request.Context()

	senderID := getLoginID(c)
	senderType := getUserType(c)

	if !ws.GlobalCrossHub.AllowMessage(senderID, enums.LoginTypeEnum(senderType)) {
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

	records := make([]imModel.Message, len(param.ReceiverIDs))
	for i, rid := range param.ReceiverIDs {
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
	if err := db.DB.WithContext(ctx).Create(&records).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("发送消息失败: "+err.Error(), 500))
		return
	}

	for _, rec := range records {
		conv := imModel.Conversation{
			ID:       rec.ConversationID,
			FromID:   rec.SenderID,
			FromType: rec.SenderType,
			ToID:     rec.ReceiverID,
			ToType:   rec.ReceiverType,
			LastMsg:  rec.Content,
		}
		db.DB.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"last_msg", "last_time", "updated_at",
			}),
		}).Create(&conv)
	}

	for i, rid := range param.ReceiverIDs {
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
		if receiverType == string(enums.LoginTypeConsumer) {
			ws.GlobalCrossHub.SendToConsumer(rid, msg, records[i].ID)
		} else {
			ws.GlobalCrossHub.SendToUser(rid, msg, records[i].ID)
		}
	}
}

// ==================== MessagePage ====================

func MessagePage(c *gin.Context, param *MessagePageParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c)

	if param.Current < 1 {
		param.Current = 1
	}
	if param.Size < 1 {
		param.Size = 10
	}
	if param.Size > 100 {
		param.Size = 100
	}

	query := db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("(sender_id = ? OR receiver_id = ?) AND (deleted_by != ? OR deleted_by IS NULL)", userID, userID, userID)
	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	var total int64
	query.Count(&total)

	var records []imModel.Message
	query.Order("created_at DESC").Limit(param.Size).Offset((param.Current - 1) * param.Size).Find(&records)

	vos := make([]MessageVO, len(records))
	for i, e := range records {
		vos[i] = *ImMessageToMessageVO(&e)
	}
	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

// ==================== MessageUnreadCount ====================

func MessageUnreadCount(c *gin.Context) {
	userID := getLoginID(c)

	var count int64
	db.DB.Model(&imModel.Message{}).Where("receiver_id = ? AND status = ?", userID, "unread").Count(&count)
	result.Success(c, UnreadCountVO{Count: count})
}

// ==================== MessageDetail ====================

func MessageDetail(c *gin.Context) {
	id := c.Query("id")

	var entity imModel.Message
	if err := db.DB.First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Success(c, nil)
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询消息失败: "+err.Error(), 500))
		return
	}
	result.Success(c, ImMessageToMessageVO(&entity))
}

// ==================== MessageMarkRead ====================

func MessageMarkRead(c *gin.Context, param *utils.IdParam) {
	if err := db.DB.Model(&imModel.Message{}).Where("id = ?", param.ID).
		Updates(map[string]interface{}{"status": "read"}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("标记已读失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageMarkConversationRead ====================

func MessageMarkConversationRead(c *gin.Context) {
	receiverID := getLoginID(c)

	var param ConversationReadParam
	if err := c.ShouldBindJSON(&param); err != nil {
		return
	}

	if err := db.DB.Model(&imModel.Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND status = ?", param.ConversationID, receiverID, "unread").
		Updates(map[string]interface{}{"status": "read"}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("标记已读失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageMarkAllRead ====================

func MessageMarkAllRead(c *gin.Context) {
	receiverID := getLoginID(c)

	if err := db.DB.Model(&imModel.Message{}).
		Where("receiver_id = ? AND status = ?", receiverID, "unread").
		Updates(map[string]interface{}{"status": "read"}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("标记全部已读失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageRemove (soft-delete) ====================

func MessageRemove(c *gin.Context, param *DeleteParam) {
	userID := getLoginID(c)

	if len(param.IDs) == 0 {
		return
	}
	if err := db.DB.Model(&imModel.Message{}).
		Where("id IN ? AND (sender_id = ? OR receiver_id = ?)", param.IDs, userID, userID).
		Update("deleted_by", userID).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除消息失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageRecall (within 5 min) ====================

func MessageRecall(c *gin.Context, param *RecallParam) {
	userID := getLoginID(c)
	userType := getUserType(c)

	var msg imModel.Message
	if err := db.DB.First(&msg, "id = ?", param.MessageID).Error; err != nil {
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
	if err := db.DB.Model(&imModel.Message{}).Where("id = ?", param.MessageID).
		Updates(map[string]interface{}{
			"content":  "消息已被撤回",
			"msg_type": imModel.MsgTypeSystem,
		}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("撤回失败: "+err.Error(), 500))
		return
	}
}

// ==================== MessageForward ====================

func MessageForward(c *gin.Context, param *ForwardParam) {
	var original imModel.Message
	if err := db.DB.First(&original, "id = ?", param.MessageID).Error; err != nil {
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

	if param.Size < 1 {
		param.Size = 20
	}
	if param.Size > 100 {
		param.Size = 100
	}

	query := db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("(sender_id = ? OR receiver_id = ?) AND content LIKE ?", userID, userID, "%"+param.Keyword+"%")
	if param.Cursor != "" {
		if t, err := utils.ParseDateTime(param.Cursor); err == nil {
			query = query.Where("created_at < ?", t)
		}
	}
	var records []imModel.Message
	query.Order("created_at DESC").Limit(param.Size + 1).Find(&records)

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
