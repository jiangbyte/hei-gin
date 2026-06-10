package group

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"


	"gorm.io/gorm/clause"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/storage"
	ws "hei-gin/plugins/plugin-im/ws"

	imModel "hei-gin/plugins/plugin-im/model"
)

// ==================== Create ====================

func Create(ctx context.Context, ownerID, ownerType string, p *CreateParam) *imModel.Group {
	if p.Name == "" {
		panic(exception.NewBusinessError("群名称不能为空", 400))
	}
	if len(p.Name) > 100 {
		panic(exception.NewBusinessError("群名称不能超过100个字符", 400))
	}
	if ownerID == "" {
		panic(exception.NewBusinessError("用户未登录", 401))
	}

	groupType := GroupTypeMixed
	if ownerType == UserTypeConsumer {
		groupType = GroupTypeConsumerOnly
	}

	now := time.Now()
	group := imModel.Group{
		ID:         utils.GenerateID(),
		Name:       p.Name,
		Avatar:     p.Avatar,
		OwnerID:    ownerID,
		OwnerType:  ownerType,
		GroupType:  groupType,
		MaxMembers: 200,
		Status:     GroupNormal,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(&group).Error; err != nil {
		panic(exception.NewBusinessError("创建群失败", 500))
	}

	ownerMember := imModel.GroupMember{
		ID: utils.GenerateID(), GroupID: group.ID,
		UserID: ownerID, UserType: ownerType,
		Role: RoleOwner, JoinedAt: &now, Status: MemberActive,
	}
	if err := tx.Create(&ownerMember).Error; err != nil {
		panic(exception.NewBusinessError("添加群主失败", 500))
	}

	if len(p.MemberIDs) > 0 {
		validateMemberType(groupType, p.MemberType)

		var existingCount int64
		tx.Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id IN ? AND user_type = ? AND status = ?",
				group.ID, p.MemberIDs, p.MemberType, MemberActive).
			Count(&existingCount)
		if existingCount > 0 {
			panic(exception.NewBusinessError("部分成员已在群中", 400))
		}

		var currentCount int64
		tx.Model(&imModel.GroupMember{}).Where("group_id = ? AND status = ?", group.ID, MemberActive).Count(&currentCount)
		if int(currentCount)+len(p.MemberIDs) > group.MaxMembers {
			panic(exception.NewBusinessError(fmt.Sprintf("群成员数量不能超过%d人", group.MaxMembers), 400))
		}

		batch := make([]imModel.GroupMember, 0, len(p.MemberIDs))
		sysBatch := make([]imModel.GroupMessage, 0, len(p.MemberIDs))
		for _, uid := range p.MemberIDs {
			if uid == ownerID {
				continue
			}
			batch = append(batch, imModel.GroupMember{
				ID: utils.GenerateID(), GroupID: group.ID,
				UserID: uid, UserType: p.MemberType,
				Role: RoleMember, JoinedAt: &now, Status: MemberActive,
			})
			extra := imModel.MsgExtraSystem{Action: "join", UserID: uid, UserType: p.MemberType}
			extraBytes, _ := json.Marshal(extra)
			sysBatch = append(sysBatch, imModel.GroupMessage{
				ID: utils.GenerateID(), GroupID: group.ID,
				SenderID: ownerID, SenderType: ownerType,
				Content: "欢迎加入群聊", Extra: string(extraBytes),
				MsgType: imModel.MsgTypeSystem, CreatedAt: &now,
			})
		}
		if len(batch) > 0 {
			if err := tx.Create(&batch).Error; err != nil {
				panic(exception.NewBusinessError("邀请成员失败", 500))
			}
		}
		if len(sysBatch) > 0 {
			if err := tx.Create(&sysBatch).Error; err != nil {
				panic(exception.NewBusinessError("发送系统消息失败", 500))
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("创建群失败", 500))
	}
	return &group
}

// ==================== Update ====================

func Update(ctx context.Context, operatorID, operatorType string, p *UpdateParam) {
	if p.GroupID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}
	_, member := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)

	updates := map[string]interface{}{"updated_at": time.Now()}
	if p.Name != nil {
		if len(*p.Name) > 100 {
			panic(exception.NewBusinessError("群名称不能超过100个字符", 400))
		}
		updates["name"] = *p.Name
	}
	if p.Avatar != nil {
		updates["avatar"] = *p.Avatar
	}
	if p.Notice != nil && member.Role == RoleOwner {
		updates["notice"] = *p.Notice
	}

	if err := db.DB.WithContext(ctx).Model(&imModel.Group{}).Where("id = ?", p.GroupID).Updates(updates).Error; err != nil {
		panic(exception.NewBusinessError("修改群信息失败", 500))
	}
}

