package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"hei-gin/sdk/config"
)

var Redis *redis.Client

func InitRedis() error {
	cfg := config.C.Redis
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	poolSize := cfg.MaxConnections
	if poolSize <= 0 {
		poolSize = 100
	}
	dialTimeout := time.Duration(cfg.SocketConnectTimeout) * time.Second
	if dialTimeout <= 0 {
		dialTimeout = 5 * time.Second
	}
	socketTimeout := time.Duration(cfg.SocketTimeout) * time.Second
	if socketTimeout <= 0 {
		socketTimeout = 3 * time.Second
	}
	poolTimeout := time.Duration(cfg.PoolTimeout) * time.Second
	if poolTimeout <= 0 {
		poolTimeout = 4 * socketTimeout
	}
	healthCheckInterval := time.Duration(cfg.HealthCheckInterval) * time.Second
	if healthCheckInterval <= 0 {
		healthCheckInterval = 30 * time.Second
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        cfg.Password,
		DB:              cfg.Database,
		PoolSize:        poolSize,
		MinIdleConns:    poolSize / 10,
		DialTimeout:     dialTimeout,
		ReadTimeout:     socketTimeout,
		WriteTimeout:    socketTimeout,
		PoolTimeout:     poolTimeout,
		ConnMaxIdleTime: healthCheckInterval,
		ConnMaxLifetime: 30 * time.Minute,
		MaxRetries:      retryCount(cfg.RetryOnTimeout),
		MinRetryBackoff: 50 * time.Millisecond,
		MaxRetryBackoff: 500 * time.Millisecond,
	})

	ctx := context.Background()
	if err := Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	log.Printf("[Database] Redis connection verified, pool_size=%d", poolSize)
	return nil
}

func retryCount(enabled bool) int {
	if enabled {
		return 3
	}
	return 0
}

func CloseRedis() {
	if Redis != nil {
		Redis.Close()
		Redis = nil
	}
}
