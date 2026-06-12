package session

import (
	"context"
	"fmt"
	"strings"
	"time"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

func (s *Service) Analysis(c *gin.Context) *SessionAnalysisResult {
	ctx := c.Request.Context()
	bStats, _ := auth.GetSessionStats(ctx, string(enums.LoginTypeBusiness))
	cStats, _ := auth.GetSessionStats(ctx, string(enums.LoginTypeConsumer))

	maxTokenCount := bStats.MaxTokenCount
	if cStats.MaxTokenCount > maxTokenCount {
		maxTokenCount = cStats.MaxTokenCount
	}

	return &SessionAnalysisResult{
		TotalCount:        bStats.TotalCount + cStats.TotalCount,
		MaxTokenCount:     maxTokenCount,
		OneHourNewlyAdded: bStats.OneHourNewlyAdded + cStats.OneHourNewlyAdded,
		ProportionOfBAndC: fmt.Sprintf("%d/%d", bStats.TotalCount, cStats.TotalCount),
	}
}

func (s *Service) Page(c *gin.Context, p *SessionPageParam) {
	ctx := c.Request.Context()
	current, size := normalizePage(p.Current, p.Size)

	keyword := strings.TrimSpace(p.Keyword)
	infos, total, err := s.listConsumerSessionInfos(ctx, keyword, current, size)
	if err != nil {
		infos = []auth.SessionInfo{}
		total = 0
	}

	users := s.repo.LoadUsers(ctx, sessionUserIDs(infos))
	rows := make([]*SessionPageResult, 0, len(infos))
	for _, info := range infos {
		row := &SessionPageResult{
			UserID:                info.UserID,
			TokenCount:            info.TokenCount,
			SessionCreateTime:     info.SessionCreateTime,
			SessionTimeout:        formatTimeout(info.SessionTimeoutSeconds),
			SessionTimeoutSeconds: info.SessionTimeoutSeconds,
		}
		if user := users[info.UserID]; user != nil {
			row.Username = user.Username
			row.Nickname = user.Nickname
			row.Avatar = user.Avatar
			row.Status = &user.Status
			row.LastLoginIP = user.LastLoginIP
		}
		rows = append(rows, row)
	}
	result.PageDataResult(c, rows, total, current, size)
}

func (s *Service) listConsumerSessionInfos(ctx context.Context, keyword string, current, size int) ([]auth.SessionInfo, int64, error) {
	if keyword == "" {
		return auth.ListSessionInfos(ctx, string(enums.LoginTypeConsumer), current, size, "")
	}
	userIDs := s.repo.FindUserIDs(ctx, keyword, maxKeywordCandidates)
	return auth.ListSessionInfosByUserIDs(ctx, string(enums.LoginTypeConsumer), userIDs, current, size)
}

const maxKeywordCandidates = 1000

func (s *Service) Exit(c *gin.Context, userID string) {
	auth.Consumer.KickoutWithContext(c.Request.Context(), userID)
}

func (s *Service) TokenList(c *gin.Context, userID string) []*SessionTokenResult {
	tokens, err := auth.GetSessionTokens(c.Request.Context(), string(enums.LoginTypeConsumer), userID)
	if err != nil || len(tokens) == 0 {
		return []*SessionTokenResult{}
	}

	results := make([]*SessionTokenResult, 0, len(tokens))
	for _, token := range tokens {
		results = append(results, &SessionTokenResult{
			Token:          token.Token,
			CreatedAt:      token.CreatedAt,
			Timeout:        formatTimeout(token.TimeoutSeconds),
			TimeoutSeconds: token.TimeoutSeconds,
			DeviceType:     token.DeviceType,
			DeviceID:       token.DeviceID,
		})
	}
	return results
}

func (s *Service) ExitToken(c *gin.Context, userID, token string) {
	auth.Consumer.KickoutTokenWithContext(c.Request.Context(), userID, token)
}

func (s *Service) Chart(c *gin.Context) *SessionChartData {
	ctx := c.Request.Context()
	days := lastNDays(7)

	cStats, _ := auth.GetSessionStats(ctx, string(enums.LoginTypeConsumer))
	cDaily := auth.GetSessionDailyCounts(ctx, string(enums.LoginTypeConsumer), days)

	series := make([]int, len(days))
	for i, day := range days {
		series[i] = cDaily[day]
	}

	return &SessionChartData{
		BarChart: BarChartData{
			Days: days,
			Series: []CategorySeries{
				{Name: "新增在线数", Data: series},
			},
		},
		PieChart: PieChartData{
			Data: []CategoryTotal{
				{Category: string(enums.LoginTypeConsumer), Total: cStats.TotalCount},
			},
		},
	}
}

func sessionUserIDs(infos []auth.SessionInfo) []string {
	ids := make([]string, 0, len(infos))
	seen := make(map[string]struct{}, len(infos))
	for _, info := range infos {
		if info.UserID == "" {
			continue
		}
		if _, ok := seen[info.UserID]; ok {
			continue
		}
		seen[info.UserID] = struct{}{}
		ids = append(ids, info.UserID)
	}
	return ids
}

func normalizePage(current, size int) (int, int) {
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return current, size
}

func formatTimeout(seconds int) string {
	if seconds < 0 {
		return "已过期"
	}
	if seconds == 0 {
		return "永久"
	}
	if seconds < 60 {
		return fmt.Sprintf("剩余 %d秒", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("剩余 %d分钟", seconds/60)
	}
	if seconds < 86400 {
		return fmt.Sprintf("剩余 %d小时%d分钟", seconds/3600, (seconds%3600)/60)
	}
	return fmt.Sprintf("剩余 %d天%d小时", seconds/86400, (seconds%86400)/3600)
}

func lastNDays(n int) []string {
	days := make([]string, n)
	now := time.Now()
	for i := 0; i < n; i++ {
		days[i] = now.AddDate(0, 0, -(n - 1 - i)).Format("2006-01-02")
	}
	return days
}
