// internal/modules/sys/config/register.go 模块自注册。
//
// Author: Charlie

package config

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 sys.config 模块。
func init() {
	module.Register("sys.config", 50, func(d *module.Deps) module.Module {
		return New(d)
	})
}
