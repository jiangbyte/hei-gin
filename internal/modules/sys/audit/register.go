// internal/modules/sys/audit/register.go 模块自注册。
//
// Author: Charlie

package audit

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œ sys.audit æ¨¡å—ã€‚
func init() {
	module.Register("sys.audit", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
