package health

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 internal.health 模块。
func init() {
	module.Register("internal.health", 5, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
