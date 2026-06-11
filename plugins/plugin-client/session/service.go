package session

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/constants"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/db"
	"hei-gin/sdk/result"
	cliUser "hei-gin/plugins/plugin-client/user"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var svcCtx = context.Background()

func scanKeys(ctx context.Context, redis *redis.Client, pattern string) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		batch, nextCursor, err := redis.Scan(ctx, cursor, pattern, 200).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, batch...)
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return keys, nil
}

// ===== Analysis =====

func SessionAnalysis(c *gin.Context) *SessionAnalysisResult {
	bKeys, _ := scanKeys(svcCtx, db.Redis, constants.SESSION_PREFIX_BUSINESS+"*")
	cKeys, _ := scanKeys(svcCtx, db.Redis, constants.SESSION_PREFIX_CONSUMER+"*")

	bTotal, bNewly, bMax := countTokens(svcCtx, db.Redis, bKeys, constants.TOKEN_PREFIX_BUSINESS)
	cTotal, cNewly, cMax := countTokens(svcCtx, db.Redis, cKeys, constants.TOKEN_PREFIX_CONSUMER)

	maxTokenCount := bMax
	if cMax > maxTokenCount {
		maxTokenCount = cMax
	}

	return &SessionAnalysisResult{
		TotalCount:        bTotal + cTotal,
		MaxTokenCount:     maxTokenCount,
		OneHourNewlyAdded: bNewly + cNewly,
		ProportionOfBAndC: fmt.Sprintf("%d/%d", bTotal, cTotal),
	}
}

func countTokens(ctx context.Context, redis *redis.Client, sessionKeys []string, tokenPrefix string) (total, oneHourNewlyAdded, maxPerUser int) {
	userTokenCounts := make(map[string]int)
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	for _, sessionKey := range sessionKeys {
		parts := strings.Split(sessionKey, ":")
		userID := parts[len(parts)-1]

		tokens, err := redis.SMembers(ctx, sessionKey).Result()
		if err != nil {
			continue
		}
		userTokenCounts[userID] = len(tokens)

		for _, token := range tokens {
			total++
			tokenKey := tokenPrefix + token
			data, err := redis.Get(ctx, tokenKey).Result()
			if err != nil {
				continue
			}
			var tokenData map[string]any
			if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
				continue
			}
			createdAtStr, _ := tokenData["created_at"].(string)
			if createdAtStr != "" {
				createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
				if err == nil && createdAt.After(oneHourAgo) {
					oneHourNewlyAdded++
				}
			}
		}
	}

	for _, count := range userTokenCounts {
		if count > maxPerUser {
			maxPerUser = count
		}
	}
	return
}

func countDaily(ctx context.Context, redis *redis.Client, sessionKeys []string, tokenPrefix string) map[string]int {
	daily := make(map[string]int)
	for _, sessionKey := range sessionKeys {
		tokens, err := redis.SMembers(ctx, sessionKey).Result()
		if err != nil {
			continue
		}
		for _, token := range tokens {
			tokenKey := tokenPrefix + token
			data, err := redis.Get(ctx, tokenKey).Result()
			if err != nil {
				continue
			}
			var tokenData map[string]any
			if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
				continue
			}
			createdAtStr, _ := tokenData["created_at"].(string)
			if createdAtStr != "" {
				createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
				if err == nil {
					daily[createdAt.Format("2006-01-02")]++
				}
			}
		}
	}
	return daily
}

// ===== Page =====

func SessionPage(c *gin.Context, p *SessionPageParam) {
	sessions, err := collectSessions(svcCtx, db.Redis, constants.SESSION_PREFIX_CONSUMER, constants.TOKEN_PREFIX_CONSUMER, p.Keyword)
	if err != nil || sessions == nil {
		sessions = []*SessionPageResult{}
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].SessionCreateTime > sessions[j].SessionCreateTime
	})

	total := len(sessions)
	if p.Current < 1 {
		p.Current = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}

	start := (p.Current - 1) * p.Size
	var pageRecords []*SessionPageResult
	if start >= total {
		pageRecords = []*SessionPageResult{}
	} else {
		end := start + p.Size
		if end > total {
			end = total
		}
		pageRecords = sessions[start:end]
	}
	result.PageDataResult(c, pageRecords, int64(total), p.Current, p.Size)
}

