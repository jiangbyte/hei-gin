# SDK 分层依赖规则

本文档定义 `sdk` 当前四层结构下的依赖方向，目的是避免后续继续出现跨层污染。

---

## 当前分层

```text
sdk/kernel
sdk/infra
sdk/web
sdk/shared
```

含义：

- `kernel`
  应用运行时与装配内核
- `infra`
  基础设施能力
- `web`
  HTTP / Gin 协议层能力
- `shared`
  跨层共享的稳定能力与契约

---

## 允许的依赖方向

允许：

- `kernel -> infra`
- `kernel -> web`
- `kernel -> shared`
- `web -> infra`
- `web -> shared`
- `infra -> shared`

禁止：

- `infra -> web`
- `shared -> infra`
- `shared -> web`
- `shared -> kernel`

原因：

- `infra` 应保持协议无关，不能反向依赖 HTTP 层
- `shared` 应保持最稳定、最轻，不承载运行时和协议细节

---

## 各层职责约束

`sdk/kernel`

- 负责应用启动、插件生命周期、注册中心
- 可以组合 `infra` 和 `web`
- 不承载业务模块逻辑

`sdk/infra`

- 负责 DB、Redis、Storage、Scheduler、EventBus
- 不直接处理 Gin `Context`
- 不直接输出 HTTP 响应
- 不依赖 `sdk/web`

`sdk/web`

- 负责 `result`、`middleware`、`exception`
- 可以依赖 `infra`，例如限流中间件访问 Redis
- 不承载具体业务 API

`sdk/shared`

- 只放稳定契约和跨层通用能力
- 推荐内容：
  - `contracts`
  - 少量 panic-safe / helper 能力
- 不放 HTTP、DB、插件注册逻辑

---

## 业务模块约束

`plugins/...` 可以依赖：

- `sdk/kernel/registry`
- `sdk/infra/*`
- `sdk/web/*`
- `sdk/shared/*`

但不要：

- 把业务 API 放回 `sdk`
- 把业务 service/repository 放回 `sdk`
- 在 `sdk` 中引入具体插件业务概念

---

## 当前已收敛的一处典型案例

`eventbus` 原先通过 `sdk/web/middleware.GoSafe` 做 goroutine panic 保护，这会形成 `infra -> web` 依赖。

现在已调整为：

- 通用安全执行能力放到 `sdk/shared/safe`
- `sdk/infra/eventbus` 依赖 `sdk/shared/safe`

这类问题后续都应按同样思路处理：把真正跨层通用的能力抽到 `shared`，而不是让下层反向依赖上层。
