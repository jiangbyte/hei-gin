package position

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œæ¨¡å— iam.positionã€‚
func init() {
	module.Register("iam.position", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