// ==================== Dissolve ====================

func Dissolve(ctx context.Context, operatorID string, groupID string) {
	if groupID == "" || operatorID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", groupID).Error; err != nil {
		panic(exception.NewBusinessError("群不存在", 400))
	}
	if group.OwnerID != operatorID {
		panic(exception.NewBusinessError("仅群主可解散群", 403))
	}

	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Model(&imModel.Group{}).Where("id = ?", groupID).
		Updates(map[string]interface{}{"status": GroupDissolved, "updated_at": time.Now()}).Error; err != nil {
		panic(exception.NewBusinessError("解散群失败", 500))
	}

	if err := tx.Model(&imModel.GroupMember{}).Where("group_id = ? AND status = ?", groupID, MemberActive).
		Update("status", MemberLeft).Error; err != nil {
		panic(exception.NewBusinessError("解散群失败", 500))
	}

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("解散群失败", 500))
	}
}

// ==================== Invite ====================

func Invite(ctx context.Context, operatorID, operatorType string, p *InviteParam) {
	if len(p.UserIDs) == 0 {
		return
	}
	if p.GroupID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", p.GroupID).Error; err != nil {
		panic(exception.NewBusinessError("群不存在", 400))
	}
	if group.Status != GroupNormal {
		panic(exception.NewBusinessError("群已解散", 400))
	}
	validateMemberType(group.GroupType, p.UserType)

	// Verify operator is owner or admin
	checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)

	// Single batch check for existing members
	var existingIDs []string
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id IN ? AND user_type = ? AND status = ?",
			p.GroupID, p.UserIDs, p.UserType, MemberActive).
		Pluck("user_id", &existingIDs)
	if len(existingIDs) > 0 {
		panic(exception.NewBusinessError(fmt.Sprintf("用户 %v 已在群中", existingIDs), 400))
	}

	// Check max members
	var currentCount int64
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND status = ?", p.GroupID, MemberActive).Count(&currentCount)
	if int(currentCount)+len(p.UserIDs) > group.MaxMembers {
		panic(exception.NewBusinessError(fmt.Sprintf("群成员数量不能超过%d人", group.MaxMembers), 400))
	}

	now := time.Now()
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	batch := make([]imModel.GroupMember, 0, len(p.UserIDs))
	sysBatch := make([]imModel.GroupMessage, 0, len(p.UserIDs))
	for _, uid := range p.UserIDs {
		batch = append(batch, imModel.GroupMember{
			ID: utils.GenerateID(), GroupID: p.GroupID,
			UserID: uid, UserType: p.UserType,
			Role: RoleMember, JoinedAt: &now, Status: MemberActive,
		})
		extra := imModel.MsgExtraSystem{Action: "join", UserID: uid, UserType: p.UserType}
		extraBytes, _ := json.Marshal(extra)
		sysBatch = append(sysBatch, imModel.GroupMessage{
			ID: utils.GenerateID(), GroupID: p.GroupID,
			SenderID: operatorID, SenderType: operatorType,
			Content: "欢迎加入群聊", Extra: string(extraBytes),
			MsgType: imModel.MsgTypeSystem, CreatedAt: &now,
		})
	}
	if err := tx.Create(&batch).Error; err != nil {
		panic(exception.NewBusinessError("邀请成员失败", 500))
	}
	if err := tx.Create(&sysBatch).Error; err != nil {
		panic(exception.NewBusinessError("发送系统消息失败", 500))
	}

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("邀请成员失败", 500))
	}
}

// ==================== Join ====================

