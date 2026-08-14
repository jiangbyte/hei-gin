// internal/modules/iam/group/register.go 模块自注册。
//
// Author: Charlie

package group

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œæ¨¡å— iam.groupã€‚
func init() {
	module.Register("iam.group", 40, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
