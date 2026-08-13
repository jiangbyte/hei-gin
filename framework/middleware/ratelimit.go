package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"hei-gin/framework/core/response"
)

// RateLimit 基于 Redis INCR + EXPIRE 的固定窗口限流。
//
// Author: Charlie
func RateLimit(rdb *redis.Client, keyPrefix string, permits int, windowSeconds int) gin.HandlerFunc {
	if permits <= 0 {
		permits = 60
	}
	if windowSeconds <= 0 {
		windowSeconds = 60
	}
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s:%s", keyPrefix, ip)
		ctx := c.Request.Context()
		n, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if n == 1 {
			_ = rdb.Expire(ctx, key, time.Duration(windowSeconds)*time.Second).Err()
		}
		if n > int64(permits) {
			response.Fail(c, http.StatusTooManyRequests, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
