// Package admin 管理端用户中心模块注册（对齐 hei-boot/fastapi profile.admin）。
//
// Author: Charlie
package admin

import (
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/profile"
	"hei-gin/internal/modules/profile/identity"
	"hei-gin/internal/modules/sys/file"
)

// init 自注册 profile.admin。
func init() {
	module.Register("profile.admin", 70, func(d *module.Deps) module.Module {
		idSvc := identity.FromDeps(d)
		s := profile.NewService(d.DB, d.Redis, d.Notify, d.Storage, file.FromDeps(d), d.Runtime,
			d.Audit, security.AccountAdmin, profile.ProfileTableAdmin, "admin", idSvc)
		return module.Module{
			Name:   "profile.admin",
			Order:  70,
			Models: []any{&profile.AdminProfileModel{}},
			Routes: []module.RouteRegistrar{s.AdminRoutes},
		}
	})
}
