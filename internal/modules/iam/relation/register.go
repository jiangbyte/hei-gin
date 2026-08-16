// internal/modules/iam/relation/register.go 模块自注册。
//
// Author: Charlie

package relation

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册模块 iam.relation。
func init() {
	module.Register("iam.relation", 40, func(d *module.Deps) module.Module {
		return New(d)
	})
}
