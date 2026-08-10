package relation

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册模块 iam.relation。
func init() {
	module.Register("iam.relation", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
