package group

import (
	"encoding/json"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	imModel "hei-gin/plugins/plugin-im/model"
	ws "hei-gin/plugins/plugin-im/ws"

	"github.com/gin-gonic/gin"
)

// ── Group management ─────────────────────────────────────────────

// ── Group queries ─────────────────────────────────────────────

func (s *Service) MyGroups(c *gin.Context) []GroupVO {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if userID == "" {
		return nil
	}

	members := s.repo.ListMyGroupMemberships(ctx, userID, userType)
	if len(members) == 0 {
		return nil
	}

	groupIDs := make([]string, len(members))
	for i, m := range members {
		groupIDs[i] = m.GroupID
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

	result := make([]GroupVO, 0, len(groups))
	for _, g := range groups {
		vo := *ImGroupToGroupVO(&g)
		vo.MemberCount = countMap[g.ID]
		vo.UnreadCount = unreadMap[g.ID]
		if l, ok := lastMap[g.ID]; ok {
			vo.LastContent = l.Content
			vo.LastTime = utils.FormatDateTime(l.CreatedAt)
		}
		result = append(result, vo)
	}
	return result
}

// ==================== GroupDetail ====================

func (s *Service) Detail(c *gin.Context) *GroupVO {
	ctx := c.Request.Context()
	groupID := c.Query("group_id")

	if groupID == "" {
		return nil
	}

	group, err := s.repo.FindGroupByID(ctx, groupID)
	if err != nil {
		return nil
	}
	count := s.repo.CountActiveMembers(ctx, groupID)
	vo := ImGroupToGroupVO(group)
	vo.MemberCount = int(count)
	return vo
}

// ==================== GroupSearchGroups ====================

func (s *Service) SearchGroups(c *gin.Context, keyword string, limit int) []GroupVO {
	ctx := c.Request.Context()
	if keyword == "" || limit < 1 {
		return nil
	}
	if limit > 50 {
		limit = 50
	}
	like := "%" + keyword + "%"

	groups := s.repo.SearchGroups(ctx, like, limit)
	if len(groups) == 0 {
		return nil
	}

	groupIDs := make([]string, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.ID
	}

	counts := s.repo.CountActiveMembersByGroupIDs(ctx, groupIDs)
	countMap := make(map[string]int, len(counts))
	for _, c := range counts {
		countMap[c.GroupID] = c.Count
	}

	result := make([]GroupVO, 0, len(groups))
	for _, g := range groups {
		vo := *ImGroupToGroupVO(&g)
		vo.MemberCount = countMap[g.ID]
		result = append(result, vo)
	}
	return result
}

// ==================== GroupMembers ====================

func (s *Service) Members(c *gin.Context) []MemberVO {
	ctx := c.Request.Context()
	groupID := c.Query("group_id")

	if groupID == "" {
		return nil
	}

	members := s.repo.ListActiveMembers(ctx, groupID)
	if len(members) == 0 {
		return nil
	}
	result := make([]MemberVO, len(members))
	for i, m := range members {
		result[i] = *ImGroupMemberToMemberVO(&m)
	}
	return result
}

// ==================== GroupMessages ====================

func (s *Service) Messages(c *gin.Context, groupID, cursor string, size int) ([]MessageVO, bool) {
	ctx := c.Request.Context()
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

// ==================== GroupSearchMessages ====================

func (s *Service) SearchMessages(c *gin.Context, groupID, keyword, cursor string, size int) ([]MessageVO, bool) {
	ctx := c.Request.Context()
	if groupID == "" || keyword == "" {
		return nil, false
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	msgs := s.repo.SearchGroupMessages(ctx, groupID, keyword, cursor, size)
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

// ==================== GroupPendingJoinRequests ====================

func (s *Service) PendingJoinRequests(c *gin.Context) []imModel.GroupJoinRequest {
	ctx := c.Request.Context()
	groupID := c.Query("group_id")

	return s.repo.ListPendingJoinRequests(ctx, groupID)
}

// ==================== GroupHandleJoinRequest ====================

func (s *Service) HandleJoinRequest(c *gin.Context, p *HandleJoinRequestParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	req, err := s.repo.FindPendingJoinRequest(ctx, p.RequestID)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("申请不存在或已处理", 400))
		return
	}

	if _, _, err := s.checkOwnerOrAdmin(ctx, req.GroupID, operatorID, operatorType); err != nil {
		result.WriteError(c, err)
		return
	}

	var member *imModel.GroupMember
	var joinMsg *imModel.GroupMessage
	if p.Action == "approved" {
		member = &imModel.GroupMember{
			ID:       utils.GenerateID(),
			GroupID:  req.GroupID,
			UserID:   req.UserID,
			UserType: req.UserType,
			Role:     RoleMember,
			Status:   MemberActive,
		}

		extra := imModel.MsgExtraSystem{Action: "join", UserID: req.UserID, UserType: req.UserType}
		extraBytes, _ := json.Marshal(extra)
		joinMsg = &imModel.GroupMessage{
			ID: utils.GenerateID(), GroupID: req.GroupID,
			SenderID: req.UserID, SenderType: req.UserType,
			Content: "加入了群聊", Extra: string(extraBytes),
			MsgType: imModel.MsgTypeSystem,
		}
	}

	if err := s.repo.HandleJoinRequest(ctx, p.RequestID, operatorID, p.Action, member, joinMsg); err != nil {
		result.WriteError(c, exception.NewBusinessError("处理失败: "+err.Error(), 500))
		return
	}

	msg := map[string]interface{}{
		"group_id": req.GroupID,
		"status":   p.Action,
		"action":   "join_request_result",
	}
	if req.UserType == string(enums.LoginTypeConsumer) {
		ws.GlobalCrossHub.SendToConsumer(req.UserID, ws.Message{Type: "group_event", Payload: msg})
	} else {
		ws.GlobalCrossHub.SendToUser(req.UserID, ws.Message{Type: "group_event", Payload: msg})
	}
}

// ==================== GroupMarkConversationRead ====================

// Messages is a backward-compatible wrapper used by the message package.
// MyGroupConversations is a backward-compatible wrapper used by the message package.
