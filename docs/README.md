# HEI Gin 文档索引

以根目录 [README.md](../README.md) 为约定与架构的主文档。

| 文档 / 目录 | 说明 |
|-------------|------|
| [README.md](../README.md) | 架构、启动、登录 RSA、二次开发、多模块、业务分层、全局 stringly |
| [config.example.yaml](../config.example.yaml) | 启动配置样例（库名 `hei_gin`） |
| [go.work](../go.work) | Maven-like reactor：framework + modules/* + app |
| [framework/](../framework/) | `hei-gin/framework`，可修改的运行时 |
| [framework/core/stringly](../framework/core/stringly/) | 对齐 boot 的全局 stringly JSON |
| [framework/core/bind](../framework/core/bind/) | `bind.JSON`（入参走 stringly） |
| [modules/](../modules/) | 业务 module（各有独立 `go.mod`） |
| [modules/iam/account](../modules/iam/account/) | 业务包分层样板 |
| [app/](../app/) | 组装根：`cmd` + `internal/app` + `modules/all` |
| [migrations/](../migrations/) | goose SQL（仓根） |
| [scripts/migrate.sh](../scripts/migrate.sh) | 迁移 |
| [web/](../web/) | 前端 |

## 包职责

| 路径 / module | 职责 |
|---------------|------|
| `hei-gin/framework` | 配置、安全、中间件、Module 注册表、DB/Redis/存储、stringly/bind |
| `hei-gin/modules/*` | 粗粒度业务；同包 `param`/`result`/`repo`/`service`/`handler`；`init` 自注册 |
| `hei-gin/modules/shared` | 跨业务共享 |
| `hei-gin/app` | 装配根 + cmd；`internal/modules/all` 汇总 blank import |

## 约定摘要

- **JSON**：标量 bool/数字线上为字符串；对象与数组不变。业务字段用原生类型，禁止再加 Wire* 包裹。
- **分层**：handler → service → repo；handler 不直接 `db.`。
- **同步**：整仓 Git 合并上游，不是依赖坐标升级 framework。
