package group

import (
	"context"
	"time"

	"gorm.io/gorm/clause"
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
)

type groupCountRow struct {
	GroupID string
	Count   int
}

type groupLastMessageRow struct {
	GroupID   string
	Content   string
	CreatedAt time.Time
}

type groupUnreadRow struct {
	GroupID string
	Count   int64
}

type groupRecipientRow struct {
	UserID   string
	UserType string
}

func FindGroupByID(ctx context.Context, groupID string) (*imModel.Group, error) {
	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", groupID).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func FindActiveMember(ctx context.Context, groupID, userID, userType string) (*imModel.GroupMember, error) {
	var member imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, MemberActive).
		First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func CreateGroupWithMembers(ctx context.Context, group *imModel.Group, owner *imModel.GroupMember, members []imModel.GroupMember, sysMsgs []imModel.GroupMessage, memberIDs []string, memberType string) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Create(group).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(owner).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(memberIDs) > 0 {
		var existingCount int64
		tx.Model(&imModel.GroupMember{}).
			Where("group_id = ? AND user_id IN ? AND user_type = ? AND status = ?",
				group.ID, memberIDs, memberType, MemberActive).
			Count(&existingCount)
		if existingCount > 0 {
			tx.Rollback()
			return exception.NewBusinessError("EXISTS_MEMBERS", 400)
		}
		var currentCount int64
		tx.Model(&imModel.GroupMember{}).Where("group_id = ? AND status = ?", group.ID, MemberActive).Count(&currentCount)
		if int(currentCount)+len(memberIDs) > group.MaxMembers {
			tx.Rollback()
			return exception.NewBusinessError("EXCEED_MAX", 400)
		}
		if len(members) > 0 {
			if err := tx.Create(&members).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		if len(sysMsgs) > 0 {
			if err := tx.Create(&sysMsgs).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

func UpdateGroup(ctx context.Context, groupID string, updates map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&imModel.Group{}).Where("id = ?", groupID).Updates(updates).Error
}

func DissolveGroup(ctx context.Context, groupID string) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Model(&imModel.Group{}).Where("id = ?", groupID).Update("status", GroupDissolved).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&imModel.GroupMember{}).Where("group_id = ? AND status = ?", groupID, MemberActive).Update("status", MemberLeft).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func ListMyGroupMemberships(ctx context.Context, userID, userType string) []imModel.GroupMember {
	var rows []imModel.GroupMember
	db.DB.WithContext(ctx).
		Select("group_id, joined_at").
		Where("user_id = ? AND user_type = ? AND status = ?", userID, userType, MemberActive).
		Find(&rows)
	return rows
}

func ListMyGroupIDs(ctx context.Context, userID, userType string) []string {
	var ids []string
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("user_id = ? AND user_type = ? AND status = ?", userID, userType, MemberActive).
		Pluck("group_id", &ids)
	return ids
}

func ListGroupsByIDs(ctx context.Context, groupIDs []string) []imModel.Group {
	var groups []imModel.Group
	db.DB.WithContext(ctx).Where("id IN ? AND status = ?", groupIDs, GroupNormal).Find(&groups)
	return groups
}

func CountActiveMembersByGroupIDs(ctx context.Context, groupIDs []string) []groupCountRow {
	var rows []groupCountRow
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Select("group_id, COUNT(*) as count").
		Where("group_id IN ? AND status = ?", groupIDs, MemberActive).
		Group("group_id").Scan(&rows)
	return rows
}

func ListLastMessagesByGroupIDs(ctx context.Context, groupIDs []string) []groupLastMessageRow {
	var rows []groupLastMessageRow
	lastSubQ := db.DB.WithContext(ctx).Table("im_group_message").
		Select("group_id, MAX(created_at) as max_ct").
		Where("group_id IN ?", groupIDs).
		Group("group_id")
	db.DB.WithContext(ctx).Table("im_group_message g2").
		Select("g2.group_id, g2.content, g2.created_at").
		Joins("INNER JOIN (?) g1 ON g1.group_id = g2.group_id AND g1.max_ct = g2.created_at", lastSubQ).
		Scan(&rows)
	return rows
}

func CountUnreadByGroupIDs(ctx context.Context, groupIDs []string, userID, userType string) []groupUnreadRow {
	var rows []groupUnreadRow
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
		Scan(&rows)
	return rows
}

func CountActiveMembers(ctx context.Context, groupID string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).Where("group_id = ? AND status = ?", groupID, MemberActive).Count(&count)
	return count
}

