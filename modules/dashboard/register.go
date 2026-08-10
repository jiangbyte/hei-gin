package dashboard

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 dashboard 模块。
func init() {
	module.Register("dashboard", 80, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
