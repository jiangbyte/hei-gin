package group

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/clause"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/storage"

	imModel "hei-gin/plugins/plugin-im/model"
	ws "hei-gin/plugins/plugin-im/ws"

	"github.com/gin-gonic/gin"
)

func getLoginID(c *gin.Context) string {
	path := c.Request.URL.Path
	if len(path) > 8 && path[:8] == "/api/v1/c" {
		return auth.Consumer.GetLoginID(c)
	}
	return auth.GetLoginID(c)
}

func getUserType(c *gin.Context) string {
	path := c.Request.URL.Path
	if len(path) > 8 && path[:8] == "/api/v1/c" {
		return string(enums.LoginTypeConsumer)
	}
	return string(enums.LoginTypeBusiness)
}

// ==================== GroupCreate ====================

func GroupCreate(c *gin.Context) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	var p CreateParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.WriteError(c, exception.NewBusinessError("参数错误: "+err.Error(), 400))
		return
	}

	if p.Name == "" {
		result.WriteError(c, exception.NewBusinessError("群名称不能为空", 400))
		return
	}
	if len(p.Name) > 100 {
		result.WriteError(c, exception.NewBusinessError("群名称不能超过100个字符", 400))
		return
	}
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return
	}

	groupType := GroupTypeMixed
	if userType == string(enums.LoginTypeConsumer) {
		groupType = GroupTypeConsumerOnly
	}

	group := imModel.Group{
		ID:         utils.GenerateID(),
		Name:       p.Name,
		Avatar:     p.Avatar,
		OwnerID:    userID,
		OwnerType:  userType,
		GroupType:  groupType,
		MaxMembers: 200,
		Status:     GroupNormal,
	}

	tx := db.DB.WithContext(ctx).Begin()

	if err := tx.Create(&group).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("创建群失败: "+err.Error(), 500))
		return
	}

	ownerMember := imModel.GroupMember{
		ID: utils.GenerateID(), GroupID: group.ID,
		UserID: userID, UserType: userType,
		Role: RoleOwner, Status: MemberActive,
	}
	if err := tx.Create(&ownerMember).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("添加群主失败: "+err.Error(), 500))
		return
	}

	if len(p.MemberIDs) > 0 {
		if err := validateMemberType(groupType, p.MemberType); err != nil {
			tx.Rollback()
			result.WriteError(c, err)
			return
		}

		var existingCount int64
		tx.Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id IN ? AND user_type = ? AND status = ?",
				group.ID, p.MemberIDs, p.MemberType, MemberActive).
			Count(&existingCount)
		if existingCount > 0 {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError("部分成员已在群中", 400))
			return
		}

		var currentCount int64
		tx.Model(&imModel.GroupMember{}).Where("group_id = ? AND status = ?", group.ID, MemberActive).Count(&currentCount)
		if int(currentCount)+len(p.MemberIDs) > group.MaxMembers {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError(fmt.Sprintf("群成员数量不能超过%d人", group.MaxMembers), 400))
			return
		}

		batch := make([]imModel.GroupMember, 0, len(p.MemberIDs))
		sysBatch := make([]imModel.GroupMessage, 0, len(p.MemberIDs))
		for _, uid := range p.MemberIDs {
			if uid == userID {
				continue
			}
			batch = append(batch, imModel.GroupMember{
				ID: utils.GenerateID(), GroupID: group.ID,
				UserID: uid, UserType: p.MemberType,
				Role: RoleMember, Status: MemberActive,
			})
			extra := imModel.MsgExtraSystem{Action: "join", UserID: uid, UserType: p.MemberType}
			extraBytes, _ := json.Marshal(extra)
			sysBatch = append(sysBatch, imModel.GroupMessage{
				ID: utils.GenerateID(), GroupID: group.ID,
				SenderID: userID, SenderType: userType,
				Content: "欢迎加入群聊", Extra: string(extraBytes),
				MsgType: imModel.MsgTypeSystem,
			})
		}
		if len(batch) > 0 {
			if err := tx.Create(&batch).Error; err != nil {
				tx.Rollback()
				result.WriteError(c, exception.NewBusinessError("邀请成员失败: "+err.Error(), 500))
				return
			}
		}
		if len(sysBatch) > 0 {
			if err := tx.Create(&sysBatch).Error; err != nil {
				tx.Rollback()
				result.WriteError(c, exception.NewBusinessError("发送系统消息失败: "+err.Error(), 500))
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("创建群失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupUpdate ====================

func GroupUpdate(c *gin.Context, p *UpdateParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}
	_, member, err := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
	if err != nil {
		result.WriteError(c, err)
		return
	}

	updates := map[string]interface{}{}
	if p.Name != nil {
		if len(*p.Name) > 100 {
			result.WriteError(c, exception.NewBusinessError("群名称不能超过100个字符", 400))
			return
		}
		updates["name"] = *p.Name
	}
	if p.Avatar != nil {
		updates["avatar"] = *p.Avatar
	}
	if p.Notice != nil && member.Role == RoleOwner {
		updates["notice"] = *p.Notice
	}

	if len(updates) == 0 {
		return
	}

	if err := db.DB.WithContext(ctx).Model(&imModel.Group{}).Where("id = ?", p.GroupID).Updates(updates).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("修改群信息失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupDissolve ====================

func GroupDissolve(c *gin.Context, p *DissolveParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)

	if p.GroupID == "" || operatorID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", p.GroupID).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("群不存在", 400))
		return
	}
	if group.OwnerID != operatorID {
		result.WriteError(c, exception.NewBusinessError("仅群主可解散群", 403))
		return
	}

	tx := db.DB.WithContext(ctx).Begin()

	if err := tx.Model(&imModel.Group{}).Where("id = ?", p.GroupID).
		Update("status", GroupDissolved).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("解散群失败: "+err.Error(), 500))
		return
	}

	if err := tx.Model(&imModel.GroupMember{}).Where("group_id = ? AND status = ?", p.GroupID, MemberActive).
		Update("status", MemberLeft).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("解散群失败: "+err.Error(), 500))
		return
	}

	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("解散群失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupInvite ====================

func GroupInvite(c *gin.Context, p *InviteParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if len(p.UserIDs) == 0 {
		return
	}
	if p.GroupID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", p.GroupID).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("群不存在", 400))
		return
	}
	if group.Status != GroupNormal {
		result.WriteError(c, exception.NewBusinessError("群已解散", 400))
		return
	}
	if err := validateMemberType(group.GroupType, p.UserType); err != nil {
		result.WriteError(c, err)
		return
	}

	if _, _, err := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType); err != nil {
		result.WriteError(c, err)
		return
	}

	var existingIDs []string
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id IN ? AND user_type = ? AND status = ?",
			p.GroupID, p.UserIDs, p.UserType, MemberActive).
		Pluck("user_id", &existingIDs)
	if len(existingIDs) > 0 {
		result.WriteError(c, exception.NewBusinessError(fmt.Sprintf("用户 %v 已在群中", existingIDs), 400))
		return
	}

	var currentCount int64
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND status = ?", p.GroupID, MemberActive).Count(&currentCount)
	if int(currentCount)+len(p.UserIDs) > group.MaxMembers {
		result.WriteError(c, exception.NewBusinessError(fmt.Sprintf("群成员数量不能超过%d人", group.MaxMembers), 400))
		return
	}

	tx := db.DB.WithContext(ctx).Begin()

	batch := make([]imModel.GroupMember, 0, len(p.UserIDs))
	sysBatch := make([]imModel.GroupMessage, 0, len(p.UserIDs))
	for _, uid := range p.UserIDs {
		batch = append(batch, imModel.GroupMember{
			ID: utils.GenerateID(), GroupID: p.GroupID,
			UserID: uid, UserType: p.UserType,
			Role: RoleMember, Status: MemberActive,
		})
		extra := imModel.MsgExtraSystem{Action: "join", UserID: uid, UserType: p.UserType}
		extraBytes, _ := json.Marshal(extra)
		sysBatch = append(sysBatch, imModel.GroupMessage{
			ID: utils.GenerateID(), GroupID: p.GroupID,
			SenderID: operatorID, SenderType: operatorType,
			Content: "欢迎加入群聊", Extra: string(extraBytes),
			MsgType: imModel.MsgTypeSystem,
		})
	}
	if err := tx.Create(&batch).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("邀请成员失败: "+err.Error(), 500))
		return
	}
	if err := tx.Create(&sysBatch).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("发送系统消息失败: "+err.Error(), 500))
		return
	}

	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("邀请成员失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupJoin ====================

func GroupJoin(c *gin.Context, p *JoinOrLeaveParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if p.GroupID == "" || userID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", p.GroupID).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("群不存在", 400))
		return
	}
	if group.Status != GroupNormal {
		result.WriteError(c, exception.NewBusinessError("群已解散", 400))
		return
	}
	if err := validateMemberType(group.GroupType, userType); err != nil {
		result.WriteError(c, err)
		return
	}

	var existing int64
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			p.GroupID, userID, userType, MemberActive).Count(&existing)
	if existing > 0 {
		result.WriteError(c, exception.NewBusinessError("已在群中", 400))
		return
	}

	var pending int64
	db.DB.WithContext(ctx).Model(&imModel.GroupJoinRequest{}).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			p.GroupID, userID, userType, "pending").Count(&pending)
	if pending > 0 {
		result.WriteError(c, exception.NewBusinessError("已发送过入群申请，请等待审核", 400))
		return
	}

	if err := db.DB.WithContext(ctx).Create(&imModel.GroupJoinRequest{
		ID:       utils.GenerateID(),
		GroupID:  p.GroupID,
		UserID:   userID,
		UserType: userType,
		Status:   "pending",
	}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("申请加入失败: "+err.Error(), 500))
		return
	}

	var members []imModel.GroupMember
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND (role = ? OR role = ?) AND status = ?",
			p.GroupID, RoleOwner, RoleAdmin, MemberActive).
		Find(&members)
	for _, m := range members {
		payload := map[string]interface{}{
			"group_id":  p.GroupID,
			"user_id":   userID,
			"user_type": userType,
			"action":    "join_request",
		}
		if m.UserType == string(enums.LoginTypeConsumer) {
			ws.GlobalCrossHub.SendToConsumer(m.UserID, ws.Message{Type: "group_event", Payload: payload})
		} else {
			ws.GlobalCrossHub.SendToUser(m.UserID, ws.Message{Type: "group_event", Payload: payload})
		}
	}
}

