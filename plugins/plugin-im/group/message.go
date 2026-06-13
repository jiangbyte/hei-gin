package group

import (
	"context"
	"time"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	imModel "hei-gin/plugins/plugin-im/model"
	ws "hei-gin/plugins/plugin-im/ws"

	"github.com/gin-gonic/gin"
)

// ── Group messaging ─────────────────────────────────────────────

func (s *Service) SendMessage(c *gin.Context, p *SendMessageParam) {
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

	member, err := s.repo.FindActiveMember(ctx, p.GroupID, senderID, senderType)
	if err != nil {
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
	if err := s.repo.CreateMessage(ctx, &msg); err != nil {
		result.WriteError(c, exception.NewBusinessError("发送消息失败: "+err.Error(), 500))
		return
	}

	memberIDs := s.repo.ListRecipientMembers(ctx, p.GroupID, senderID, senderType)

	msgPayload := buildPushPayload(&msg)
	if runtime := ws.Runtime(); runtime != nil {
		for _, m := range memberIDs {
			if m.UserType == string(enums.LoginTypeConsumer) {
				runtime.SendToConsumer(m.UserID, ws.Message{Type: "group_message", Payload: msgPayload})
			} else {
				runtime.SendToUser(m.UserID, ws.Message{Type: "group_message", Payload: msgPayload})
			}
		}
	}
}

// ==================== GroupRecallMessage ====================

func (s *Service) RecallMessage(c *gin.Context, p *RecallMessageParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" || p.MessageID == "" || operatorID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	msg, err := s.repo.FindMessageByID(ctx, p.MessageID, p.GroupID)
	if err != nil {
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

	if err := s.repo.RecallMessage(ctx, p.MessageID); err != nil {
		result.WriteError(c, exception.NewBusinessError("撤回消息失败: "+err.Error(), 500))
		return
	}

	memberIDs := s.repo.ListRecipientMembers(ctx, p.GroupID, operatorID, operatorType)

	msg.Content = "消息已被撤回"
	msg.MsgType = imModel.MsgTypeSystem
	recallPayload := buildRecallPayload(msg, operatorID, operatorType)
	if runtime := ws.Runtime(); runtime != nil {
		for _, m := range memberIDs {
			if m.UserType == string(enums.LoginTypeConsumer) {
				runtime.SendToConsumer(m.UserID, ws.Message{Type: "group_message", Payload: recallPayload})
			} else {
				runtime.SendToUser(m.UserID, ws.Message{Type: "group_message", Payload: recallPayload})
			}
		}
	}
}

// ==================== GroupMarkRead ====================

func (s *Service) MarkRead(c *gin.Context, p *MarkReadParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if p.GroupID == "" || userID == "" {
		return
	}

	_ = s.repo.UpsertMessageRead(ctx, &imModel.GroupMessageRead{
		MessageID: p.MessageID, GroupID: p.GroupID,
		ID:     utils.GenerateID(),
		UserID: userID, UserType: userType,
	})
}

// ==================== GroupMuteMember ====================

func (s *Service) MarkConversationRead(c *gin.Context) {
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

	lastMessageID := s.repo.FindLastMessageID(ctx, p.GroupID)
	if lastMessageID == "" {
		return
	}

	_ = s.repo.UpsertMessageRead(ctx, &imModel.GroupMessageRead{
		MessageID: lastMessageID, GroupID: p.GroupID,
		ID:     utils.GenerateID(),
		UserID: userID, UserType: userType,
	})
}

// ==================== Backward-compatible wrappers ====================

func Messages(ctx context.Context, groupID, cursor string, size int) ([]MessageVO, bool) {
	return DefaultModule.service.listMessages(ctx, groupID, cursor, size)
}

func (s *Service) listMessages(ctx context.Context, groupID, cursor string, size int) ([]MessageVO, bool) {
	if groupID == "" {
		return nil, false
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	msgs := s.repo.ListGroupMessages(ctx, groupID, cursor, size)
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

func (s *Service) MyGroupConversations(userID, userType string) []*ConversationVO {
	return s.MyGroupConversationsWithContext(context.Background(), userID, userType)
}

func (s *Service) MyGroupConversationsWithContext(ctx context.Context, userID, userType string) []*ConversationVO {
	if userID == "" {
		return nil
	}

	groupIDs := s.repo.ListMyGroupIDs(ctx, userID, userType)
	if len(groupIDs) == 0 {
		return nil
	}
	groups := s.repo.ListGroupsByIDs(ctx, groupIDs)
	counts := s.repo.CountActiveMembersByGroupIDs(ctx, groupIDs)
	countMap := make(map[string]int, len(counts))
	for _, c := range counts {
		countMap[c.GroupID] = c.Count
	}
	lastMsgs := s.repo.ListLastMessagesByGroupIDs(ctx, groupIDs)
	lastMap := make(map[string]groupLastMessageRow, len(lastMsgs))
	for _, l := range lastMsgs {
		lastMap[l.GroupID] = l
	}
	unreads := s.repo.CountUnreadByGroupIDs(ctx, groupIDs, userID, userType)
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
func (s *Service) MarkConversationReadWithContext(ctx context.Context, groupID, userID, userType string) {
	if groupID == "" || userID == "" {
		return
	}

	lastMessageID := s.repo.FindLastMessageID(ctx, groupID)
	if lastMessageID == "" {
		return
	}

	_ = s.repo.UpsertMessageRead(ctx, &imModel.GroupMessageRead{
		MessageID: lastMessageID, GroupID: groupID,
		ID:     utils.GenerateID(),
		UserID: userID, UserType: userType,
	})
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
