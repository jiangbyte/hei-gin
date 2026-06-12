# SDK 结构调整说明

本文档说明 `sdk` 在中大型项目下的推荐定位，以及本仓库当前采用的收敛方向。

---

## 定位

`sdk` 应作为基础设施与运行时支撑层，不承载具体业务模块。

适合放在 `sdk` 的内容：

- 应用启动与生命周期
- 插件装配与运行时注册
- 路由注册机制
- DB、Redis、Storage、Scheduler、EventBus
- 通用 Middleware、Result、Exception
- 稳定的 Contracts、Constants、Enums、Utils

不适合继续放入 `sdk` 的内容：

- 具体业务模块的 `api`
- 具体业务模块的 `service/repository/model`
- 强依赖某个插件上下文的业务逻辑

---

## 当前收敛方向

为了避免 `sdk` 继续变成“大杂烩公共目录”，当前已收敛为四块：

```text
sdk/kernel/app
sdk/kernel/plugin
sdk/kernel/registry
sdk/infra/db
sdk/infra/storage
sdk/infra/scheduler
sdk/infra/eventbus
sdk/web/result
sdk/web/middleware
sdk/web/exception
sdk/shared/contracts
```

职责划分：

- `sdk/kernel/app`
  应用启动、HTTP Server、Swagger、路由总装配
- `sdk/kernel/plugin`
  插件生命周期注册与调度
- `sdk/kernel/registry`
  路由注册、中间件注册、权限注册
- `sdk/infra`
  基础设施能力，包括 DB、Redis、Storage、Scheduler、EventBus
- `sdk/web`
  面向 HTTP / Gin 的协议层支撑，包括响应、异常、中间件
- `sdk/shared`
  跨层共享但不绑定 HTTP 的稳定契约类型

---

## 当前结果

旧路径兼容层已经移除：

- `sdk/app`
- `sdk/plugin`
- `sdk/registry`

当前代码统一依赖：

- `hei-gin/sdk/kernel/app`
- `hei-gin/sdk/kernel/plugin`
- `hei-gin/sdk/kernel/registry`
- `hei-gin/sdk/infra/db`
- `hei-gin/sdk/infra/storage`
- `hei-gin/sdk/infra/scheduler`
- `hei-gin/sdk/infra/eventbus`
- `hei-gin/sdk/web/result`
- `hei-gin/sdk/web/middleware`
- `hei-gin/sdk/web/exception`
- `hei-gin/sdk/shared/contracts`

---

## 业务模块约束

业务模块仍应留在 `plugins/...`，例如：

```text
plugins/plugin-sys/banner
plugins/plugin-sys/config
plugins/plugin-im/message
```

模块内部维持：

- `module.go`
- `repository.go`
- `service.go`
- `api/v1/api.go`

不要把业务模块 API 搬进 `sdk`。

---

## 后续建议

后续可以继续渐进收敛：

1. 新代码继续统一使用 `sdk/kernel/*`
2. 继续把 `sdk` 中其他职责按 `kernel / infra / web / shared` 收拢
3. 避免把业务模块 API、service、repository 再放回 `sdk`
