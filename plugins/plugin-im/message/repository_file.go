package message

import (
	"context"

	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/db"
)

func CreateFile(ctx context.Context, record *imModel.ImFile) error {
	return db.DB.WithContext(ctx).Create(record).Error
}

func FindFileByKey(ctx context.Context, bucket, fileKey string) (*imModel.ImFile, error) {
	var entity imModel.ImFile
	if err := db.DB.WithContext(ctx).
		First(&entity, "bucket = ? AND file_key = ?", bucket, fileKey).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