func Leave(ctx context.Context, userID, userType string, groupID string) {
	if groupID == "" || userID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	var member imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, MemberActive).First(&member).Error; err != nil {
		panic(exception.NewBusinessError("不在群中", 400))
	}
	if member.Role == RoleOwner {
		panic(exception.NewBusinessError("群主不能退群，请先转让群", 400))
	}

	now := time.Now()
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Model(&member).Update("status", MemberLeft).Error; err != nil {
		panic(exception.NewBusinessError("退群失败", 500))
	}

	extra := imModel.MsgExtraSystem{Action: "leave", UserID: userID, UserType: userType}
	extraBytes, _ := json.Marshal(extra)
	if err := tx.Create(&imModel.GroupMessage{
		ID: utils.GenerateID(), GroupID: groupID,
		SenderID: userID, SenderType: userType,
		Content: "退出了群聊", Extra: string(extraBytes),
		MsgType: imModel.MsgTypeSystem, CreatedAt: &now,
	}).Error; err != nil {
		panic(exception.NewBusinessError("发送系统消息失败", 500))
	}

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("退群失败", 500))
	}
}

// ==================== Kick ====================

func Kick(ctx context.Context, operatorID, operatorType string, p *KickParam) {
	if p.GroupID == "" || p.UserID == "" || p.UserType == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	group, operatorMember := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
	if p.UserID == group.OwnerID {
		panic(exception.NewBusinessError("不能踢出群主", 400))
	}
	if operatorMember.Role == RoleAdmin {
		var target imModel.GroupMember
		if err := db.DB.WithContext(ctx).
			Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
				p.GroupID, p.UserID, p.UserType, MemberActive).First(&target).Error; err != nil {
			panic(exception.NewBusinessError("成员不存在或已离开", 400))
		}
		if target.Role != RoleMember {
			panic(exception.NewBusinessError("不能踢出管理员", 403))
		}
	}

	now := time.Now()
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
		Update("status", MemberKicked).Error; err != nil {
		panic(exception.NewBusinessError("踢出失败", 500))
	}

	extra := imModel.MsgExtraSystem{Action: "kick", UserID: p.UserID, UserType: p.UserType, OperatorID: operatorID}
	extraBytes, _ := json.Marshal(extra)
	if err := tx.Create(&imModel.GroupMessage{
		ID: utils.GenerateID(), GroupID: p.GroupID,
		SenderID: operatorID, SenderType: operatorType,
		Content: "被移出群聊", Extra: string(extraBytes),
		MsgType: imModel.MsgTypeSystem, CreatedAt: &now,
	}).Error; err != nil {
		panic(exception.NewBusinessError("发送系统消息失败", 500))
	}

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("踢出失败", 500))
	}
}

// ==================== SetRole / Transfer ====================

func SetRole(ctx context.Context, operatorID string, p *SetRoleParam) {
	if p.GroupID == "" || p.UserID == "" || p.Role == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", p.GroupID).Error; err != nil {
		panic(exception.NewBusinessError("群不存在", 400))
	}
	if group.OwnerID != operatorID {
		panic(exception.NewBusinessError("仅群主可设置角色", 403))
	}

	switch p.Role {
	case RoleAdmin:
		if err := db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
				p.GroupID, p.UserID, p.UserType, MemberActive).
			Update("role", RoleAdmin).Error; err != nil {
			panic(exception.NewBusinessError("设置管理员失败", 500))
		}
	case RoleOwner:
		now := time.Now()
		tx := db.DB.WithContext(ctx).Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			}
		}()

		if err := tx.Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, group.OwnerID, group.OwnerType).
			Update("role", RoleAdmin).Error; err != nil {
			panic(exception.NewBusinessError("转让群失败", 500))
		}
		if err := tx.Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
			Update("role", RoleOwner).Error; err != nil {
			panic(exception.NewBusinessError("转让群失败", 500))
		}
		if err := tx.Model(&imModel.Group{}).Where("id = ?", p.GroupID).
			Updates(map[string]interface{}{"owner_id": p.UserID, "owner_type": p.UserType, "updated_at": now}).Error; err != nil {
			panic(exception.NewBusinessError("转让群失败", 500))
		}
		if err := tx.Commit().Error; err != nil {
			panic(exception.NewBusinessError("转让群失败", 500))
		}
	default:
		panic(exception.NewBusinessError("不支持的角色", 400))
	}
}

// ==================== Send Message ====================

