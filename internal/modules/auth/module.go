// internal/modules/auth/module.go 模块自注册（对齐 hei-boot auth 模块）。
//
// Author: Charlie

package auth

import (
	"hei-gin/internal/framework/platform/module"
)

const accountFinderKey = "account_finder"

// init 自注册 auth 模块；order 25 保证 iam.account（20）先 Provide account_finder。
func init() {
	module.Register("auth", 25, func(d *module.Deps) module.Module {
		var finder AccountFinder
		if v, ok := d.Service(accountFinderKey); ok {
			finder = v.(AccountFinder)
		}
		return New(d, finder)
	})
}
