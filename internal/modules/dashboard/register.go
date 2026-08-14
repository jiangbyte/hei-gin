// internal/modules/dashboard/register.go 模块自注册。
//
// Author: Charlie

package dashboard

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 dashboard 模块。
func init() {
	module.Register("dashboard", 80, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
