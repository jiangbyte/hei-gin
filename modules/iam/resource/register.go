package resource

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册模块 iam.resource。
func init() {
	module.Register("iam.resource", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
