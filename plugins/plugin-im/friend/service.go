package friend

import (
	"context"
	"time"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"
	"hei-gin/plugins/plugin-im/ws"
	imModel "hei-gin/plugins/plugin-im/model"

	sysUser "hei-gin/plugins/plugin-sys/user"
	cliUser "hei-gin/plugins/plugin-client/user"
)

// ==================== Send Friend Request ====================

func SendRequest(ctx context.Context, senderID, senderType string, p *SendRequestParam) {
	if senderID == "" || p.ReceiverID == "" || p.ReceiverType == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}
	if senderID == p.ReceiverID && senderType == p.ReceiverType {
		panic(exception.NewBusinessError("不能添加自己为好友", 400))
	}

	// Check existing friendship
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?) OR "+
			"(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?)",
			senderID, senderType, p.ReceiverID, p.ReceiverType,
			p.ReceiverID, p.ReceiverType, senderID, senderType).
		Count(&count)
	if count > 0 {
		panic(exception.NewBusinessError("已经是好友了", 400))
	}

	// Check pending request
	var existing int64
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("sender_id = ? AND sender_type = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
			senderID, senderType, p.ReceiverID, p.ReceiverType, "pending").
		Count(&existing)
	if existing > 0 {
		panic(exception.NewBusinessError("已发送过好友请求，请等待回复", 400))
	}

	now := time.Now()
	req := &imModel.FriendRequest{
		ID:           utils.GenerateID(),
		SenderID:     senderID,
		SenderType:   senderType,
		ReceiverID:   p.ReceiverID,
		ReceiverType: p.ReceiverType,
		Remark:       p.Remark,
		Status:       "pending",
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	if err := db.DB.WithContext(ctx).Create(req).Error; err != nil {
		panic(exception.NewBusinessError("发送好友请求失败", 500))
	}

	// WS push to receiver
	payload := map[string]interface{}{
		"request_id":  req.ID,
		"sender_id":   senderID,
		"sender_type": senderType,
		"remark":      p.Remark,
		"action":      "friend_request",
	}
	msg := ws.Message{Type: "friend_request", Payload: payload}
	if p.ReceiverType == "CONSUMER" {
		ws.GlobalCrossHub.SendToConsumer(p.ReceiverID, msg)
	} else {
		ws.GlobalCrossHub.SendToUser(p.ReceiverID, msg)
	}
}

// ==================== Accept Friend Request ====================

func AcceptRequest(ctx context.Context, userID, userType string, p *HandleRequestParam) {
	if userID == "" || p.RequestID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	var req imModel.FriendRequest
	if err := db.DB.WithContext(ctx).First(&req, "id = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
		p.RequestID, userID, userType, "pending").Error; err != nil {
		panic(exception.NewBusinessError("好友请求不存在或已处理", 400))
	}

	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	now := time.Now()
	tx.Model(&req).Updates(map[string]interface{}{"status": "accepted", "updated_at": now})

	// Create bidirectional friendship records
	pair1 := imModel.Friendship{
		ID: utils.GenerateID(), UserID: req.ReceiverID, UserType: req.ReceiverType,
		FriendID: req.SenderID, FriendType: req.SenderType, CreatedAt: &now,
	}
	pair2 := imModel.Friendship{
		ID: utils.GenerateID(), UserID: req.SenderID, UserType: req.SenderType,
		FriendID: req.ReceiverID, FriendType: req.ReceiverType, CreatedAt: &now,
	}
	if err := tx.Create(&pair1).Error; err != nil {
		panic(exception.NewBusinessError("添加好友失败", 500))
	}
	if err := tx.Create(&pair2).Error; err != nil {
		panic(exception.NewBusinessError("添加好友失败", 500))
	}

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("添加好友失败", 500))
	}

	// WS push to sender
	payload := map[string]interface{}{
		"request_id":  req.ID,
		"receiver_id":   userID,
		"receiver_type": userType,
		"action":      "friend_request_accepted",
	}
	msg := ws.Message{Type: "friend_request", Payload: payload}
	if req.SenderType == "CONSUMER" {
		ws.GlobalCrossHub.SendToConsumer(req.SenderID, msg)
	} else {
		ws.GlobalCrossHub.SendToUser(req.SenderID, msg)
	}
}

// ==================== Reject Friend Request ====================

