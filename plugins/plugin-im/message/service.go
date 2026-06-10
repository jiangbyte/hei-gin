package message

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/plugins/plugin-im/ws"
	imModel "hei-gin/plugins/plugin-im/model"

	"github.com/gin-gonic/gin"
)

func toVO(e *imModel.Message) *MessageVO {
	v := &MessageVO{
		ID: e.ID, ConversationID: e.ConversationID,
		Content: e.Content, MsgType: e.MsgType, Extra: e.Extra,
		SenderID: e.SenderID, SenderType: e.SenderType,
		ReceiverID: e.ReceiverID, ReceiverType: e.ReceiverType,
		Status: e.Status,
		CreatedAt: utils.FormatDateTimePtr(e.CreatedAt),
		UpdatedAt: utils.FormatDateTimePtr(e.UpdatedAt),
	}
	if e.ReadAt != nil {
		s := utils.FormatDateTime(*e.ReadAt)
		v.ReadAt = &s
	}
	return v
}

func toVOList(records []imModel.Message) []MessageVO {
	r := make([]MessageVO, len(records))
	for i, e := range records {
		r[i] = *toVO(&e)
	}
	return r
}


// ==================== Send ====================

func Send(c *gin.Context, param *MessageSendParam, senderID string, senderType string) []string {

	ctx := c.Request.Context()
	now := time.Now()

	if !ws.GlobalCrossHub.AllowMessage(senderID, enums.LoginTypeEnum(senderType)) {
		panic(exception.NewBusinessError("发送消息过于频繁，请稍后重试", 429))
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
			CreatedAt:      &now,
			UpdatedAt:      &now,
		}
	}
	if err := db.DB.WithContext(ctx).Create(&records).Error; err != nil {
		panic(exception.NewBusinessError("发送消息失败: "+err.Error(), 500))
	}

	// Update im_conversation cache for each receiver
	for _, rec := range records {
		conv := imModel.Conversation{
			ID:        rec.ConversationID,
			FromID:   rec.SenderID,
			FromType: rec.SenderType,
			ToID:   rec.ReceiverID,
			ToType: rec.ReceiverType,
			LastMsg:   rec.Content,
			LastTime:  &now,
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		db.DB.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"last_msg", "last_time", "updated_at",
			}),
		}).Create(&conv)
	}

	// Send WS notifications
	for i, rid := range param.ReceiverIDs {
		cid := records[i].ConversationID
		msg := ws.Message{
			Type: ws.MsgNewMessage,
			Payload: ws.NewMessagePayload{
				MessageID:      records[i].ID,
				ConversationID: cid,
				Content:        param.Content,
				MsgType:        msgType,
				Extra:          param.Extra,
				SenderID:       senderID,
				SenderType:     senderType,
				CreatedAt:      utils.FormatDateTime(now),
			},
		}
		if receiverType == string(enums.LoginTypeConsumer) {
			ws.GlobalCrossHub.SendToConsumer(rid, msg, records[i].ID)
		} else {
			ws.GlobalCrossHub.SendToUser(rid, msg, records[i].ID)
		}
	}
	convIDs := make([]string, len(records))
	for i := range records {
		convIDs[i] = records[i].ConversationID
	}
	return convIDs
}

// ==================== Page ====================

func Page(c *gin.Context, userID string, param *MessagePageParam) {
	ctx := c.Request.Context()
	if param.Current < 1 {
		param.Current = 1
	}
	if param.Size < 1 {
		param.Size = 10
	}
	if param.Size > 100 {
		param.Size = 100
	}

	query := db.DB.WithContext(ctx).Model(&imModel.Message{}).Where("(sender_id = ? OR receiver_id = ?) AND (deleted_by != ? OR deleted_by IS NULL)", userID, userID, userID)
	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	var total int64
	query.Count(&total)

	var records []imModel.Message
	query.Order("created_at DESC").Limit(param.Size).Offset((param.Current - 1) * param.Size).Find(&records)
	result.PageDataResult(c, toVOList(records), total, param.Current, param.Size)
}

// ==================== UnreadCount ====================

func UnreadCount(userID string) int64 {
	var count int64
	db.DB.Model(&imModel.Message{}).Where("receiver_id = ? AND status = ?", userID, "unread").Count(&count)
	return count
}

