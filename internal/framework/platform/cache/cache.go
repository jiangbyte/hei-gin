// Package cache å°è£… Redis å®¢æˆ·ç«¯æ‰“å¼€ä¸Žé…ç½®å˜æ›´é¢‘é“å¸¸é‡ã€‚
//
// Author: Charlie
package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"hei-gin/internal/framework/core/config"
)

// Open æŒ‰é…ç½®è§£æž URL å¹¶ Ping æ ¡éªŒåŽè¿”å›ž Redis å®¢æˆ·ç«¯ã€‚
func Open(cfg config.RedisConfig) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	if cfg.MaxConnections > 0 {
		opt.PoolSize = cfg.MaxConnections
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return rdb, nil
}

// ConfigChangedChannel ä¸ºä¸šåŠ¡é…ç½®å˜æ›´ Pub/Sub é¢‘é“ã€‚
const ConfigChangedChannel = "hei:config:changed"
