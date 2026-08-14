// internal/modules/dashboard/register.go 模块自注册。
//
// Author: Charlie

package dashboard

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œ dashboard æ¨¡å—ã€‚
func init() {
	module.Register("dashboard", 80, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