// ==================== Detail ====================

func Detail(id string) *MessageVO {
	var entity imModel.Message
	if err := db.DB.First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(exception.NewBusinessError("查询消息失败: "+err.Error(), 500))
	}
	return toVO(&entity)
}

// ==================== MarkRead ====================

func MarkRead(id string) {
	now := time.Now()
	if err := db.DB.Model(&imModel.Message{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "read", "read_at": &now}).Error; err != nil {
		panic(exception.NewBusinessError("标记已读失败: "+err.Error(), 500))
	}
}

func MarkConversationRead(receiverID string, conversationID string) {
	now := time.Now()
	if err := db.DB.Model(&imModel.Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND status = ?", conversationID, receiverID, "unread").
		Updates(map[string]interface{}{"status": "read", "read_at": &now}).Error; err != nil {
		panic(exception.NewBusinessError("标记已读失败: "+err.Error(), 500))
	}
}

func MarkAllRead(receiverID string) {
	now := time.Now()
	if err := db.DB.Model(&imModel.Message{}).
		Where("receiver_id = ? AND status = ?", receiverID, "unread").
		Updates(map[string]interface{}{"status": "read", "read_at": &now}).Error; err != nil {
		panic(exception.NewBusinessError("标记全部已读失败: "+err.Error(), 500))
	}
}

// ==================== Remove (soft-delete) ====================

func Remove(userID string, ids []string) {
	if len(ids) == 0 {
		return
	}
	if err := db.DB.Model(&imModel.Message{}).
		Where("id IN ? AND (sender_id = ? OR receiver_id = ?)", ids, userID, userID).
		Update("deleted_by", userID).Error; err != nil {
		panic(exception.NewBusinessError("删除消息失败: "+err.Error(), 500))
	}
}

// ==================== Recall (within 5 min) ====================

func Recall(userID string, userType string, param *RecallParam) {
	var msg imModel.Message
	if err := db.DB.First(&msg, "id = ?", param.MessageID).Error; err != nil {
		panic(exception.NewBusinessError("消息不存在", 400))
	}
	if msg.SenderID != userID || msg.SenderType != userType {
		panic(exception.NewBusinessError("只能撤回自己的消息", 403))
	}
	if msg.CreatedAt != nil && time.Since(*msg.CreatedAt) > 5*time.Minute {
		panic(exception.NewBusinessError("超过5分钟，无法撤回", 400))
	}
	now := time.Now()
	if err := db.DB.Model(&imModel.Message{}).Where("id = ?", param.MessageID).
		Updates(map[string]interface{}{
			"content":    "消息已被撤回",
			"msg_type":   imModel.MsgTypeSystem,
			"updated_at": &now,
		}).Error; err != nil {
		panic(exception.NewBusinessError("撤回失败: "+err.Error(), 500))
	}
}

// ==================== Forward ====================

func Forward(c *gin.Context, userID string, userType string, param *ForwardParam) {
	var original imModel.Message
	if err := db.DB.First(&original, "id = ?", param.MessageID).Error; err != nil {
		panic(exception.NewBusinessError("消息不存在", 400))
	}

	sendParam := &MessageSendParam{
		Content:      original.Content,
		MsgType:      original.MsgType,
		Extra:        original.Extra,
		ReceiverIDs:  param.TargetIDs,
		ReceiverType: param.TargetType,
	}
	Send(c, sendParam, userID, userType)
}

// ==================== Search ====================

func Search(c *gin.Context, userID string, param *SearchParam) ([]MessageVO, bool) {
	ctx := c.Request.Context()
	if param.Size < 1 {
		param.Size = 20
	}
	if param.Size > 100 {
		param.Size = 100
	}

	query := db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("(sender_id = ? OR receiver_id = ?) AND content LIKE ?", userID, userID, "%"+param.Keyword+"%")
	if param.Cursor != "" {
		if t, err := utils.ParseDateTimeLocal(param.Cursor); err == nil {
			query = query.Where("created_at < ?", t)
		}
	}
	var records []imModel.Message
	query.Order("created_at DESC").Limit(param.Size + 1).Find(&records)

	hasMore := len(records) > param.Size
	if hasMore {
		records = records[:param.Size]
	}
	return toVOList(records), hasMore
}