func SearchGroups(ctx context.Context, like string, limit int) []imModel.Group {
	var groups []imModel.Group
	db.DB.WithContext(ctx).
		Where("name LIKE ? AND status = ?", like, GroupNormal).
		Limit(limit).
		Find(&groups)
	return groups
}

func ListActiveMembers(ctx context.Context, groupID string) []imModel.GroupMember {
	var rows []imModel.GroupMember
	db.DB.WithContext(ctx).
		Where("group_id = ? AND status = ?", groupID, MemberActive).
		Order("FIELD(role, 'owner', 'admin', 'member'), joined_at ASC").
		Find(&rows)
	return rows
}

func CreateMessage(ctx context.Context, msg *imModel.GroupMessage) error {
	return db.DB.WithContext(ctx).Create(msg).Error
}

func FindMessageByID(ctx context.Context, messageID, groupID string) (*imModel.GroupMessage, error) {
	var msg imModel.GroupMessage
	if err := db.DB.WithContext(ctx).Where("id = ? AND group_id = ?", messageID, groupID).First(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

func RecallMessage(ctx context.Context, messageID string) error {
	return db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).Where("id = ?", messageID).
		Updates(map[string]interface{}{"content": "消息已被撤回", "msg_type": imModel.MsgTypeSystem}).Error
}

func ListRecipientMembers(ctx context.Context, groupID, excludeUserID, excludeUserType string) []groupRecipientRow {
	var rows []groupRecipientRow
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Select("user_id, user_type").
		Where("group_id = ? AND status = ? AND NOT (user_id = ? AND user_type = ?)",
			groupID, MemberActive, excludeUserID, excludeUserType).
		Find(&rows)
	return rows
}

func UpsertMessageRead(ctx context.Context, record *imModel.GroupMessageRead) error {
	return db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}, {Name: "user_id"}, {Name: "user_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"read_at", "group_id"}),
	}).Create(record).Error
}

func FindLastMessageID(ctx context.Context, groupID string) string {
	var row struct{ ID string }
	if err := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).
		Select("id").
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Limit(1).
		Scan(&row).Error; err != nil {
		return ""
	}
	return row.ID
}

func ListGroupMessages(ctx context.Context, groupID, cursor string, size int) []imModel.GroupMessage {
	q := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).Where("group_id = ?", groupID)
	if cursor != "" {
		if t, err := utils.ParseDateTime(cursor); err == nil {
			q = q.Where("created_at < ?", t)
		}
	}
	order := "created_at DESC"
	if cursor != "" {
		order = "created_at ASC"
	}
	var rows []imModel.GroupMessage
	q.Order(order).Limit(size + 1).Find(&rows)
	return rows
}

func SearchGroupMessages(ctx context.Context, groupID, keyword, cursor string, size int) []imModel.GroupMessage {
	q := db.DB.WithContext(ctx).Model(&imModel.GroupMessage{}).
		Where("group_id = ? AND content LIKE ? AND msg_type != ?", groupID, "%"+keyword+"%", imModel.MsgTypeSystem)
	if cursor != "" {
		if t, err := utils.ParseDateTime(cursor); err == nil {
			q = q.Where("created_at < ?", t)
		}
	}
	order := "created_at DESC"
	if cursor != "" {
		order = "created_at ASC"
	}
	var rows []imModel.GroupMessage
	q.Order(order).Limit(size + 1).Find(&rows)
	return rows
}

func ListPendingJoinRequests(ctx context.Context, groupID string) []imModel.GroupJoinRequest {
	var rows []imModel.GroupJoinRequest
	db.DB.WithContext(ctx).Model(&imModel.GroupJoinRequest{}).
		Where("group_id = ? AND status = ?", groupID, "pending").
		Order("created_at DESC").Find(&rows)
	return rows
}

