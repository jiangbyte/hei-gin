// internal/modules/iam/group/register.go 模块自注册。
//
// Author: Charlie

package group

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册模块 iam.group。
func init() {
	module.Register("iam.group", 40, func(d *module.Deps) module.Module {
		return New(d)
	})
}