func SendMessage(ctx context.Context, senderID, senderType string, p *SendMessageParam) *MessageVO {
	if p.GroupID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}
	if p.Content == "" && p.MsgType == imModel.MsgTypeText {
		panic(exception.NewBusinessError("消息不能为空", 400))
	}
	if len(p.Content) > 5000 {
		panic(exception.NewBusinessError("消息内容不能超过5000个字符", 400))
	}

	var member imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			p.GroupID, senderID, senderType, MemberActive).
		First(&member).Error; err != nil {
		panic(exception.NewBusinessError("不在群中", 400))
	}
	if member.MutedUntil != nil && member.MutedUntil.After(time.Now()) {
		panic(exception.NewBusinessError("你已被禁言", 403))
	}

	msgType := p.MsgType
	if msgType == "" {
		msgType = imModel.MsgTypeText
	}

	now := time.Now()
	msg := imModel.GroupMessage{
		ID: utils.GenerateID(), GroupID: p.GroupID,
		SenderID: senderID, SenderType: senderType,
		Content: p.Content, Extra: p.Extra, MsgType: msgType,
		ReplyTo: p.ReplyTo, CreatedAt: &now,
	}
	if err := db.DB.WithContext(ctx).Create(&msg).Error; err != nil {
		panic(exception.NewBusinessError("发送消息失败", 500))
	}

	// Batch get all active member IDs excluding sender (single query, avoid N+1)
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
		if m.UserType == UserTypeConsumer {
			ws.GlobalCrossHub.SendToConsumer(m.UserID, ws.Message{Type: "group_message", Payload: msgPayload})
		} else {
			ws.GlobalCrossHub.SendToUser(m.UserID, ws.Message{Type: "group_message", Payload: msgPayload})
		}
	}

	fileURL := ""
	if msg.MsgType == "IMAGE" || msg.MsgType == "FILE" {
		fileURL = resolveFileURL(msg.Content, msg.Extra)
	}
	return &MessageVO{
		ID: msg.ID, SenderID: msg.SenderID, SenderType: msg.SenderType,
		Content: msg.Content, Extra: msg.Extra, MsgType: msg.MsgType,
		ReplyTo: msg.ReplyTo, FileURL: fileURL,
		CreatedAt: utils.FormatDateTimePtr(msg.CreatedAt),
	}
}

// ==================== Group List ====================

func MyGroups(ctx context.Context, userID, userType string) []GroupVO {
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

	// Single batch: get all groups
	var groups []imModel.Group
	db.DB.WithContext(ctx).Where("id IN ? AND status = ?", groupIDs, GroupNormal).Find(&groups)

	// Single batch: count active members per group
	type cnt struct{ GroupID string; Count int }
	var counts []cnt
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Select("group_id, COUNT(*) as count").
		Where("group_id IN ? AND status = ?", groupIDs, MemberActive).
		Group("group_id").Scan(&counts)
	countMap := make(map[string]int, len(counts))
	for _, c := range counts {
		countMap[c.GroupID] = c.Count
	}

	// Single batch: get last message per group (using subquery)
	type lm struct{ GroupID string; Content string; CreatedAt time.Time }
	// Single batch: get last message per group using GORM subquery (no raw SQL)
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

	type uc struct{ GroupID string; Count int64 }
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
	for _, u := range unreads {
		unreadMap[u.GroupID] = u.Count
	}
	for _, l := range lastMsgs {
		lastMap[l.GroupID] = l
	}

	result := make([]GroupVO, 0, len(groups))
	for _, g := range groups {
		vo := GroupVO{
			ID: g.ID, Name: g.Name, Avatar: g.Avatar,
			OwnerID: g.OwnerID, OwnerType: g.OwnerType,
			GroupType: g.GroupType, Notice: g.Notice,
			MemberCount: countMap[g.ID],
			UnreadCount:  unreadMap[g.ID],
		}
		if l, ok := lastMap[g.ID]; ok {
			vo.LastContent = l.Content
			vo.LastTime = utils.FormatDateTime(l.CreatedAt)
		}
		result = append(result, vo)
	}
	return result
}

// ==================== Detail ====================

func Detail(ctx context.Context, groupID string) *GroupVO {
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
	return &GroupVO{
		ID: group.ID, Name: group.Name, Avatar: group.Avatar,
		OwnerID: group.OwnerID, OwnerType: group.OwnerType,
		GroupType: group.GroupType, Notice: group.Notice,
		MemberCount: int(count),
	}
}