func RejectRequest(ctx context.Context, userID, userType string, p *HandleRequestParam) {
	if userID == "" || p.RequestID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	result := db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("id = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
			p.RequestID, userID, userType, "pending").
		Updates(map[string]interface{}{"status": "rejected", "updated_at": time.Now()})
	if result.RowsAffected == 0 {
		panic(exception.NewBusinessError("好友请求不存在或已处理", 400))
	}
}

// ==================== Friend List ====================

func FriendList(ctx context.Context, userID, userType string) []FriendVO {
	var friendships []imModel.Friendship
	db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&friendships)

	if len(friendships) == 0 {
		return nil
	}

	businessIDs, consumerIDs := []string{}, []string{}
	for _, f := range friendships {
		switch f.FriendType {
		case "BUSINESS":
			businessIDs = append(businessIDs, f.FriendID)
		case "CONSUMER":
			consumerIDs = append(consumerIDs, f.FriendID)
		}
	}

	nicknameMap := make(map[string]string, len(friendships))
	avatarMap := make(map[string]string, len(friendships))

	if len(businessIDs) > 0 {
		var users []sysUser.SysUser
		db.DB.WithContext(ctx).Model(&sysUser.SysUser{}).Where("id IN ?", businessIDs).Find(&users)
		for _, u := range users {
			k := "BUSINESS:" + u.ID
			if u.Nickname != nil {
				nicknameMap[k] = *u.Nickname
			}
			if u.Avatar != nil {
				avatarMap[k] = *u.Avatar
			}
		}
	}
	if len(consumerIDs) > 0 {
		var users []cliUser.ClientUser
		db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).Where("id IN ?", consumerIDs).Find(&users)
		for _, u := range users {
			k := "CONSUMER:" + u.ID
			if u.Nickname != nil {
				nicknameMap[k] = *u.Nickname
			}
			if u.Avatar != nil {
				avatarMap[k] = *u.Avatar
			}
		}
	}

	result := make([]FriendVO, 0, len(friendships))
	for _, f := range friendships {
		k := f.FriendType + ":" + f.FriendID
		result = append(result, FriendVO{
			UserID:   f.FriendID,
			UserType: f.FriendType,
			Nickname: nicknameMap[k],
			Avatar:   avatarMap[k],
			Remark:   f.Remark,
			AddedAt:  utils.FormatDateTimePtr(f.CreatedAt),
		})
	}
	return result
}

// ==================== Friend Requests (incoming + outgoing) ====================

func PendingRequests(ctx context.Context, userID, userType string) ([]FriendRequestVO, []FriendRequestVO) {
	var incoming []imModel.FriendRequest
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("receiver_id = ? AND receiver_type = ? AND status = ?", userID, userType, "pending").
		Order("created_at DESC").Find(&incoming)

	var outgoing []imModel.FriendRequest
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("sender_id = ? AND sender_type = ? AND status = ?", userID, userType, "pending").
		Order("created_at DESC").Find(&outgoing)

	return toVOList(incoming), toVOList(outgoing)
}

// ==================== Remove Friend ====================

func RemoveFriend(ctx context.Context, userID, userType, friendID, friendType string) {
	if userID == "" || friendID == "" || friendType == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	r1 := tx.Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
		userID, userType, friendID, friendType).Delete(&imModel.Friendship{})
	r2 := tx.Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
		friendID, friendType, userID, userType).Delete(&imModel.Friendship{})

	if r1.RowsAffected == 0 && r2.RowsAffected == 0 {
		panic(exception.NewBusinessError("好友关系不存在", 400))
	}

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("删除好友失败", 500))
	}
}

// ==================== Search Users ====================

