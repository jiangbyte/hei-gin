# 生产化整改计划

本文档用于记录 `hei-gin` 当前距离“生产级”还有哪些差距，并给出一份可执行的整改顺序。

目标不是一次性“重写成大厂架构”，而是在当前代码基础上，按优先级分批收敛，先解决上线阻断项，再补齐稳定性、可观测性和可维护性。

---

## 当前结论

当前项目：

- 已具备中小型项目可运行能力
- 已有较清晰的插件化和模块分层基础
- 已有 MySQL / Redis / WebSocket / 基础中间件 / Swagger / 迁移入口

但仍不能直接认定为“生产级”。

主要原因：

- 存在若干上线阻断风险
- 启动与配置策略不够严格
- 全局状态与 `init()` 副作用偏重
- 健康检查、观测、验证体系不完整
- 测试与 CI 对多模块 workspace 的覆盖还不够可靠

---

## 生产级判定标准

后续整改以以下标准为目标：

- 配置缺失或关键依赖异常时，服务必须明确启动失败
- 核心链路不能依赖脆弱的全局对象隐式初始化
- HTTP / WS / DB / Redis 具备基本安全边界
- 有可落地的健康检查、日志、指标、错误定位手段
- 有稳定的迁移策略和回归验证策略
- 有明确的模块依赖边界，避免继续回退到“大 service / 大全局状态”

---

## P0 上线阻断项

P0 表示不建议带着这些问题直接上正式生产。

### 1. 消息推送链路对 `GlobalCrossHub` 依赖不一致

问题：

- `plugin-im` 中部分逻辑直接调用 `ws.GlobalCrossHub.*`
- 只有少数路径做了判空保护
- 一旦初始化顺序、工具命令、测试路径、局部运行方式发生变化，存在 panic 风险

涉及位置：

- `plugins/plugin-im/friend/service.go`
- `plugins/plugin-im/group/message.go`
- `plugins/plugin-im/group/member.go`
- `plugins/plugin-im/group/query.go`
- `plugins/plugin-im/broadcast/service.go`
- `plugins/plugin-im/message/service.go`

整改要求：

- 禁止业务代码直接依赖裸全局 `GlobalCrossHub`
- 收敛出统一的推送访问面，例如：
  - `type PushGateway interface`
  - 或 `ws.Service / ws.Runtime`
- 至少做到：
  - 初始化失败时行为明确
  - 单机场景、Redis 缺失场景、测试场景都不 panic
  - 所有消息发送路径使用一致的空值与降级策略

验收标准：

- 全仓搜索不再出现业务层直接访问 `GlobalCrossHub.*`
- `plugin-im` 在无 Redis / 测试环境下不因推送链路崩溃

### 2. 配置缺失时静默使用空配置继续启动

问题：

- `sdk/config/config.go` 中 `FindAndLoad()` 在找不到配置文件时会创建空配置并继续运行
- 这会把启动错误转成运行时隐患

涉及位置：

- `sdk/config/config.go`
- `sdk/kernel/app/app.go`

整改要求：

- 生产模式下，关键配置缺失必须启动失败
- 至少校验：
  - `app.host`
  - `app.port`
  - `db.host / db.port / db.user / db.database`
  - Redis 若为必需，也应强校验
- 可以保留开发模式默认值，但必须显式区分

建议做法：

- 引入 `Validate()` 方法
- `FindAndLoad()` 只负责加载
- `app.Run()` 在启动前执行严格校验

验收标准：

- 配置文件缺失或关键字段缺失时，服务明确退出
- 不再存在“空配置也能把服务拉起来”的情况

### 3. WebSocket 安全边界过松

问题：

- `CheckOrigin` 当前直接返回 `true`
- 连接数限流依赖 `X-Forwarded-For` / `X-Real-IP`，但没有可信代理边界

涉及位置：

- `plugins/plugin-im/ws/hub.go`

整改要求：

- `CheckOrigin` 必须改成基于配置的白名单校验
- 代理头只在可信反向代理场景下使用
- 未配置可信代理时，应回退到 `RemoteAddr`

建议做法：

- 新增 WS 配置项：
  - `allowed_origins`
  - `trusted_proxies`
- 对 IP 获取逻辑做统一封装

验收标准：

- 未授权来源不能建立 WS 连接
- 伪造 `X-Forwarded-For` 不能绕过连接数限制

### 4. 多模块测试验证策略不完整

问题：

- 仓库使用 `go.work`
- 根模块执行 `go test ./...` 不能等价覆盖所有子模块
- 当前验证容易给出“看起来都过了”的错觉

涉及位置：

- `go.work`
- 根 `go.mod`
- 各插件子模块 `go.mod`

整改要求：

- 明确 workspace 级测试入口
- CI 不能只跑根模块测试

建议做法：

- 在 `Makefile` 中新增统一命令，例如：
  - `test-root`
  - `test-sdk`
  - `test-plugin-sys`
  - `test-plugin-client`
  - `test-plugin-im`
  - `test-all`
- 在 CI 中逐模块执行

验收标准：

- 新提交必须经过 workspace 全模块测试
- 任何模块回归都能在 CI 中暴露

---

