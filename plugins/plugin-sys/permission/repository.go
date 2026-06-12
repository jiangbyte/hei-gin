package permission

import (
	"context"

	"hei-gin/sdk/constants"

	"github.com/redis/go-redis/v9"
)

type repository struct {
	rdb redis.UniversalClient
}

func (r *repository) Cache(ctx context.Context) (string, error) {
	return r.rdb.Get(ctx, constants.PERMISSION_CACHE_KEY).Result()
}
