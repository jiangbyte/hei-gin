// Package portal 门户用户中心模块注册（对齐 hei-boot/fastapi profile.portal）。
//
// Author: Charlie
package portal

import (
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/profile"
	"hei-gin/internal/modules/sys/file"
)

// init 自注册 profile.portal。
func init() {
	module.Register("profile.portal", 71, func(d *module.Deps) module.Module {
		s := profile.NewService(d.DB, d.Redis, d.Notify, d.Storage, file.FromDeps(d), d.Runtime,
			d.Audit, security.AccountPortal, profile.ProfileTablePortal, "portal")
		return module.Module{
			Name:   "profile.portal",
			Order:  71,
			Models: []any{&profile.PortalProfileModel{}},
			Routes: []module.RouteRegistrar{s.PortalRoutes},
		}
	})
}
