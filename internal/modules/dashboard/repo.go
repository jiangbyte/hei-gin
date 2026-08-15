// internal/modules/dashboard/repo.go 持久化仓储（聚合 SQL，对齐 hei-boot DashboardStatsMapper）。
//
// Author: Charlie

package dashboard

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Repo 仪表盘统计持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Aggregates 聚合计数（对齐 hei-boot aggregateCounts）。
func (r *Repo) Aggregates(ctx context.Context, dayStart time.Time) map[string]int64 {
	out := make(map[string]int64, 14)
	var row struct {
		AccountTotal    int64 `gorm:"column:account_total"`
		AccountEnabled  int64 `gorm:"column:account_enabled"`
		AccountDisabled int64 `gorm:"column:account_disabled"`
		AccountTodayNew int64 `gorm:"column:account_today_new"`
		FileTotal       int64 `gorm:"column:file_total"`
		StorageBytes    int64 `gorm:"column:storage_bytes"`
		RoleCount       int64 `gorm:"column:role_count"`
		DeptCount       int64 `gorm:"column:dept_count"`
		GroupCount      int64 `gorm:"column:group_count"`
		MenuCount       int64 `gorm:"column:menu_count"`
		AuditTotal      int64 `gorm:"column:audit_total"`
		AuditFailed     int64 `gorm:"column:audit_failed"`
		FeedbackPending int64 `gorm:"column:feedback_pending"`
	}
	if err := r.with(ctx).Raw(`SELECT (SELECT COUNT(*) FROM sys_account) AS account_total, (SELECT COUNT(*) FROM sys_account WHERE account_status = 'ENABLED') AS account_enabled, (SELECT COUNT(*) FROM sys_account WHERE account_status = 'DISABLED') AS account_disabled, (SELECT COUNT(*) FROM sys_account WHERE created_at >= ?) AS account_today_new, (SELECT COUNT(*) FROM sys_file) AS file_total, (SELECT COALESCE(SUM(size), 0) FROM sys_file) AS storage_bytes, (SELECT COUNT(*) FROM sys_role) AS role_count, (SELECT COUNT(*) FROM sys_dept) AS dept_count, (SELECT COUNT(*) FROM sys_group) AS group_count, (SELECT COUNT(*) FROM sys_resource WHERE resource_type = 'MENU' AND status = 'ENABLED') AS menu_count, (SELECT COUNT(*) FROM sys_operation_audit_log WHERE created_at >= ?) AS audit_total, (SELECT COUNT(*) FROM sys_operation_audit_log WHERE created_at >= ? AND success = false) AS audit_failed, (SELECT COUNT(*) FROM sys_feedback WHERE status = 'PENDING') AS feedback_pending`, dayStart, dayStart, dayStart).Scan(&row).Error; err != nil {
		return out
	}
	out["account_total"] = row.AccountTotal
	out["account_enabled"] = row.AccountEnabled
	out["account_disabled"] = row.AccountDisabled
	out["account_today_new"] = row.AccountTodayNew
	out["file_total"] = row.FileTotal
	out["storage_bytes"] = row.StorageBytes
	out["role_count"] = row.RoleCount
	out["dept_count"] = row.DeptCount
	out["group_count"] = row.GroupCount
	out["menu_count"] = row.MenuCount
	out["audit_total"] = row.AuditTotal
	out["audit_failed"] = row.AuditFailed
	out["feedback_pending"] = row.FeedbackPending
	return out
}

// AccountByType 按账号类型分组统计（对齐 hei-boot accountByType）。
func (r *Repo) AccountByType(ctx context.Context) []StatusItem {
	var rows []StatusItem
	if err := r.with(ctx).Raw(`SELECT account_type AS name, COUNT(*) AS value FROM sys_account GROUP BY account_type ORDER BY value DESC`).Scan(&rows).Error; err != nil {
		return []StatusItem{}
	}
	if rows == nil {
		rows = []StatusItem{}
	}
	return rows
}

// FileTypeShare 按内容类型统计文件数量（前 8；对齐 hei-boot fileTypeShare）。
func (r *Repo) FileTypeShare(ctx context.Context) []StatusItem {
	var rows []StatusItem
	if err := r.with(ctx).Raw(`SELECT COALESCE(content_type, 'unknown') AS name, COUNT(*) AS value FROM sys_file GROUP BY content_type ORDER BY value DESC LIMIT 8`).Scan(&rows).Error; err != nil {
		return []StatusItem{}
	}
	if rows == nil {
		rows = []StatusItem{}
	}
	return rows
}

// DailyCountRow 每日计数行。
type DailyCountRow struct {
	Day string `gorm:"column:day"`
	Cnt int64  `gorm:"column:cnt"`
}

// AccountDailyCounts 自 since 起每日新增账号数（对齐 hei-boot accountDailyCounts）。
func (r *Repo) AccountDailyCounts(ctx context.Context, since time.Time) []DailyCountRow {
	var rows []DailyCountRow
	if err := r.with(ctx).Raw(`SELECT TO_CHAR(created_at::date, 'YYYY-MM-DD') AS day, COUNT(*) AS cnt FROM sys_account WHERE created_at >= ? GROUP BY created_at::date ORDER BY day`, since).Scan(&rows).Error; err != nil {
		return nil
	}
	return rows
}

// AuditDailyCounts 自 since 起每日审计日志数（对齐 hei-boot auditDailyCounts）。
func (r *Repo) AuditDailyCounts(ctx context.Context, since time.Time) []DailyCountRow {
	var rows []DailyCountRow
	if err := r.with(ctx).Raw(`SELECT TO_CHAR(created_at::date, 'YYYY-MM-DD') AS day, COUNT(*) AS cnt FROM sys_operation_audit_log WHERE created_at >= ? GROUP BY created_at::date ORDER BY day`, since).Scan(&rows).Error; err != nil {
		return nil
	}
	return rows
}