// ==================== Search Groups ====================

func SearchGroups(ctx context.Context, keyword string, limit int) []GroupVO {
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

	type cnt struct{ GroupID string; Count int }
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
		result = append(result, GroupVO{
			ID:          g.ID,
			Name:        g.Name,
			Avatar:      g.Avatar,
			OwnerID:     g.OwnerID,
			OwnerType:   g.OwnerType,
			GroupType:   g.GroupType,
			Notice:      g.Notice,
			MemberCount: countMap[g.ID],
		})
	}
	return result
}

// ==================== Members List ====================

func Members(ctx context.Context, groupID string) []MemberVO {
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
		result[i] = MemberVO{
			UserID: m.UserID, UserType: m.UserType,
			Role: m.Role, Nickname: m.Nickname,
			IsMuted: m.MutedUntil != nil && m.MutedUntil.After(time.Now()),
			JoinedAt: utils.FormatDateTimePtr(m.JoinedAt),
		}
	}
	return result
}

// ==================== Messages (cursor pagination) ====================

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
		result[i] = MessageVO{
			ID: m.ID, SenderID: m.SenderID, SenderType: m.SenderType,
			Content: m.Content, Extra: m.Extra, MsgType: m.MsgType,
			ReplyTo: m.ReplyTo, FileURL: fileURL,
			CreatedAt: utils.FormatDateTimePtr(m.CreatedAt),
		}
	}
	return result, hasMore
}

// ==================== MarkRead ====================


// ==================== Search Messages ====================

func SearchMessages(ctx context.Context, groupID, keyword string, cursor string, size int) ([]MessageVO, bool) {
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
		result[i] = MessageVO{
			ID: m.ID, SenderID: m.SenderID, SenderType: m.SenderType,
			Content: m.Content, Extra: m.Extra, MsgType: m.MsgType,
			ReplyTo: m.ReplyTo, FileURL: fileURL,
			CreatedAt: utils.FormatDateTimePtr(m.CreatedAt),
		}
	}
	return result, hasMore
}

func MarkRead(ctx context.Context, groupID, userID, userType, messageID string) {
	if groupID == "" || userID == "" {
		return
	}
	now := time.Now()
	_ = db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}, {Name: "user_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"read_at", "group_id"}),
	}).Create(&imModel.GroupMessageRead{
		MessageID: messageID, GroupID: groupID,
		ID: utils.GenerateID(),
		UserID: userID, UserType: userType, ReadAt: &now,
	}).Error
}

// ==================== Recall Message ====================

func RecallMessage(ctx context.Context, groupID, messageID, operatorID, operatorType string) {
	if groupID == "" || messageID == "" || operatorID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	var msg imModel.GroupMessage
	if err := db.DB.WithContext(ctx).Where("id = ? AND group_id = ?", messageID, groupID).First(&msg).Error; err != nil {
		panic(exception.NewBusinessError("消息不存在", 400))
	}
	if msg.SenderID != operatorID || msg.SenderType != operatorType {
		panic(exception.NewBusinessError("只能撤回自己的消息", 403))
	}
	if msg.CreatedAt.Add(5 * time.Minute).Before(time.Now()) {
		panic(exception.NewBusinessError("只能撤回5分钟内的消息", 400))
	}
	if msg.MsgType == imModel.MsgTypeSystem {
		panic(exception.NewBusinessError("系统消息不能撤回", 400))
	}

	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).Where("id = ?", messageID).
		Updates(map[string]interface{}{"content": "消息已被撤回", "msg_type": imModel.MsgTypeSystem}).Error; err != nil {
		panic(exception.NewBusinessError("撤回消息失败", 500))
	}

	// Push recall notification to all members EXCEPT sender
	var memberIDs []struct{ UserID string; UserType string }
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Select("user_id, user_type").
		Where("group_id = ? AND status = ? AND NOT (user_id = ? AND user_type = ?)",
			groupID, MemberActive, operatorID, operatorType).
		Find(&memberIDs)

	msg.Content = "消息已被撤回"
	msg.MsgType = imModel.MsgTypeSystem
	recallPayload := buildRecallPayload(&msg, operatorID, operatorType)
	for _, m := range memberIDs {
		if m.UserType == UserTypeConsumer {
			ws.GlobalCrossHub.SendToConsumer(m.UserID, ws.Message{Type: "group_message", Payload: recallPayload})
		} else {
			ws.GlobalCrossHub.SendToUser(m.UserID, ws.Message{Type: "group_message", Payload: recallPayload})
		}
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

// ==================== Mute / Unmute ====================

func MuteMember(ctx context.Context, operatorID, operatorType string, p *KickParam, duration time.Duration) {
	if p.GroupID == "" || p.UserID == "" || p.UserType == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	group, operator := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
	if p.UserID == group.OwnerID {
		panic(exception.NewBusinessError("不能禁言群主", 400))
	}
	if operator.Role == RoleAdmin {
		var target imModel.GroupMember
		if err := db.DB.WithContext(ctx).
			Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
				p.GroupID, p.UserID, p.UserType, MemberActive).First(&target).Error; err != nil {
			panic(exception.NewBusinessError("成员不存在", 400))
		}
		if target.Role != RoleMember {
			panic(exception.NewBusinessError("不能禁言管理员", 403))
		}
	}
	until := time.Now().Add(duration)
	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
		Update("muted_until", &until).Error; err != nil {
		panic(exception.NewBusinessError("禁言失败", 500))
	}
}

