// internal/modules/sys/audit/repo.go 持久化仓储。
//
// Author: Charlie

package audit

import (
	"context"

	"time"

	"gorm.io/gorm"
)

// Repo 审计日志持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Page 分页查询审计日志。
func (r *Repo) Page(ctx context.Context, q PageParam) (rows []OperationLog, total int64, err error) {
	cur, size := q.Normalize()
	db := r.with(ctx).Model(&OperationLog{})
	if q.Module != "" {
		db = db.Where("module = ?", q.Module)
	}
	if q.Action != "" {
		db = db.Where("action = ?", q.Action)
	} else if q.ExcludeAction != "" {
		db = db.Where("action <> ?", q.ExcludeAction)
	}
	if q.AccountID != "" {
		db = db.Where("account_id = ?", q.AccountID)
	}
	if q.Success != nil {
		db = db.Where("success = ?", *q.Success)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*OperationLog, error) {
	var row OperationLog
	if err := r.with(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// DB 返回底层 *gorm.DB（供 Job 直查）。
func (r *Repo) DB() *gorm.DB { return r.db }

const maxAuditCleanupRounds = 100

// CleanupExpiredLoginLogs 按保留天数分批删除 login/logout 审计日志。
func (r *Repo) CleanupExpiredLoginLogs(ctx context.Context, retentionDays, batchSize int) (int, error) {
	return r.cleanupExpired(ctx, retentionDays, batchSize, true)
}

// CleanupExpiredOperationLogs 按保留天数分批删除非 login/logout 操作审计日志。
func (r *Repo) CleanupExpiredOperationLogs(ctx context.Context, retentionDays, batchSize int) (int, error) {
	return r.cleanupExpired(ctx, retentionDays, batchSize, false)
}

func (r *Repo) cleanupExpired(ctx context.Context, retentionDays, batchSize int, loginLogs bool) (int, error) {
	if retentionDays <= 0 {
		return 0, nil
	}
	limit := batchSize
	if limit <= 0 {
		limit = 1000
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -retentionDays)
	total := 0
	for round := 0; round < maxAuditCleanupRounds; round++ {
		var ids []string
		q := r.with(ctx).Model(&OperationLog{}).
			Select("id").
			Where("created_at < ?", cutoff).
			Order("created_at ASC").
			Limit(limit)
		if loginLogs {
			q = q.Where("action IN ?", []string{"login", "logout"})
		} else {
			q = q.Where("action IS NULL OR action NOT IN ?", []string{"login", "logout"})
		}
		if err := q.Pluck("id", &ids).Error; err != nil {
			return total, err
		}
		if len(ids) == 0 {
			break
		}
		res := r.with(ctx).Where("id IN ?", ids).Delete(&OperationLog{})
		if res.Error != nil {
			return total, res.Error
		}
		total += int(res.RowsAffected)
		if len(ids) < limit {
			break
		}
	}
	return total, nil
}
