package position

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册模块 iam.position。
func init() {
	module.Register("iam.position", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
