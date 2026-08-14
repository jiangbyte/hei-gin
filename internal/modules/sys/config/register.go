// internal/modules/sys/config/register.go 模块自注册。
//
// Author: Charlie

package config

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 sys.config 模块。
func init() {
	module.Register("sys.config", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
