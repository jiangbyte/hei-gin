// Package shared å‘ä¸šåŠ¡æ¨¡å—æ³¨å…¥çš„ä¾èµ–è§†å›¾ï¼ˆå­—æ®µåŒ framework/module.Depsï¼‰ã€‚
package shared

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/notify"
	"hei-gin/internal/framework/platform/storage"
)

// Deps ä¸šåŠ¡æ¨¡å—æž„é€ æ—¶ä½¿ç”¨çš„è¿è¡Œæ—¶ä¾èµ–ã€‚
//
// Author: Charlie
type Deps struct {
	Cfg      *config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	Sessions *security.SessionStore
	Perms    *security.PermissionRegistry
	Storage  *storage.Manager
	Audit    *audit.Queue
	Notify   *notify.Facade
}

// FromModule ä»Ž framework æ³¨å†Œè¡¨ Deps è½¬æ¢ã€‚
func FromModule(d *module.Deps) *Deps {
	return &Deps{
		Cfg:      d.Cfg,
		DB:       d.DB,
		Redis:    d.Redis,
		Sessions: d.Sessions,
		Perms:    d.Perms,
		Storage:  d.Storage,
		Audit:    d.Audit,
		Notify:   d.Notify,
	}
}

// AccountFinderKey ä¾› auth ä»Ž Deps æœåŠ¡è¢‹å–å‡ºè´¦å·æŸ¥æ‰¾å®žçŽ°ã€‚
const AccountFinderKey = "account_finder"
