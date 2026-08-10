package notice

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 message.notice 模块。
func init() {
	module.Register("message.notice", 60, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
