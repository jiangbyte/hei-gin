package middleware

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/constants"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

const maxNoRepeatBodyBytes int64 = 1 << 20

// NoRepeat returns a middleware that prevents duplicate submissions within the given interval (in milliseconds).
// It uses Redis to store a hash of the request params keyed by userID + IP + URL path.
func NoRepeat(interval int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID (try CONSUMER first, fallback to BUSINESS)
		clientAuth := auth.Consumer
		userID := clientAuth.GetLoginIDDefaultNull(c)
		if userID == "" {
			userID = auth.GetLoginIDDefaultNull(c)
		}

		// Get client IP
		ip := utils.GetClientIP(c)

		// Hash request params
		phash := paramsHash(c)
		cacheKey := constants.NO_REPEAT_PREFIX + ip + ":" + userID + ":" + c.Request.URL.Path + ":" + phash

		// Check Redis
		redisClient := db.Redis
		if redisClient != nil {
			ctx := c.Request.Context()
			cacheTTL := time.Duration(interval) * time.Millisecond
			if cacheTTL <= 0 {
				cacheTTL = time.Second
			}
			ok, err := redisClient.SetNX(ctx, cacheKey, "1", cacheTTL).Result()
			if err == nil && !ok {
				remaining := int64(cacheTTL / time.Second)
				if remaining < 1 {
					remaining = 1
				}
				c.Abort()
				c.JSON(200, gin.H{
					"code":    429,
					"message": "请求过于频繁，请" + strconv.FormatInt(remaining, 10) + "秒后再试",
					"success": false,
				})
				return
			}
		}

		c.Next()
	}
}

// paramsHash generates a deterministic hash from the request's query, form, and body parameters.
func paramsHash(c *gin.Context) string {
	params := make(map[string]interface{})

	// Collect query parameters
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 {
			params[k] = v[0]
		} else {
			params[k] = v
		}
	}

	// Collect form parameters (for POST/PUT/PATCH)
	if c.Request.Method != "GET" {
		if strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/") {
			return hashParams(params)
		}
		_ = c.Request.ParseForm()
		for k, v := range c.Request.PostForm {
			if len(v) == 1 {
				params[k] = v[0]
			} else {
				params[k] = v
			}
		}
	}

	// Read request body and restore it for downstream handlers (Gin v1.12.0 GetRawData does not restore Body)
	if c.Request.ContentLength < 0 || c.Request.ContentLength > maxNoRepeatBodyBytes {
		return hashParams(params)
	}
	bodyReader := http.MaxBytesReader(c.Writer, c.Request.Body, maxNoRepeatBodyBytes)
	if body, err := io.ReadAll(bodyReader); err == nil {
		if len(body) > 0 {
			params["_body"] = string(body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return hashParams(params)
}

func hashParams(params map[string]interface{}) string {
	jsonBytes, _ := json.Marshal(params)
	h := fnv.New64a()
	_, _ = h.Write(jsonBytes)
	return strconv.FormatUint(h.Sum64(), 16)
}
