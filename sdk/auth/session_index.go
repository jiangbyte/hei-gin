package auth

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionInfo struct {
	UserID                string
	Username              string
	SessionCreateTime     string
	SessionTimeoutSeconds int
	TokenCount            int
}

type SessionStats struct {
	TotalCount        int
	OneHourNewlyAdded int
	MaxTokenCount     int
}

type SessionTokenInfo struct {
	Token          string
	CreatedAt      string
	TimeoutSeconds int
	DeviceType     string
	DeviceID       string
}

const sessionCleanupBatchSize int64 = 500

type parsedTokenData struct {
	CreatedAt       time.Time
	CreatedAtString string
	Username        string
	DeviceType      string
	DeviceID        string
}

func (t *baseAuthTool) getSessionIndexKey() string {
	return "hei:auth:" + string(t.realmID) + ":session:index"
}

func (t *baseAuthTool) getSessionExpiryKey() string {
	return "hei:auth:" + string(t.realmID) + ":session:expiry"
}

func (t *baseAuthTool) getSessionCountKey() string {
	return "hei:auth:" + string(t.realmID) + ":session:counts"
}

func (t *baseAuthTool) getTokenCreatedIndexKey() string {
	return "hei:auth:" + string(t.realmID) + ":token:created"
}

func (t *baseAuthTool) getTokenExpiryIndexKey() string {
	return "hei:auth:" + string(t.realmID) + ":token:expiry"
}

func (t *baseAuthTool) getTokenOwnerKey() string {
	return "hei:auth:" + string(t.realmID) + ":token:owners"
}

func (t *baseAuthTool) trackLoginSession(ctx context.Context, userID, token string, createdAt time.Time, ttl time.Duration) {
	redisClient := t.getRedis()
	if redisClient == nil || userID == "" || token == "" {
		return
	}
	count, _ := redisClient.SCard(ctx, t.getSessionKey(userID)).Result()
	if count < 1 {
		count = 1
	}
	expiresAt := createdAt.Add(ttl).Unix()
	pipe := redisClient.Pipeline()
	pipe.ZAddGT(ctx, t.getSessionIndexKey(), redis.Z{Score: float64(createdAt.Unix()), Member: userID})
	pipe.ZAddGT(ctx, t.getSessionExpiryKey(), redis.Z{Score: float64(expiresAt), Member: userID})
	pipe.ZAddGT(ctx, t.getSessionCountKey(), redis.Z{Score: float64(count), Member: userID})
	pipe.ZAdd(ctx, t.getTokenCreatedIndexKey(), redis.Z{Score: float64(createdAt.Unix()), Member: token})
	pipe.ZAdd(ctx, t.getTokenExpiryIndexKey(), redis.Z{Score: float64(expiresAt), Member: token})
	pipe.HSet(ctx, t.getTokenOwnerKey(), token, userID)
	_, _ = pipe.Exec(ctx)
}

func (t *baseAuthTool) untrackToken(ctx context.Context, userID, token string) {
	redisClient := t.getRedis()
	if redisClient == nil || token == "" {
		return
	}
	if userID == "" {
		if owner, err := redisClient.HGet(ctx, t.getTokenOwnerKey(), token).Result(); err == nil {
			userID = owner
		}
	}
	pipe := redisClient.Pipeline()
	if userID != "" {
		pipe.SRem(ctx, t.getSessionKey(userID), token)
	}
	pipe.ZRem(ctx, t.getTokenCreatedIndexKey(), token)
	pipe.ZRem(ctx, t.getTokenExpiryIndexKey(), token)
	pipe.HDel(ctx, t.getTokenOwnerKey(), token)
	_, _ = pipe.Exec(ctx)
	if userID != "" {
		t.refreshSessionIndexes(ctx, userID)
	}
}

