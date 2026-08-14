// internal/modules/iam/client/register.go 模块自注册。
//
// Author: Charlie

package client

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init 自注册模块 iam.client。
func init() {
	module.Register("iam.client", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
