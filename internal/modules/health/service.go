package health

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// Service å¥åº·æ£€æŸ¥æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewService æž„é€ æœåŠ¡ã€‚
func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{db: db, redis: rdb}
}

// New æž„å»º internal.health æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Redis)
	return module.Module{
		Name:   "internal.health",
		Order:  1,
		Routes: []module.RouteRegistrar{s.registerRoutes},
	}
}

// Live å­˜æ´»æ£€æŸ¥ã€‚
func (s *Service) Live() LiveResult {
	return LiveResult{Status: "live"}
}

// Ready å°±ç»ªæ£€æŸ¥ã€‚
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
