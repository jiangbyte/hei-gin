package auth

import (
	"testing"

	"hei-gin/sdk/infra/db"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupPermissionScanTest(t *testing.T) {
	t.Helper()
	mr := miniredis.RunT(t)
	db.Redis = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	permissionRegistry = nil
	t.Cleanup(func() {
		permissionRegistry = nil
		if db.Redis != nil {
			_ = db.Redis.Close()
			db.Redis = nil
		}
	})
}

func TestBuildModuleTreeUsesEntryModuleOrCodePrefix(t *testing.T) {
	tree := buildModuleTree([]PermissionEntry{
		{Code: "sys:user:view", Name: "view"},
		{Code: "custom:export", Module: "custom", Name: "export"},
	})

	if _, ok := tree["sys:user"]["sys:user:view"]; !ok {
		t.Fatalf("missing sys:user:view in tree: %#v", tree)
	}
	if _, ok := tree["custom"]["custom:export"]; !ok {
		t.Fatalf("missing custom:export in tree: %#v", tree)
	}
}

func TestRunPermissionScanAndReadBackModulesAndPermissions(t *testing.T) {
	setupPermissionScanTest(t)

	RegisterPermission(PermissionEntry{Code: "sys:user:view", Name: "查看用户"})
	RegisterPermission(PermissionEntry{Code: "sys:user:edit", Name: "编辑用户"})
	RegisterPermission(PermissionEntry{Code: "sys:role:view", Name: "查看角色"})

	if err := RunPermissionScan(); err != nil {
		t.Fatalf("run permission scan: %v", err)
	}

	modules, err := GetModulesFromRedis()
	if err != nil {
		t.Fatalf("get modules: %v", err)
	}
	if len(modules) != 2 || modules[0] != "sys:role" || modules[1] != "sys:user" {
		t.Fatalf("modules = %#v", modules)
	}

	perms, err := GetPermissionsByModuleFromRedis("sys:user")
	if err != nil {
		t.Fatalf("get perms: %v", err)
	}
	if len(perms) != 2 {
		t.Fatalf("perms = %#v", perms)
	}
}
