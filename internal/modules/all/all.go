// Package all 汇总 blank import 全部内置业务模块，触发 init 自注册。
//
// 二次开发推荐：整仓 Git 合并上游；cmd 保留对本包的 blank import，
// 自有扩展在自己的 main 里额外 _ import，减少与官方 all 的合并冲突。
// 复杂场景直接改本仓 framework/（非外部依赖升级模型）。
//
// Author: Charlie
package all

import (
	_ "hei-gin/internal/modules/auth"
	_ "hei-gin/internal/modules/biz/cg_test_activity"
	_ "hei-gin/internal/modules/biz/cg_test_catalog"
	_ "hei-gin/internal/modules/biz/cg_test_knowledge_category"
	_ "hei-gin/internal/modules/biz/cg_test_order"
	_ "hei-gin/internal/modules/dashboard"
	_ "hei-gin/internal/modules/health"
	_ "hei-gin/internal/modules/iam/account"
	_ "hei-gin/internal/modules/iam/client"
	_ "hei-gin/internal/modules/iam/dept"
	_ "hei-gin/internal/modules/iam/group"
	_ "hei-gin/internal/modules/iam/permission"
	_ "hei-gin/internal/modules/iam/position"
	_ "hei-gin/internal/modules/iam/relation"
	_ "hei-gin/internal/modules/iam/resource"
	_ "hei-gin/internal/modules/iam/role"
	_ "hei-gin/internal/modules/message/feedback"
	_ "hei-gin/internal/modules/message/notice"
	_ "hei-gin/internal/modules/sys/audit"
	_ "hei-gin/internal/modules/sys/banner"
	_ "hei-gin/internal/modules/sys/codegen"
	_ "hei-gin/internal/modules/sys/config"
	_ "hei-gin/internal/modules/sys/dict"
	_ "hei-gin/internal/modules/sys/file"
	_ "hei-gin/internal/modules/sys/weak_password"
	_ "hei-gin/internal/modules/user/admin"
	_ "hei-gin/internal/modules/user/portal"
)