func (t *baseAuthTool) refreshSessionIndexes(ctx context.Context, userID string) {
	redisClient := t.getRedis()
	if redisClient == nil || userID == "" {
		return
	}
	tokens, err := redisClient.SMembers(ctx, t.getSessionKey(userID)).Result()
	if err != nil || len(tokens) == 0 {
		t.removeSessionIndexes(ctx, userID)
		return
	}

	liveTokens, tokenData := t.loadTokenData(ctx, tokens)
	staleTokens := diffTokens(tokens, liveTokens)
	if len(staleTokens) > 0 {
		pipe := redisClient.Pipeline()
		pipe.SRem(ctx, t.getSessionKey(userID), stringArgs(staleTokens)...)
		pipe.ZRem(ctx, t.getTokenCreatedIndexKey(), stringArgs(staleTokens)...)
		pipe.ZRem(ctx, t.getTokenExpiryIndexKey(), stringArgs(staleTokens)...)
		pipe.HDel(ctx, t.getTokenOwnerKey(), staleTokens...)
		_, _ = pipe.Exec(ctx)
	}
	if len(liveTokens) == 0 {
		t.removeSessionIndexes(ctx, userID)
		return
	}

	maxCreated := time.Now()
	for i, token := range liveTokens {
		if data, ok := tokenData[token]; ok && (i == 0 || data.CreatedAt.After(maxCreated)) {
			maxCreated = data.CreatedAt
		}
	}
	ttl, _ := redisClient.TTL(ctx, t.getSessionKey(userID)).Result()
	if ttl <= 0 {
		t.removeSessionIndexes(ctx, userID)
		return
	}
	expiresAt := time.Now().Add(ttl).Unix()
	pipe := redisClient.Pipeline()
	pipe.ZAdd(ctx, t.getSessionIndexKey(), redis.Z{Score: float64(maxCreated.Unix()), Member: userID})
	pipe.ZAdd(ctx, t.getSessionExpiryKey(), redis.Z{Score: float64(expiresAt), Member: userID})
	pipe.ZAdd(ctx, t.getSessionCountKey(), redis.Z{Score: float64(len(liveTokens)), Member: userID})
	for _, token := range liveTokens {
		pipe.HSet(ctx, t.getTokenOwnerKey(), token, userID)
	}
	_, _ = pipe.Exec(ctx)
}

func (t *baseAuthTool) removeSessionIndexes(ctx context.Context, userID string) {
	redisClient := t.getRedis()
	if redisClient == nil || userID == "" {
		return
	}
	pipe := redisClient.Pipeline()
	pipe.ZRem(ctx, t.getSessionIndexKey(), userID)
	pipe.ZRem(ctx, t.getSessionExpiryKey(), userID)
	pipe.ZRem(ctx, t.getSessionCountKey(), userID)
	_, _ = pipe.Exec(ctx)
}

func (t *baseAuthTool) removeTokenIndexes(ctx context.Context, tokens []string) {
	redisClient := t.getRedis()
	if redisClient == nil || len(tokens) == 0 {
		return
	}
	pipe := redisClient.Pipeline()
	pipe.ZRem(ctx, t.getTokenCreatedIndexKey(), stringArgs(tokens)...)
	pipe.ZRem(ctx, t.getTokenExpiryIndexKey(), stringArgs(tokens)...)
	pipe.HDel(ctx, t.getTokenOwnerKey(), tokens...)
	_, _ = pipe.Exec(ctx)
}

func (t *baseAuthTool) updateTokenExpiryIndex(ctx context.Context, userID, token string, ttl time.Duration) {
	redisClient := t.getRedis()
	if redisClient == nil || userID == "" {
		return
	}
	expiresAt := time.Now().Add(ttl).Unix()
	pipe := redisClient.Pipeline()
	if token != "" {
		pipe.ZAdd(ctx, t.getTokenExpiryIndexKey(), redis.Z{Score: float64(expiresAt), Member: token})
		pipe.HSet(ctx, t.getTokenOwnerKey(), token, userID)
	}
	pipe.ZAdd(ctx, t.getSessionExpiryKey(), redis.Z{Score: float64(expiresAt), Member: userID})
	_, _ = pipe.Exec(ctx)
}

