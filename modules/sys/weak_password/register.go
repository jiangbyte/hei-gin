package weakpassword

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 sys.weak_password 模块。
func init() {
	module.Register("sys.weak_password", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
