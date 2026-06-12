package middleware

import (
	"net/http"
	"strconv"

	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

const defaultRateLimitWindow = 10
const defaultRateLimitMax = 30

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
