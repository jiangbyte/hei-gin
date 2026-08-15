# HEI Gin 文档索引

以根目录 [README.md](../README.md) 为约定与架构的主文档。

| 文档 / 目录 | 说明 |
|-------------|------|
| [README.md](../README.md) | 架构、启动、登录 RSA、SnailJob、二次开发、单体、业务分层、全局 stringly |
| [config.example.yaml](../configs/config.example.yaml) | 启动配置样例（库名 `hei_gin`） |
| [go.mod](../go.mod) | 唯一 Go module：`hei-gin` |
| [framework/](../internal/framework/) | 可修改的运行时包 |
| [framework/core/stringly](../internal/framework/core/stringly/) | 全局 stringly JSON |
| [framework/core/bind](../internal/framework/core/bind/) | `bind.JSON`（入参走 stringly） |
| [framework/platform/snailjob](../internal/framework/platform/snailjob/) | API 内嵌 SnailJob 执行器 |
| [modules/](../internal/modules/) | 业务包目录（同属根 module） |
| [modules/iam/account](../internal/modules/iam/account/) | 业务包分层样板 |
| [modules/profile](../internal/modules/profile/) | 用户中心（admin/portal 共享服务） |
| [app/](../internal/app/) | 装配根：唯一运行入口 `cmd/api`；`modules/all` |
| [scripts/db.sql](../scripts/db.sql) | 建表 + seed（与 hei-boot schema 对齐） |
| [web/](../web/) | 前端（admin） |

## 包职责

| 路径 | 职责 |
|------|------|
| `hei-gin/internal/framework/...` | 配置、安全、中间件、Module 注册表、DB/Redis/存储、stringly/bind、SnailJob |
| `hei-gin/internal/modules/...` | 业务；同包 `param`/`result`/`repo`/`service`/`handler`；`init` 自注册 |
| `hei-gin/internal/modules/shared` | 跨业务共享 |
| `hei-gin/internal/app/...` | 装配根：`cmd` 入口加载、基础设施、模块钩子 |

## 约定摘要

- **单体**：一个 `go.mod`，仓根 `go run ./cmd/api` 直接启动。
- **JSON**：标量 bool/数字线上为字符串；对象与数组不变。业务字段用原生类型，禁止再加 Wire* 包裹。
- **分层**：handler → service → repo；handler 不直接 `db.`。
- **调度**：模块 `Jobs` → SnailJob 客户端；配置键 `snail_job`。
- **同步**：整仓 Git 合并上游，不是依赖坐标升级 framework。