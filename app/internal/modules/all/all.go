// Package all 汇总 blank import 全部内置业务模块，触发 init 自注册。
//
// 二次开发推荐：整仓 Git 合并上游；cmd 保留对本包的 blank import，
// 自有扩展在自己的 main 里额外 _ import，减少与官方 all 的合并冲突。
// 复杂场景直接改本仓 framework/（非外部依赖升级模型）。
package all

import (
	_ "hei-gin/modules/auth"
	_ "hei-gin/modules/biz/cg_test_activity"
	_ "hei-gin/modules/biz/cg_test_catalog"
	_ "hei-gin/modules/biz/cg_test_knowledge_category"
	_ "hei-gin/modules/biz/cg_test_order"
	_ "hei-gin/modules/dashboard"
	_ "hei-gin/modules/health"
	_ "hei-gin/modules/iam/account"
	_ "hei-gin/modules/iam/client"
	_ "hei-gin/modules/iam/dept"
	_ "hei-gin/modules/iam/group"
	_ "hei-gin/modules/iam/permission"
	_ "hei-gin/modules/iam/position"
	_ "hei-gin/modules/iam/relation"
	_ "hei-gin/modules/iam/resource"
	_ "hei-gin/modules/iam/role"
	_ "hei-gin/modules/message/feedback"
	_ "hei-gin/modules/message/notice"
	_ "hei-gin/modules/sys/audit"
	_ "hei-gin/modules/sys/banner"
	_ "hei-gin/modules/sys/codegen"
	_ "hei-gin/modules/sys/config"
	_ "hei-gin/modules/sys/dict"
	_ "hei-gin/modules/sys/file"
	_ "hei-gin/modules/sys/weak_password"
	_ "hei-gin/modules/user/admin"
	_ "hei-gin/modules/user/portal"
)
