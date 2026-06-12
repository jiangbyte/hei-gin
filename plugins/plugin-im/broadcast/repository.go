package broadcast

import (
	"context"

	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/db"
	"hei-gin/sdk/utils"
)

func Create(ctx context.Context, entity *imModel.Broadcast) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func Page(ctx context.Context, cursor string, size int) []imModel.Broadcast {
	q := db.DB.WithContext(ctx).Model(&imModel.Broadcast{})
	if cursor != "" {
		if t, err := utils.ParseDateTime(cursor); err == nil {
			q = q.Where("created_at < ?", t)
		}
	}
	var records []imModel.Broadcast
	q.Order("created_at DESC").Limit(size + 1).Find(&records)
	return records
}

func ListLatest(ctx context.Context, size int) []imModel.Broadcast {
	var records []imModel.Broadcast
	db.DB.WithContext(ctx).Model(&imModel.Broadcast{}).Order("created_at DESC").Limit(size).Find(&records)
	return records
}

func ListReads(ctx context.Context, userID, userType string) []imModel.BroadcastRead {
	var rows []imModel.BroadcastRead
	db.DB.WithContext(ctx).Model(&imModel.BroadcastRead{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&rows)
	return rows
}

func MarkRead(ctx context.Context, broadcastID, userID, userType string) {
	_ = db.DB.WithContext(ctx).Where("broadcast_id = ? AND user_id = ? AND user_type = ?", broadcastID, userID, userType).
		FirstOrCreate(&imModel.BroadcastRead{
			BroadcastID: broadcastID,
			ID:          utils.GenerateID(),
			UserID:      userID,
			UserType:    userType,
		})
}

func FindByID(ctx context.Context, id string) (*imModel.Broadcast, error) {
	var entity imModel.Broadcast
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
