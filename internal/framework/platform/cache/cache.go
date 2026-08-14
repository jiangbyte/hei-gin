// Package cache 封装 Redis 客户端打开与配置变更频道常量。
//
// Author: Charlie
package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"hei-gin/internal/framework/core/config"
)

// Open 按配置解析 URL 并 Ping 校验后返回 Redis 客户端。
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

// ConfigChangedChannel 为业务配置变更 Pub/Sub 频道。
const ConfigChangedChannel = "hei:config:changed"