func UnmuteMember(ctx context.Context, operatorID, operatorType string, p *KickParam) {
	if p.GroupID == "" || p.UserID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
		Update("muted_until", nil).Error; err != nil {
		panic(exception.NewBusinessError("解除禁言失败", 500))
	}
}

// ==================== Helpers ====================

func validateMemberType(groupType, userType string) {
	if groupType == GroupTypeConsumerOnly && userType != UserTypeConsumer {
		panic(exception.NewBusinessError("该群仅限C端用户", 403))
	}
}

func checkOwnerOrAdmin(ctx context.Context, groupID, userID, userType string) (*imModel.Group, *imModel.GroupMember) {
	if groupID == "" || userID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", groupID).Error; err != nil {
		panic(exception.NewBusinessError("群不存在", 400))
	}
	var member imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, MemberActive).
		First(&member).Error; err != nil {
		panic(exception.NewBusinessError("不在群中", 400))
	}
	if member.Role != RoleOwner && member.Role != RoleAdmin {
		panic(exception.NewBusinessError("无权限", 403))
	}
	return &group, &member
}

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

// MarkConversationRead marks an entire group conversation as read by recording
// the latest message as the user's last-read position (upserts on user_id+user_type+message_id).

// MarkConversationRead marks an entire group conversation as read by recording
// the latest message as the user's last-read position (upserts on user_id+user_type+message_id).
func MarkConversationRead(ctx context.Context, groupID, userID, userType string) {
	if groupID == "" || userID == "" {
		return
	}

	// Find the latest message in the group
	type lastMsg struct {
		ID string
	}
	var lm lastMsg
	err := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).
		Select("id").
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Limit(1).
		Scan(&lm).Error
	if err != nil || lm.ID == "" {
		return
	}

	now := time.Now()
	_ = db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}, {Name: "user_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"read_at", "group_id"}),
	}).Create(&imModel.GroupMessageRead{
		MessageID: lm.ID, GroupID: groupID,
		ID: utils.GenerateID(),
		UserID: userID, UserType: userType, ReadAt: &now,
	}).Error
}

// ==================== Join (rewritten: creates join request) ====================