func (t *baseAuthTool) cleanupExpiredIndexes(ctx context.Context) {
	redisClient := t.getRedis()
	if redisClient == nil {
		return
	}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	expiredUsers, _ := redisClient.ZRangeByScore(ctx, t.getSessionExpiryKey(), &redis.ZRangeBy{
		Min: "-inf", Max: now, Count: sessionCleanupBatchSize,
	}).Result()
	if len(expiredUsers) > 0 {
		pipe := redisClient.Pipeline()
		pipe.ZRem(ctx, t.getSessionIndexKey(), stringArgs(expiredUsers)...)
		pipe.ZRem(ctx, t.getSessionExpiryKey(), stringArgs(expiredUsers)...)
		pipe.ZRem(ctx, t.getSessionCountKey(), stringArgs(expiredUsers)...)
		_, _ = pipe.Exec(ctx)
	}

	expiredTokens, _ := redisClient.ZRangeByScore(ctx, t.getTokenExpiryIndexKey(), &redis.ZRangeBy{
		Min: "-inf", Max: now, Count: sessionCleanupBatchSize,
	}).Result()
	if len(expiredTokens) == 0 {
		return
	}
	owners, _ := redisClient.HMGet(ctx, t.getTokenOwnerKey(), expiredTokens...).Result()
	usersToRefresh := make(map[string]struct{})
	pipe := redisClient.Pipeline()
	for i, token := range expiredTokens {
		if i < len(owners) {
			if owner, ok := owners[i].(string); ok && owner != "" {
				pipe.SRem(ctx, t.getSessionKey(owner), token)
				usersToRefresh[owner] = struct{}{}
			}
		}
	}
	pipe.ZRem(ctx, t.getTokenCreatedIndexKey(), stringArgs(expiredTokens)...)
	pipe.ZRem(ctx, t.getTokenExpiryIndexKey(), stringArgs(expiredTokens)...)
	pipe.HDel(ctx, t.getTokenOwnerKey(), expiredTokens...)
	_, _ = pipe.Exec(ctx)
	for userID := range usersToRefresh {
		t.refreshSessionIndexes(ctx, userID)
	}
}

func (t *baseAuthTool) listSessionInfos(ctx context.Context, current, size int, keyword string) ([]SessionInfo, int64, error) {
	redisClient := t.getRedis()
	if redisClient == nil {
		return nil, 0, nil
	}
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	t.cleanupExpiredIndexes(ctx)

	userIDs, total, err := t.pageIndexedUserIDs(ctx, current, size, keyword)
	if err != nil {
		return nil, 0, err
	}
	infos, err := t.hydrateSessionInfos(ctx, userIDs)
	if err != nil {
		return nil, 0, err
	}
	return infos, total, nil
}

func (t *baseAuthTool) listSessionInfosByUserIDs(ctx context.Context, userIDs []string, current, size int) ([]SessionInfo, int64, error) {
	redisClient := t.getRedis()
	if redisClient == nil || len(userIDs) == 0 {
		return []SessionInfo{}, 0, nil
	}
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	t.cleanupExpiredIndexes(ctx)

	deduped := uniqueNonEmptyStrings(userIDs)
	if len(deduped) == 0 {
		return []SessionInfo{}, 0, nil
	}

	const chunkSize = 1000
	type scoredUser struct {
		UserID string
		Score  float64
	}
	matches := make([]scoredUser, 0, len(deduped))
	for start := 0; start < len(deduped); start += chunkSize {
		end := start + chunkSize
		if end > len(deduped) {
			end = len(deduped)
		}
		chunk := deduped[start:end]
		scores, err := redisClient.ZMScore(ctx, t.getSessionIndexKey(), chunk...).Result()
		if err != nil {
			return nil, 0, err
		}
		for i, score := range scores {
			if score == 0 || math.IsNaN(score) {
				continue
			}
			matches = append(matches, scoredUser{UserID: chunk[i], Score: score})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score == matches[j].Score {
			return matches[i].UserID < matches[j].UserID
		}
		return matches[i].Score > matches[j].Score
	})

	total := int64(len(matches))
	start := (current - 1) * size
	if start >= len(matches) {
		return []SessionInfo{}, total, nil
	}
	end := start + size
	if end > len(matches) {
		end = len(matches)
	}

	pageIDs := make([]string, 0, end-start)
	for _, row := range matches[start:end] {
		pageIDs = append(pageIDs, row.UserID)
	}
	infos, err := t.hydrateSessionInfos(ctx, pageIDs)
	if err != nil {
		return nil, 0, err
	}
	return infos, total, nil
}

