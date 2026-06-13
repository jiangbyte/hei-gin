# 注册与装配现状

当前项目已经完成以下阶段：

- 第一阶段：注册中心增强
  - plugin / route / middleware / model / seed 支持去重、冻结、快照、测试重置
- 第二阶段：显式装配入口
  - `sdk` builtin plugin 改为显式 `RegisterPlugin()`
  - `plugin-sys` / `plugin-client` / `plugin-im` 改为显式：
    - `RegisterPlugin()`
    - `RegisterRoutes()`
    - `RegisterMigrations()`
- 第三阶段：顶层 bootstrap 显式调用
  - `main.go` 不再依赖顶层 plugin blank import 触发注册
  - `cmd/migrate/main.go` 不再依赖业务包 blank import 触发模型/seed 注册

---

## 当前原则

1. 注册调用可以在显式装配函数内发生
2. 不再依赖“仅导入包就自动完成业务装配”
3. `init()` 可以保留给非装配用途，但不应继续承担 route / model / seed / plugin 注册职责
4. 顶层启动链应明确声明装配顺序

---

## 当前装配入口

### `sdk` builtin

- `auth.RegisterPlugin()`
- `captcha.RegisterPlugin()`
- `utils.RegisterPlugin()`
- `scheduler.RegisterPlugin()`

由 [sdk/kernel/app/app.go](/mnt/e/projects/mine/hei/hei-gin/sdk/kernel/app/app.go:1) 显式调用。

### 业务 plugin

#### `plugin-sys`

- `plugin_sys.RegisterPlugin()`
- `plugin_sys.RegisterRoutes()`
- `plugin_sys.RegisterMigrations()`

#### `plugin-client`

- `plugin_client.RegisterPlugin()`
- `plugin_client.RegisterRoutes()`
- `plugin_client.RegisterMigrations()`

#### `plugin-im`

- `plugin_im.RegisterPlugin()`
- `plugin_im.RegisterRoutes()`
- `plugin_im.RegisterMigrations()`

由 [main.go](/mnt/e/projects/mine/hei/hei-gin/main.go:1) 和 [cmd/migrate/main.go](/mnt/e/projects/mine/hei/hei-gin/cmd/migrate/main.go:1) 显式调用。

---

## 后续约束

新增模块时，优先按以下模式：

1. 在子模块 `api/v1` 中暴露 `Register()`
2. 在 plugin 根包中集中调用这些 `Register()`
3. 在 plugin 根包中集中声明 `RegisterMigrations()`
4. 由顶层 bootstrap 显式调用 plugin 的装配入口

不再推荐：

- 新增 `imports.go` 做 blank import 聚合
- 在 `init()` 中直接 `plugin.Register(...)`
- 在 `init()` 中直接 `registry.RegisterRoute(...)`
- 在 `init()` 中直接 `db.RegisterModel(...)`
- 在 `init()` 中直接 `db.RegisterSeed(...)`
