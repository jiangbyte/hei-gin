// internal/modules/sys/file/repo.go 持久化仓储。
//
// Author: Charlie

package file

import (
	"context"

	"gorm.io/gorm"
)

// Repo 文件元数据持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建文件记录。
func (r *Repo) Create(ctx context.Context, row *File) error {
	return r.with(ctx).Create(row).Error
}

// Update 按 ID 更新。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&File{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.with(ctx).Where("id IN ?", ids).Delete(&File{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*File, error) {
	var row File
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// FindByObjectName 按对象名查询。
func (r *Repo) FindByObjectName(ctx context.Context, objectName string) (*File, error) {
	var row File
	if err := r.with(ctx).Where("object_name = ?", objectName).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ListByIDs 批量查询。
func (r *Repo) ListByIDs(ctx context.Context, ids []string) ([]File, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []File
	err := r.with(ctx).Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

// LoadAccountNames 批量按账号 ID 解析主登录标识（ACCOUNT 身份），供 created_name/updated_name 回填。
func (r *Repo) LoadAccountNames(ctx context.Context, ids []string) map[string]string {
	out := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return out
	}
	var rows []struct {
		AccountID  string `gorm:"column:account_id"`
		Identifier string `gorm:"column:identifier"`
	}
	if err := r.with(ctx).Table("sys_account_identity").
		Select("account_id", "identifier").
		Where("account_id IN ? AND identity_type = ?", ids, "ACCOUNT").
		Find(&rows).Error; err != nil {
		return out
	}
	for _, row := range rows {
		if _, ok := out[row.AccountID]; !ok {
			out[row.AccountID] = row.Identifier
		}
	}
	return out
}

// Page 分页查询（对齐 hei-boot FileServiceImpl.page 过滤条件）。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []File, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&File{})
	if q.OriginalName != "" {
		db = db.Where("original_name ILIKE ?", "%"+q.OriginalName+"%")
	}
	if q.ObjectName != "" {
		db = db.Where("object_name ILIKE ?", "%"+q.ObjectName+"%")
	}
	if q.StorageProvider != "" {
		db = db.Where("storage_provider = ?", q.StorageProvider)
	}
	if q.ContentType != "" {
		db = db.Where("content_type ILIKE ?", "%"+q.ContentType+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at DESC, id DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
