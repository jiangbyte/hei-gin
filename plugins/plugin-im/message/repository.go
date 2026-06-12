package message

import (
	"context"

	"gorm.io/gorm/clause"
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/utils"
)

type repository struct{}

func (r *repository) CreateMessages(ctx context.Context, records []imModel.Message) error {
	return db.DB.WithContext(ctx).Create(&records).Error
}

func (r *repository) UpsertConversations(ctx context.Context, conversations []imModel.Conversation) error {
	if len(conversations) == 0 {
		return nil
	}
	return db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"last_msg", "last_time", "updated_at",
		}),
	}).Create(&conversations).Error
}

func (r *repository) Page(ctx context.Context, userID, userType, status string, current, size int) ([]imModel.Message, int64) {
	query := db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("((sender_id = ? AND sender_type = ?) OR (receiver_id = ? AND receiver_type = ?)) AND (deleted_by != ? OR deleted_by IS NULL)",
			userID, userType, userID, userType, userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	query.Count(&total)
	var records []imModel.Message
	query.Order("created_at DESC").Limit(size).Offset((current - 1) * size).Find(&records)
	return records, total
}

func (r *repository) CountUnread(ctx context.Context, userID, userType string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("receiver_id = ? AND receiver_type = ? AND status = ?", userID, userType, "unread").
		Count(&count)
	return count
}

func (r *repository) FindOwnedByID(ctx context.Context, id, userID, userType string) (*imModel.Message, error) {
	var entity imModel.Message
	if err := db.DB.WithContext(ctx).
		Where("id = ? AND ((sender_id = ? AND sender_type = ?) OR (receiver_id = ? AND receiver_type = ?))",
			id, userID, userType, userID, userType).
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*imModel.Message, error) {
	var entity imModel.Message
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository) MarkRead(ctx context.Context, id, userID, userType string) error {
	return db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("id = ? AND receiver_id = ? AND receiver_type = ?", id, userID, userType).
		Updates(map[string]interface{}{"status": "read"}).Error
}

func (r *repository) MarkConversationRead(ctx context.Context, conversationID, userID, userType string) error {
	return db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
			conversationID, userID, userType, "unread").
		Updates(map[string]interface{}{"status": "read"}).Error
}

func (r *repository) MarkAllRead(ctx context.Context, userID, userType string) error {
	return db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("receiver_id = ? AND receiver_type = ? AND status = ?", userID, userType, "unread").
		Updates(map[string]interface{}{"status": "read"}).Error
}

func (r *repository) SoftDelete(ctx context.Context, ids []string, userID, userType string) error {
	return db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("id IN ? AND ((sender_id = ? AND sender_type = ?) OR (receiver_id = ? AND receiver_type = ?))",
			ids, userID, userType, userID, userType).
		Update("deleted_by", userID).Error
}

func (r *repository) Recall(ctx context.Context, messageID string) error {
	return db.DB.WithContext(ctx).Model(&imModel.Message{}).Where("id = ?", messageID).
		Updates(map[string]interface{}{"content": "消息已被撤回", "msg_type": imModel.MsgTypeSystem}).Error
}

func (r *repository) Search(ctx context.Context, userID, userType, keyword, cursor string, size int) []imModel.Message {
	query := db.DB.WithContext(ctx).Model(&imModel.Message{}).
		Where("((sender_id = ? AND sender_type = ?) OR (receiver_id = ? AND receiver_type = ?)) AND content LIKE ?",
			userID, userType, userID, userType, "%"+keyword+"%")
	if cursor != "" {
		if t, err := utils.ParseDateTime(cursor); err == nil {
			query = query.Where("created_at < ?", t)
		}
	}
	var records []imModel.Message
	query.Order("created_at DESC").Limit(size + 1).Find(&records)
	return records
}
