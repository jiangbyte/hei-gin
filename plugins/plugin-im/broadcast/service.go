package broadcast

import (
	"time"

	"gorm.io/gorm"
	imModel "hei-gin/plugins/plugin-im/model"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"
	"hei-gin/plugins/plugin-im/ws"
)

// ==================== Send ====================

func Send(senderID string, p *SendBroadcastParam) {
	scope := p.Scope
	if scope == "" {
		scope = "ALL"
	}

	now := time.Now()
	if err := db.DB.Create(&imModel.Broadcast{
		ID:        utils.GenerateID(),
		Title:     p.Title,
		Content:   p.Content,
		Scope:     scope,
		SenderID:  senderID,
		CreatedAt: &now,
		UpdatedAt: &now,
	}).Error; err != nil {
		panic(exception.NewBusinessError("发送通知失败", 500))
	}

	// WS broadcast
	payload := map[string]interface{}{
		"title":   p.Title,
		"content": p.Content,
		"scope":   scope,
		"action":  "broadcast",
	}
	msg := ws.Message{Type: "broadcast", Payload: payload}
	switch scope {
	case "ALL":
		ws.GlobalCrossHub.BroadcastAll(msg)
	case "BUSINESS":
		ws.GlobalCrossHub.BroadcastBusiness(msg)
	case "CONSUMER":
		ws.GlobalCrossHub.BroadcastConsumers(msg)
	}
}

// ==================== List (admin) ====================

func List(cursor string, size int) ([]BroadcastVO, bool) {
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	q := db.DB.Model(&imModel.Broadcast{})
	if cursor != "" {
		if t, err := utils.ParseDateTime(cursor); err == nil {
			q = q.Where("created_at < ?", t)
		}
	}
	var records []imModel.Broadcast
	q.Order("created_at DESC").Limit(size + 1).Find(&records)

	hasMore := len(records) > size
	if hasMore {
		records = records[:size]
	}

	result := make([]BroadcastVO, len(records))
	for i, b := range records {
		result[i] = BroadcastVO{
			ID: b.ID, Title: b.Title, Content: b.Content,
			Scope: b.Scope, SenderID: b.SenderID,
			CreatedAt: utils.FormatDateTimePtr(b.CreatedAt),
		}
	}
	return result, hasMore
}

// ==================== Unread List ====================

func UnreadList(userID, userType string) ([]BroadcastVO, bool) {
	var records []imModel.Broadcast
	db.DB.Model(&imModel.Broadcast{}).Order("created_at DESC").Limit(50).Find(&records)

	// Check read status
	var readRecords []imModel.BroadcastRead
	db.DB.Model(&imModel.BroadcastRead{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&readRecords)
	readMap := make(map[string]*time.Time)
	for _, r := range readRecords {
		readMap[r.BroadcastID] = r.ReadAt
	}

	result := make([]BroadcastVO, 0, len(records))
	for _, b := range records {
		readAt, read := readMap[b.ID]
		vo := BroadcastVO{
			ID: b.ID, Title: b.Title,
			Content: b.Content,
			Scope: b.Scope, Read: read,
			CreatedAt: utils.FormatDateTimePtr(b.CreatedAt),
		}
		if read && readAt != nil {
			s := utils.FormatDateTime(*readAt)
			vo.ReadAt = s
		}
		result = append(result, vo)
	}
	return result, false
}

// ==================== Mark Read ====================

func MarkRead(userID, userType, broadcastID string) {
	now := time.Now()
	_ = db.DB.Where("broadcast_id = ? AND user_id = ? AND user_type = ?", broadcastID, userID, userType).
		FirstOrCreate(&imModel.BroadcastRead{
			BroadcastID: broadcastID,
			ID: utils.GenerateID(),
			UserID:      userID,
			UserType:    userType,
			ReadAt:      &now,
		})
}

// ==================== Detail ====================

func Detail(id string) *BroadcastVO {
	var b imModel.Broadcast
	if err := db.DB.First(&b, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(exception.NewBusinessError("查询通知失败", 500))
	}
	return &BroadcastVO{
		ID: b.ID, Title: b.Title, Content: b.Content,
		Scope: b.Scope, SenderID: b.SenderID,
		CreatedAt: utils.FormatDateTimePtr(b.CreatedAt),
	}
}
