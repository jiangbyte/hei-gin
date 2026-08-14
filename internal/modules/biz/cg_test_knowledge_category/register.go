// internal/modules/biz/cg_test_knowledge_category/register.go 模块自注册。
//
// Author: Charlie

package cg_test_knowledge_category

import (
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// init è‡ªæ³¨å†Œ biz.cg_test_knowledge_category æ¨¡å—ã€‚
func init() {
	module.Register("biz.cg_test_knowledge_category", 90, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
