package cg_test_catalog

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œ biz.cg_test_catalog æ¨¡å—ã€‚
func init() {
	module.Register("biz.cg_test_catalog", 90, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
