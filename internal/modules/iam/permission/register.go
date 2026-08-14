// internal/modules/iam/permission/register.go 模块自注册。
//
// Author: Charlie

package permission

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册模块 iam.permission。
func init() {
	module.Register("iam.permission", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