func (t *baseAuthTool) pageIndexedUserIDs(ctx context.Context, current, size int, keyword string) ([]string, int64, error) {
	redisClient := t.getRedis()
	if keyword != "" {
		score, err := redisClient.ZScore(ctx, t.getSessionIndexKey(), keyword).Result()
		if err == redis.Nil {
			return []string{}, 0, nil
		}
		if err != nil {
			return nil, 0, err
		}
		if score == 0 {
			return []string{}, 0, nil
		}
		if current > 1 {
			return []string{}, 1, nil
		}
		return []string{keyword}, 1, nil
	}
	total, _ := redisClient.ZCard(ctx, t.getSessionIndexKey()).Result()
	start := int64((current - 1) * size)
	userIDs, err := redisClient.ZRevRange(ctx, t.getSessionIndexKey(), start, start+int64(size)-1).Result()
	return userIDs, total, err
}

func (t *baseAuthTool) hydrateSessionInfos(ctx context.Context, userIDs []string) ([]SessionInfo, error) {
	redisClient := t.getRedis()
	if redisClient == nil || len(userIDs) == 0 {
		return []SessionInfo{}, nil
	}
	pipe := redisClient.Pipeline()
	memberCmds := make([]*redis.StringSliceCmd, len(userIDs))
	ttlCmds := make([]*redis.DurationCmd, len(userIDs))
	for i, userID := range userIDs {
		memberCmds[i] = pipe.SMembers(ctx, t.getSessionKey(userID))
		ttlCmds[i] = pipe.TTL(ctx, t.getSessionKey(userID))
	}
	_, _ = pipe.Exec(ctx)

	tokenKeys := make([]string, 0)
	userTokens := make(map[string][]string, len(userIDs))
	for i, userID := range userIDs {
		tokens, err := memberCmds[i].Result()
		if err != nil || len(tokens) == 0 {
			t.removeSessionIndexes(ctx, userID)
			continue
		}
		userTokens[userID] = tokens
		for _, token := range tokens {
			tokenKeys = append(tokenKeys, t.getTokenKey(token))
		}
	}

	tokenData := make(map[string]parsedTokenData, len(tokenKeys))
	if len(tokenKeys) > 0 {
		values, _ := redisClient.MGet(ctx, tokenKeys...).Result()
		for i, val := range values {
			raw, ok := val.(string)
			if !ok || i >= len(tokenKeys) {
				continue
			}
			token := strings.TrimPrefix(tokenKeys[i], t.getTokenKey(""))
			if data, ok := parseTokenData(raw); ok {
				tokenData[token] = data
			}
		}
	}

	infos := make([]SessionInfo, 0, len(userIDs))
	for i, userID := range userIDs {
		tokens := userTokens[userID]
		if len(tokens) == 0 {
			continue
		}
		liveTokens := make([]string, 0, len(tokens))
		staleTokens := make([]string, 0)
		var minCreated time.Time
		var maxCreated time.Time
		username := ""
		for _, token := range tokens {
			data, ok := tokenData[token]
			if !ok {
				staleTokens = append(staleTokens, token)
				continue
			}
			liveTokens = append(liveTokens, token)
			if minCreated.IsZero() || data.CreatedAt.Before(minCreated) {
				minCreated = data.CreatedAt
			}
			if maxCreated.IsZero() || data.CreatedAt.After(maxCreated) {
				maxCreated = data.CreatedAt
			}
			if username == "" {
				username = data.Username
			}
		}
		if len(staleTokens) > 0 {
			t.removeStaleTokens(ctx, userID, staleTokens)
		}
		if len(liveTokens) == 0 {
			t.removeSessionIndexes(ctx, userID)
			continue
		}
		timeoutSeconds := -1
		if ttl, err := ttlCmds[i].Result(); err == nil {
			timeoutSeconds = int(ttl.Seconds())
		}
		if timeoutSeconds <= 0 {
			t.removeSessionIndexes(ctx, userID)
			continue
		}
		t.updateSessionSummary(ctx, userID, maxCreated, len(liveTokens), time.Duration(timeoutSeconds)*time.Second, liveTokens)
		infos = append(infos, SessionInfo{
			UserID:                userID,
			Username:              username,
			SessionCreateTime:     minCreated.Format("2006-01-02 15:04:05"),
			SessionTimeoutSeconds: timeoutSeconds,
			TokenCount:            len(liveTokens),
		})
	}
	return infos, nil
}

