// internal/modules/sys/notice/register.go 模块自注册。
//
// Author: Charlie

package notice

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 sys.notice 模块。
func init() {
	module.Register("sys.notice", 60, func(d *module.Deps) module.Module {
		return New(d)
	})
}
