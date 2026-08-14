package weakpassword

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œ sys.weak_password æ¨¡å—ã€‚
func init() {
	module.Register("sys.weak_password", 50, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
