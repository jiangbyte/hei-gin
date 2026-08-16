// internal/modules/internal/health/register.go 模块自注册。
//
// Author: Charlie

package health

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 internal.health 模块。
func init() {
	module.Register("internal.health", 5, func(d *module.Deps) module.Module {
		return New(d)
	})
}
