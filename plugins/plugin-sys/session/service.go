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
	"hei-gin/sdk/db"
	"hei-gin/sdk/result"
	userModel "hei-gin/plugins/plugin-sys/user"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

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

func Analysis(c *gin.Context) *SessionAnalysisResult {
	ctx := c.Request.Context()
	bKeys, _ := scanKeys(ctx, db.Redis, constants.SESSION_PREFIX_BUSINESS+"*")
	cKeys, _ := scanKeys(ctx, db.Redis, constants.SESSION_PREFIX_CONSUMER+"*")

	bTotal, bNewly, bMax := countTokens(ctx, db.Redis, bKeys, constants.TOKEN_PREFIX_BUSINESS)
	cTotal, cNewly, cMax := countTokens(ctx, db.Redis, cKeys, constants.TOKEN_PREFIX_CONSUMER)

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
					day := createdAt.Format("2006-01-02")
					daily[day]++
				}
			}
		}
	}
	return daily
}

func Page(c *gin.Context, param *SessionPageParam) {
	ctx := c.Request.Context()
	sessions, err := collectSessions(ctx, db.Redis, constants.SESSION_PREFIX_BUSINESS, constants.TOKEN_PREFIX_BUSINESS, param.Keyword)
	if err != nil || sessions == nil {
		sessions = []*SessionPageResult{}
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].SessionCreateTime > sessions[j].SessionCreateTime
	})

	total := len(sessions)
	current := param.Current
	if current < 1 {
		current = 1
	}
	size := param.Size
	if size < 1 {
		size = 10
	}

	start := (current - 1) * size
	var pageRecords []*SessionPageResult
	if start >= total {
		pageRecords = []*SessionPageResult{}
	} else {
		end := start + size
		if end > total {
			end = total
		}
		pageRecords = sessions[start:end]
	}

	result.PageDataResult(c, pageRecords, int64(total), current, size)
}

func collectSessions(ctx context.Context, redis *redis.Client, sessionPrefix, tokenPrefix, keyword string) ([]*SessionPageResult, error) {
	sessionKeys, err := scanKeys(ctx, redis, sessionPrefix+"*")
	if err != nil {
		return nil, err
	}

	var result []*SessionPageResult
	userCache := make(map[string]*userModel.SysUser)

	for _, sessionKey := range sessionKeys {
		parts := strings.Split(sessionKey, ":")
		userID := parts[len(parts)-1]

		if keyword != "" && !strings.Contains(userID, keyword) {
			continue
		}

		// Session key is a Redis SET (stores token members via SAdd), NOT a String.
		// Do NOT redis.Get() it — that causes WRONGTYPE error.
		// Instead, read token data to derive session info.
		tokens, err := redis.SMembers(ctx, sessionKey).Result()
		if err != nil {
			continue
		}
		tokenCount := len(tokens)

		sessionCreateTime := ""
		username := ""
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
			ct, _ := tokenData["created_at"].(string)
			if ct != "" && (sessionCreateTime == "" || ct < sessionCreateTime) {
				sessionCreateTime = ct
			}
			if ext, ok := tokenData["extra"].(map[string]any); ok {
				if un, _ := ext["username"].(string); un != "" && username == "" {
					username = un
				}
			}
			break
		}

		ttl, err := redis.TTL(ctx, sessionKey).Result()
		timeoutSeconds := -1
		if err == nil {
			timeoutSeconds = int(ttl.Seconds())
		}

		user, ok := userCache[userID]
		if !ok {
			var u userModel.SysUser
			err := db.DB.WithContext(ctx).First(&u, "id = ?", userID).Error
			if err == nil {
				user = &u
				userCache[userID] = user
			}
		}

		sr := &SessionPageResult{
			UserID:                userID,
			TokenCount:            tokenCount,
			SessionCreateTime:     sessionCreateTime,
			SessionTimeout:        formatTimeout(timeoutSeconds),
			SessionTimeoutSeconds: timeoutSeconds,
		}
		if username != "" {
			sr.Username = &username
		}
		if user != nil {
			sr.Nickname = user.Nickname
			sr.Avatar = user.Avatar
			sr.Status = user.Status
			sr.LastLoginIP = user.LastLoginIP
			if user.LastLoginAt != nil {
				sr.LastLoginTime = user.LastLoginAt.Format("2006-01-02 15:04:05")
			}
		}

		result = append(result, sr)
	}

	return result, nil
}

func Exit(c *gin.Context, param *SessionExitParam) {
	auth.Kickout(param.UserID)
}

func TokenList(c *gin.Context, userID string) []*SessionTokenResult {
	sessionKey := constants.SESSION_PREFIX_BUSINESS + userID
	ctx := c.Request.Context()
	tokens, err := db.Redis.SMembers(ctx, sessionKey).Result()
	if err != nil || len(tokens) == 0 {
		return []*SessionTokenResult{}
	}

	var results []*SessionTokenResult
	for _, token := range tokens {
		tokenKey := constants.TOKEN_PREFIX_BUSINESS + token
		data, err := db.Redis.Get(ctx, tokenKey).Result()
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

		ttl, err := db.Redis.TTL(ctx, tokenKey).Result()
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

func ExitToken(c *gin.Context, param *SessionExitTokenParam) {
	auth.KickoutToken(param.UserID, param.Token)
}

func ChartData(c *gin.Context) *SessionChartData {
	ctx := c.Request.Context()
	bKeys, _ := scanKeys(ctx, db.Redis, constants.SESSION_PREFIX_BUSINESS+"*")
	cKeys, _ := scanKeys(ctx, db.Redis, constants.SESSION_PREFIX_CONSUMER+"*")

	bTotal, _, _ := countTokens(ctx, db.Redis, bKeys, constants.TOKEN_PREFIX_BUSINESS)
	cTotal, _, _ := countTokens(ctx, db.Redis, cKeys, constants.TOKEN_PREFIX_CONSUMER)

	bDaily := countDaily(ctx, db.Redis, bKeys, constants.TOKEN_PREFIX_BUSINESS)
	cDaily := countDaily(ctx, db.Redis, cKeys, constants.TOKEN_PREFIX_CONSUMER)

	days := lastNDays(7)
	bSeries := make([]int, 7)
	cSeries := make([]int, 7)
	for i, day := range days {
		bSeries[i] = bDaily[day]
		cSeries[i] = cDaily[day]
	}

	return &SessionChartData{
		BarChart: BarChartData{
			Days: days,
			Series: []CategorySeries{
				{Name: "BUSINESS", Data: bSeries},
				{Name: "CONSUMER", Data: cSeries},
			},
		},
		PieChart: PieChartData{
			Data: []CategoryTotal{
				{Category: "BUSINESS", Total: bTotal},
				{Category: "CONSUMER", Total: cTotal},
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