// ==================== GroupLeave ====================

func GroupLeave(c *gin.Context, p *JoinOrLeaveParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if p.GroupID == "" || userID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	var member imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			p.GroupID, userID, userType, MemberActive).First(&member).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("不在群中", 400))
		return
	}
	if member.Role == RoleOwner {
		result.WriteError(c, exception.NewBusinessError("群主不能退群，请先转让群", 400))
		return
	}

	tx := db.DB.WithContext(ctx).Begin()

	if err := tx.Model(&member).Update("status", MemberLeft).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("退群失败: "+err.Error(), 500))
		return
	}

	extra := imModel.MsgExtraSystem{Action: "leave", UserID: userID, UserType: userType}
	extraBytes, _ := json.Marshal(extra)
	if err := tx.Create(&imModel.GroupMessage{
		ID: utils.GenerateID(), GroupID: p.GroupID,
		SenderID: userID, SenderType: userType,
		Content: "退出了群聊", Extra: string(extraBytes),
		MsgType: imModel.MsgTypeSystem,
	}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("发送系统消息失败: "+err.Error(), 500))
		return
	}

	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("退群失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupKick ====================

func GroupKick(c *gin.Context, p *KickParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" || p.UserID == "" || p.UserType == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	group, operatorMember, err := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
	if err != nil {
		result.WriteError(c, err)
		return
	}
	if p.UserID == group.OwnerID {
		result.WriteError(c, exception.NewBusinessError("不能踢出群主", 400))
		return
	}
	if operatorMember.Role == RoleAdmin {
		var target imModel.GroupMember
		if err := db.DB.WithContext(ctx).
			Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
				p.GroupID, p.UserID, p.UserType, MemberActive).First(&target).Error; err != nil {
			result.WriteError(c, exception.NewBusinessError("成员不存在或已离开", 400))
			return
		}
		if target.Role != RoleMember {
			result.WriteError(c, exception.NewBusinessError("不能踢出管理员", 403))
			return
		}
	}

	tx := db.DB.WithContext(ctx).Begin()

	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
		Update("status", MemberKicked).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("踢出失败: "+err.Error(), 500))
		return
	}

	extra := imModel.MsgExtraSystem{Action: "kick", UserID: p.UserID, UserType: p.UserType, OperatorID: operatorID}
	extraBytes, _ := json.Marshal(extra)
	if err := tx.Create(&imModel.GroupMessage{
		ID: utils.GenerateID(), GroupID: p.GroupID,
		SenderID: operatorID, SenderType: operatorType,
		Content: "被移出群聊", Extra: string(extraBytes),
		MsgType: imModel.MsgTypeSystem,
	}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("发送系统消息失败: "+err.Error(), 500))
		return
	}

	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("踢出失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupSetRole ====================

func GroupSetRole(c *gin.Context, p *SetRoleParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)

	if p.GroupID == "" || p.UserID == "" || p.Role == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", p.GroupID).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("群不存在", 400))
		return
	}
	if group.OwnerID != operatorID {
		result.WriteError(c, exception.NewBusinessError("仅群主可设置角色", 403))
		return
	}

	switch p.Role {
	case RoleAdmin:
		if err := db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
				p.GroupID, p.UserID, p.UserType, MemberActive).
			Update("role", RoleAdmin).Error; err != nil {
			result.WriteError(c, exception.NewBusinessError("设置管理员失败: "+err.Error(), 500))
			return
		}
	case RoleOwner:
		now := time.Now()
		tx := db.DB.WithContext(ctx).Begin()

		if err := tx.Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, group.OwnerID, group.OwnerType).
			Update("role", RoleAdmin).Error; err != nil {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError("转让群失败: "+err.Error(), 500))
			return
		}
		if err := tx.Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
			Update("role", RoleOwner).Error; err != nil {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError("转让群失败: "+err.Error(), 500))
			return
		}
		if err := tx.Model(&imModel.Group{}).Where("id = ?", p.GroupID).
			Updates(map[string]interface{}{"owner_id": p.UserID, "owner_type": p.UserType, "updated_at": now}).Error; err != nil {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError("转让群失败: "+err.Error(), 500))
			return
		}
		if err := tx.Commit().Error; err != nil {
			result.WriteError(c, exception.NewBusinessError("转让群失败: "+err.Error(), 500))
			return
		}
	default:
		result.WriteError(c, exception.NewBusinessError("不支持的角色", 400))
		return
	}
}

