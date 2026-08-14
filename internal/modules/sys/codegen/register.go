// internal/modules/sys/codegen/register.go 模块自注册。
//
// Author: Charlie

package codegen

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 sys.codegen 模块。
func init() {
	module.Register("sys.codegen", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
