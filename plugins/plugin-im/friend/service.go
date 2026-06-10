package friend

import (
	imModel "hei-gin/plugins/plugin-im/model"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
	"hei-gin/plugins/plugin-im/ws"

	sysUser "hei-gin/plugins/plugin-sys/user"
	cliUser "hei-gin/plugins/plugin-client/user"

	"github.com/gin-gonic/gin"
)

func getLoginID(c *gin.Context, userType string) string {
	if userType == string(enums.LoginTypeConsumer) {
		return auth.Consumer.GetLoginID(c)
	}
	return auth.GetLoginID(c)
}

// ==================== FriendSendRequest ====================

func FriendSendRequest(c *gin.Context, userType string, p *SendRequestParam) {
	ctx := c.Request.Context()
	senderID := getLoginID(c, userType)

	if p.ReceiverID == "" || p.ReceiverType == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}
	if senderID == p.ReceiverID && userType == p.ReceiverType {
		result.WriteError(c, exception.NewBusinessError("不能添加自己为好友", 400))
		return
	}

	// Check existing friendship
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?) OR "+
			"(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?)",
			senderID, userType, p.ReceiverID, p.ReceiverType,
			p.ReceiverID, p.ReceiverType, senderID, userType).
		Count(&count)
	if count > 0 {
		result.WriteError(c, exception.NewBusinessError("已经是好友了", 400))
		return
	}

	// Check pending request
	var existing int64
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("sender_id = ? AND sender_type = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
			senderID, userType, p.ReceiverID, p.ReceiverType, "pending").
		Count(&existing)
	if existing > 0 {
		result.WriteError(c, exception.NewBusinessError("已发送过好友请求，请等待回复", 400))
		return
	}

	req := &imModel.FriendRequest{
		ID:           utils.GenerateID(),
		SenderID:     senderID,
		SenderType:   userType,
		ReceiverID:   p.ReceiverID,
		ReceiverType: p.ReceiverType,
		Remark:       p.Remark,
		Status:       "pending",
	}
	if err := db.DB.WithContext(ctx).Create(req).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("发送好友请求失败: "+err.Error(), 500))
		return
	}

	// WS push to receiver
	payload := map[string]interface{}{
		"request_id":  req.ID,
		"sender_id":   senderID,
		"sender_type": userType,
		"remark":      p.Remark,
		"action":      "friend_request",
	}
	msg := ws.Message{Type: "friend_request", Payload: payload}
	if p.ReceiverType == string(enums.LoginTypeConsumer) {
		ws.GlobalCrossHub.SendToConsumer(p.ReceiverID, msg)
	} else {
		ws.GlobalCrossHub.SendToUser(p.ReceiverID, msg)
	}
}

// ==================== FriendAcceptRequest ====================

func FriendAcceptRequest(c *gin.Context, userType string, p *HandleRequestParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	if p.RequestID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	var req imModel.FriendRequest
	if err := db.DB.WithContext(ctx).First(&req, "id = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
		p.RequestID, userID, userType, "pending").Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("好友请求不存在或已处理", 400))
		return
	}

	tx := db.DB.WithContext(ctx).Begin()

	if err := tx.Model(&req).Update("status", "accepted").Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("添加好友失败: "+err.Error(), 500))
		return
	}

	// Create bidirectional friendship records
	pair1 := imModel.Friendship{
		ID: utils.GenerateID(), UserID: req.ReceiverID, UserType: req.ReceiverType,
		FriendID: req.SenderID, FriendType: req.SenderType,
	}
	pair2 := imModel.Friendship{
		ID: utils.GenerateID(), UserID: req.SenderID, UserType: req.SenderType,
		FriendID: req.ReceiverID, FriendType: req.ReceiverType,
	}
	if err := tx.Create(&pair1).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("添加好友失败: "+err.Error(), 500))
		return
	}
	if err := tx.Create(&pair2).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("添加好友失败: "+err.Error(), 500))
		return
	}

	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加好友失败: "+err.Error(), 500))
		return
	}

	// WS push to sender
	payload := map[string]interface{}{
		"request_id":    req.ID,
		"receiver_id":   userID,
		"receiver_type": userType,
		"action":        "friend_request_accepted",
	}
	msg := ws.Message{Type: "friend_request", Payload: payload}
	if req.SenderType == string(enums.LoginTypeConsumer) {
		ws.GlobalCrossHub.SendToConsumer(req.SenderID, msg)
	} else {
		ws.GlobalCrossHub.SendToUser(req.SenderID, msg)
	}
}

// ==================== FriendRejectRequest ====================

func FriendRejectRequest(c *gin.Context, userType string, p *HandleRequestParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	if p.RequestID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	res := db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("id = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
			p.RequestID, userID, userType, "pending").
		Update("status", "rejected")
	if res.RowsAffected == 0 {
		result.WriteError(c, exception.NewBusinessError("好友请求不存在或已处理", 400))
		return
	}
	if res.Error != nil {
		result.WriteError(c, exception.NewBusinessError("拒绝好友请求失败: "+res.Error.Error(), 500))
		return
	}
}

// ==================== FriendList ====================

