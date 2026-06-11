package group

import (
	"context"
	"time"

	"gorm.io/gorm/clause"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	imModel "hei-gin/plugins/plugin-im/model"
	ws "hei-gin/plugins/plugin-im/ws"

	"github.com/gin-gonic/gin"
)

// ── Group messaging ─────────────────────────────────────────────

func GroupSendMessage(c *gin.Context, p *SendMessageParam) {
	ctx := c.Request.Context()
	senderID := getLoginID(c)
	senderType := getUserType(c)

	if p.GroupID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}
	if p.Content == "" && p.MsgType == "" {
		p.MsgType = imModel.MsgTypeText
	}
	if p.Content == "" && p.MsgType == imModel.MsgTypeText {
		result.WriteError(c, exception.NewBusinessError("消息不能为空", 400))
		return
	}
	if len(p.Content) > 5000 {
		result.WriteError(c, exception.NewBusinessError("消息内容不能超过5000个字符", 400))
		return
	}

	var member imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			p.GroupID, senderID, senderType, MemberActive).
		First(&member).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("不在群中", 400))
		return
	}
	if member.MutedUntil != nil && member.MutedUntil.After(time.Now()) {
		result.WriteError(c, exception.NewBusinessError("你已被禁言", 403))
		return
	}

	msgType := p.MsgType
	if msgType == "" {
		msgType = imModel.MsgTypeText
	}

	msg := imModel.GroupMessage{
		ID: utils.GenerateID(), GroupID: p.GroupID,
		SenderID: senderID, SenderType: senderType,
		Content: p.Content, Extra: p.Extra, MsgType: msgType,
		ReplyTo: p.ReplyTo,
	}
	if err := db.DB.WithContext(ctx).Create(&msg).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("发送消息失败: "+err.Error(), 500))
		return
	}

	var memberIDs []struct {
		UserID   string
		UserType string
	}
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Select("user_id, user_type").
		Where("group_id = ? AND status = ? AND NOT (user_id = ? AND user_type = ?)",
			p.GroupID, MemberActive, senderID, senderType).
		Find(&memberIDs)

	msgPayload := buildPushPayload(&msg)
	for _, m := range memberIDs {
		if m.UserType == string(enums.LoginTypeConsumer) {
			ws.GlobalCrossHub.SendToConsumer(m.UserID, ws.Message{Type: "group_message", Payload: msgPayload})
		} else {
			ws.GlobalCrossHub.SendToUser(m.UserID, ws.Message{Type: "group_message", Payload: msgPayload})
		}
	}
}

// ==================== GroupRecallMessage ====================

func GroupRecallMessage(c *gin.Context, p *RecallMessageParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" || p.MessageID == "" || operatorID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	var msg imModel.GroupMessage
	if err := db.DB.WithContext(ctx).Where("id = ? AND group_id = ?", p.MessageID, p.GroupID).First(&msg).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("消息不存在", 400))
		return
	}
	if msg.SenderID != operatorID || msg.SenderType != operatorType {
		result.WriteError(c, exception.NewBusinessError("只能撤回自己的消息", 403))
		return
	}
	if msg.CreatedAt != nil && msg.CreatedAt.Add(5*time.Minute).Before(time.Now()) {
		result.WriteError(c, exception.NewBusinessError("只能撤回5分钟内的消息", 400))
		return
	}
	if msg.MsgType == imModel.MsgTypeSystem {
		result.WriteError(c, exception.NewBusinessError("系统消息不能撤回", 400))
		return
	}

	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).Where("id = ?", p.MessageID).
		Updates(map[string]interface{}{"content": "消息已被撤回", "msg_type": imModel.MsgTypeSystem}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("撤回消息失败: "+err.Error(), 500))
		return
	}

	var memberIDs []struct {
		UserID   string
		UserType string
	}
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Select("user_id, user_type").
		Where("group_id = ? AND status = ? AND NOT (user_id = ? AND user_type = ?)",
			p.GroupID, MemberActive, operatorID, operatorType).
		Find(&memberIDs)

	msg.Content = "消息已被撤回"
	msg.MsgType = imModel.MsgTypeSystem
	recallPayload := buildRecallPayload(&msg, operatorID, operatorType)
	for _, m := range memberIDs {
		if m.UserType == string(enums.LoginTypeConsumer) {
			ws.GlobalCrossHub.SendToConsumer(m.UserID, ws.Message{Type: "group_message", Payload: recallPayload})
		} else {
			ws.GlobalCrossHub.SendToUser(m.UserID, ws.Message{Type: "group_message", Payload: recallPayload})
		}
	}
}

