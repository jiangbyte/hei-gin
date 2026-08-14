// internal/modules/user/portal/register.go 模块自注册。
//
// Author: Charlie

package portal

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 user.portal 模块。
func init() {
	module.Register("user.portal", 70, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
