// internal/modules/profile/register.go 模块自注册（admin/portal 共享一套服务实现）。
//
// Author: Charlie

package profile

import (
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// NewAdmin 构建 profile.admin 模块。
func NewAdmin(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Redis, d.Notify, d.Storage, security.AccountAdmin, ProfileTableAdmin, "admin")
	return module.Module{
		Name:   "profile.admin",
		Order:  70,
		Models: []any{&AdminProfileModel{}},
		Routes: []module.RouteRegistrar{s.adminRoutes},
	}
}

// NewPortal 构建 profile.portal 模块。
func NewPortal(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Redis, d.Notify, d.Storage, security.AccountPortal, ProfileTablePortal, "portal")
	return module.Module{
		Name:   "profile.portal",
		Order:  71,
		Models: []any{&PortalProfileModel{}},
		Routes: []module.RouteRegistrar{s.registerRoutes},
	}
}

// init 自注册 profile.admin / profile.portal。
func init() {
	module.Register("profile.admin", 70, func(d *module.Deps) module.Module {
		return NewAdmin(shared.FromModule(d))
	})
	module.Register("profile.portal", 71, func(d *module.Deps) module.Module {
		return NewPortal(shared.FromModule(d))
	})
}