type SearchResult struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func SearchUsers(ctx context.Context, keyword string, limit int) []SearchResult {
	if keyword == "" || limit < 1 {
		return nil
	}
	if limit > 50 {
		limit = 50
	}

	like := "%" + keyword + "%"
	results := make([]SearchResult, 0, limit)

	var sysUsers []sysUser.SysUser
	db.DB.WithContext(ctx).Model(&sysUser.SysUser{}).
		Where("username LIKE ? OR nickname LIKE ?", like, like).
		Limit(limit).Find(&sysUsers)
	for _, u := range sysUsers {
		nickname := ""
		if u.Nickname != nil {
			nickname = *u.Nickname
		}
		avatar := ""
		if u.Avatar != nil {
			avatar = *u.Avatar
		}
		results = append(results, SearchResult{
			UserID: u.ID, UserType: "BUSINESS",
			Nickname: nickname, Avatar: avatar,
		})
	}

	if len(results) < limit {
		remaining := limit - len(results)
		var cliUsers []cliUser.ClientUser
		db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).
			Where("username LIKE ? OR nickname LIKE ?", like, like).
			Limit(remaining).Find(&cliUsers)
		for _, u := range cliUsers {
			nickname := ""
			if u.Nickname != nil {
				nickname = *u.Nickname
			}
			avatar := ""
			if u.Avatar != nil {
				avatar = *u.Avatar
			}
			results = append(results, SearchResult{
				UserID: u.ID, UserType: "CONSUMER",
				Nickname: nickname, Avatar: avatar,
			})
		}
	}

	return results
}

// ==================== Block / Unblock / BlockList ====================

func BlockUser(ctx context.Context, userID, userType, blockedID, blockedType string) {
	if userID == "" || blockedID == "" || blockedType == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}
	if userID == blockedID && userType == blockedType {
		panic(exception.NewBusinessError("不能拉黑自己", 400))
	}

	// Check if already blocked
	var existing int64
	db.DB.WithContext(ctx).Model(&imModel.FriendBlock{}).
		Where("user_id = ? AND user_type = ? AND blocked_id = ? AND blocked_type = ?",
			userID, userType, blockedID, blockedType).Count(&existing)
	if existing > 0 {
		panic(exception.NewBusinessError("已经拉黑了该用户", 400))
	}

	now := time.Now()
	if err := db.DB.WithContext(ctx).Create(&imModel.FriendBlock{
		ID:          utils.GenerateID(),
		UserID:      userID,
		UserType:    userType,
		BlockedID:   blockedID,
		BlockedType: blockedType,
		CreatedAt:   &now,
	}).Error; err != nil {
		panic(exception.NewBusinessError("拉黑失败", 500))
	}

	// Also remove friendship if exists
	db.DB.WithContext(ctx).
		Where("(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?) OR "+
			"(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?)",
			userID, userType, blockedID, blockedType,
			blockedID, blockedType, userID, userType).
		Delete(&imModel.Friendship{})
}

func UnblockUser(ctx context.Context, userID, userType, blockedID, blockedType string) {
	if userID == "" || blockedID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}

	result := db.DB.WithContext(ctx).
		Where("user_id = ? AND user_type = ? AND blocked_id = ? AND blocked_type = ?",
			userID, userType, blockedID, blockedType).
		Delete(&imModel.FriendBlock{})
	if result.RowsAffected == 0 {
		panic(exception.NewBusinessError("未拉黑该用户", 400))
	}
}

func BlockList(ctx context.Context, userID, userType string) []BlockVO {
	var blocks []imModel.FriendBlock
	db.DB.WithContext(ctx).Model(&imModel.FriendBlock{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&blocks)

	result := make([]BlockVO, len(blocks))
	for i, b := range blocks {
		result[i] = BlockVO{
			BlockedID:   b.BlockedID,
			BlockedType: b.BlockedType,
			CreatedAt:   utils.FormatDateTimePtr(b.CreatedAt),
		}
	}
	return result
}

// ==================== Update Friend Remark ====================

func UpdateFriendRemark(ctx context.Context, userID, userType, friendID, friendType, remark string) {
	if userID == "" || friendID == "" {
		panic(exception.NewBusinessError("参数错误", 400))
	}
	if err := db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
			userID, userType, friendID, friendType).
		Update("remark", remark).Error; err != nil {
		panic(exception.NewBusinessError("修改备注失败", 500))
	}
}

// ==================== Helpers ====================

func toVOList(requests []imModel.FriendRequest) []FriendRequestVO {
	r := make([]FriendRequestVO, len(requests))
	for i, req := range requests {
		r[i] = FriendRequestVO{
			ID: req.ID, SenderID: req.SenderID, SenderType: req.SenderType,
			ReceiverID: req.ReceiverID, ReceiverType: req.ReceiverType,
			Remark: req.Remark, Status: req.Status,
			CreatedAt: utils.FormatDateTimePtr(req.CreatedAt),
		}
	}
	return r
}
