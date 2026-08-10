// Package shared 向业务模块注入的依赖视图（与 framework/module.Deps 字段对齐）。
package shared

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/framework/core/config"
	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/module"
	"hei-gin/framework/platform/storage"
)

// Deps 业务模块构造时使用的运行时依赖。
//
// Author: Charlie
type Deps struct {
	Cfg      *config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	Sessions *security.SessionStore
	Perms    *security.PermissionRegistry
	Storage  *storage.Manager
}

// FromModule 从 framework 注册表 Deps 转换。
func FromModule(d *module.Deps) *Deps {
	return &Deps{
		Cfg:      d.Cfg,
		DB:       d.DB,
		Redis:    d.Redis,
		Sessions: d.Sessions,
		Perms:    d.Perms,
		Storage:  d.Storage,
	}
}

// AccountFinderKey 供 auth 从 Deps 服务袋取出账号查找实现。
const AccountFinderKey = "account_finder"
