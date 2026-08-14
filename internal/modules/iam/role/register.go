// internal/modules/iam/role/register.go 模块自注册。
//
// Author: Charlie

package role

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册模块 iam.role。
func init() {
	module.Register("iam.role", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
