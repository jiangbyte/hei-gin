// internal/modules/user/admin/register.go 模块自注册。
//
// Author: Charlie

package admin

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œ user.admin æ¨¡å—ã€‚
func init() {
	module.Register("user.admin", 70, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