func (t *baseAuthTool) updateSessionSummary(ctx context.Context, userID string, maxCreated time.Time, tokenCount int, ttl time.Duration, tokens []string) {
	redisClient := t.getRedis()
	if redisClient == nil || userID == "" {
		return
	}
	pipe := redisClient.Pipeline()
	pipe.ZAdd(ctx, t.getSessionIndexKey(), redis.Z{Score: float64(maxCreated.Unix()), Member: userID})
	pipe.ZAdd(ctx, t.getSessionExpiryKey(), redis.Z{Score: float64(time.Now().Add(ttl).Unix()), Member: userID})
	pipe.ZAdd(ctx, t.getSessionCountKey(), redis.Z{Score: float64(tokenCount), Member: userID})
	for _, token := range tokens {
		pipe.HSet(ctx, t.getTokenOwnerKey(), token, userID)
	}
	_, _ = pipe.Exec(ctx)
}

func (t *baseAuthTool) removeStaleTokens(ctx context.Context, userID string, tokens []string) {
	redisClient := t.getRedis()
	if redisClient == nil || len(tokens) == 0 {
		return
	}
	pipe := redisClient.Pipeline()
	if userID != "" {
		pipe.SRem(ctx, t.getSessionKey(userID), stringArgs(tokens)...)
	}
	pipe.ZRem(ctx, t.getTokenCreatedIndexKey(), stringArgs(tokens)...)
	pipe.ZRem(ctx, t.getTokenExpiryIndexKey(), stringArgs(tokens)...)
	pipe.HDel(ctx, t.getTokenOwnerKey(), tokens...)
	_, _ = pipe.Exec(ctx)
}

func (t *baseAuthTool) sessionStats(ctx context.Context) (SessionStats, error) {
	redisClient := t.getRedis()
	if redisClient == nil {
		return SessionStats{}, nil
	}
	t.cleanupExpiredIndexes(ctx)
	total, _ := redisClient.ZCard(ctx, t.getTokenCreatedIndexKey()).Result()
	oneHourAgo := strconv.FormatInt(time.Now().Add(-time.Hour).Unix(), 10)
	oneHour, _ := redisClient.ZCount(ctx, t.getTokenCreatedIndexKey(), oneHourAgo, "+inf").Result()

	var maxTokenCount int
	rows, _ := redisClient.ZRevRangeWithScores(ctx, t.getSessionCountKey(), 0, 0).Result()
	if len(rows) > 0 {
		maxTokenCount = int(rows[0].Score)
	}
	return SessionStats{TotalCount: int(total), OneHourNewlyAdded: int(oneHour), MaxTokenCount: maxTokenCount}, nil
}

func (t *baseAuthTool) sessionDailyCounts(ctx context.Context, days []string) map[string]int {
	redisClient := t.getRedis()
	result := make(map[string]int, len(days))
	if redisClient == nil {
		return result
	}
	t.cleanupExpiredIndexes(ctx)
	for _, day := range days {
		start, err := time.ParseInLocation("2006-01-02", day, time.Local)
		if err != nil {
			continue
		}
		end := start.AddDate(0, 0, 1).Add(-time.Second)
		count, _ := redisClient.ZCount(ctx, t.getTokenCreatedIndexKey(),
			strconv.FormatInt(start.Unix(), 10),
			strconv.FormatInt(end.Unix(), 10),
		).Result()
		result[day] = int(count)
	}
	return result
}

