package log

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Page(ctx context.Context, p *LogPageParam) ([]SysLog, int64) {
	q := r.db.WithContext(ctx).Model(&SysLog{})
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	if p.ExeStatus != "" {
		q = q.Where("exe_status = ?", p.ExeStatus)
	}
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ? OR op_user LIKE ? OR op_ip LIKE ?", like, like, like)
	}
	var total int64
	q.Count(&total)
	var rows []SysLog
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysLog, error) {
	var e SysLog
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *SysLog) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysLog{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysLog{}).Error
}

func (r *repository) DeleteByCategory(ctx context.Context, category string) error {
	return r.db.WithContext(ctx).Where("category = ?", category).Delete(&SysLog{}).Error
}

func (r *repository) ListByCategoriesSince(ctx context.Context, categories []string, since time.Time) []SysLog {
	var records []SysLog
	r.db.WithContext(ctx).Where("category IN ?", categories).Where("op_time >= ?", since).Find(&records)
	return records
}

func (r *repository) CountByCategory(ctx context.Context, category string) int64 {
	var count int64
	r.db.WithContext(ctx).Model(&SysLog{}).Where("category = ?", category).Count(&count)
	return count
}
