// internal/modules/iam/dept/register.go 模块自注册。
//
// Author: Charlie

package dept

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œæ¨¡å— iam.deptã€‚
func init() {
	module.Register("iam.dept", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
