// internal/modules/workspace/register.go 模块自注册。
//
// Author: Charlie

package workspace

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 workspace 模块。
func init() {
	module.Register("workspace", 50, func(d *module.Deps) module.Module {
		return New(d)
	})
}
