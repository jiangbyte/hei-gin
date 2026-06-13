package analyze

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/plugins/plugin-sys/shared"
	logModel "hei-gin/plugins/plugin-sys/log"
)

type repository struct {
	db *gorm.DB
}

type monthlyCount struct {
	Month string
	Count int
}

type orgCount struct {
	OrgID string
	Count int
}

type orgNameRow struct {
	ID   string
	Name string
}

type catCount struct {
	Category string
	Count    int
}

func (r *repository) Page(ctx context.Context, p *logModel.LogPageParam) ([]logModel.SysLog, int64) {
	q := r.db.WithContext(ctx).Model(&logModel.SysLog{})
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	if p.Keyword != "" {
		kw := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ? OR op_user LIKE ? OR op_ip LIKE ?", kw, kw, kw)
	}
	var total int64
	q.Count(&total)
	var rows []logModel.SysLog
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) CountLogsByCategory(ctx context.Context, category string) int64 {
	var count int64
	r.db.WithContext(ctx).Model(&logModel.SysLog{}).Where("category = ?", category).Count(&count)
	return count
}

func (r *repository) CountLogsByCategoryAndStatus(ctx context.Context, category, status string) int64 {
	var count int64
	r.db.WithContext(ctx).Model(&logModel.SysLog{}).Where("category = ?", category).Where("exe_status = ?", status).Count(&count)
	return count
}

func (r *repository) CountLogsByCategoryBetween(ctx context.Context, category string, start, end interface{}) int64 {
	var count int64
	r.db.WithContext(ctx).Model(&logModel.SysLog{}).Where("category = ?", category).Where("op_time >= ? AND op_time < ?", start, end).Count(&count)
	return count
}

func (r *repository) CountTable(ctx context.Context, table string) int64 {
	var count int64
	r.db.WithContext(ctx).Table(table).Count(&count)
	return count
}

func (r *repository) CountTableByStatus(ctx context.Context, table, status string) int64 {
	var count int64
	r.db.WithContext(ctx).Table(table).Where("status = ?", status).Count(&count)
	return count
}

func (r *repository) MonthlyTrend(ctx context.Context, table string) []monthlyCount {
	var rows []monthlyCount
	r.db.WithContext(ctx).Table(table).
		Select("DATE_FORMAT(created_at, '%Y-%m') AS month, COUNT(*) AS count").
		Where("created_at IS NOT NULL").
		Group("month").
		Order("month ASC").
		Limit(12).
		Find(&rows)
	return rows
}

func (r *repository) OrgUserCounts(ctx context.Context) []orgCount {
	var rows []orgCount
	r.db.WithContext(ctx).Table("sys_user").
		Select("org_id, COUNT(*) AS count").
		Where("org_id IS NOT NULL AND org_id != ''").
		Group("org_id").
		Find(&rows)
	return rows
}

func (r *repository) OrgNames(ctx context.Context, orgIDs []string) []orgNameRow {
	var rows []orgNameRow
	r.db.WithContext(ctx).Table("sys_org").Select("id, name").Where("id IN ?", orgIDs).Find(&rows)
	return rows
}

func (r *repository) RoleCategoryCounts(ctx context.Context) []catCount {
	var rows []catCount
	r.db.WithContext(ctx).Table("sys_role").
		Select("category, COUNT(*) AS count").
		Group("category").
		Find(&rows)
	return rows
}

func (r *repository) ActiveStatus() string {
	return shared.UserStatusActive
}
