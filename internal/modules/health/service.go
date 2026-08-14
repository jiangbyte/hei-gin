// internal/modules/health/service.go 业务服务。
//
// Author: Charlie

package health

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service 健康检查服务。
//
// Author: Charlie
type Service struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewService 构造服务。
func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{db: db, redis: rdb}
}

// New 构建 internal.health 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Redis)
	return module.Module{
		Name:   "internal.health",
		Order:  1,
		Routes: []module.RouteRegistrar{s.registerRoutes},
	}
}

// Live 存活检查。
func (s *Service) Live() LiveResult {
	return LiveResult{Status: "live"}
}

// Ready 就绪检查。
func (s *Service) Ready(ctx context.Context) (ReadyResult, bool) {
	out := ReadyResult{Status: "ready"}
	out.Checks.Database = CheckItem{Enabled: true}
	out.Checks.Redis = CheckItem{Enabled: true}

	if sqlDB, err := s.db.DB(); err != nil {
		detail := err.Error()
		out.Checks.Database.Detail = &detail
	} else if err := sqlDB.PingContext(ctx); err != nil {
		detail := err.Error()
		out.Checks.Database.Detail = &detail
	} else {
		out.Checks.Database.OK = true
		d := "connection ok"
		out.Checks.Database.Detail = &d
	}

	if err := s.redis.Ping(ctx).Err(); err != nil {
		detail := err.Error()
		out.Checks.Redis.Detail = &detail
	} else {
		out.Checks.Redis.OK = true
		d := "connection ok"
		out.Checks.Redis.Detail = &d
	}

	ok := out.Checks.Database.OK && out.Checks.Redis.OK
	if !ok {
		out.Status = "not_ready"
	}
	return out, ok
}