// ==================== GroupTransferOwner ====================

func GroupTransferOwner(c *gin.Context, p *TransferOwnerParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", p.GroupID).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("群不存在", 400))
		return
	}
	if group.OwnerID != operatorID {
		result.WriteError(c, exception.NewBusinessError("仅群主可转让群", 403))
		return
	}

	var newOwner imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			p.GroupID, p.NewOwnerID, p.NewOwnerType, MemberActive).First(&newOwner).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("新群主不在群中", 400))
		return
	}

	tx := db.DB.WithContext(ctx).Begin()

	if err := tx.Model(&imModel.Group{}).Where("id = ?", p.GroupID).
		Updates(map[string]interface{}{"owner_id": p.NewOwnerID, "owner_type": p.NewOwnerType}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("转让失败: "+err.Error(), 500))
		return
	}

	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, operatorID, group.OwnerType).
		Update("role", RoleAdmin).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("转让失败: "+err.Error(), 500))
		return
	}

	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.NewOwnerID, p.NewOwnerType).
		Update("role", RoleOwner).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("转让失败: "+err.Error(), 500))
		return
	}

	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("转让失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupSetMemberNickname ====================

func GroupSetMemberNickname(c *gin.Context, p *SetNicknameParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if _, _, err := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType); err != nil {
		result.WriteError(c, err)
		return
	}

	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
		Update("nickname", p.Nickname).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("设置昵称失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupSendMessage ====================

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

	var memberIDs []struct{ UserID string; UserType string }
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
		ID: utils.GenerateID(),
		UserID: userID, UserType: userType,
	}).Error
}

