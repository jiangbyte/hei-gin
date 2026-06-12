package message

import (
	"context"
	"time"

	cliUser "hei-gin/plugins/plugin-client/user"
	imModel "hei-gin/plugins/plugin-im/model"
	sysUser "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/utils"
)

type conversationMessageRow struct {
	ConversationID string
	SenderID       string
	SenderType     string
	ReceiverID     string
	ReceiverType   string
	Content        string
	CreatedAt      time.Time
	Status         string
}

type conversationUnreadRow struct {
	ConversationID string
	Count          int64
}

func (r *repository) ListConversationLatest(ctx context.Context, currentUserID, userType string) []conversationMessageRow {
	queryDB := db.DB.WithContext(ctx)
	subQuery := queryDB.Table("im_message").
		Select("conversation_id, MAX(created_at) as max_ct").
		Where("((sender_id = ? AND sender_type = ?) OR (receiver_id = ? AND receiver_type = ?)) AND (deleted_by IS NULL OR deleted_by != ?)",
			currentUserID, userType, currentUserID, userType, currentUserID).
		Group("conversation_id")

	var rows []conversationMessageRow
	queryDB.Table("im_message m").
		Select("m.conversation_id, m.sender_id, m.sender_type, m.receiver_id, m.receiver_type, m.content, m.created_at, m.status").
		Joins("INNER JOIN (?) latest ON latest.conversation_id = m.conversation_id AND latest.max_ct = m.created_at", subQuery).
		Order("m.created_at DESC").
		Scan(&rows)
	return rows
}

func (r *repository) CountConversationUnread(ctx context.Context, convIDs []string, currentUserID, userType string) []conversationUnreadRow {
	var rows []conversationUnreadRow
	if len(convIDs) == 0 {
		return rows
	}
	db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Select("conversation_id, COUNT(*) as count").
		Where("conversation_id IN ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
			convIDs, currentUserID, userType, "unread").
		Group("conversation_id").
		Scan(&rows)
	return rows
}

func (r *repository) ListConversationMessages(ctx context.Context, conversationID, currentUserID, userType, cursor string, size int) []imModel.Message {
	query := db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("conversation_id = ? AND ((sender_id = ? AND sender_type = ?) OR (receiver_id = ? AND receiver_type = ?)) AND (deleted_by != ? OR deleted_by IS NULL)",
			conversationID, currentUserID, userType, currentUserID, userType, currentUserID)
	if cursor != "" {
		if t, err := utils.ParseDateTime(cursor); err == nil {
			query = query.Where("created_at < ?", t)
		}
	}
	order := "created_at DESC"
	if cursor != "" {
		order = "created_at ASC"
	}
	var rows []imModel.Message
	query.Order(order).Limit(size + 1).Find(&rows)
	return rows
}

func (r *repository) FindBusinessUsers(ctx context.Context, ids []string) []sysUser.SysUser {
	var users []sysUser.SysUser
	if len(ids) == 0 {
		return users
	}
	db.DB.WithContext(ctx).Model(&sysUser.SysUser{}).Where("id IN ?", ids).Find(&users)
	return users
}

func (r *repository) FindConsumerUsers(ctx context.Context, ids []string) []cliUser.ClientUser {
	var users []cliUser.ClientUser
	if len(ids) == 0 {
		return users
	}
	db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).Where("id IN ?", ids).Find(&users)
	return users
}

func (r *repository) FindBusinessUser(ctx context.Context, id string) (*sysUser.SysUser, error) {
	var user sysUser.SysUser
	if err := db.DB.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindConsumerUser(ctx context.Context, id string) (*cliUser.ClientUser, error) {
	var user cliUser.ClientUser
	if err := db.DB.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
