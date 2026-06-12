package permission

import (
	"context"

	"hei-gin/sdk/constants"
	"hei-gin/sdk/db"
)

func Cache(ctx context.Context) (string, error) {
	return db.Redis.Get(ctx, constants.PERMISSION_CACHE_KEY).Result()
}
