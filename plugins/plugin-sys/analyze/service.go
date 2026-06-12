package analyze

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	logModel "hei-gin/plugins/plugin-sys/log"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

// Server start time, used for computing uptime/runtime display
var serverStartTime = time.Now()

// getServerIP attempts to detect a non-loopback IPv4 address
func getServerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, mins)
	}
	return fmt.Sprintf("%d分钟", mins)
}

// ===== Page =====

func (s *Service) Page(c *gin.Context, p *logModel.LogPageParam) {
	ctx := c.Request.Context()
	if p.Current < 1 {
		p.Current = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}

	rows, total := s.repo.Page(ctx, p)
	result.PageDataResult(c, rows, total, p.Current, p.Size)
}

// ===== LoginAnalysis =====

func (s *Service) LoginAnalysis(c *gin.Context) *LogAnalysisData {
	ctx := c.Request.Context()
	todayStart := time.Now().Truncate(24 * time.Hour)
	tomorrowStart := todayStart.Add(24 * time.Hour)

	loginTotal := s.repo.CountLogsByCategory(ctx, "LOGIN")
	failedTotal := s.repo.CountLogsByCategoryAndStatus(ctx, "LOGIN", "FAIL")
	loginToday := s.repo.CountLogsByCategoryBetween(ctx, "LOGIN", todayStart, tomorrowStart)

	log.Printf("[Analyze] Login stats: total=%d, failed=%d, today=%d", loginTotal, failedTotal, loginToday)

	return &LogAnalysisData{
		LoginTotal:  int(loginTotal),
		LoginFailed: int(failedTotal),
		LoginToday:  int(loginToday),
	}
}

// ===== LogAnalysis =====

func (s *Service) LogAnalysis(c *gin.Context) *LogAnalysisData {
	ctx := c.Request.Context()
	todayStart := time.Now().Truncate(24 * time.Hour)
	tomorrowStart := todayStart.Add(24 * time.Hour)

	logTotal := s.repo.CountTable(ctx, "sys_log")
	exceptionTotal := s.repo.CountLogsByCategory(ctx, "EXCEPTION")
	exceptionToday := s.repo.CountLogsByCategoryBetween(ctx, "EXCEPTION", todayStart, tomorrowStart)

	return &LogAnalysisData{
		LogTotal:       int(logTotal),
		LogException:   int(exceptionTotal),
		ExceptionToday: int(exceptionToday),
	}
}

// ===== Dashboard =====

func (s *Service) Dashboard(c *gin.Context) *DashboardVO {
	ctx := c.Request.Context()
	stats := DashboardStats{}

	stats.TotalUsers = s.repo.CountTable(ctx, "sys_user")
	stats.ActiveUsers = s.repo.CountTableByStatus(ctx, "sys_user", s.repo.ActiveStatus())
	stats.TotalRoles = s.repo.CountTable(ctx, "sys_role")
	stats.TotalOrgs = s.repo.CountTable(ctx, "sys_org")
	stats.TotalConfigs = s.repo.CountTable(ctx, "sys_config")
	stats.TotalNotices = s.repo.CountTable(ctx, "sys_notice")

	clientStats := ClientStats{}
	clientStats.TotalUsers = s.repo.CountTable(ctx, "client_user")
	clientStats.ActiveUsers = s.repo.CountTableByStatus(ctx, "client_user", s.repo.ActiveStatus())

	// User growth trend: monthly registration counts over the last 12 months
	userTrend := s.getMonthlyTrend(ctx, "sys_user")
	clientTrend := s.getMonthlyTrend(ctx, "client_user")

	// Org user distribution
	orgDist := s.getOrgUserDistribution(ctx)

	// Role category distribution
	roleDist := s.getRoleCategoryDistribution(ctx)

	sysInfo := SysInfo{
		OsName:   runtime.GOOS,
		ServerIP: getServerIP(),
		RunTime:  fmt.Sprintf("已运行 %s", formatDuration(time.Since(serverStartTime))),
	}

	return &DashboardVO{
		Stats:                    stats,
		ClientStats:              clientStats,
		UserTrend:                userTrend,
		ClientTrend:              clientTrend,
		OrgUserDistribution:      orgDist,
		RoleCategoryDistribution: roleDist,
		SysInfo:                  sysInfo,
	}
}

func (s *Service) getMonthlyTrend(ctx context.Context, table string) []TrendItem {
	rows := s.repo.MonthlyTrend(ctx, table)
	result := make([]TrendItem, len(rows))
	for i, r := range rows {
		result[i] = TrendItem{Month: r.Month, Count: r.Count}
	}
	if result == nil {
		result = []TrendItem{}
	}
	return result
}

func (s *Service) getOrgUserDistribution(ctx context.Context) []OrgUserDistribution {
	rows := s.repo.OrgUserCounts(ctx)
	orgIDs := make([]string, len(rows))
	for i, r := range rows {
		orgIDs[i] = r.OrgID
	}
	orgNames := make(map[string]string)
	if len(orgIDs) > 0 {
		orgRows := s.repo.OrgNames(ctx, orgIDs)
		for _, o := range orgRows {
			orgNames[o.ID] = o.Name
		}
	}
	result := make([]OrgUserDistribution, 0, len(rows))
	for _, r := range rows {
		name := orgNames[r.OrgID]
		if name == "" {
			name = "未分配"
		}
		result = append(result, OrgUserDistribution{Name: name, Count: r.Count})
	}
	if result == nil {
		result = []OrgUserDistribution{}
	}
	return result
}

func (s *Service) getRoleCategoryDistribution(ctx context.Context) []CategoryDistribution {
	rows := s.repo.RoleCategoryCounts(ctx)
	result := make([]CategoryDistribution, len(rows))
	for i, r := range rows {
		result[i] = CategoryDistribution{Category: r.Category, Count: r.Count}
	}
	if result == nil {
		result = []CategoryDistribution{}
	}
	return result
}