func collectSessions(ctx context.Context, redis *redis.Client, sessionPrefix, tokenPrefix, keyword string) ([]*SessionPageResult, error) {
	sessionKeys, err := scanKeys(ctx, redis, sessionPrefix + "*")
	if err != nil {
		return nil, err
	}

	var result []*SessionPageResult
	userCache := make(map[string]*cliUser.ClientUser)

	for _, sessionKey := range sessionKeys {
		parts := strings.Split(sessionKey, ":")
		userID := parts[len(parts)-1]

		if keyword != "" && !strings.Contains(userID, keyword) {
			continue
		}

		// Session key is a Redis SET (stores token members via SAdd), NOT a String.
		// Do NOT redis.Get() it — that causes WRONGTYPE error.
		// Instead, read token data to derive session info from the earliest token.
		tokens, err := redis.SMembers(ctx, sessionKey).Result()
		if err != nil || len(tokens) == 0 {
			continue
		}

		var minCreatedAt time.Time
		hasCreatedAt := false
		for _, token := range tokens {
			tokenKey := tokenPrefix + token
			data, err := redis.Get(ctx, tokenKey).Result()
			if err != nil {
				continue
			}
			var tokenData map[string]any
			if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
				continue
			}
			createdAtStr, _ := tokenData["created_at"].(string)
			if createdAtStr != "" {
				createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
				if err == nil && (!hasCreatedAt || createdAt.Before(minCreatedAt)) {
					minCreatedAt = createdAt
					hasCreatedAt = true
				}
			}
		}

		user, ok := userCache[userID]
		if !ok {
			var u cliUser.ClientUser
			if err := db.DB.First(&u, "id = ?", userID).Error; err != nil {
				userCache[userID] = nil
			} else {
				user = &u
				userCache[userID] = user
			}
		}

		item := &SessionPageResult{
			UserID:     userID,
			TokenCount: len(tokens),
		}
		if user != nil {
			item.Username = user.Username
			item.Nickname = user.Nickname
			item.Avatar = user.Avatar
			item.Status = &user.Status
			item.LastLoginIP = user.LastLoginIP
		}
		if hasCreatedAt {
			item.SessionCreateTime = minCreatedAt.Format("2006-01-02 15:04:05")
		}
		result = append(result, item)
	}
	return result, nil
}

// ===== Exit =====

func SessionExit(c *gin.Context, userID string) {
	auth.Consumer.Kickout(userID)
}

// ===== TokenList =====

func SessionTokenList(c *gin.Context, userID string) []*SessionTokenResult {
	sessionKey := constants.SESSION_PREFIX_CONSUMER + userID
	tokens, err := db.Redis.SMembers(svcCtx, sessionKey).Result()
	if err != nil || len(tokens) == 0 {
		return []*SessionTokenResult{}
	}

	var results []*SessionTokenResult
	for _, token := range tokens {
		tokenKey := constants.TOKEN_PREFIX_CONSUMER + token
		data, err := db.Redis.Get(svcCtx, tokenKey).Result()
		if err != nil {
			continue
		}
		var tokenData map[string]any
		if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
			continue
		}

		createdAt, _ := tokenData["created_at"].(string)
		extra, _ := tokenData["extra"].(map[string]any)
		deviceType, _ := extra["device_type"].(string)
		deviceID, _ := extra["device_id"].(string)

		ttl, err := db.Redis.TTL(svcCtx, tokenKey).Result()
		timeoutSeconds := -1
		if err == nil {
			timeoutSeconds = int(ttl.Seconds())
		}

		results = append(results, &SessionTokenResult{
			Token: token, CreatedAt: createdAt,
			Timeout: formatTimeout(timeoutSeconds), TimeoutSeconds: timeoutSeconds,
			DeviceType: deviceType, DeviceID: deviceID,
		})
	}
	return results
}

// ===== ExitToken =====

func SessionExitToken(c *gin.Context, userID, token string) {
	auth.Consumer.KickoutToken(userID, token)
}

// ===== ChartData =====

func SessionChart(c *gin.Context) *SessionChartData {
	cKeys, _ := scanKeys(svcCtx, db.Redis, constants.SESSION_PREFIX_CONSUMER+"*")
	cTotal, _, _ := countTokens(svcCtx, db.Redis, cKeys, constants.TOKEN_PREFIX_CONSUMER)
	cDaily := countDaily(svcCtx, db.Redis, cKeys, constants.TOKEN_PREFIX_CONSUMER)

	days := lastNDays(7)
	series := make([]int, 7)
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
				{Category: string(enums.LoginTypeConsumer), Total: cTotal},
			},
		},
	}
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
