// internal/modules/sys/file/register.go 模块自注册。
//
// Author: Charlie

package file

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 sys.file 模块。
func init() {
	module.Register("sys.file", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
