// internal/modules/iam/resource/register.go 模块自注册。
//
// Author: Charlie

package resource

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册模块 iam.resource。
func init() {
	module.Register("iam.resource", 40, func(d *module.Deps) module.Module {
		return New(d)
	})
}
