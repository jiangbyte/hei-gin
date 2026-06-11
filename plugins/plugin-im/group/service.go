package group

import (
	"encoding/json"
	"fmt"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	imModel "hei-gin/plugins/plugin-im/model"

	"github.com/gin-gonic/gin"
)

// ── Group management ─────────────────────────────────────────────

// ── Group lifecycle ─────────────────────────────────────────────

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
