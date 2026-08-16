// internal/modules/iam/role/register.go 模块自注册。
//
// Author: Charlie

package role

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册模块 iam.role。
func init() {
	module.Register("iam.role", 40, func(d *module.Deps) module.Module {
		return New(d)
	})
}
