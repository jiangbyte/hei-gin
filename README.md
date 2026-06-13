# Hei Gin

> 本项目一直由我独立设计和维护。最初是为了沉淀一套自己在实际业务里反复使用的 Go 后端基础框架，后续也在真实项目中不断裁剪、补强、演进。
>
> 它不是为了追求最小化，而是为了把后台系统、双端认证、权限体系、文件存储、在线会话、IM、日志与生产化基础能力，整理成一套可持续扩展的工程底座。

<img width="120" src="vitepress/docs/public/logo.svg">

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.25+-brightgreen.svg)
![Gin](https://img.shields.io/badge/Gin-1.12+-blue.svg)
![GORM](https://img.shields.io/badge/GORM-1.25+-blue.svg)

## 简介

**Hei Gin** 是一个基于 **Go + Gin + GORM** 的后端开发框架，主要面向后台管理系统、双端 API 服务、权限驱动型业务系统，以及带实时消息能力的场景。

它目前已经实现了一批常用基础能力，这些能力主要分布在“认证、权限、会话、插件装配、运行治理”几个方向：

- 认证与鉴权
- RBAC 与权限范围
- 在线会话与 ACL 刷新
- 文件上传与多存储后端
- 操作日志与链路追踪
- 图形验证码与基础安全能力
- WebSocket IM 与跨实例消息投递
- 插件化业务组织
- 数据库迁移与种子数据
- 健康检查、指标与装配快照

在线文档：<https://jiangbyte.github.io/hei-gin/>（暂时未维护）

## 预览

![](./docs/readme/login.png)

![](./docs/readme/dashboard.png)

![](./docs/readme/home.png)

## 项目背景

Hei Gin 不是从“做一个通用教程项目”出发，而是从真实业务开发的痛点出发逐步沉淀出来的：

- 认证和权限每个项目都要重新搭一遍
- 文件上传、日志记录、在线会话这类能力经常重复实现
- 业务不断变大后，代码很容易退化成大包、大 service、大量隐式依赖
- 很多项目早期能跑，但到后期缺少注册治理、健康检查、观测和会话治理，维护成本会迅速上升

因此，这个项目的目标一直比较明确：

- 不是只解决 CRUD
- 不是只解决启动
- 而是整理一套在实际项目里反复使用的后端工程基础

## 核心特性

### 认证与鉴权

- **双端认证体系**
  - B 端与 C 端分离，统一通过 `auth.Business` 和 `auth.Consumer` 访问
- **Token 会话化**
  - 登录后将角色、权限、Scope 信息写入 claims / session，减少请求时重复查询
- **角色与权限校验**
  - 支持登录校验、权限校验、角色校验
- **Claims / Session 职责分层**
  - claims 负责请求期读取，session 负责在线态管理、刷新与失效控制
- **ACL 刷新机制**
  - 用户、角色权限变更后支持刷新在线会话 ACL
- **显式 Realm 模型**
  - 不再依赖裸字符串 login type，改为显式 realm 常量

### 权限体系

- **RBAC**
  - 用户 -> 角色 -> 权限
- **用户直授权限**
  - 支持用户直接挂权限，不仅依赖角色
- **权限扫描与缓存**
  - 启动时扫描权限，并写入 Redis
- **权限范围与 Scope**
  - Scope 信息支持写入 token session，并随会话一起刷新
- **权限通配符匹配**
  - 支持单级和多级权限匹配
- **对外权限 Provider**
  - 对外提供 `GetPermissionList`、`GetRoleList` 等接口

### 数据权限

- 支持常见数据范围控制模型，包括：
  - `ALL`
  - `SELF`
  - `ORG`
  - `ORG_AND_BELOW`
  - `CUSTOM_ORG`
  - `GROUP`
  - `GROUP_AND_BELOW`
  - `CUSTOM_GROUP`

### 在线会话治理

- B 端、C 端独立会话体系
- 在线用户与会话查询
- 强制下线
- 权限变更后的会话 ACL 刷新
- Session 数据与鉴权逻辑在 auth 模块内处理
- 会话可用于在线查询、强制下线、ACL 刷新与状态核验

### 安全能力

- **SM2**
  - 支持国密 SM2 加解密
- **SM3**
  - 支持摘要与日志签名
- **bcrypt**
  - 密码安全哈希
- **验证码**
  - B 端、C 端独立验证码服务
- **防重复提交**
  - `NoRepeat` 中间件
- **基础登录边界**
  - 公共路径、双端分流、会话校验在认证中间件中处理
- **中间件收口**
  - `sdk/auth/middleware` 处理登录态、权限收集、权限校验与上下文写入

### 文件存储

- 统一文件存储抽象
- 支持：
  - Local
  - MinIO
  - S3
- 支持文件上传与业务文件访问能力
- 支持按配置切换存储后端

### WebSocket / IM

- 内置 IM 插件
- 支持：
  - 单聊
  - 好友关系
  - 群聊
  - 广播
  - 在线状态
- 支持跨实例消息投递
- 支持连接限流、消息限流、去重、心跳、在线人数统计
- WebSocket 安全边界已支持：
  - `allowed_origins`
  - `trusted_proxies`

### 日志与追踪

- **操作日志**
  - 统一记录业务操作
- **日志持久化扩展**
  - 通过 provider / persistence 机制落库
- **trace_id**
  - 全链路透传
- **统一响应结构**
  - 便于问题排查与前后端对接
- **可观测性基础**
  - 当前已接入 HTTP、DB、Redis、WebSocket 基础指标

### 注册与装配

- **插件化组织**
  - `plugin-sys`
  - `plugin-client`
  - `plugin-im`
- **显式装配**
  - 顶层显式调用：
    - `RegisterPlugin()`
    - `RegisterRoutes()`
    - `RegisterMigrations()`
- **注册中心治理**
  - 支持去重
  - 支持冻结
  - 支持快照
  - 支持测试重置
- **装配审计**
  - 调试模式下支持 `/debug/registry`

### 生产化基础

- 严格配置校验
- 数据库与 Redis readiness 检查
- Prometheus 指标输出
- 装配摘要日志
- workspace 级测试验证
- WebSocket 安全边界与代理信任控制
- debug 模式装配快照查看

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
| 可观测性 | Prometheus metrics + health endpoints |

## 当前架构特点

### Go Workspace + 多模块组织

项目不是单一大模块，而是由根应用、`sdk` 和多个业务插件模块共同组成。

### 三大业务插件

- `plugin-sys`
  - 系统管理插件，覆盖用户、角色、组织、资源、权限、日志、公告、文件、会话等
- `plugin-client`
  - C 端用户、认证、会话能力
- `plugin-im`
  - 好友、群组、消息、广播、WebSocket IM 能力

### 从副作用注册收口到显式装配

项目当前以显式装配为主。  
主入口会明确装配 plugin、route 与 migration，而不是只靠 blank import 完成。

### 注册中心已经具备治理能力

当前注册体系不再只是简单地“append 到切片”：

- plugin 注册可去重、冻结、快照
- route / middleware 注册可去重、冻结、快照
- model / seed 注册可去重、冻结、快照

这部分能力主要用于排查装配状态和支持测试。

## 适用场景

- 后台管理系统
- 双端 API 服务
- 需要权限、会话、日志、文件统一底座的业务系统
- 带实时消息、通知、站内信、在线状态的业务场景
- 希望在早期就保留清晰工程边界的项目

## 运行要求

- Go 1.25+
- MySQL 8+
- Redis 6+

项目默认使用 `config.yaml` 运行，关键配置缺失时会直接启动失败。

## 快速开始

### 1. 准备配置

基于 `config.example.yaml` 创建本地配置文件：

- `config.example.yaml` -> `config.yaml`

至少需要正确配置：

- `app.host`
- `app.port`
- `db.host / db.port / db.user / db.database`
- `redis.host / redis.port`
- `token.expire_seconds / token.token_name`

如果启用 WebSocket 并准备上线，建议同时关注：

- `ws.allowed_origins`
- `ws.trusted_proxies`

### 2. 初始化数据库

执行迁移：

```bash
go run cmd/migrate/main.go
```

仅执行结构迁移：

```bash
go run cmd/migrate/main.go -skip-seed
```

### 3. 启动服务

```bash
go run main.go
```

默认基础入口：

- `/`
- `/health/live`
- `/health/ready`
- `/metrics`

当 `app.debug: true` 时，还会开放：

- `/debug/registry`

## 开发方式

常用命令：

- `go run main.go`
- `go run cmd/migrate/main.go`
- `make test`
- `make lint`
- `make build`

如果本地安装了 `air`，也可以使用热重载开发。

## 生产化现状

当前项目已经补上一些偏运行阶段会直接用到的基础项：

- 严格配置校验
- 健康检查
- 指标暴露
- 会话 ACL 刷新
- 显式注册装配
- 注册冻结与装配快照
- workspace 级测试验证
- WebSocket origin / trusted proxy 安全边界
- 基础可观测性与调试入口治理

如果继续往生产环境长期使用推进，仍建议补齐：

- 更完善的指标面板与告警
- 迁移版本记录与执行历史
- 更完整的 CI / staticcheck / 生成校验

## 目录概览

按职责理解即可：

- `sdk/`
  - 框架基础设施与通用能力
- `plugins/plugin-sys/`
  - 系统管理插件
- `plugins/plugin-client/`
  - 客户端业务插件
- `plugins/plugin-im/`
  - 即时消息插件
- `cmd/`
  - 迁移、代码生成等命令行入口
- `vitepress/`
  - 文档站

## 相关项目

- [Hei Boot](https://github.com/jiangbyte/hei-boot)
- [Hei FastAPI](https://github.com/jiangbyte/hei-fastapi)
- [Hei Admin Vue](https://github.com/jiangbyte/hei-admin-vue)

## 协议

[MIT License](LICENSE)
