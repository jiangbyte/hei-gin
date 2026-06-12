# 项目结构

Hei Gin 项目采用 **Go Workspace 多模块架构**，将 SDK 框架、API 接口定义、业务插件和应用入口组织为独立的 Go 模块，保持了良好的关注点分离和依赖管理。

## 多模块架构概览

```
hei-gin/                   (根模块)
├── main.go                # 应用入口：导入插件，调用 sdk/app.Run()
├── go.mod                 # 根模块依赖 sdk
├── go.work                # Workspace 定义：聚合所有子模块
│
├── sdk/                   # 框架 SDK（独立 Go 模块：hei-gin/sdk）
├── sdk/contracts/         # 公共契约（位于 SDK 内）
├── plugins/plugin-sys/    # 系统管理插件（独立 Go 模块）
├── plugins/plugin-client/ # 客户端插件（独立 Go 模块）
├── plugins/plugin-im/     # IM 插件（独立 Go 模块）
├── app/                   # 应用组装层（独立 Go 模块：hei-gin/app）
└── cmd/                   # 命令行工具
```

## 完整目录树

```
hei-gin/
├── main.go                          # 应用入口，初始化并启动服务
├── config.example.yaml              # 配置文件示例
├── go.mod / go.sum / go.work        # Go 模块依赖 + Workspace
│
├── sdk/                             # 框架 SDK（独立模块）
│   ├── go.mod / go.sum              # SDK 自身依赖
│   ├── app/
│   │   ├── app.go                   # 初始化流程编排
│   │   └── health.go                # 健康检查 "/" 处理器
│   ├── auth/
│   │   ├── base_auth.go             # 认证基础逻辑（Business + Consumer）
│   │   ├── auth_tool.go             # B 端认证包级便捷函数
│   │   ├── permission_interface.go  # 权限查询接口 + 默认实现
│   │   ├── permission_matcher.go    # 权限通配符匹配器
│   │   ├── permission_scan.go       # 权限自动扫描与 Redis 缓存
│   │   ├── permission_tool.go       # 权限/角色查询门面
│   │   ├── module.go                # 认证模块初始化
│   │   ├── middleware/
│   │   │   ├── check_login.go       # B 端登录验证中间件
│   │   │   ├── check_permission.go  # B 端 + C 端权限检查中间件
│   │   │   ├── check_role.go        # B 端角色检查中间件
│   │   │   └── norepeat.go          # 防重复提交中间件
│   │   └── pojo/
│   ├── captcha/                     # 图形验证码
│   ├── config/                      # 配置加载（YAML）
│   ├── constants/                   # 系统常量和 Redis 键前缀
│   ├── crud/                        # 通用 CRUD 帮助函数
│   ├── db/                          # 数据库初始化（GORM + go-redis）
│   ├── enums/                       # 枚举定义
│   ├── eventbus/                    # 内存事件总线
│   ├── exception/                   # BusinessError 业务异常
│   ├── log/                         # 操作日志（SysLog 中间件 + 工具函数）
│   ├── middleware/                  # 全局中间件
│   │   ├── trace.go                # 链路追踪
│   │   ├── auth_check.go           # 认证路由分流
│   │   ├── recovery.go             # 全局异常恢复
│   │   ├── cors.go                 # CORS 跨域
│   │   └── ratelimit.go            # API 限流
│   ├── module/                      # 模块生命周期管理
│   ├── pojo/                        # 通用 POJO
│   ├── registry/                    # 注册中心（路由/中间件/权限）
│   │   ├── route.go        # 路由注册
│   │   ├── middleware.go    # 中间件注册
│   │   └── perm.go          # Perm() / ClientPerm() 权限快捷注册
│   ├── result/                      # 统一响应格式
│   ├── scheduler/                   # 定时调度（cron）
│   ├── storage/                     # 文件存储抽象（Engine 接口）
│   └── utils/                       # 工具函数（雪花 ID、加密、IP 等）
│
│   ├── contracts/                   # 公共契约（权限/事件/日志接口）
│
├── plugins/                         # 业务插件目录
│   ├── plugin-sys/                  # 系统管理插件（B 端）
│   │   ├── go.mod                   # 依赖 sdk
│   │   ├── plugin.go               # 插件入口
│   │   ├── imports.go              # 导入所有子模块
│   │   ├── persistence.go          # 日志持久化
│   │   ├── provider/               # 权限/用户 Provider
│   │   ├── auth/                   # 认证（captcha/sm2/username）
│   │   ├── banner/                 # Banner 管理
│   │   ├── config/                 # 系统配置
│   │   ├── dict/                   # 数据字典
│   │   ├── file/                   # 文件管理
│   │   ├── group/                  # 用户组
│   │   ├── home/                   # 首页仪表盘
│   │   ├── log/                    # 操作日志查询
│   │   ├── notice/                 # 通知公告
│   │   ├── org/                    # 组织架构
│   │   ├── permission/             # 权限查询
│   │   ├── position/               # 职位管理
│   │   ├── resource/               # 资源管理（菜单/按钮）
│   │   ├── role/                   # 角色管理
│   │   ├── session/                # 会话管理
│   │   ├── user/                   # 用户管理
│   │   ├── analyze/                # 系统分析
│   │   └── ...                     # 子模块数据模型 + Model 定义
│   │
│   ├── plugin-client/              # 客户端插件（C 端）
│   │   ├── go.mod
│   │   ├── plugin.go
│   │   ├── auth/                   # C 端认证
│   │   ├── session/                # C 端会话管理
│   │   └── user/                   # C 端用户管理
│   │
│   └── plugin-im/                  # IM 插件（WebSocket）
│       ├── go.mod
│       ├── plugin.go
│       ├── ws/                     # WebSocket 核心
│       │   ├── hub.go             # 本地连接管理
│       │   ├── client.go          # WS 客户端管理
│       │   ├── cross_hub.go       # Redis 跨实例消息桥
│       │   ├── config.go          # WS 配置
│       │   └── message.go         # 消息协议
│       ├── friend/                # 好友管理
│       ├── group/                 # 群组管理
│       ├── message/               # 消息管理
│       ├── broadcast/             # 广播消息
│       └── model/                 # IM 数据模型
│
├── app/                            # 应用组装层（独立模块）
│   ├── main.go                     # 应用启动入口
│   └── go.mod                      # 依赖所有插件 + sdk
│
├── cmd/                            # 命令行工具
│   ├── migrate/main.go             # DB 迁移 + 种子数据
│   └── codegen/main.go             # 代码生成器
│
├── scripts/                        # 辅助脚本
├── docs/                           # 设计规划文档
├── uploads/                        # 本地上传文件存储目录
└── vitepress/                      # VitePress 在线文档
```