func Join(ctx context.Context, userID, userType string, groupID string) {
	if groupID == "" || userID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", groupID).Error; err != nil {
		panic(exception.NewBusinessError("群不存在", 400))
	}
	if group.Status != GroupNormal {
		panic(exception.NewBusinessError("群已解散", 400))
	}
	validateMemberType(group.GroupType, userType)

	// Check if already a member
	var existing int64
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, MemberActive).Count(&existing)
	if existing > 0 {
		panic(exception.NewBusinessError("已在群中", 400))
	}

	// Check pending join request
	var pending int64
	db.DB.WithContext(ctx).Model(&imModel.GroupJoinRequest{}).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, "pending").Count(&pending)
	if pending > 0 {
		panic(exception.NewBusinessError("已发送过入群申请，请等待审核", 400))
	}

	now := time.Now()
	if err := db.DB.WithContext(ctx).Create(&imModel.GroupJoinRequest{
		ID:        utils.GenerateID(),
		GroupID:   groupID,
		UserID:    userID,
		UserType:  userType,
		Status:    "pending",
		CreatedAt: &now,
		UpdatedAt: &now,
	}).Error; err != nil {
		panic(exception.NewBusinessError("申请加入失败", 500))
	}

	// Notify group admins via WS
	var members []imModel.GroupMember
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND (role = ? OR role = ?) AND status = ?", groupID, RoleOwner, RoleAdmin, MemberActive).
		Find(&members)
	for _, m := range members {
		payload := map[string]interface{}{
			"group_id":  groupID,
			"user_id":   userID,
			"user_type": userType,
			"action":    "join_request",
		}
		if m.UserType == UserTypeConsumer {
			ws.GlobalCrossHub.SendToConsumer(m.UserID, ws.Message{Type: "group_event", Payload: payload})
		} else {
			ws.GlobalCrossHub.SendToUser(m.UserID, ws.Message{Type: "group_event", Payload: payload})
		}
	}
}

// ==================== Handle Join Request ====================

func HandleJoinRequest(ctx context.Context, operatorID, operatorType string, p *HandleJoinRequestParam) {
	var req imModel.GroupJoinRequest
	if err := db.DB.WithContext(ctx).First(&req, "id = ? AND status = ?", p.RequestID, "pending").Error; err != nil {
		panic(exception.NewBusinessError("申请不存在或已处理", 400))
	}

	// Verify operator is admin/owner of the group
	checkOwnerOrAdmin(ctx, req.GroupID, operatorID, operatorType)

	now := time.Now()
	updates := map[string]interface{}{
		"status":     p.Action,
		"handled_by": operatorID,
		"updated_at": &now,
	}
	if err := db.DB.WithContext(ctx).Model(&imModel.GroupJoinRequest{}).Where("id = ?", p.RequestID).Updates(updates).Error; err != nil {
		panic(exception.NewBusinessError("处理失败", 500))
	}

	if p.Action == "approved" {
		// Actually add member
		joinedAt := time.Now()
		if err := db.DB.WithContext(ctx).Create(&imModel.GroupMember{
			ID:       utils.GenerateID(),
			GroupID:  req.GroupID,
			UserID:   req.UserID,
			UserType: req.UserType,
			Role:     RoleMember,
			JoinedAt: &joinedAt,
			Status:   MemberActive,
		}).Error; err != nil {
			panic(exception.NewBusinessError("添加成员失败", 500))
		}

		// System message
		extra := imModel.MsgExtraSystem{Action: "join", UserID: req.UserID, UserType: req.UserType}
		extraBytes, _ := json.Marshal(extra)
		db.DB.WithContext(ctx).Create(&imModel.GroupMessage{
			ID: utils.GenerateID(), GroupID: req.GroupID,
			SenderID: req.UserID, SenderType: req.UserType,
			Content: "加入了群聊", Extra: string(extraBytes),
			MsgType: imModel.MsgTypeSystem, CreatedAt: &joinedAt,
		})
	}

	// Notify applicant
	msg := map[string]interface{}{
		"group_id": req.GroupID,
		"status":   p.Action,
		"action":   "join_request_result",
	}
	if req.UserType == UserTypeConsumer {
		ws.GlobalCrossHub.SendToConsumer(req.UserID, ws.Message{Type: "group_event", Payload: msg})
	} else {
		ws.GlobalCrossHub.SendToUser(req.UserID, ws.Message{Type: "group_event", Payload: msg})
	}
}

// ==================== Pending Join Requests ====================

func PendingJoinRequests(ctx context.Context, groupID string) []imModel.GroupJoinRequest {
	var requests []imModel.GroupJoinRequest
	db.DB.WithContext(ctx).Model(&imModel.GroupJoinRequest{}).
		Where("group_id = ? AND status = ?", groupID, "pending").
		Order("created_at DESC").Find(&requests)
	return requests
}

// ==================== Transfer Owner ====================

