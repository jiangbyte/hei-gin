package group

import (
	"encoding/json"
	"fmt"
	"time"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	imModel "hei-gin/plugins/plugin-im/model"
	ws "hei-gin/plugins/plugin-im/ws"

	"github.com/gin-gonic/gin"
)

// ── Group management ─────────────────────────────────────────────

// ── Member management ─────────────────────────────────────────────

func (s *Service) Invite(c *gin.Context, p *InviteParam) {
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

	group, err := s.repo.FindGroupByID(ctx, p.GroupID)
	if err != nil {
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

	if _, _, err := s.checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType); err != nil {
		result.WriteError(c, err)
		return
	}

	existingIDs := s.repo.FindExistingMemberIDs(ctx, p.GroupID, p.UserIDs, p.UserType)
	if len(existingIDs) > 0 {
		result.WriteError(c, exception.NewBusinessError(fmt.Sprintf("用户 %v 已在群中", existingIDs), 400))
		return
	}

	currentCount := s.repo.CountActiveMembersByGroup(ctx, p.GroupID)
	if int(currentCount)+len(p.UserIDs) > group.MaxMembers {
		result.WriteError(c, exception.NewBusinessError(fmt.Sprintf("群成员数量不能超过%d人", group.MaxMembers), 400))
		return
	}

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
	if err := s.repo.InviteMembers(ctx, batch, sysBatch); err != nil {
		result.WriteError(c, exception.NewBusinessError("邀请成员失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupJoin ====================

func (s *Service) Join(c *gin.Context, p *JoinOrLeaveParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if p.GroupID == "" || userID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	group, err := s.repo.FindGroupByID(ctx, p.GroupID)
	if err != nil {
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

	existing := s.repo.CountMemberExists(ctx, p.GroupID, userID, userType)
	if existing > 0 {
		result.WriteError(c, exception.NewBusinessError("已在群中", 400))
		return
	}

	pending := s.repo.CountPendingJoin(ctx, p.GroupID, userID, userType)
	if pending > 0 {
		result.WriteError(c, exception.NewBusinessError("已发送过入群申请，请等待审核", 400))
		return
	}

	if err := s.repo.CreateJoinRequest(ctx, &imModel.GroupJoinRequest{
		ID:       utils.GenerateID(),
		GroupID:  p.GroupID,
		UserID:   userID,
		UserType: userType,
		Status:   "pending",
	}); err != nil {
		result.WriteError(c, exception.NewBusinessError("申请加入失败: "+err.Error(), 500))
		return
	}

	members := s.repo.ListManagers(ctx, p.GroupID)
	if runtime := ws.Runtime(); runtime != nil {
		for _, m := range members {
			payload := map[string]interface{}{
				"group_id":  p.GroupID,
				"user_id":   userID,
				"user_type": userType,
				"action":    "join_request",
			}
			if m.UserType == string(enums.LoginTypeConsumer) {
				runtime.SendToConsumer(m.UserID, ws.Message{Type: "group_event", Payload: payload})
			} else {
				runtime.SendToUser(m.UserID, ws.Message{Type: "group_event", Payload: payload})
			}
		}
	}
}

// ==================== GroupLeave ====================

func (s *Service) Leave(c *gin.Context, p *JoinOrLeaveParam) {
	ctx := c.Request.Context()
	userID := getLoginID(c)
	userType := getUserType(c)

	if p.GroupID == "" || userID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	member, err := s.repo.FindActiveMember(ctx, p.GroupID, userID, userType)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("不在群中", 400))
		return
	}
	if member.Role == RoleOwner {
		result.WriteError(c, exception.NewBusinessError("群主不能退群，请先转让群", 400))
		return
	}

	extra := imModel.MsgExtraSystem{Action: "leave", UserID: userID, UserType: userType}
	extraBytes, _ := json.Marshal(extra)
	msg := &imModel.GroupMessage{
		ID: utils.GenerateID(), GroupID: p.GroupID,
		SenderID: userID, SenderType: userType,
		Content: "退出了群聊", Extra: string(extraBytes),
		MsgType: imModel.MsgTypeSystem,
	}
	if err := s.repo.LeaveGroup(ctx, member, msg); err != nil {
		result.WriteError(c, exception.NewBusinessError("退群失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupKick ====================

func (s *Service) Kick(c *gin.Context, p *KickParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" || p.UserID == "" || p.UserType == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	group, operatorMember, err := s.checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
	if err != nil {
		result.WriteError(c, err)
		return
	}
	if p.UserID == group.OwnerID {
		result.WriteError(c, exception.NewBusinessError("不能踢出群主", 400))
		return
	}
	if operatorMember.Role == RoleAdmin {
		target, err := s.repo.FindActiveMember(ctx, p.GroupID, p.UserID, p.UserType)
		if err != nil {
			result.WriteError(c, exception.NewBusinessError("成员不存在或已离开", 400))
			return
		}
		if target.Role != RoleMember {
			result.WriteError(c, exception.NewBusinessError("不能踢出管理员", 403))
			return
		}
	}

	extra := imModel.MsgExtraSystem{Action: "kick", UserID: p.UserID, UserType: p.UserType, OperatorID: operatorID}
	extraBytes, _ := json.Marshal(extra)
	msg := &imModel.GroupMessage{
		ID: utils.GenerateID(), GroupID: p.GroupID,
		SenderID: operatorID, SenderType: operatorType,
		Content: "被移出群聊", Extra: string(extraBytes),
		MsgType: imModel.MsgTypeSystem,
	}
	if err := s.repo.KickMember(ctx, p.GroupID, p.UserID, p.UserType, msg); err != nil {
		result.WriteError(c, exception.NewBusinessError("踢出失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupSetRole ====================

func (s *Service) SetRole(c *gin.Context, p *SetRoleParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)

	if p.GroupID == "" || p.UserID == "" || p.Role == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	group, err := s.repo.FindGroupByID(ctx, p.GroupID)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("群不存在", 400))
		return
	}
	if group.OwnerID != operatorID {
		result.WriteError(c, exception.NewBusinessError("仅群主可设置角色", 403))
		return
	}

	switch p.Role {
	case RoleAdmin:
		if err := s.repo.UpdateMemberRole(ctx, p.GroupID, p.UserID, p.UserType, RoleAdmin); err != nil {
			result.WriteError(c, exception.NewBusinessError("设置管理员失败: "+err.Error(), 500))
			return
		}
	case RoleOwner:
		if err := s.repo.TransferOwnerAndRole(ctx, p.GroupID, group.OwnerID, group.OwnerType, p.UserID, p.UserType); err != nil {
			result.WriteError(c, exception.NewBusinessError("转让群失败: "+err.Error(), 500))
			return
		}
	default:
		result.WriteError(c, exception.NewBusinessError("不支持的角色", 400))
		return
	}
}

// ==================== GroupTransferOwner ====================

func (s *Service) TransferOwner(c *gin.Context, p *TransferOwnerParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)

	group, err := s.repo.FindGroupByID(ctx, p.GroupID)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("群不存在", 400))
		return
	}
	if group.OwnerID != operatorID {
		result.WriteError(c, exception.NewBusinessError("仅群主可转让群", 403))
		return
	}

	if _, err := s.repo.FindActiveMember(ctx, p.GroupID, p.NewOwnerID, p.NewOwnerType); err != nil {
		result.WriteError(c, exception.NewBusinessError("新群主不在群中", 400))
		return
	}
	if err := s.repo.TransferOwner(ctx, p.GroupID, p.NewOwnerID, p.NewOwnerType, operatorID, group.OwnerType); err != nil {
		result.WriteError(c, exception.NewBusinessError("转让失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupSetMemberNickname ====================

func (s *Service) SetMemberNickname(c *gin.Context, p *SetNicknameParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if _, _, err := s.checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType); err != nil {
		result.WriteError(c, err)
		return
	}

	if err := s.repo.UpdateMemberNickname(ctx, p.GroupID, p.UserID, p.UserType, p.Nickname); err != nil {
		result.WriteError(c, exception.NewBusinessError("设置昵称失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupSendMessage ====================

func (s *Service) MuteMember(c *gin.Context, p *MuteParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" || p.UserID == "" || p.UserType == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	group, operator, err := s.checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
	if err != nil {
		result.WriteError(c, err)
		return
	}
	if p.UserID == group.OwnerID {
		result.WriteError(c, exception.NewBusinessError("不能禁言群主", 400))
		return
	}
	if operator.Role == RoleAdmin {
		target, err := s.repo.FindActiveMember(ctx, p.GroupID, p.UserID, p.UserType)
		if err != nil {
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
	if err := s.repo.UpdateMutedUntil(ctx, p.GroupID, p.UserID, p.UserType, &until); err != nil {
		result.WriteError(c, exception.NewBusinessError("禁言失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupUnmuteMember ====================

func (s *Service) UnmuteMember(c *gin.Context, p *UnmuteParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" || p.UserID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	if _, _, err := s.checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType); err != nil {
		result.WriteError(c, err)
		return
	}
	if err := s.repo.UpdateMutedUntil(ctx, p.GroupID, p.UserID, p.UserType, nil); err != nil {
		result.WriteError(c, exception.NewBusinessError("解除禁言失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupMyGroups ====================
