// internal/modules/sys/codegen/register.go 模块自注册。
//
// Author: Charlie

package codegen

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 sys.codegen 模块。
func init() {
	module.Register("sys.codegen", 50, func(d *module.Deps) module.Module {
		return New(d)
	})
}
