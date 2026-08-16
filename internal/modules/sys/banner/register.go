// internal/modules/sys/banner/register.go 模块自注册。
//
// Author: Charlie

package banner

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 sys.banner 模块。
func init() {
	module.Register("sys.banner", 50, func(d *module.Deps) module.Module {
		return New(d)
	})
}
