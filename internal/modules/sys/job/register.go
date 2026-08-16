// internal/modules/sys/job/register.go 模块自注册。
//
// Author: Charlie

package job

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册 sys.job 模块。
func init() {
	module.Register("sys.job", 30, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
