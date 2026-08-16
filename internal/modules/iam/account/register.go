// internal/modules/iam/account/register.go 模块自注册。
//
// Author: Charlie

package account

import (
	"hei-gin/internal/framework/platform/module"
)

// AccountFinderKey 供 auth 从 Deps 服务袋取出账号查找实现。
const AccountFinderKey = "account_finder"

// init 自注册模块 iam.account。
func init() {
	module.Register("iam.account", 20, New)
}
