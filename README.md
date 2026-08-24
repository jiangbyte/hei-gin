# HEI Gin

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-4169E1?logo=postgresql&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-Supported-4479A1?logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Supported-DC382D?logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-Apache_2.0-blue)
![Version](https://img.shields.io/badge/version-1.1.0--beta-orange)

**HEI Gin** 是一套面向中后台场景的 Go / Gin 工程化脚手架：单个进程同时提供 **Admin** 与 **Portal** 双端 API，统一认证、权限、运维与消息能力，并与 [hei-boot](https://github.com/jiangbyte/hei-boot)、[hei-fastapi](https://github.com/jiangbyte/hei-fastapi) 等姊妹后端保持契约一致。

> 当前版本：`1.1.0-beta` · 协议：[Apache License 2.0](LICENSE)

## 目录

- [功能特性](#功能特性)
- [前端姊妹项目](#前端姊妹项目)
- [技术栈](#技术栈)
- [工程结构](#工程结构)
- [快速开始](#快速开始)
- [默认账号](#默认账号)
- [姊妹项目](#姊妹项目)
- [License](#license)

## 功能特性

API 前缀统一为 `/api/v1/admin/*` 与 `/api/v1/portal/*`，常见中后台能力按模块划分如下：

| 模块 | 说明 |
| --- | --- |
| 双端账号体系 | ADMIN / PORTAL 独立会话（Redis 不透明 Token）；密码 RSA 传输、验证码登录、失败锁定与限流；可配置三方 OAuth 登录 |
| RBAC 权限 | 账号 / 角色 / 部门 / 用户组 / 岗位；菜单、按钮与 API 资源授权；在线会话踢出 |
| 系统管理 | 字典、动态配置（`sys_config`，敏感项可加密）、Banner、公告 / 通知、意见反馈、弱口令库 |
| 对象存储 | S3 兼容存储（MinIO / RustFS / 阿里云 OSS 等），引擎与凭证走运行时配置，直链或预签名访问 |
| 运维能力 | 操作审计与告警、登录日志、运营工作台概览、内置任务调度（`sys_job`；DB 扫描 + Redis 锁 + cron） |
| 代码生成 | 单表 / 树表 / 主子表方案，预览与 ZIP 下载（含前端与菜单权限 SQL） |
| 实名认证 | 工单提交与审核、敏感字段加密存储（对齐 hei-boot） |
| 业务扩展 | `internal/modules/biz` 示例模块，可按同样模式横向扩展 |

## 前端姊妹项目

| 项目 | 说明 |
| --- | --- |
| [**hei-admin**](https://github.com/jiangbyte/hei-admin) | Vue 3 管理端，对接 `/api/v1/admin/*` |
| [**hei-portal**](https://github.com/jiangbyte/hei-portal) | React 门户，对接 `/api/v1/portal/*` |
| [**hei-admin-uniapp**](https://github.com/jiangbyte/hei-admin-uniapp) | uni-app 管理端移动端 |

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 后端 | Go 1.25+ · Gin · 单 module 单体（`cmd/api`） |
| 持久化 | PostgreSQL / MySQL · GORM（按 `db.driver` 或 DSN 推断） |
| 缓存 / 会话 | Redis（go-redis）· 不透明会话 Token |
| 配置 | Viper（`config.yaml`）+ 运行时 `sys_config` |
| 文档 | swaggo · OpenAPI `/v3/api-docs` · Swagger UI `/swagger-ui/index.html` |
| 其他 | AWS SDK v2（S3）· zap · robfig/cron · snowflake |

## 工程结构

```text
hei-gin/
├── cmd/api/                  # 唯一可启动入口
├── internal/
│   ├── app/                  # 装配根（基础设施 + 模块钩子）
│   ├── framework/            # 运行时（config / security / middleware / storage / gojob …）
│   └── modules/              # 业务模块（auth / iam / sys / profile / workspace / biz …）
├── configs/config.example.yaml
└── scripts/hei_gin.sql       # MySQL 全量建表、种子数据与表/列 COMMENT
```

`scripts/` 目录：

| 文件 | 用途 |
| --- | --- |
| `hei_gin.sql` | MySQL 全量建表、种子数据与表/列 `COMMENT`（`sys_job.handler` 为 Gin 原生 key） |

## 快速开始

### 环境要求

- Go **1.25+**
- MySQL 8+（演示种子）、Redis
- PostgreSQL 亦可（需自行准备 DDL / 数据，配置见 `configs/config.example.yaml`）

### 1. 初始化数据库

**维护原则：** 在 `scripts/hei_gin.sql`（MySQL 8）直接维护表结构、种子数据与表/列 `COMMENT`。与 `hei_boot` 表结构对齐；`sys_job.handler` 使用 Gin 栈标识（如 `sys_job_sample`），非 Boot 的 Java 全限定类名。

**MySQL 8+（本地运行默认）：**

```bash
mysql -u root -p -e "CREATE DATABASE hei_gin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p hei_gin < scripts/hei_gin.sql
```

在 [`configs/config.example.yaml`](configs/config.example.yaml) 中配置 `db.driver` / `db.url`（也可只写 URL，由 scheme 推断）：

```yaml
db:
  driver: mysql
  url: mysql://root:123456@127.0.0.1:3306/hei_gin?charset=utf8mb4&parseTime=true&loc=Local
```

开发配置可复制为根目录 `config.yaml`，按需修改 `db` / `redis` / `auth` 等项。

> 本地 / 演示环境以 SQL 脚本全量初始化。`scripts/hei_gin.sql` 已包含全部表/列中文注释。

### 2. 启动后端

```bash
cp configs/config.example.yaml config.yaml
# 按需修改 db / redis / app 等
go run ./cmd/api
```

| 项 | 地址 |
| --- | --- |
| API | http://127.0.0.1:8000 |
| 接口文档（Swagger UI） | http://127.0.0.1:8000/swagger-ui/index.html |
| OpenAPI JSON | http://127.0.0.1:8000/v3/api-docs |

### 3. 启动前端（可选）

前端为独立仓库，默认将 `/api` 代理到本后端 `http://127.0.0.1:8000`：

```bash
# 管理端 → http://127.0.0.1:5173
git clone https://github.com/jiangbyte/hei-admin.git && cd hei-admin
pnpm install && pnpm dev

# 门户 → http://127.0.0.1:5174
git clone https://github.com/jiangbyte/hei-portal.git && cd hei-portal
pnpm install && pnpm dev
```

详见 [hei-admin](https://github.com/jiangbyte/hei-admin) / [hei-portal](https://github.com/jiangbyte/hei-portal) 各仓库 README。

## 默认账号

| 端 | 前端仓库 | 地址 | 账号 | 密码 | 说明 |
| --- | --- | --- | --- | --- | --- |
| Admin | [hei-admin](https://github.com/jiangbyte/hei-admin) | http://127.0.0.1:5173 | `superadmin` | `123456` | 超级管理员（`*:*:*`） |

> 仅供本地演示。部署后请修改默认密码，并更换配置加密密钥、对象存储凭证等敏感项。更多演示账号与内容种子已写入 `scripts/hei_gin.sql`。

## 姊妹项目

| 项目 | 说明 | 协议 |
| --- | --- | --- |
| [**hei-boot**](https://github.com/jiangbyte/hei-boot) | Spring Boot 脚手架 | Apache License 2.0 |
| [**hei-gin**](https://github.com/jiangbyte/hei-gin) | Go / Gin 后端（本仓库） | Apache License 2.0 |
| [**hei-fastapi**](https://github.com/jiangbyte/hei-fastapi) | FastAPI 后端 | Apache License 2.0 |
| [**hei-admin**](https://github.com/jiangbyte/hei-admin) | Vue 3 管理端前端 | Apache License 2.0 |
| [**hei-portal**](https://github.com/jiangbyte/hei-portal) | React 门户前端 | Apache License 2.0 |
| [**hei-admin-uniapp**](https://github.com/jiangbyte/hei-admin-uniapp) | uni-app 管理端移动端 | Apache License 2.0 |

## License

本项目基于 [Apache License 2.0](LICENSE) 开源。完整条款见 [LICENSE](LICENSE)，版权声明见 [NOTICE](NOTICE).
