package middleware

import (
	"net/http"
	"strconv"

	"hei-gin/sdk/db"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"
)

// defaultRateLimitWindow is the default rate limit window in seconds.
const defaultRateLimitWindow = 10

// defaultRateLimitMax is the default max requests per window per user.
const defaultRateLimitMax = 30

// RateLimiter returns a Gin middleware that limits requests per user per window.
// Uses Redis INCR + EXPIRE for distributed rate limiting across instances.
//
// Parameters:
//   - endpointKey: unique identifier for the endpoint group
//   - window: time window in seconds (0 = use default 10s)
//   - maxRequests: max requests per window (0 = use default 30)
func RateLimiter(endpointKey string, window int, maxRequests int) gin.HandlerFunc {
	win := window
	if win <= 0 {
		win = defaultRateLimitWindow
	}
	max := maxRequests
	if max <= 0 {
		max = defaultRateLimitMax
	}

	return func(c *gin.Context) {
		// Extract user ID: prefer login_id from context, fallback to client IP
		userID := ""
		if v, exists := c.Get("login_id"); exists {
			if s, ok := v.(string); ok {
				userID = s
			}
		}
		if userID == "" {
			userID = c.ClientIP()
		}

		key := "ratelimit:api:" + endpointKey + ":" + userID

		rdb := db.Redis
		if rdb == nil {
			c.Next()
			return
		}

		// Single Redis round-trip: INCR (creates key if not exists) + EXPIRE on first creation
		// The Lua script ensures atomicity: increment and set expiry only on first creation.
		script := `
			local key = KEYS[1]
			local window = tonumber(ARGV[1])
			local max = tonumber(ARGV[2])
			local current = redis.call("INCR", key)
			if current == 1 then
				redis.call("EXPIRE", key, window)
			end
			return current
		`
		val, err := rdb.Eval(c.Request.Context(), script, []string{key}, win, max).Result()
		if err != nil {
			// Redis error — allow through to avoid blocking legitimate traffic
			c.Next()
			return
		}

		var count int64
		switch v := val.(type) {
		case int64:
			count = v
		case float64:
			count = int64(v)
		case string:
			count, _ = strconv.ParseInt(v, 10, 64)
		default:
			count = 0
		}
		if count > int64(max) {
			c.Abort()
			result.Failure(c, "请求过于频繁，请稍后重试", http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
