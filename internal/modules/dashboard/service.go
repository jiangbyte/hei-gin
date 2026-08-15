// internal/modules/dashboard/service.go 业务服务（对齐 hei-boot DashboardServiceImpl）。
//
// Author: Charlie

package dashboard

import (
	"context"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service 仪表盘服务。
//
// Author: Charlie
type Service struct {
	repo     *Repo
	sessions *security.SessionStore
}

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 dashboard 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	s.sessions = d.Sessions
	return module.Module{
		Name:   "dashboard",
		Order:  50,
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Overview 概览统计（对齐 hei-boot overview：汇总 + 账号分布 + IAM + 今日运维 + 七日趋势 + 文件类型）。
func (s *Service) Overview(ctx context.Context) OverviewResult {
	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	since := dayStart.AddDate(0, 0, -6)
	counts := s.repo.Aggregates(ctx, dayStart)

	out := OverviewResult{
		Summary: SummaryResult{
			AccountTotal:   counts["account_total"],
			OnlineSessions: s.onlineSessionCount(ctx),
			FileTotal:      counts["file_total"],
			StorageBytes:   counts["storage_bytes"],
		},
		Accounts: AccountsResult{
			Enabled:  counts["account_enabled"],
			Disabled: counts["account_disabled"],
			TodayNew: counts["account_today_new"],
			ByType:   s.repo.AccountByType(ctx),
		},
		IAM: IAMResult{
			RoleCount:  counts["role_count"],
			DeptCount:  counts["dept_count"],
			GroupCount: counts["group_count"],
			MenuCount:  counts["menu_count"],
		},
		OpsToday: OpsTodayResult{
			AuditTotal:      counts["audit_total"],
			AuditFailed:     counts["audit_failed"],
			FeedbackPending: counts["feedback_pending"],
		},
		Trends: TrendsResult{
			AccountTrend: s.buildTrend(s.repo.AccountDailyCounts(ctx, since), since, "accounts"),
			AuditTrend:   s.buildTrend(s.repo.AuditDailyCounts(ctx, since), since, "audits"),
		},
		Files: FilesResult{
			ByContentType: s.repo.FileTypeShare(ctx),
		},
	}
	return out
}

// buildTrend 补齐近七日趋势点（date 取 MM-DD；对齐 hei-boot buildTrend）。
func (s *Service) buildTrend(rows []DailyCountRow, since time.Time, typ string) []TrendPoint {
	byDay := make(map[string]int64, len(rows))
	for _, row := range rows {
		byDay[row.Day] = row.Cnt
	}
	points := make([]TrendPoint, 0, 7)
	for i := 0; i < 7; i++ {
		d := since.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		points = append(points, TrendPoint{
			Date:  d.Format("01-02"),
			Type:  typ,
			Value: byDay[key],
		})
	}
	return points
}

// onlineSessionCount 在线会话数（经 SessionStore 账号索引统计；对齐 hei-boot StpKit searchSessionId 估算）。
func (s *Service) onlineSessionCount(ctx context.Context) int64 {
	if s.sessions == nil {
		return 0
	}
	accountIDs, err := s.sessions.ListAccountIDs(ctx)
	if err != nil {
		return 0
	}
	var total int64
	for _, id := range accountIDs {
		tokens, err2 := s.sessions.ListTokensForAccount(ctx, id)
		if err2 == nil {
			total += int64(len(tokens))
		}
	}
	return total
}