func (t *baseAuthTool) sessionTokens(ctx context.Context, userID string) ([]SessionTokenInfo, error) {
	redisClient := t.getRedis()
	if redisClient == nil || userID == "" {
		return []SessionTokenInfo{}, nil
	}
	tokens, err := redisClient.SMembers(ctx, t.getSessionKey(userID)).Result()
	if err != nil || len(tokens) == 0 {
		return []SessionTokenInfo{}, nil
	}

	tokenKeys := make([]string, len(tokens))
	for i, token := range tokens {
		tokenKeys[i] = t.getTokenKey(token)
	}
	values, _ := redisClient.MGet(ctx, tokenKeys...).Result()

	pipe := redisClient.Pipeline()
	ttlCmds := make([]*redis.DurationCmd, len(tokens))
	for i, key := range tokenKeys {
		ttlCmds[i] = pipe.TTL(ctx, key)
	}
	_, _ = pipe.Exec(ctx)

	result := make([]SessionTokenInfo, 0, len(tokens))
	staleTokens := make([]string, 0)
	for i, token := range tokens {
		raw, ok := values[i].(string)
		if !ok {
			staleTokens = append(staleTokens, token)
			continue
		}
		data, ok := parseTokenData(raw)
		if !ok {
			staleTokens = append(staleTokens, token)
			continue
		}
		timeoutSeconds := -1
		if ttl, err := ttlCmds[i].Result(); err == nil {
			timeoutSeconds = int(ttl.Seconds())
		}
		result = append(result, SessionTokenInfo{
			Token:          token,
			CreatedAt:      data.CreatedAtString,
			TimeoutSeconds: timeoutSeconds,
			DeviceType:     data.DeviceType,
			DeviceID:       data.DeviceID,
		})
	}
	if len(staleTokens) > 0 {
		t.removeStaleTokens(ctx, userID, staleTokens)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})
	t.refreshSessionIndexes(ctx, userID)
	return result, nil
}

func parseTokenData(raw string) (parsedTokenData, bool) {
	var tokenData map[string]any
	if err := json.Unmarshal([]byte(raw), &tokenData); err != nil {
		return parsedTokenData{}, false
	}
	createdAtStr, _ := tokenData["created_at"].(string)
	createdAt := time.Now()
	if createdAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			createdAt = t
		}
	}
	result := parsedTokenData{CreatedAt: createdAt, CreatedAtString: createdAtStr}
	if extra, ok := tokenData["extra"].(map[string]any); ok {
		result.Username, _ = extra["username"].(string)
		result.DeviceType, _ = extra["device_type"].(string)
		result.DeviceID, _ = extra["device_id"].(string)
	}
	return result, true
}

func diffTokens(tokens []string, live []string) []string {
	liveSet := make(map[string]struct{}, len(live))
	for _, token := range live {
		liveSet[token] = struct{}{}
	}
	stale := make([]string, 0)
	for _, token := range tokens {
		if _, ok := liveSet[token]; !ok {
			stale = append(stale, token)
		}
	}
	return stale
}

func stringArgs(values []string) []any {
	args := make([]any, len(values))
	for i, value := range values {
		args[i] = value
	}
	return args
}

func uniqueNonEmptyStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func (t *baseAuthTool) loadTokenData(ctx context.Context, tokens []string) ([]string, map[string]parsedTokenData) {
	redisClient := t.getRedis()
	if redisClient == nil || len(tokens) == 0 {
		return nil, nil
	}
	keys := make([]string, len(tokens))
	for i, token := range tokens {
		keys[i] = t.getTokenKey(token)
	}
	values, _ := redisClient.MGet(ctx, keys...).Result()
	live := make([]string, 0, len(tokens))
	dataMap := make(map[string]parsedTokenData, len(tokens))
	for i, value := range values {
		raw, ok := value.(string)
		if !ok {
			continue
		}
		if data, ok := parseTokenData(raw); ok {
			live = append(live, tokens[i])
			dataMap[tokens[i]] = data
		}
	}
	return live, dataMap
}
