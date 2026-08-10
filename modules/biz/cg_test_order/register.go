package cg_test_order

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 biz.cg_test_order 模块。
func init() {
	module.Register("biz.cg_test_order", 90, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
