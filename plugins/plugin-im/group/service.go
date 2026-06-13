package group

import (
	"encoding/json"
	"fmt"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	imModel "hei-gin/plugins/plugin-im/model"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

// ── Group management ─────────────────────────────────────────────

// ── Group lifecycle ─────────────────────────────────────────────

func (s *Service) Create(c *gin.Context) {
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
	if userType == string(auth.ConsumerID) {
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

	ownerMember := imModel.GroupMember{
		ID: utils.GenerateID(), GroupID: group.ID,
		UserID: userID, UserType: userType,
		Role: RoleOwner, Status: MemberActive,
	}

	memberIDs := make([]string, 0, len(p.MemberIDs))
	batch := make([]imModel.GroupMember, 0, len(p.MemberIDs))
	sysBatch := make([]imModel.GroupMessage, 0, len(p.MemberIDs))
	if len(p.MemberIDs) > 0 {
		if err := validateMemberType(groupType, p.MemberType); err != nil {
			result.WriteError(c, err)
			return
		}
		for _, uid := range p.MemberIDs {
			if uid == userID {
				continue
			}
			memberIDs = append(memberIDs, uid)
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
	}

	if err := s.repo.CreateGroupWithMembers(ctx, &group, &ownerMember, batch, sysBatch, memberIDs, p.MemberType); err != nil {
		switch err.Error() {
		case "EXISTS_MEMBERS":
			result.WriteError(c, exception.NewBusinessError("部分成员已在群中", 400))
		case "EXCEED_MAX":
			result.WriteError(c, exception.NewBusinessError(fmt.Sprintf("群成员数量不能超过%d人", group.MaxMembers), 400))
		default:
			result.WriteError(c, exception.NewBusinessError("创建群失败: "+err.Error(), 500))
		}
		return
	}
}

// ==================== GroupUpdate ====================

func (s *Service) Update(c *gin.Context, p *UpdateParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)
	operatorType := getUserType(c)

	if p.GroupID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}
	_, member, err := s.checkOwnerOrAdmin(ctx, p.GroupID, operatorID, operatorType)
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

	if err := s.repo.UpdateGroup(ctx, p.GroupID, updates); err != nil {
		result.WriteError(c, exception.NewBusinessError("修改群信息失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupDissolve ====================

func (s *Service) Dissolve(c *gin.Context, p *DissolveParam) {
	ctx := c.Request.Context()
	operatorID := getLoginID(c)

	if p.GroupID == "" || operatorID == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}

	group, err := s.repo.FindGroupByID(ctx, p.GroupID)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("群不存在", 400))
		return
	}
	if group.OwnerID != operatorID {
		result.WriteError(c, exception.NewBusinessError("仅群主可解散群", 403))
		return
	}
	if err := s.repo.DissolveGroup(ctx, p.GroupID); err != nil {
		result.WriteError(c, exception.NewBusinessError("解散群失败: "+err.Error(), 500))
		return
	}
}

// ==================== GroupInvite ====================
