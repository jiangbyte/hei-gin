// internal/modules/iam/position/register.go 模块自注册。
//
// Author: Charlie

package position

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册模块 iam.position。
func init() {
	module.Register("iam.position", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
