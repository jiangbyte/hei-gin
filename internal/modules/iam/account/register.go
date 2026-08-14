package account

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œæ¨¡å— iam.accountã€‚
func init() {
	module.Register("iam.account", 20, func(d *module.Deps) module.Module {
		svc := NewService(d.DB)
		d.Provide(shared.AccountFinderKey, svc)
		return New(shared.FromModule(d))
	})
}
