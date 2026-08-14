// internal/modules/sys/dict/register.go 模块自注册。
//
// Author: Charlie

package dict

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 sys.dict 模块。
func init() {
	module.Register("sys.dict", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