// ==================== GroupMuteMember ====================

func GroupMuteMember(c *gin.Context, p *MuteParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" || p.UserID == "" || p.UserType == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	group, operator, err := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
	if err != nil {
		result.WriteError(c, err)
		return
	}
	if p.UserID == group.OwnerID {
		result.WriteError(c, exception.NewBusinessError("不能禁言群主", 400))
		return
	}
	if operator.Role == RoleAdmin {
		var target imModel.GroupMember
		if err := db.DB.WithContext(ctx).
			Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
				p.GroupID, p.UserID, p.UserType, MemberActive).First(&target).Error; err != nil {
			result.WriteError(c, exception.NewBusinessError("成员不存在", 400))
			return
		}
		if target.Role != RoleMember {
			result.WriteError(c, exception.NewBusinessError("不能禁言管理员", 403))
			return
		}
	}

	duration := 60
	if p.Duration > 0 {
		duration = p.Duration
	}
	until := time.Now().Add(time.Duration(duration) * time.Minute)
	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
		Update("muted_until", &until).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("禁言失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupUnmuteMember ====================

func GroupUnmuteMember(c *gin.Context, p *UnmuteParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" || p.UserID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	if _, _, err := checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType); err != nil {
		result.WriteError(c, err)
		return
	}
	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", p.GroupID, p.UserID, p.UserType).
		Update("muted_until", nil).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("解除禁言失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupMyGroups ====================

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

	type lm struct{ GroupID string; Content string; CreatedAt time.Time }
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
		ID: utils.GenerateID(),
		UserID: userID, UserType: userType,
	}).Error
}


// ==================== Backward-compatible wrappers ====================


// Messages is a backward-compatible wrapper used by the message package.
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

// MyGroupConversations is a backward-compatible wrapper used by the message package.
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
		ID: utils.GenerateID(),
		UserID: userID, UserType: userType,
	}).Error
}

// ==================== Helpers ====================

func validateMemberType(groupType, userType string) error {
	if groupType == GroupTypeConsumerOnly && userType != string(enums.LoginTypeConsumer) {
		return exception.NewBusinessError("该群仅限C端用户", 403)
	}
	return nil
}

func checkOwnerOrAdmin(ctx context.Context, groupID, userID, userType string) (*imModel.Group, *imModel.GroupMember, error) {
	if groupID == "" || userID == "" {
		return nil, nil, exception.NewBusinessError("参数错误", 400)
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", groupID).Error; err != nil {
		return nil, nil, exception.NewBusinessError("群不存在", 400)
	}
	var member imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, MemberActive).
		First(&member).Error; err != nil {
		return nil, nil, exception.NewBusinessError("不在群中", 400)
	}
	if member.Role != RoleOwner && member.Role != RoleAdmin {
		return nil, nil, exception.NewBusinessError("无权限", 403)
	}
	return &group, &member, nil
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
