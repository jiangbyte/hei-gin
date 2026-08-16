// internal/modules/sys/audit/register.go 模块自注册。
//
// Author: Charlie

package audit

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 sys.audit 模块。
func init() {
	module.Register("sys.audit", 50, func(d *module.Deps) module.Module {
		return New(d)
	})
}
