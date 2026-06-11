package broadcast

import (
	"gorm.io/gorm"
	imModel "hei-gin/plugins/plugin-im/model"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
	"hei-gin/plugins/plugin-im/ws"

	"github.com/gin-gonic/gin"
)

// ==================== BroadcastSend ====================

func BroadcastSend(c *gin.Context, p *SendBroadcastParam) {
	ctx := c.Request.Context()
	senderID := auth.GetLoginID(c)

	if p.Scope == "" {
		p.Scope = "ALL"
	}

	e := imModel.Broadcast{
		ID:       utils.GenerateID(),
		Title:    p.Title,
		Content:  p.Content,
		Scope:    p.Scope,
		SenderID: senderID,
	}
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("发送通知失败: "+err.Error(), 500))
		return
	}

	// WS broadcast
	payload := map[string]interface{}{
		"title":   p.Title,
		"content": p.Content,
		"scope":   p.Scope,
		"action":  "broadcast",
	}
	msg := ws.Message{Type: "broadcast", Payload: payload}
	switch p.Scope {
	case "ALL":
		ws.GlobalCrossHub.BroadcastAll(msg)
	case string(enums.LoginTypeBusiness):
		ws.GlobalCrossHub.BroadcastBusiness(msg)
	case string(enums.LoginTypeConsumer):
		ws.GlobalCrossHub.BroadcastConsumers(msg)
	}
}

// ==================== BroadcastPage (admin) ====================

func BroadcastPage(c *gin.Context, cursor string, size int) ([]BroadcastVO, bool) {
	ctx := c.Request.Context()
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	q := db.DB.WithContext(ctx).Model(&imModel.Broadcast{})
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

	vos := make([]BroadcastVO, len(records))
	for i, b := range records {
		vos[i] = *ImBroadcastToBroadcastVO(&b)
	}
	return vos, hasMore
}

// ==================== BroadcastUnreadList ====================

func BroadcastUnreadList(c *gin.Context, userType string) []BroadcastVO {
	ctx := c.Request.Context()
	var userID string
	if userType == string(enums.LoginTypeConsumer) {
		userID = auth.Consumer.GetLoginID(c)
	} else {
		userID = auth.GetLoginID(c)
	}

	var records []imModel.Broadcast
	db.DB.WithContext(ctx).Model(&imModel.Broadcast{}).Order("created_at DESC").Limit(50).Find(&records)

	var readRecords []imModel.BroadcastRead
	db.DB.WithContext(ctx).Model(&imModel.BroadcastRead{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&readRecords)
	readMap := make(map[string]bool)
	for _, r := range readRecords {
		readMap[r.BroadcastID] = true
	}

	vos := make([]BroadcastVO, 0, len(records))
	for _, b := range records {
		vo := *ImBroadcastToBroadcastVO(&b)
		if _, ok := readMap[b.ID]; ok {
			vo.Read = true
		}
		vos = append(vos, vo)
	}
	return vos
}

// ==================== BroadcastMarkRead ====================

func BroadcastMarkRead(c *gin.Context, p *ReadParam) {
	ctx := c.Request.Context()
	userID := auth.GetLoginID(c)

	_ = db.DB.WithContext(ctx).Where("broadcast_id = ? AND user_id = ? AND user_type = ?", p.BroadcastID, userID, string(enums.LoginTypeBusiness)).
		FirstOrCreate(&imModel.BroadcastRead{
			BroadcastID: p.BroadcastID,
			ID:          utils.GenerateID(),
			UserID:      userID,
			UserType:    string(enums.LoginTypeBusiness),
		})
}

// ==================== BroadcastMarkReadConsumer ====================

func BroadcastMarkReadConsumer(c *gin.Context, p *ReadParam) {
	ctx := c.Request.Context()
	userID := auth.Consumer.GetLoginID(c)

	_ = db.DB.WithContext(ctx).Where("broadcast_id = ? AND user_id = ? AND user_type = ?", p.BroadcastID, userID, string(enums.LoginTypeConsumer)).
		FirstOrCreate(&imModel.BroadcastRead{
			BroadcastID: p.BroadcastID,
			ID:          utils.GenerateID(),
			UserID:      userID,
			UserType:    string(enums.LoginTypeConsumer),
		})
}

// ==================== BroadcastDetail ====================

func BroadcastDetail(c *gin.Context, id string) *BroadcastVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var b imModel.Broadcast
	if err := db.DB.WithContext(ctx).First(&b, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("通知不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询通知失败: "+err.Error(), 500))
		return nil
	}
	return ImBroadcastToBroadcastVO(&b)
}
