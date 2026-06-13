package provider

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/constants"
	"hei-gin/sdk/infra/db"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupProviderRedisTest(t *testing.T) {
	t.Helper()
	mr := miniredis.RunT(t)
	db.Redis = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() {
		if db.Redis != nil {
			_ = db.Redis.Close()
			db.Redis = nil
		}
	})
}

func TestGetAllPermissionsFromRedisSupportsFlatList(t *testing.T) {
	setupProviderRedisTest(t)

	payload, _ := json.Marshal([]string{"sys:user:view", "sys:user:edit"})
	if err := db.Redis.Set(context.Background(), constants.PERMISSION_CACHE_KEY, payload, 0).Err(); err != nil {
		t.Fatalf("set redis: %v", err)
	}

	perms, err := (&PermissionProvider{}).getAllPermissionsFromRedis(context.Background())
	if err != nil {
		t.Fatalf("get all permissions: %v", err)
	}
	want := []string{"sys:user:view", "sys:user:edit"}
	if !reflect.DeepEqual(perms, want) {
		t.Fatalf("perms = %#v, want %#v", perms, want)
	}
}

func TestGetAllPermissionsFromRedisSupportsModuleTree(t *testing.T) {
	setupProviderRedisTest(t)

	tree := map[string]map[string]auth.PermissionEntry{
		"sys:user": {
			"sys:user:edit": {Code: "sys:user:edit", Name: "编辑用户"},
			"sys:user:view": {Code: "sys:user:view", Name: "查看用户"},
		},
	}
	payload, _ := json.Marshal(tree)
	if err := db.Redis.Set(context.Background(), constants.PERMISSION_CACHE_KEY, payload, 0).Err(); err != nil {
		t.Fatalf("set redis: %v", err)
	}

	perms, err := (&PermissionProvider{}).getAllPermissionsFromRedis(context.Background())
	if err != nil {
		t.Fatalf("get all permissions: %v", err)
	}
	want := []string{"sys:user:edit", "sys:user:view"}
	if !reflect.DeepEqual(perms, want) {
		t.Fatalf("perms = %#v, want %#v", perms, want)
	}
}
