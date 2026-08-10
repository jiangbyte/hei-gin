package portal

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 user.portal 模块。
func init() {
	module.Register("user.portal", 70, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