// ==================== GroupMarkRead ====================

func GroupMarkRead(c *gin.Context, p *MarkReadParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if p.GroupID == "" || userID == "" {
		return
	}

	_ = db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}, {Name: "user_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"read_at", "group_id"}),
	}).Create(&imModel.GroupMessageRead{
		MessageID: p.MessageID, GroupID: p.GroupID,
		ID:     utils.GenerateID(),
		UserID: userID, UserType: userType,
	}).Error
}

// ==================== GroupMuteMember ====================

func GroupMarkConversationRead(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	var p struct{ GroupID string }
	if err := c.ShouldBindJSON(&p); err != nil {
		return
	}

	if p.GroupID == "" || userID == "" {
		return
	}

	var lm struct{ ID string }
	err := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).
		Select("id").
		Where("group_id = ?", p.GroupID).
		Order("created_at DESC").
		Limit(1).
		Scan(&lm).Error
	if err != nil || lm.ID == "" {
		return
	}

	_ = db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}, {Name: "user_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"read_at", "group_id"}),
	}).Create(&imModel.GroupMessageRead{
		MessageID: lm.ID, GroupID: p.GroupID,
		ID:     utils.GenerateID(),
		UserID: userID, UserType: userType,
	}).Error
}

// ==================== Backward-compatible wrappers ====================

func Messages(ctx context.Context, groupID, cursor string, size int) ([]MessageVO, bool) {
	if groupID == "" {
		return nil, false
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	q := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).Where("group_id = ?", groupID)
	if cursor != "" {
		if t, err := utils.ParseDateTime(cursor); err == nil {
			q = q.Where("created_at < ?", t)
		}
	}
	var msgs []imModel.GroupMessage
	order := "created_at DESC"
	if cursor != "" {
		order = "created_at ASC"
	}
	q.Order(order).Limit(size + 1).Find(&msgs)
	hasMore := len(msgs) > size
	if hasMore {
		msgs = msgs[:size]
	}

	result := make([]MessageVO, len(msgs))
	for i, m := range msgs {
		fileURL := ""
		if m.MsgType == "IMAGE" || m.MsgType == "FILE" {
			fileURL = resolveFileURL(m.Content, m.Extra)
		}
		result[i] = *ImGroupMessageToMessageVO(&m, fileURL)
	}
	return result, hasMore
}

func MyGroupConversations(userID, userType string) []*ConversationVO {
	return MyGroupConversationsWithContext(context.Background(), userID, userType)
}