## P1 生产稳定性项

P1 表示不一定阻断首发，但应在正式放量前尽快补齐。

### 1. 健康检查过浅

当前 `/` 仅返回字符串与版本，不能反映真实可用性。

整改要求：

- 增加：
  - `liveness`
  - `readiness`
  - 可选 `startup`
- `readiness` 至少检查：
  - MySQL ping
  - Redis ping
  - 关键插件初始化状态

建议路径：

- `/health/live`
- `/health/ready`

验收标准：

- K8s / 网关 / 运维系统可以基于 `readiness` 做流量摘除

### 2. 启动与关闭流程缺少更细粒度降级与状态控制

问题：

- 当前大量 `log.Fatalf`
- 适合简单启动失败，但不利于未来做嵌入式运行、自动化控制、局部降级

整改方向：

- `Run()` 内部改为显式返回 error 或分阶段启动状态
- `plugin.StartAll()` / `plugin.StopAll()` 的失败策略明确化
- 对调度器、WS、Redis 订阅等长生命周期组件做状态跟踪

### 3. 全局注册体系依赖 `init()` 副作用

问题：

- 路由、插件、模型、中间件注册大量依赖 blank import 和 `init()`
- 对大型项目和多二进制演进不友好

整改方向：

- 短期保留 `init()` 模式
- 中期逐步收敛为显式装配入口
- 至少增加：
  - 注册去重
  - 启动后冻结
  - 注册结果可观测

建议：

- 为 plugin / route / middleware / model 注册中心增加只读快照能力

### 4. 业务层错误处理风格不统一

当前存在三种风格混用：

- `result.WriteError(...)`
- 返回 `error`
- `panic(BusinessError)` + Recovery

整改方向：

- 模块内统一一套主风格
- 当前项目建议继续保留“service 可直接写响应”，但要统一约束

建议规则：

- handler 只做参数绑定和响应收口
- service 内统一使用 `result.WriteError`
- 尽量减少业务 panic

### 5. 可观测性不足

当前虽然有日志与 trace 中间件，但还不足以支撑生产排障。

建议补齐：

- Prometheus metrics
- 请求耗时统计
- DB 慢查询日志
- Redis 耗时统计
- WS 在线人数、连接失败数、消息投递数
- 插件启动耗时

建议最先补的指标：

- HTTP QPS / latency / 5xx
- DB ping / pool usage
- Redis ping / errors
- WS current connections / rejected connections

### 6. 迁移体系偏轻量

当前迁移更偏向模型注册与统一执行，适合中小项目早期，但生产要求更高。

整改方向：

- 明确迁移版本号
- 记录已执行历史
- 支持灰度环境重复执行
- 避免不可控的结构漂移

建议：

- 中短期先保留现有机制
- 但增加迁移执行记录表和版本说明

---

## P2 中长期治理项

P2 不是当前阻断项，但如果目标是中大型项目长期演进，应尽早规划。

### 1. 收紧全局默认模块实例的使用范围

当前大量使用：

- `DefaultModule`
- `defaultHandler`

这在模块内装配没问题，但不应继续扩散到复杂业务子流程中。

方向：

- 模块外跨模块调用，优先通过明确 Service 暴露面
- 避免越来越多包内 helper 通过全局默认实例取 repo

### 2. 降低业务层对 Gin Context 的耦合

当前项目采用“service 可接收 `*gin.Context`”是务实方案，但长期看会限制：

- 纯业务测试
- 异步任务复用
- 非 HTTP 场景调用

方向：

- 新增重要核心流程的 context-only 版本
- 把纯业务计算与 HTTP 响应逐步拆开

### 3. 对 `sdk` 继续分层收紧

虽然 `sdk/kernel / infra / web / shared` 已经比之前清晰，但仍需要长期约束：

- 不把业务 API 再放进 `sdk`
- 不让 `infra -> web`
- 不让 `shared` 承担运行时职责

### 4. CI / 静态检查 / 质量门禁补齐

建议逐步补齐：

- `go test` 全模块
- `go vet`
- `staticcheck`
- swagger 生成校验
- codegen 产物一致性校验
- 文档样例编译检查

---

## 建议整改顺序

建议按下面顺序落地：

1. 先做 P0
2. 再补健康检查和 CI
3. 再补观测性
4. 最后治理 `init()`、全局状态、Gin 耦合

更具体的执行顺序：

1. 收敛 `GlobalCrossHub` 访问面
2. 配置严格校验
3. WS origin / proxy / IP 信任边界
4. `test-all` 与 CI
5. `/health/live`、`/health/ready`
6. metrics 与关键运行指标
7. 迁移版本化
8. 注册中心显式化 / 冻结化

---

## 建议验收方式

每完成一项整改，至少做下面几类验证：

- 单元测试
- 模块级集成测试
- 启动失败场景测试
- 配置缺失场景测试
- Redis / DB 不可用场景测试
- WS 连接与限流测试

建议为生产化整改单独建立验收清单，例如：

- `docs/production-readiness-checklist.md`

后续每次上线前逐项勾选：

- 配置校验通过
- 迁移通过
- workspace 全量测试通过
- readiness 检查通过
- 关键指标可观测
