package message

import (
	"context"

	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/infra/db"
)

func (r *repository) CreateFile(ctx context.Context, record *imModel.ImFile) error {
	return db.DB.WithContext(ctx).Create(record).Error
}

func (r *repository) FindFileByKey(ctx context.Context, bucket, fileKey string) (*imModel.ImFile, error) {
	var entity imModel.ImFile
	if err := db.DB.WithContext(ctx).
		First(&entity, "bucket = ? AND file_key = ?", bucket, fileKey).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