func MyGroupConversationsWithContext(ctx context.Context, userID, userType string) []*ConversationVO {
	if userID == "" {
		return nil
	}

	var members []imModel.GroupMember
	queryDB := db.DB.WithContext(ctx)
	queryDB.
		Select("group_id").
		Where("user_id = ? AND user_type = ? AND status = ?", userID, userType, MemberActive).
		Find(&members)
	if len(members) == 0 {
		return nil
	}

	groupIDs := make([]string, len(members))
	for i, m := range members {
		groupIDs[i] = m.GroupID
	}

	var groups []imModel.Group
	queryDB.Where("id IN ? AND status = ?", groupIDs, GroupNormal).Find(&groups)

	type cnt struct {
		GroupID string
		Count   int
	}
	var counts []cnt
	queryDB.Model(&imModel.GroupMember{}).
		Select("group_id, COUNT(*) as count").
		Where("group_id IN ? AND status = ?", groupIDs, MemberActive).
		Group("group_id").Scan(&counts)
	countMap := make(map[string]int, len(counts))
	for _, c := range counts {
		countMap[c.GroupID] = c.Count
	}

	type lm struct {
		GroupID   string
		Content   string
		CreatedAt time.Time
	}
	var lastMsgs []lm
	lastSubQ := queryDB.Table("im_group_message").
		Select("group_id, MAX(created_at) as max_ct").
		Where("group_id IN ?", groupIDs).
		Group("group_id")
	queryDB.Table("im_group_message g2").
		Select("g2.group_id, g2.content, g2.created_at").
		Joins("INNER JOIN (?) g1 ON g1.group_id = g2.group_id AND g1.max_ct = g2.created_at", lastSubQ).
		Scan(&lastMsgs)
	lastMap := make(map[string]lm, len(lastMsgs))
	for _, l := range lastMsgs {
		lastMap[l.GroupID] = l
	}

	type uc struct {
		GroupID string
		Count   int64
	}
	var unreads []uc
	readSubQ := queryDB.Table("im_group_message_read").
		Select("group_id, MAX(read_at) as max_read").
		Where("user_id = ? AND user_type = ?", userID, userType).
		Group("group_id")
	queryDB.Table("im_group_message gm").
		Select("gm.group_id, COUNT(*) as count").
		Joins("LEFT JOIN (?) gr ON gr.group_id = gm.group_id", readSubQ).
		Where("gm.group_id IN ?", groupIDs).
		Where("gm.created_at > COALESCE(gr.max_read, ?)", "1970-01-01 00:00:00").
		Group("gm.group_id").
		Scan(&unreads)
	unreadMap := make(map[string]int64, len(unreads))
	for _, u := range unreads {
		unreadMap[u.GroupID] = u.Count
	}

	result := make([]*ConversationVO, 0, len(groups))
	for _, g := range groups {
		vo := &ConversationVO{
			GroupID:     g.ID,
			GroupName:   g.Name,
			GroupAvatar: g.Avatar,
			MemberCount: countMap[g.ID],
			UnreadCount: unreadMap[g.ID],
		}
		if l, ok := lastMap[g.ID]; ok {
			vo.LastContent = l.Content
			vo.LastTime = utils.FormatDateTime(l.CreatedAt)
		}
		result = append(result, vo)
	}
	return result
}

// MarkConversationRead is a backward-compatible wrapper used by the message package.
func MarkConversationRead(ctx context.Context, groupID, userID, userType string) {
	if groupID == "" || userID == "" {
		return
	}

	var lm struct{ ID string }
	err := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).
		Select("id").
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Limit(1).
		Scan(&lm).Error
	if err != nil || lm.ID == "" {
		return
	}

	_ = db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}, {Name: "user_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"read_at", "group_id"}),
	}).Create(&imModel.GroupMessageRead{
		MessageID: lm.ID, GroupID: groupID,
		ID:     utils.GenerateID(),
		UserID: userID, UserType: userType,
	}).Error
}

// ==================== Helpers ====================
func buildPushPayload(msg *imModel.GroupMessage) map[string]interface{} {
	return map[string]interface{}{
		"message_id":  msg.ID,
		"group_id":    msg.GroupID,
		"sender_id":   msg.SenderID,
		"sender_type": msg.SenderType,
		"content":     msg.Content,
		"extra":       msg.Extra,
		"msg_type":    msg.MsgType,
		"reply_to":    msg.ReplyTo,
		"created_at":  utils.FormatDateTimePtr(msg.CreatedAt),
	}
}
func buildRecallPayload(msg *imModel.GroupMessage, recallerID, recallerType string) map[string]interface{} {
	return map[string]interface{}{
		"message_id":  msg.ID,
		"group_id":    msg.GroupID,
		"sender_id":   msg.SenderID,
		"sender_type": msg.SenderType,
		"content":     msg.Content,
		"msg_type":    msg.MsgType,
		"recalled_by": recallerID,
		"created_at":  utils.FormatDateTimePtr(msg.CreatedAt),
		"action":      "recalled",
	}
}
