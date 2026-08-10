package config

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 sys.config 模块。
func init() {
	module.Register("sys.config", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
