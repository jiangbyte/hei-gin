// internal/modules/iam/client/register.go 模块自注册。
//
// Author: Charlie

package client

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册模块 iam.client。
func init() {
	module.Register("iam.client", 40, func(d *module.Deps) module.Module {
		return New(d)
	})
}
