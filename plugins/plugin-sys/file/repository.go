package file

import (
	"context"
	"encoding/json"
	"time"

	"hei-gin/sdk/db"
)

func SaveChunkState(ctx context.Context, uploadID string, state chunkUploadState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return db.Redis.SetEx(ctx, chunkStateKey(uploadID), data, 24*time.Hour).Err()
}

func LoadChunkState(ctx context.Context, uploadID string) (*chunkUploadState, error) {
	data, err := db.Redis.Get(ctx, chunkStateKey(uploadID)).Bytes()
	if err != nil {
		return nil, err
	}
	var state chunkUploadState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func DeleteChunkState(ctx context.Context, uploadID string) error {
	return db.Redis.Del(ctx, chunkStateKey(uploadID)).Err()
}

func Page(ctx context.Context, p *FilePageParam) ([]SysFile, int64) {
	query := db.DB.WithContext(ctx).Model(&SysFile{})
	if p.Engine != "" {
		query = query.Where("engine = ?", p.Engine)
	}
	if p.Bucket != "" {
		query = query.Where("bucket = ?", p.Bucket)
	}
	if p.Keyword != "" {
		kw := "%" + p.Keyword + "%"
		query = query.Where("name LIKE ? OR name LIKE ?", kw, kw)
	}
	var total int64
	query.Count(&total)
	var rows []SysFile
	query.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func FindByID(ctx context.Context, id string) (*SysFile, error) {
	var entity SysFile
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func ListByIDs(ctx context.Context, ids []string) []SysFile {
	var files []SysFile
	if len(ids) == 0 {
		return files
	}
	db.DB.WithContext(ctx).Where("id IN ?", ids).Find(&files)
	return files
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysFile{}).Error
}

func Create(ctx context.Context, entity *SysFile) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func FindByKey(ctx context.Context, bucket, fileKey string) (*SysFile, error) {
	var entity SysFile
	if err := db.DB.WithContext(ctx).
		First(&entity, "bucket = ? AND file_key = ?", bucket, fileKey).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
