// internal/modules/sys/feedback/register.go 模块自注册。
//
// Author: Charlie

package feedback

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 sys.feedback 模块。
func init() {
	module.Register("sys.feedback", 60, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
