# Hei Gin

> 本项目最初来自一次公司内部媒体项目重构过程中的后端框架沉淀，由我主导设计和实现。后续从内部分支中剥离出来，结合自己在实际业务里反复使用的能力进行整理、裁剪和补强，最终以独立仓库的形式开源。
>
> 它不是为了做一个最小示例，而是希望把后台系统、双端认证、权限体系、文件存储、在线会话、IM、日志和插件化组织方式，整理成一套可以继续扩展的 Go 后端工程底座。

<img width="120" src="vitepress/docs/public/logo.svg">

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.25+-brightgreen.svg)
![Gin](https://img.shields.io/badge/Gin-1.12+-blue.svg)
![GORM](https://img.shields.io/badge/GORM-1.25+-blue.svg)

## 简介

**Hei Gin** 是 HEI 快速开发框架的 Go 单体应用版本，基于 **Go + Gin + GORM** 构建，采用 Go Workspace 多模块组织方式。项目面向后台管理系统、双端 API 服务、权限驱动型业务系统，以及带实时消息、通知、在线状态能力的业务场景。

当前仓库由根应用、`sdk` 和三个业务插件组成：

- `sdk`：配置、认证、权限、中间件、数据库、Redis、存储、日志、事件、调度、健康检查与指标等通用能力。
- `plugins/plugin-sys`：B 端系统管理插件，包含用户、角色、组织、职位、用户组、资源、权限、字典、配置、公告、Banner、文件、日志、会话和首页能力。
- `plugins/plugin-client`：C 端用户插件，包含客户端用户、登录注册、验证码、会话和个人资料能力。
- `plugins/plugin-im`：IM 插件，包含 WebSocket、单聊、好友、群聊、广播、文件消息、在线状态和跨实例消息投递。

在线文档：<https://jiangbyte.github.io/hei-gin/>（暂时未维护）

## 预览

![](./docs/readme/login.png)

![](./docs/readme/dashboard.png)

![](./docs/readme/home.png)

## 项目背景

Hei Gin 不是从“做一个通用教程项目”出发，而是从真实业务开发中反复出现的问题沉淀出来的：

- 认证、权限、菜单、角色和数据范围每个后台项目都需要重新搭建。
- 文件上传、操作日志、验证码、在线会话、强制下线等能力经常重复实现。
- B 端和 C 端用户模型、认证入口、权限边界容易混在一起。
- 业务不断变大后，代码很容易退化成大包、大 service 和大量隐式依赖。
- 模块装配、路由注册、权限发现、初始化顺序如果没有治理，后期排查成本会越来越高。

因此这个项目的目标不是只解决 CRUD，也不是只解决服务启动，而是整理一套在实际项目里可复用的后端工程基础。

## 功能概览

### 认证与会话

- B 端、C 端独立认证体系，对应 `auth.Business` 和 `auth.Consumer`。
- Token 会话存储在 Redis，服务端可主动失效。
- 登录后写入角色、权限和 Scope 快照，减少请求期间重复查询。
- 支持在线会话列表、用户维度会话查询、Token 查询、强制下线。
- 用户或角色权限变化后，可刷新在线会话 ACL。
- WebSocket 连接在升级阶段通过 token 完成用户识别。

### 权限与数据范围

- RBAC 模型：用户 -> 角色 -> 权限。
- 支持用户直授权限，不完全依赖角色。
- 支持角色授权资源、权限授权和用户授权。
- 支持权限扫描、权限模块查询和按模块查询权限。
- 支持 `*` 单级与 `**` 多级通配符权限匹配。
- 支持常见数据范围模型：
  - `ALL`
  - `SELF`
  - `ORG`
  - `ORG_AND_BELOW`
  - `CUSTOM_ORG`
  - `GROUP`
  - `GROUP_AND_BELOW`
  - `CUSTOM_GROUP`

### 系统管理插件

`plugin-sys` 提供后台管理常用模块：

- 用户管理：创建、修改、删除、详情、分页、角色授权、权限授权、个人资料、头像、密码修改。
- 角色管理：角色分页、授权权限、授权资源、刷新角色相关在线会话 ACL。
- 组织、职位、用户组：支持树形结构和基础维护。
- 资源与模块：管理菜单、按钮、模块和资源树。
- 字典与配置：系统字典、业务字典、系统配置项。
- 日志：操作日志记录、查询、删除、图表数据。
- 文件：普通上传、分片上传、下载、文件记录管理。
- 公告、Banner：首页展示与公开查询。
- 首页：统计信息、快捷操作、常用资源。
- 会话：在线用户、会话分析、Token 管理和强制下线。

### C 端用户插件

`plugin-client` 提供面向客户端用户的基础能力：

- C 端用户名密码登录、注册、登出。
- C 端验证码和 SM2 公钥接口。
- C 端用户分页、详情、创建、修改、删除。
- 当前用户信息、资料更新、头像更新、密码修改。
- C 端会话统计、在线会话分页、Token 查询与强制退出。

### WebSocket / IM 插件

`plugin-im` 提供即时消息相关能力：

- WebSocket 长连接接入。
- B 端和 C 端用户连接分流。
- 在线状态、在线人数、心跳、连接数限制。
- 单聊消息发送、查询、搜索、撤回、转发、删除、已读。
- 会话列表、会话消息、会话已读、获取或创建会话。
- 好友申请、同意、拒绝、删除、拉黑、取消拉黑、好友搜索。
- 群聊创建、成员邀请、退出、解散、禁言、成员列表。
- 群消息发送、撤回、搜索、已读、群会话列表。
- 广播消息、已读记录和广播列表。
- IM 文件上传和文件消息。
- Redis 跨实例消息投递、消息去重、消息限流。

### 文件存储

- 统一 `storage.Engine` 抽象。
- 支持 Local、MinIO、S3。
- 支持流式上传。
- 支持分片上传初始化、上传分片、完成、取消。
- 上传文件写入系统文件记录，便于后续查询和管理。

### 日志、追踪与响应

- `trace_id` 链路标识，支持从请求头读取或自动生成。
- Recovery 中间件统一处理 panic 和业务错误。
- 统一响应结构，包含 `code`、`message`、`data`、`success`、`trace_id`。
- 操作日志通过持久化接口写入系统日志表。
- 暴露 HTTP、DB、Redis、WebSocket 相关 Prometheus 指标。
- 提供健康检查和 readiness 检查。

### 插件与装配

- 根应用显式装配插件：
  - `RegisterPlugin()`
  - `RegisterRoutes()`
  - `RegisterMigrations()`
- 插件注册、路由注册、迁移注册都有去重、冻结和快照能力。
- `app.debug: true` 时可查看 `/debug/registry` 装配快照。
- 插件可以独立组织模型、仓储、服务、API 和迁移。

## 技术栈

| 类别 | 说明 |
| --- | --- |
| 核心框架 | Go 1.25+, Gin 1.12+ |
| ORM | GORM |
| 数据库 | MySQL |
| 缓存 | Redis |
| 存储 | Local / MinIO / S3 |
| 安全 | Token / SM2 / SM3 / bcrypt |
| 调度 | robfig/cron |
| WebSocket | gorilla/websocket |
| 指标 | Prometheus |
| 文档 | Swagger / VitePress |

## 项目结构

```text
hei-gin
├── main.go                         # 应用入口，显式装配插件并启动
├── config.yaml                     # 本地运行配置
├── config.example.yaml             # 配置示例
├── scripts/
│   └── hei.sql                     # 初始化 SQL
├── cmd/
│   ├── migrate/                    # 数据迁移入口
│   └── codegen/                    # 代码生成入口
├── sdk/                            # 通用 SDK 模块
│   ├── auth/                       # 认证、会话、权限、Scope、权限扫描
│   ├── captcha/                    # 图形验证码
│   ├── config/                     # 配置加载与校验
│   ├── infra/                      # DB、Redis、存储、事件、调度
│   ├── kernel/                     # app、plugin、registry
│   ├── log/                        # 操作日志与日志工具
│   ├── observability/              # Prometheus 指标
│   ├── utils/                      # 加密、雪花 ID、时间、IP、UA 等工具
│   └── web/                        # 中间件、响应、异常
├── plugins/
│   ├── plugin-sys/                 # B 端系统管理插件
│   ├── plugin-client/              # C 端用户插件
│   └── plugin-im/                  # IM 插件
├── docs/                           # Swagger 生成文件与 README 图片
└── vitepress/                      # 文档站
```

## 路由约定

| 路径 | 说明 |
| --- | --- |
| `/api/v1/public/b/*` | B 端公开接口，如登录、验证码、SM2 公钥 |
| `/api/v1/public/c/*` | C 端公开接口，如登录、注册、验证码 |
| `/api/v1/sys/*` | B 端系统管理接口 |
| `/api/v1/client/*` | B 端管理 C 端用户/会话接口 |
| `/api/v1/c/*` | C 端登录后接口 |
| `/api/v1/sys/im/*` | B 端 IM 接口 |
| `/api/v1/c/im/*` | C 端 IM 接口 |
| `/uploads/:bucket/:file_key` | 本地文件访问入口 |
| `/health/live` | 存活检查 |
| `/health/ready` | 依赖检查 |
| `/metrics` | Prometheus 指标 |
| `/swagger/index.html` | Swagger UI，取决于配置是否启用 |

## 快速开始

### 1. 环境要求

- Go 1.25+
- MySQL 8+
- Redis 6+

### 2. 准备配置

复制配置示例：

```bash
cp config.example.yaml config.yaml
```

至少需要配置：

- `app.host`
- `app.port`
- `db.host`
- `db.port`
- `db.user`
- `db.database`
- `redis.host`
- `redis.port`
- `token.expire_seconds`
- `token.token_name`
- `snowflake.instance`

### 3. 初始化数据库

执行迁移和种子数据：

```bash
go run cmd/migrate/main.go
```

只执行表结构迁移：

```bash
go run cmd/migrate/main.go -skip-seed
```

也可以按需导入 SQL：

```bash
mysql -u root -p hei < scripts/hei.sql
```

### 4. 启动服务

```bash
go run main.go
```

默认入口：

- `/`
- `/health/live`
- `/health/ready`
- `/metrics`

开启 Swagger 后可访问：

- `/swagger/index.html`

## 常用命令

```bash
# 启动
go run main.go

# 数据迁移
go run cmd/migrate/main.go

# 跳过种子数据
go run cmd/migrate/main.go -skip-seed

# 测试
make test

# 构建
make build

# 代码检查
make lint
```

如果本地安装了 `air`，也可以使用热重载开发。

## API 响应格式

普通响应：

```json
{
  "code": 200,
  "message": "请求成功",
  "data": {},
  "success": true,
  "trace_id": "trace-id"
}
```

分页响应：

```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "records": [],
    "total": 100,
    "current": 1,
    "size": 20,
    "pages": 5
  },
  "success": true,
  "trace_id": "trace-id"
}
```

## 模块开发约定

业务模块通常按垂直切片组织：

```text
<module>/
├── model.go        # 数据模型
├── params.go       # 请求与响应结构
├── repository.go   # 数据访问
├── service.go      # 业务逻辑
├── module.go       # 模块组装
└── api/v1/api.go   # 路由和 Handler
```

插件入口通常负责三类事情：

- 注册插件生命周期。
- 注册路由。
- 注册迁移模型和种子数据。

主入口 [main.go](main.go) 中显式装配当前启用的插件。

## 相关项目

- [Hei Boot](https://github.com/jiangbyte/hei-boot)
- [Hei FastAPI](https://github.com/jiangbyte/hei-fastapi)
- [Hei Admin Vue](https://github.com/jiangbyte/hei-admin-vue)

## 协议

[MIT License](LICENSE)