func FriendList(c *gin.Context, userType string) []FriendVO {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	var friendships []imModel.Friendship
	db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&friendships)

	if len(friendships) == 0 {
		return nil
	}

	var businessIDs, consumerIDs []string
	for _, f := range friendships {
		switch f.FriendType {
		case string(enums.LoginTypeBusiness):
			businessIDs = append(businessIDs, f.FriendID)
		case string(enums.LoginTypeConsumer):
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

	vos := make([]FriendVO, 0, len(friendships))
	for _, f := range friendships {
		vo := *ImFriendshipToFriendVO(&f)
		k := f.FriendType + ":" + f.FriendID
		vo.Nickname = nicknameMap[k]
		vo.Avatar = avatarMap[k]
		vos = append(vos, vo)
	}
	return vos
}

// ==================== FriendPendingRequests ====================

func FriendPendingRequests(c *gin.Context, userType string) (incoming, outgoing []FriendRequestVO) {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	var inRecs []imModel.FriendRequest
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("receiver_id = ? AND receiver_type = ? AND status = ?", userID, userType, "pending").
		Order("created_at DESC").Find(&inRecs)

	var outRecs []imModel.FriendRequest
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("sender_id = ? AND sender_type = ? AND status = ?", userID, userType, "pending").
		Order("created_at DESC").Find(&outRecs)

	incoming = make([]FriendRequestVO, len(inRecs))
	for i, r := range inRecs {
		incoming[i] = *ImFriendRequestToFriendRequestVO(&r)
	}
	outgoing = make([]FriendRequestVO, len(outRecs))
	for i, r := range outRecs {
		outgoing[i] = *ImFriendRequestToFriendRequestVO(&r)
	}
	return
}

// ==================== FriendRemove ====================

func FriendRemove(c *gin.Context, userType string, p *RemoveFriendParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	if p.FriendID == "" || p.FriendType == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	tx := db.DB.WithContext(ctx).Begin()

	r1 := tx.Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
		userID, userType, p.FriendID, p.FriendType).Delete(&imModel.Friendship{})
	if r1.Error != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除好友失败: "+r1.Error.Error(), 500))
		return
	}

	r2 := tx.Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
		p.FriendID, p.FriendType, userID, userType).Delete(&imModel.Friendship{})
	if r2.Error != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除好友失败: "+r2.Error.Error(), 500))
		return
	}

	if r1.RowsAffected == 0 && r2.RowsAffected == 0 {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("好友关系不存在", 400))
		return
	}

	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除好友失败: "+err.Error(), 500))
		return
	}

	// WS push
	payload := map[string]interface{}{
		"friend_id":   p.FriendID,
		"friend_type": p.FriendType,
		"action":      "friend_removed",
	}
	msg := ws.Message{Type: "friend", Payload: payload}
	ws.GlobalCrossHub.SendToUser(userID, msg)
}

// ==================== FriendSearch ====================

func FriendSearch(c *gin.Context, keyword string, limit int) []SearchResult {
	ctx := c.Request.Context()
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
			UserID: u.ID, UserType: string(enums.LoginTypeBusiness),
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
				UserID: u.ID, UserType: string(enums.LoginTypeConsumer),
				Nickname: nickname, Avatar: avatar,
			})
		}
	}

	return results
}

// ==================== FriendBlock ====================

func FriendBlock(c *gin.Context, userType string, p *BlockParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	if p.BlockedID == "" || p.BlockedType == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}
	if userID == p.BlockedID && userType == p.BlockedType {
		result.WriteError(c, exception.NewBusinessError("不能拉黑自己", 400))
		return
	}

	// Check if already blocked
	var existing int64
	db.DB.WithContext(ctx).Model(&imModel.FriendBlock{}).
		Where("user_id = ? AND user_type = ? AND blocked_id = ? AND blocked_type = ?",
			userID, userType, p.BlockedID, p.BlockedType).Count(&existing)
	if existing > 0 {
		result.WriteError(c, exception.NewBusinessError("已经拉黑了该用户", 400))
		return
	}

	if err := db.DB.WithContext(ctx).Create(&imModel.FriendBlock{
		ID:          utils.GenerateID(),
		UserID:      userID,
		UserType:    userType,
		BlockedID:   p.BlockedID,
		BlockedType: p.BlockedType,
	}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("拉黑失败: "+err.Error(), 500))
		return
	}

	// Also remove friendship if exists
	db.DB.WithContext(ctx).
		Where("(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?) OR "+
			"(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?)",
			userID, userType, p.BlockedID, p.BlockedType,
			p.BlockedID, p.BlockedType, userID, userType).
		Delete(&imModel.Friendship{})
}

// ==================== FriendUnblock ====================

func FriendUnblock(c *gin.Context, userType string, p *BlockParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	if p.BlockedID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	res := db.DB.WithContext(ctx).
		Where("user_id = ? AND user_type = ? AND blocked_id = ? AND blocked_type = ?",
			userID, userType, p.BlockedID, p.BlockedType).
		Delete(&imModel.FriendBlock{})
	if res.RowsAffected == 0 {
		result.WriteError(c, exception.NewBusinessError("未拉黑该用户", 400))
		return
	}
}

// ==================== FriendBlockList ====================

func FriendBlockList(c *gin.Context, userType string) []BlockVO {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	var blocks []imModel.FriendBlock
	db.DB.WithContext(ctx).Model(&imModel.FriendBlock{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&blocks)

	vos := make([]BlockVO, len(blocks))
	for i, b := range blocks {
		vos[i] = *ImFriendBlockToBlockVO(&b)
	}
	return vos
}

// ==================== FriendUpdateRemark ====================

func FriendUpdateRemark(c *gin.Context, userType string, p *RemarkParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c, userType)

	if p.FriendID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}
	if err := db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
			userID, userType, p.FriendID, p.FriendType).
		Update("remark", p.Remark).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("修改备注失败: "+err.Error(), 500))
		return
	}
}