func TransferOwner(ctx context.Context, operatorID string, p *TransferOwnerParam) {
	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", p.GroupID).Error; err != nil {
		panic(exception.NewBusinessError("群不存在", 400))
	}
	if group.OwnerID != operatorID {
		panic(exception.NewBusinessError("仅群主可转让群", 403))
	}

	// Verify new owner is a member
	var newOwner imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			p.GroupID, p.NewOwnerID, p.NewOwnerType, MemberActive).First(&newOwner).Error; err != nil {
		panic(exception.NewBusinessError("新群主不在群中", 400))
	}

	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Update group owner
	tx.Model(&imModel.Group{}).Where("id = ?", p.GroupID).
		Updates(map[string]interface{}{"owner_id": p.NewOwnerID, "owner_type": p.NewOwnerType})

	// Demote old owner to admin
	tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, operatorID, group.OwnerType).
		Update("role", RoleAdmin)

	// Promote new owner
	tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.NewOwnerID, p.NewOwnerType).
		Update("role", RoleOwner)

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("转让失败", 500))
	}
}

// ==================== Set Member Nickname ====================

func SetMemberNickname(ctx context.Context, operatorID, operatorType string, p *SetNicknameParam) {
	// Verify operator is admin/owner
	checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)

	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
		Update("nickname", p.Nickname).Error; err != nil {
		panic(exception.NewBusinessError("设置昵称失败", 500))
	}
}

// ==================== My Group Conversations (for unified conversation list) ====================

func MyGroupConversations(userID, userType string) []*ConversationVO {
	if userID == "" {
		return nil
	}

	var members []imModel.GroupMember
	db.DB.
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
	db.DB.Where("id IN ? AND status = ?", groupIDs, GroupNormal).Find(&groups)

	type cnt struct{ GroupID string; Count int }
	var counts []cnt
	db.DB.Model(&imModel.GroupMember{}).
		Select("group_id, COUNT(*) as count").
		Where("group_id IN ? AND status = ?", groupIDs, MemberActive).
		Group("group_id").Scan(&counts)
	countMap := make(map[string]int, len(counts))
	for _, c := range counts {
		countMap[c.GroupID] = c.Count
	}

	// Last message per group
	type lm struct{ GroupID string; Content string; CreatedAt time.Time }
	var lastMsgs []lm
	lastSubQ := db.DB.Table("im_group_message").
		Select("group_id, MAX(created_at) as max_ct").
		Where("group_id IN ?", groupIDs).
		Group("group_id")
	db.DB.Table("im_group_message g2").
		Select("g2.group_id, g2.content, g2.created_at").
		Joins("INNER JOIN (?) g1 ON g1.group_id = g2.group_id AND g1.max_ct = g2.created_at", lastSubQ).
		Scan(&lastMsgs)
	lastMap := make(map[string]lm, len(lastMsgs))
	for _, l := range lastMsgs {
		lastMap[l.GroupID] = l
	}

	// Unread count per group
	type uc struct{ GroupID string; Count int64 }
	var unreads []uc
	readSubQ := db.DB.Table("im_group_message_read").
		Select("group_id, MAX(read_at) as max_read").
		Where("user_id = ? AND user_type = ?", userID, userType).
		Group("group_id")
	db.DB.Table("im_group_message gm").
		Select("gm.group_id, COUNT(*) as count").
		Joins("LEFT JOIN (?) gr ON gr.group_id = gm.group_id", readSubQ).
		Where("gm.group_id IN ?", groupIDs).
		Where("gm.created_at > COALESCE(gr.max_read, '1970-01-01 00:00:00')").
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
// resolveFileURL constructs a full HTTP URL from file message content and extra.
// Used by the group package to avoid circular imports with message.ResolveFileURL.
func resolveFileURL(content, extra string) string {
	if strings.HasPrefix(content, "http") {
		return content
	}
	if content == "" {
		return ""
	}
	engine := "LOCAL"
	bucket := "DEFAULT"
	if extra != "" {
		var meta struct {
			Engine string `json:"engine"`
			Bucket string `json:"bucket"`
		}
		if err := json.Unmarshal([]byte(extra), &meta); err == nil {
			if meta.Engine != "" {
				engine = meta.Engine
			}
			if meta.Bucket != "" {
				bucket = meta.Bucket
			}
		}
	}
	return storage.GetURL(engine, bucket, content)
}
