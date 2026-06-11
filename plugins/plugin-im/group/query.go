package group

import (
	"encoding/json"
	"time"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	imModel "hei-gin/plugins/plugin-im/model"
	ws "hei-gin/plugins/plugin-im/ws"

	"github.com/gin-gonic/gin"
)

// ── Group management ─────────────────────────────────────────────

// ── Group queries ─────────────────────────────────────────────

func GroupMyGroups(c *gin.Context) []GroupVO {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if userID == "" {
		return nil
	}

	var members []imModel.GroupMember
	db.DB.WithContext(ctx).
		Select("group_id, joined_at").
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
	db.DB.WithContext(ctx).Where("id IN ? AND status = ?", groupIDs, GroupNormal).Find(&groups)

	type cnt struct {
		GroupID string
		Count   int
	}
	var counts []cnt
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
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
	lastSubQ := db.DB.WithContext(ctx).Table("im_group_message").
		Select("group_id, MAX(created_at) as max_ct").
		Where("group_id IN ?", groupIDs).
		Group("group_id")
	db.DB.WithContext(ctx).Table("im_group_message g2").
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
	readSubQ := db.DB.WithContext(ctx).Table("im_group_message_read").
		Select("group_id, MAX(read_at) as max_read").
		Where("user_id = ? AND user_type = ?", userID, userType).
		Group("group_id")
	db.DB.WithContext(ctx).Table("im_group_message gm").
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

func GroupDetail(c *gin.Context) *GroupVO {
	ctx := c.Request.Context()
	groupID := c.Query("group_id")

	if groupID == "" {
		return nil
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", groupID).Error; err != nil {
		return nil
	}
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND status = ?", groupID, MemberActive).Count(&count)
	vo := ImGroupToGroupVO(&group)
	vo.MemberCount = int(count)
	return vo
}

// ==================== GroupSearchGroups ====================

func GroupSearchGroups(c *gin.Context, keyword string, limit int) []GroupVO {
	ctx := c.Request.Context()
	if keyword == "" || limit < 1 {
		return nil
	}
	if limit > 50 {
		limit = 50
	}
	like := "%" + keyword + "%"

	var groups []imModel.Group
	db.DB.WithContext(ctx).
		Where("name LIKE ? AND status = ?", like, GroupNormal).
		Limit(limit).
		Find(&groups)
	if len(groups) == 0 {
		return nil
	}

	groupIDs := make([]string, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.ID
	}

	type cnt struct {
		GroupID string
		Count   int
	}
	var counts []cnt
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Select("group_id, COUNT(*) as count").
		Where("group_id IN ? AND status = ?", groupIDs, MemberActive).
		Group("group_id").Scan(&counts)
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

func GroupMembers(c *gin.Context) []MemberVO {
	ctx := c.Request.Context()
	groupID := c.Query("group_id")

	if groupID == "" {
		return nil
	}

	var members []imModel.GroupMember
	db.DB.WithContext(ctx).
		Where("group_id = ? AND status = ?", groupID, MemberActive).
		Order("FIELD(role, 'owner', 'admin', 'member'), joined_at ASC").
		Find(&members)
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

func GroupMessages(c *gin.Context, groupID, cursor string, size int) ([]MessageVO, bool) {
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

// ==================== GroupSearchMessages ====================

func GroupSearchMessages(c *gin.Context, groupID, keyword, cursor string, size int) ([]MessageVO, bool) {
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

	q := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).
		Where("group_id = ? AND content LIKE ? AND msg_type != ?", groupID, "%"+keyword+"%", imModel.MsgTypeSystem)
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

// ==================== GroupPendingJoinRequests ====================

func GroupPendingJoinRequests(c *gin.Context) []imModel.GroupJoinRequest {
	ctx := c.Request.Context()
	groupID := c.Query("group_id")

	var requests []imModel.GroupJoinRequest
	db.DB.WithContext(ctx).Model(&imModel.GroupJoinRequest{}).
		Where("group_id = ? AND status = ?", groupID, "pending").
		Order("created_at DESC").Find(&requests)
	return requests
}

// ==================== GroupHandleJoinRequest ====================

func GroupHandleJoinRequest(c *gin.Context, p *HandleJoinRequestParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	var req imModel.GroupJoinRequest
	if err := db.DB.WithContext(ctx).First(&req, "id = ? AND status = ?", p.RequestID, "pending").Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("申请不存在或已处理", 400))
		return
	}

	if _, _, err := checkOwnerOrAdmin(ctx, req.GroupID, operatorID, operatorType); err != nil {
		result.WriteError(c, err)
		return
	}

	if err := db.DB.WithContext(ctx).Model(&imModel.GroupJoinRequest{}).Where("id = ?", p.RequestID).
		Updates(map[string]interface{}{"status": p.Action, "handled_by": operatorID}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("处理失败: "+err.Error(), 500))
		return
	}

	if p.Action == "approved" {
		if err := db.DB.WithContext(ctx).Create(&imModel.GroupMember{
			ID:       utils.GenerateID(),
			GroupID:  req.GroupID,
			UserID:   req.UserID,
			UserType: req.UserType,
			Role:     RoleMember,
			Status:   MemberActive,
		}).Error; err != nil {
			result.WriteError(c, exception.NewBusinessError("添加成员失败: "+err.Error(), 500))
			return
		}

		extra := imModel.MsgExtraSystem{Action: "join", UserID: req.UserID, UserType: req.UserType}
		extraBytes, _ := json.Marshal(extra)
		db.DB.WithContext(ctx).Create(&imModel.GroupMessage{
			ID: utils.GenerateID(), GroupID: req.GroupID,
			SenderID: req.UserID, SenderType: req.UserType,
			Content: "加入了群聊", Extra: string(extraBytes),
			MsgType: imModel.MsgTypeSystem,
		})
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
