# Hei Gin

> 本项目一直由我独立维护。后来我在公司主导某个项目时，为方便起见，直接简化后在项目中采用了该框架，目前已在生产环境中投入使用。
> 框架完全由我设计实现，正在参考优秀开源方案并在此基础上扩展更多能力。
> 后续会持续优化。因工作繁忙，维护时间不定，欢迎提出建议、Issue 或 PR。

<img width="120" src="vitepress/docs/public/logo.svg">

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.25+-brightgreen.svg)
![Gin](https://img.shields.io/badge/Gin-1.12+-blue.svg)
![GORM](https://img.shields.io/badge/GORM-1.25+-blue.svg)

## 简介

**Hei Gin** —— Go + Gin + GORM 构建的快速开发框架，开箱即用。包含完善的权限管理（RBAC）、数据权限、认证授权、文件存储（MinIO/S3/Local）、WebSocket IM、操作日志等功能模块。框架采用前后端分离架构，支持快速搭建后台管理系统和 API 服务。

**在线文档**: [https://jiangbyte.github.io/hei-gin/](https://jiangbyte.github.io/hei-gin/)

## 预览

![](./docs/readme/login.png)

![](./docs/readme/dashboard.png)

![](./docs/readme/home.png)

## 技术栈

| 类型 | 技术 |
| --- | --- |
| 核心框架 | Go 1.25+ / Gin 1.12+ |
| ORM | GORM (gorm.io/gorm) |
| 数据库 | MySQL 8.0+ |
| 缓存 | Redis 6.0+ (go-redis) |
| 认证授权 | Token / SM2 国密加解密 / SM3 哈希 / bcrypt 密码哈希 |
| 文件存储 | MinIO / S3 (AWS SDK) / Local 本地存储 |
| 定时调度 | cron 表达式 (robfig/cron/v3) |
| IP 定位 | IP2Region 离线库 |
| 分布式 ID | Snowflake ID |
| WebSocket | gorilla/websocket + go-redis List/BRPOP 跨实例 IM |

## 核心特性

- **多模块 Go Workspace** — 采用 Go Workspace 组织 SDK、API 接口、业务插件和 App 入口，模块间依赖清晰
- **插件化架构** — 业务模块以插件形式自注册（`module.Register()`），自动发现路由、权限、中间件、定时任务、DB 模型、种子数据
- **双端认证体系** — B 端（后台管理）和 C 端（客户端）独立的两套 Token 认证，通过 `auth.Business` / `auth.Consumer` 统一访问
- **SM2 国密加密** — 登录密码传输使用国密 SM2 C1C3C2 模式加密
- **SM3 哈希** — 操作日志防篡改签名
- **bcrypt 密码哈希** — 存储密码使用 bcrypt 加盐哈希
- **RBAC 权限控制** — 用户→角色→权限 + 用户直授权限，双层模型
- **数据权限（行级）** — 支持 ALL / ORG / ORG_AND_BELOW / SELF / CUSTOM_ORG / GROUP / GROUP_AND_BELOW / CUSTOM_GROUP 等数据范围控制
- **权限自动发现** — 启动时自动扫描注册的权限并缓存到 Redis
- **权限匹配器** — 支持 `*` 单级和 `**` 多级通配符匹配
- **操作日志** — `log.SysLog` 中间件自动记录用户操作，SM3 签名防篡改
- **防重复提交** — `middleware.NoRepeat` 中间件（基于 Redis）
- **链路追踪** — 基于 `trace_id` 的全链路追踪
- **统一验证码** — B 端/C 端独立的图形验证码服务（Redis 存储，300s TTL）
- **统一响应格式** — `{code, message, data, success, trace_id}` 标准结构
- **抽象文件存储** — 统一的 `Engine` 接口，三种后端（Local/MinIO/S3），支持配置切换
- **分片上传** — 大文件分片上传、合并与校验
- **定时调度** — 支持 cron 表达式和固定间隔的后台任务，优雅关闭
- **DB 自动迁移** — `cmd/migrate` 命令行工具自动发现所有模块注册的 Model
- **种子数据** — 模块自注册种子数据，幂等执行
- **模块生命周期** — `module.Register()` 统一管理模块的 Init / Start / Stop
- **模块级配置** — 通过 `config.C.Raw` 读取模块专属配置，无需修改 `config.go`
- **在线会话管理** — B 端和 C 端独立的会话管理，支持在线用户查看和强制下线
- **雪花 ID** — 分布式 Snowflake ID 生成器
- **事件总线** — 内置内存事件总线，支持模块间解耦通信（连接/断开/消息事件）
- **通用 CRUD** — `sdk/crud` 提供通用分页、详情、删除等标准操作函数
- **跨实例 WebSocket IM** — Redis List + BRPOP 驱动的跨实例消息投递、在线状态感知、消息去重、限流、心跳检测

## 项目结构

```
hei-gin/
├── main.go                          # 应用入口（import 插件并调用 app.Run()）
├── go.mod / go.sum / go.work        # Go 模块定义 + Workspace
├── config.example.yaml              # 配置文件示例
├── config.yaml                      # 本地配置文件
│
├── sdk/                             # 框架 SDK（核心基础设施，独立 Go 模块）
│   ├── app/
│   │   ├── app.go                   # 应用工厂：初始化编排（Config→DB→Redis→Auth→Router→Server）
│   │   └── health.go                # 健康检查 Handler
│   ├── auth/                        # 认证与权限系统
│   │   ├── base_auth.go             # 认证基础实现（Business + Consumer）
│   │   ├── auth_tool.go             # B 端（BUSINESS）包级便捷函数
│   │   ├── permission_interface.go  # 权限查询接口定义 + 默认实现
│   │   ├── permission_matcher.go    # 权限通配符匹配器（* / **）
│   │   ├── permission_scan.go       # 权限自动扫描与 Redis 缓存
│   │   ├── permission_tool.go       # 权限/角色查询门面
│   │   ├── module.go                # 认证模块初始化
│   │   ├── middleware/              # 认证与权限中间件
│   │   │   ├── check_login.go              # B 端登录验证
│   │   │   ├── client_check_login.go       # C 端登录验证
│   │   │   ├── check_permission.go         # B 端权限检查
│   │   │   ├── client_check_permission.go  # C 端权限检查
│   │   │   ├── check_role.go               # B 端角色检查
│   │   │   ├── client_check_role.go        # C 端角色检查
│   │   │   └── norepeat.go                 # 防重复提交
│   │   └── pojo/                    # 认证 POJO 定义
│   ├── captcha/                     # 图形验证码生成与验证
│   ├── config/                      # 配置加载（YAML 加载 + 结构体定义）
│   ├── constants/                   # 系统常量与 Redis 键前缀定义
│   ├── crud/                        # 通用 CRUD 帮助函数（分页/详情/删除）
│   ├── db/                          # 数据库初始化（GORM + go-redis）
│   ├── enums/                       # 状态/权限/资源/分页字段枚举
│   ├── eventbus/                    # 内存事件总线（发布/订阅模式）
│   ├── exception/                   # 业务异常 BusinessError（panic-based）
│   ├── log/                         # 操作日志系统（SysLog 中间件 + 工具函数）
│   ├── middleware/                  # 全局中间件
│   │   ├── trace.go                 # 链路追踪（trace_id）
│   │   ├── auth_check.go            # 认证路由分流（B/C/Public）
│   │   ├── recovery.go              # 全局异常恢复（panic→JSON）
│   │   └── cors.go                  # CORS 跨域配置
│   ├── module/                      # 模块生命周期管理（Register/Init/Start/Stop）
│   ├── pojo/                        # 通用 POJO（日期混入、ID 参数）
│   ├── registry/                    # 注册中心（路由/中间件/权限注册与执行）
│   ├── result/                      # 统一响应格式
│   ├── scheduler/                   # 定时调度（cron 表达式 + 固定间隔）
│   ├── storage/                     # 文件存储抽象（Engine 接口 + Local/MinIO/S3 实现）
│   └── utils/                       # 工具函数（雪花 ID、加密、IP 定位、图片等）
│
├── api/                             # API 接口定义层（独立 Go 模块）
│   ├── plugin.go                    # PluginInfo / Plugin 接口定义
│   ├── auth.go                      # 认证/权限相关接口
│   ├── event.go                     # 事件总线接口定义
│   └── log.go                       # 日志持久化接口
│
├── plugins/                         # 业务插件目录（每个插件为独立 Go 模块）
│   ├── plugin-sys/                  # 系统管理插件（B 端）
│   │   ├── plugin.go                # 插件入口：注册模块、权限接口、日志持久化
│   │   ├── imports.go               # 导入所有子模块触发 init() 自注册
│   │   ├── persistence.go           # 日志持久化实现
│   │   ├── provider/                # 权限/用户 Provider 实现
│   │   ├── auth/                    # 认证模块（captcha / sm2 / username）
│   │   ├── banner/                  # Banner 管理
│   │   ├── config/                  # 系统配置管理
│   │   ├── dict/                    # 数据字典管理
│   │   ├── file/                    # 文件上传与管理
│   │   ├── group/                   # 用户组管理
│   │   ├── home/                    # 首页仪表盘
│   │   ├── log/                     # 操作日志查询
│   │   ├── notice/                  # 通知公告
│   │   ├── org/                     # 组织架构
│   │   ├── permission/              # 权限查询
│   │   ├── position/                # 职位管理
│   │   ├── resource/                # 资源管理（菜单/按钮）
│   │   ├── role/                    # 角色管理
│   │   ├── session/                 # 在线会话管理
│   │   └── user/                    # B 端用户管理
│   │
│   ├── plugin-client/               # 客户端插件（C 端）
│   │   ├── plugin.go                # 插件入口
│   │   ├── auth/                    # C 端认证（captcha / sm2 / username）
│   │   ├── session/                 # C 端会话管理
│   │   └── user/                    # C 端用户管理
│   │
│   └── plugin-im/                   # WebSocket IM 插件
│       ├── plugin.go                # 插件入口
│       ├── ws/                      # WebSocket 核心（hub / client / cross_hub）
│       ├── friend/                  # 好友管理
│       ├── group/                   # 群组管理
│       ├── message/                 # 消息管理（含消息/会话/文件）
│       ├── broadcast/               # 广播消息
│       └── model/                   # IM 数据模型
│
├── sdk/app/                         # 应用工厂（Run() 入口）
│   ├── app.go                       # 初始化编排（Config→DB→Redis→Auth→Router→Server）
│   └── health.go                    # 健康检查 Handler
│
├── cmd/                             # 命令行工具
│   ├── migrate/main.go              # DB 迁移 + 种子数据
│   └── codegen/main.go              # 代码生成器（scaffold / list / gen-imports）
│
├── Makefile                         # 开发命令（dev / build / scaffold / migrate ...）
├── .air.toml                        # air 热重载配置
│
├── scripts/                         # 辅助脚本
├── docs/                            # 设计文档
├── uploads/                         # 本地上传文件存储目录
└── vitepress/                       # 在线文档站点（VitePress）
```

## 快速开始

### 环境要求

| 依赖 | 版本要求 |
|------|----------|
| Go | 1.25 或更高版本 |
| MySQL | 8.0 或更高版本 |
| Redis | 6.0 或更高版本 |

### 1. 克隆项目

```bash
git clone https://github.com/jiangbyte/hei-gin.git
cd hei-gin
```

### 2. 创建数据库

```sql
CREATE DATABASE IF NOT EXISTS `hei-gin` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 修改配置

编辑 `config.yaml`，修改数据库和 Redis 连接信息：

```yaml
db:
  host: 127.0.0.1
  port: 3306
  user: root
  password: "your-password"
  database: hei-gin

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  database: 1
```

### 4. 运行项目

```bash
# 方式一：直接运行
go run main.go

# 方式二：热重载（推荐开发时使用）
# 安装 air：go install github.com/air-verse/air@latest
air

# 方式三：使用 Makefile
make dev
```

启动成功后访问 `http://localhost:18885/` 验证服务运行状态。

## 开发工具

### Makefile 命令

| 命令 | 说明 |
|------|------|
| `make dev` | 热重载启动（需安装 air） |
| `make build` | 编译二进制到 `bin/hei-gin` |
| `make run` | 直接运行 |
| `make scaffold name=plugin-xxx` | 创建新插件脚手架 |
| `make list-plugins` | 列出所有插件 |
| `make gen-imports` | 重新生成 `main.go` 中的空白导入 |
| `make migrate` | 执行数据库迁移 |
| `make test` | 运行测试 |
| `make lint` | 代码检查（go vet） |
| `make clean` | 清理构建产物 |

### 插件脚手架

使用代码生成器快速创建插件：

```bash
# 列出已发现的插件
go run cmd/codegen/main.go list

# 创建新插件（自动注册到 main.go）
go run cmd/codegen/main.go scaffold plugin-xxx

# 重新生成 main.go 中的插件导入
go run cmd/codegen/main.go gen-imports
```

生成的插件结构：

```
plugins/plugin-xxx/
├── plugin.go      # Module 实现 + init() 注册
├── model.go       # GORM 模型
├── migrate.go     # db.RegisterModel()
├── params.go      # 请求/响应参数
├── service.go     # 业务逻辑
└── api/v1/
    └── api.go     # 路由 + Handler
```

### 热重载（air）

hei-gin 内置 [air](https://github.com/air-verse/air) 配置，文件变化时自动重编译重启：

```bash
# 安装 air
go install github.com/air-verse/air@latest

# 启动（.air.toml 已预配置）
air
```

配置说明：

| 参数 | 值 |
|------|-----|
| 监控扩展 | `.go`, `.tpl`, `.tmpl`, `.html`, `.yaml`, `.yml` |
| 排除目录 | `.air_tmp`, `.git`, `node_modules`, `vitepress` |
| 延迟 | 1000ms（防频繁触发） |
| 异常时停止 | true |

### 数据库迁移

```bash
# 迁移 + 种子数据
go run cmd/migrate/main.go

# 仅迁移，跳过种子数据
go run cmd/migrate/main.go -skip-seed
```

## WebSocket / 站内信 IM

框架内置跨实例 WebSocket IM 系统，支持实时消息推送、在线状态感知、多实例水平扩展。

### 架构

```
┌─ Instance A ─────────────┐    ┌─ Instance B ─────────────┐
│  CrossHub                 │    │  CrossHub                 │
│  ├─ Local Hub (in-mem)    │    │  ├─ Local Hub (in-mem)    │
│  └─ Redis List BRPOP      │    │  └─ Redis List BRPOP      │
└──────────┬────────────────┘    └──────────┬────────────────┘
           │                                │
           └────────── Redis ───────────────┘
                      │
        ┌─────────────┴─────────────┐
        │  ws:user:{type}:{uid}     │ → 用户→实例映射（Set）
        │  ws:messages:{instance}   │ → 实例消息队列（List）
        │  ws:instance:{instance}   │ → 实例心跳（String + TTL）
        └───────────────────────────┘
```

### WebSocket 端点

| 路径 | 说明 |
|------|------|
| `ws://host:port/api/v1/sys/ws` | B 端（后台管理）WebSocket |
| `ws://host:port/api/v1/c/ws`  | C 端（客户端）WebSocket |

### 事件类型

| 类型 | 方向 | 说明 |
|------|------|------|
| `heartbeat` | Client → Server | 客户端心跳，30s 间隔 |
| `new_message` | Server → Client | 新消息推送 |
| `unread_count` | Server → Client | 通知前端刷新未读数 |
| `presence` | Server → Client | 用户在线/离线状态变更 |
| `online_count` | Server → Client | 在线人数广播（60s） |

### 跨实例特性

| 特性 | 实现 |
|------|------|
| **跨实例消息投递** | Redis List + BRPOP，每个实例独享消息队列 |
| **用户连接追踪** | `ws:user:{type}:{uid}` → Set（用户→实例映射） |
| **消息去重** | Redis SETNX + TTL，防止跨实例重复投递 |
| **限流** | Redis INCR 滑动窗口，默认 10s / 30 条 |
| **心跳检测** | `ws:instance:{id}` 每 15s 刷新 TTL，60s 过期 |
| **过期实例清理** | 后台协程每 5 分钟自动清理僵尸实例 |
| **在线状态广播** | 用户连接/断开时广播 `presence` 事件 |

### 配置

```yaml
ws:
  read_buffer_size: 1024            # WS 读取缓冲区（字节）
  write_buffer_size: 1024           # WS 写入缓冲区（字节）
  heartbeat_interval: 15            # 心跳发送间隔（秒）
  instance_ttl: 60                  # 实例心跳 TTL（秒），超时视为宕机
  stale_clean_interval: 5           # 过期实例清理间隔（分钟）
  rate_limit_window: 10             # 限流时间窗口（秒）
  rate_limit_max: 30                # 窗口内最大消息数
  dedup_ttl: 30                     # 消息去重 TTL（秒）
  poll_timeout: 2                   # Redis BRPOP 超时（秒）
  pong_timeout: 60                  # WS Pong 超时（秒）
  write_timeout: 10                 # WS 写入超时（秒）
  online_broadcast_interval: 60     # 在线人数广播间隔（秒）
```

实例 ID 使用 Snowflake 配置 `snowflake.instance`，生产环境每个实例需配置不同值。

### 前端连接

```typescript
const wsUrl = `ws://${host}/api/v1/sys/ws?token=${token}`
const ws = new WebSocket(wsUrl)
// 心跳：每 30s 发送 { type: "heartbeat" }
// 重连：指数退避，最多 10 次
```

## 库引用说明

使用框架 SDK 中的组件时，引用路径统一为 `hei-gin/sdk/...`：

```go
import (
    "hei-gin/sdk/config"
    "hei-gin/sdk/auth"
    authMiddleware "hei-gin/sdk/auth/middleware"
    "hei-gin/sdk/middleware"
    "hei-gin/sdk/result"
    "hei-gin/sdk/exception"
    "hei-gin/sdk/log"
    "hei-gin/sdk/module"
    "hei-gin/sdk/scheduler"
    "hei-gin/sdk/storage"
    "hei-gin/sdk/db"
    "hei-gin/sdk/crud"
    "hei-gin/sdk/utils"
)
```

## 配置管理

```go
import "hei-gin/sdk/config"

// 读取标准配置
dbHost := config.C.DB.Host
redisPort := config.C.Redis.Port

// 读取模块专属配置（无需修改 config.go）
rawVal := config.C.Raw["my_module"]
```

## 相关项目

- **[Hei Boot](https://github.com/jiangbyte/hei-boot)** — Java Spring Boot 单体版本
- **[Hei FastAPI](https://github.com/jiangbyte/hei-fastapi)** — Python FastAPI 单体版本
- **[Hei Admin Vue](https://github.com/jiangbyte/hei-admin-vue)** — Vue3 前端管理后台

## 开源协议

本项目采用 [MIT License](LICENSE) 开源协议
