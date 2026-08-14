package health

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œ internal.health æ¨¡å—ã€‚
func init() {
	module.Register("internal.health", 5, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
