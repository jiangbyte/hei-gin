package permission

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册模块 iam.permission。
func init() {
	module.Register("iam.permission", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