func FindPendingJoinRequest(ctx context.Context, requestID string) (*imModel.GroupJoinRequest, error) {
	var req imModel.GroupJoinRequest
	if err := db.DB.WithContext(ctx).First(&req, "id = ? AND status = ?", requestID, "pending").Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func HandleJoinRequest(ctx context.Context, requestID, operatorID, action string, member *imModel.GroupMember, msg *imModel.GroupMessage) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Model(&imModel.GroupJoinRequest{}).Where("id = ?", requestID).
		Updates(map[string]interface{}{"status": action, "handled_by": operatorID}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if action == "approved" {
		if err := tx.Create(member).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Create(msg).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func FindExistingMemberIDs(ctx context.Context, groupID string, userIDs []string, userType string) []string {
	var ids []string
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id IN ? AND user_type = ? AND status = ?", groupID, userIDs, userType, MemberActive).
		Pluck("user_id", &ids)
	return ids
}

func CountActiveMembersByGroup(ctx context.Context, groupID string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND status = ?", groupID, MemberActive).Count(&count)
	return count
}

func InviteMembers(ctx context.Context, members []imModel.GroupMember, sysMsgs []imModel.GroupMessage) error {
	tx := db.DB.WithContext(ctx).Begin()
	if len(members) > 0 {
		if err := tx.Create(&members).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(sysMsgs) > 0 {
		if err := tx.Create(&sysMsgs).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func CountMemberExists(ctx context.Context, groupID, userID, userType string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, MemberActive).Count(&count)
	return count
}

func CountPendingJoin(ctx context.Context, groupID, userID, userType string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.GroupJoinRequest{}).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, "pending").Count(&count)
	return count
}

func CreateJoinRequest(ctx context.Context, req *imModel.GroupJoinRequest) error {
	return db.DB.WithContext(ctx).Create(req).Error
}

func ListManagers(ctx context.Context, groupID string) []imModel.GroupMember {
	var rows []imModel.GroupMember
	db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND (role = ? OR role = ?) AND status = ?",
			groupID, RoleOwner, RoleAdmin, MemberActive).
		Find(&rows)
	return rows
}

func LeaveGroup(ctx context.Context, member *imModel.GroupMember, msg *imModel.GroupMessage) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Model(member).Update("status", MemberLeft).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(msg).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func KickMember(ctx context.Context, groupID, userID, userType string, msg *imModel.GroupMessage) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", groupID, userID, userType).
		Update("status", MemberKicked).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(msg).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func TransferOwnerAndRole(ctx context.Context, groupID, oldOwnerID, oldOwnerType, newOwnerID, newOwnerType string) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", groupID, oldOwnerID, oldOwnerType).
		Update("role", RoleAdmin).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", groupID, newOwnerID, newOwnerType).
		Update("role", RoleOwner).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&imModel.Group{}).Where("id = ?", groupID).
		Updates(map[string]interface{}{"owner_id": newOwnerID, "owner_type": newOwnerType}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func TransferOwner(ctx context.Context, groupID, newOwnerID, newOwnerType, operatorID, operatorType string) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Model(&imModel.Group{}).Where("id = ?", groupID).
		Updates(map[string]interface{}{"owner_id": newOwnerID, "owner_type": newOwnerType}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", groupID, operatorID, operatorType).
		Update("role", RoleAdmin).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", groupID, newOwnerID, newOwnerType).
		Update("role", RoleOwner).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func UpdateMemberRole(ctx context.Context, groupID, userID, userType, role string) error {
	return db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, MemberActive).
		Update("role", role).Error
}

func UpdateMemberNickname(ctx context.Context, groupID, userID, userType, nickname string) error {
	return db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", groupID, userID, userType).
		Update("nickname", nickname).Error
}

func UpdateMutedUntil(ctx context.Context, groupID, userID, userType string, until *time.Time) error {
	return db.DB.WithContext(ctx).Model(&imModel.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND user_type = ?", groupID, userID, userType).
		Update("muted_until", until).Error
}
