package account

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册模块 iam.account。
func init() {
	module.Register("iam.account", 20, func(d *module.Deps) module.Module {
		svc := NewService(d.DB)
		d.Provide(shared.AccountFinderKey, svc)
		return New(shared.FromModule(d))
	})
}
