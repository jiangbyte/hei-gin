package portal

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œ user.portal æ¨¡å—ã€‚
func init() {
	module.Register("user.portal", 70, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
