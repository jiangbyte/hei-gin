package log

import (
	"context"
	"time"

	"hei-gin/sdk/db"
)

func Page(ctx context.Context, p *LogPageParam) ([]SysLog, int64) {
	q := db.DB.WithContext(ctx).Model(&SysLog{})
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
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func FindByID(ctx context.Context, id string) (*SysLog, error) {
	var e SysLog
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func Create(ctx context.Context, entity *SysLog) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysLog{}).Where("id = ?", id).Updates(up).Error
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysLog{}).Error
}

func DeleteByCategory(ctx context.Context, category string) error {
	return db.DB.WithContext(ctx).Where("category = ?", category).Delete(&SysLog{}).Error
}

func ListByCategoriesSince(ctx context.Context, categories []string, since time.Time) []SysLog {
	var records []SysLog
	db.DB.WithContext(ctx).Where("category IN ?", categories).Where("op_time >= ?", since).Find(&records)
	return records
}

func CountByCategory(ctx context.Context, category string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&SysLog{}).Where("category = ?", category).Count(&count)
	return count
}
