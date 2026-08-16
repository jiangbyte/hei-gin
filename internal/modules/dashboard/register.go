// internal/modules/dashboard/register.go 模块自注册。
//
// Author: Charlie

package dashboard

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 dashboard 模块。
func init() {
	module.Register("dashboard", 80, func(d *module.Deps) module.Module {
		return New(d)
	})
}