## 目录设计说明

### sdk/ 框架核心

`sdk/` 是框架的 SDK 包（独立 Go 模块 `hei-gin/sdk`），包含与具体业务无关的通用框架能力。每个子包职责单一、高内聚低耦合：

- **app/**：应用生命周期管理，统一的初始化编排
- **auth/**：完整的认证授权体系，包含 B/C 双端认证、权限管理、角色检查
- **middleware/**：HTTP 全局中间件（Trace / AuthCheck / Recovery / CORS）
- **storage/**：文件存储抽象层（Engine 接口），支持 Local / MinIO / S3 三种实现
- **config/**：配置文件加载与结构体定义
- **db/**：GORM + Redis 数据库初始化
- **module/**：模块生命周期管理（Register/Init/Start/Stop）
- **registry/**：注册中心模式，支持路由、中间件自注册 + Perm() 权限快捷注册
- **result/**：统一 JSON 响应格式
- **exception/**：BusinessError 业务异常（panic-based）
- **crud/**：通用 CRUD 帮助函数，减少重复代码
- **eventbus/**：内存事件总线，模块间解耦通信
- **scheduler/**：基于 cron 的定时任务调度

### plugins/ 业务插件

`plugins/` 目录包含多个独立 Go 模块，每个插件实现一组业务功能。每个插件遵循 `Plugin` 接口约定，通过 `module.Register()` 自注册。

每个插件内部的子模块遵循统一的内部约定：

```
plugins/<plugin>/<module>/
├── api/v1/
│   └── api.go             # 路由注册 + HTTP Handler
├── service.go             # 业务逻辑
├── params.go              # 请求/响应参数
└── model.go               # GORM 数据模型
```

这种结构的好处：

- **高内聚**：一个功能的所有代码集中管理
- **低耦合**：插件之间通过 api 包定义的接口交互
- **易于裁剪**：不需要的插件可以直接从 go.work 中移除
- **便于协作**：不同团队维护不同插件

### api/ 接口定义层

`api/` 是独立的 Go 模块，仅依赖 Go 标准库。它定义插件和 SDK 之间的接口契约：

- **Plugin 接口**：插件必须实现 `Info()` 和 `Name()` 方法
- **EventBus 接口**：事件总线的发布/订阅接口
- **LogPersistenceAPI**：日志持久化接口
- **AuthAPI**：认证相关接口

### app/ 应用组装层

`app/` 是独立的 Go 模块，负责导入所有插件并调用 `sdk/app.Run()` 启动应用。它是整个应用的最上层组装点，不包含任何业务逻辑。
