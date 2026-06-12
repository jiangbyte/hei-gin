package file

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type repository struct {
	db  *gorm.DB
	rdb redis.UniversalClient
}

func (r *repository) SaveChunkState(ctx context.Context, uploadID string, state chunkUploadState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return r.rdb.SetEx(ctx, chunkStateKey(uploadID), data, 24*time.Hour).Err()
}

func (r *repository) LoadChunkState(ctx context.Context, uploadID string) (*chunkUploadState, error) {
	data, err := r.rdb.Get(ctx, chunkStateKey(uploadID)).Bytes()
	if err != nil {
		return nil, err
	}
	var state chunkUploadState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *repository) DeleteChunkState(ctx context.Context, uploadID string) error {
	return r.rdb.Del(ctx, chunkStateKey(uploadID)).Err()
}

func (r *repository) Page(ctx context.Context, p *FilePageParam) ([]SysFile, int64) {
	query := r.db.WithContext(ctx).Model(&SysFile{})
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

func (r *repository) FindByID(ctx context.Context, id string) (*SysFile, error) {
	var entity SysFile
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository) ListByIDs(ctx context.Context, ids []string) []SysFile {
	var files []SysFile
	if len(ids) == 0 {
		return files
	}
	r.db.WithContext(ctx).Where("id IN ?", ids).Find(&files)
	return files
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysFile{}).Error
}

func (r *repository) Create(ctx context.Context, entity *SysFile) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) FindByKey(ctx context.Context, bucket, fileKey string) (*SysFile, error) {
	var entity SysFile
	if err := r.db.WithContext(ctx).
		First(&entity, "bucket = ? AND file_key = ?", bucket, fileKey).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
