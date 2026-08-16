// internal/modules/iam/dept/register.go 模块自注册。
//
// Author: Charlie

package dept

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册模块 iam.dept。
func init() {
	module.Register("iam.dept", 40, func(d *module.Deps) module.Module {
		return New(d)
	})
}
