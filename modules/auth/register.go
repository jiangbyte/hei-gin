package auth

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册模块 auth。
func init() {
	module.Register("auth", 30, func(d *module.Deps) module.Module {
		var finder AccountFinder
		if v, ok := d.Service(shared.AccountFinderKey); ok {
			finder, _ = v.(AccountFinder)
		}
		return New(shared.FromModule(d), finder)
	})
}
