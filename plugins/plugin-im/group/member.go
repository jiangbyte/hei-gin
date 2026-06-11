package group

import (
	"encoding/json"
	"fmt"
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

// ── Member management ─────────────────────────────────────────────

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
			Updates(map[string]interface{}{"owner_id": p.UserID, "owner_type": p.UserType}).Error; err != nil {
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
