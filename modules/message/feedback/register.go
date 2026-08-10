package feedback

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 message.feedback 模块。
func init() {
	module.Register("message.feedback", 60, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
