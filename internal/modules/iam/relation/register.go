// internal/modules/iam/relation/register.go 模块自注册。
//
// Author: Charlie

package relation

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œæ¨¡å— iam.relationã€‚
func init() {
	module.Register("iam.relation", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
