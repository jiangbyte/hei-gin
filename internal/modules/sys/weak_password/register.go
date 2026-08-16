// internal/modules/sys/weak_password/register.go 模块自注册。
//
// Author: Charlie

package weakpassword

import (
	"hei-gin/internal/framework/platform/module"
)

// init 自注册 sys.weak_password 模块。
func init() {
	module.Register("sys.weak_password", 50, func(d *module.Deps) module.Module {
		return New(d)
	})
}
