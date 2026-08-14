// internal/modules/auth/register.go 模块自注册。
//
// Author: Charlie

package auth

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œæ¨¡å— authã€‚
func init() {
	module.Register("auth", 30, func(d *module.Deps) module.Module {
		var finder AccountFinder
		if v, ok := d.Service(shared.AccountFinderKey); ok {
			finder, _ = v.(AccountFinder)
		}
		return New(shared.FromModule(d), finder)
	})
}
