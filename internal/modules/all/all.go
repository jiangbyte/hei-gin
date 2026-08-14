// Package all æ±‡æ€» blank import å…¨éƒ¨å†…ç½®ä¸šåŠ¡æ¨¡å—ï¼Œè§¦å‘ init è‡ªæ³¨å†Œã€‚
//
// äºŒæ¬¡å¼€å‘æŽ¨èï¼šæ•´ä»“ Git åˆå¹¶ä¸Šæ¸¸ï¼›cmd ä¿ç•™å¯¹æœ¬åŒ…çš„ blank importï¼Œ
// è‡ªæœ‰æ‰©å±•åœ¨è‡ªå·±çš„ main é‡Œé¢å¤– _ importï¼Œå‡å°‘ä¸Žå®˜æ–¹ all çš„åˆå¹¶å†²çªã€‚
// å¤æ‚åœºæ™¯ç›´æŽ¥æ”¹æœ¬ä»“ framework/ï¼ˆéžå¤–éƒ¨ä¾èµ–å‡çº§æ¨¡åž‹ï¼‰ã€‚
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
