# hei-gin ↔ hei-boot 契约对齐清单

以 **hei-boot** 为 API/数据模型真源，**hei-admin** `src/api/**` 为前端验收基准。

## 已对齐模块

| 模块 | 状态 | 说明 |
|------|------|------|
| auth | 已对齐 | 登录/会话/OAuth/手机找回密码/site-footer |
| workspace | 已对齐 | 替代 dashboard；overview + shortcuts |
| profile | 已对齐 | MeResult 含 identity；移除 name |
| profile/identity | 已对齐 | 实名认证全链路 |
| sys/audit | 已对齐 | admin + portal my-page/my-detail |
| iam/* | 基本对齐 | 账号/RBAC/组织/资源 |
| sys/* | 基本对齐 | 字典/配置/文件/任务/代码生成等 |
| easyTrans | 排除 | boot 生态依赖，gin/fastapi 不实现 |

## 验收脚本

```bash
# 三栈统一入口（仓库根目录）
python scripts/compare_stacks.py

# 单栈
python hei-gin/scripts/compare-modules.py
python hei-fastapi/scripts/e2e/compare_modules.py --boot http://127.0.0.1:8000 --fastapi http://127.0.0.1:8100
```

## API 路径（新增/变更）

| 方法 | 路径 | 替代 |
|------|------|------|
| GET | `/api/v1/admin/workspace/overview` | `/dashboard/overview`（已删除） |
| GET/POST | `/api/v1/admin/workspace/shortcuts` | 新增 |
| GET | `/api/v1/public/site-footer` | 新增 |
| POST | `/api/v1/{admin\|portal}/forgot-password/phone` | 新增 |
| POST | `/api/v1/{admin\|portal}/reset-password/phone` | 新增 |
| GET | `/api/v1/{admin\|portal}/profile/identity/status` | 新增 |
| GET/POST | `/api/v1/{admin\|portal}/real-name/case/*` | 新增 |
| GET/POST | `/api/v1/admin/sys/real-name-case/*` | 新增 |
| GET/POST | `/api/v1/admin/sys/identity/*` | 新增 |
| GET | `/api/v1/{admin\|portal}/sys/audit/my-page` | 新增 |
| GET | `/api/v1/{admin\|portal}/sys/audit/my-detail` | 新增 |

## 数据库

- 全量脚本与 hei-boot 同步：`scripts/db.mysql.sql`（MySQL 主库）；`scripts/db.sql` 为 PostgreSQL 方言 CI
- 活库导出：`hei-boot/scripts/export-db-mysql.sh` → `sync-db-mysql-to-stacks.py`
- 增量迁移：`scripts/migration/20260822_real_name_identity.*.sql`
- 新增表：`profile_identity`、`real_name_case`、`real_name_case_record`、`sys_workspace_shortcut`
- `profile_user_*.name` 列已移除

## 成熟包

| 能力 | 包 |
|------|-----|
| Fernet | `github.com/fernet/fernet-go` |
| OAuth (GitHub) | `golang.org/x/oauth2` + `golang.org/x/oauth2/github` |
| 验证码 | 保留 SVG 自研（hei-admin 前端硬编码 `image/svg+xml`） |

## 验收

```bash
# 模块运行时 shape 对比（Boot:8000 vs Gin:8001）
python scripts/compare-modules.py
python scripts/compare-modules.py --module sys/audit

# OpenAPI 全量契约
python scripts/compare-contract-full.py --output scripts/reports/boot_gin_full_diff.json

go run ./scripts/e2e --base http://127.0.0.1:8000
go run ./scripts/smoke
```

模块状态看板：`docs/alignment-status.json`

前端联调：hei-admin 指向 `http://127.0.0.1:8000`。
