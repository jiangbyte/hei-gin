package codegen

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 sys.codegen 模块。
func init() {
	module.Register("sys.codegen", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
