package permission

import (
	"context"

	"hei-gin/plugins/plugin-sys/shared"

	"github.com/redis/go-redis/v9"
)

type repository struct {
	rdb redis.UniversalClient
}

func (r *repository) Cache(ctx context.Context) (string, error) {
	return r.rdb.Get(ctx, shared.PermissionCacheKey).Result()
}
