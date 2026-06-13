# Auth 重构计划

本文档用于定义 `sdk/auth` 的重构目标、边界、实施顺序和验收标准。

这次重构的目标不是局部修补，而是把当前分散在 `sdk/auth`、`sdk/auth/middleware`、`plugins/plugin-sys/session`、`plugins/plugin-sys/provider` 中的认证、鉴权、权限元数据、会话治理能力重新收拢，形成一套职责明确、性能稳定、可扩展的统一内核。

---

## 一、重构目标

### 1. 统一调用模型

对外调用统一为：

```go
auth.Business.Login(c, userID, extra)
auth.Business.IsLogin(c)
auth.Business.HasPermission(c, "sys:user:page")
auth.Business.Sessions().Page(ctx, query)
```

不再允许：

- 传裸字符串 `loginType`
- 使用默认 BUSINESS 包级函数
- 在中间件或业务中手写 `if consumer else business`

### 2. 职责明确

`sdk/auth` 统一负责：

- 认证
- 鉴权
- 权限声明与扫描
- 权限缓存
- session 管理与统计

插件层只负责：

- 提供权限/角色数据
- 提供业务侧用户资料补充
- 提供业务侧数据权限应用策略

### 3. 高性能 session 管理

把后台会话分页、token 明细、强退、趋势分析等能力并回 `sdk/auth`，避免插件层重复编排底层 Redis 访问。

### 4. 未来可扩展

底层应支持未来新增：

- `auth.Developer`
- 其他端态

新增端态应是“注册一个 realm”，而不是到处扩散分支。

---

## 二、当前问题

### 1. 多端支持是半抽象、半硬编码

当前虽然有 `baseAuthTool`，但实际调用路径仍大量依赖：

- `businessAuth`
- `Consumer`
- `ToolForLoginType()`
- `if loginType == CONSUMER`

这会导致：

- 新增端态成本高
- 容易静默回退到错误端态
- 代码横向重复

### 2. auth 职责过杂但边界又不清

当前 `sdk/auth` 同时承担：

- token/session 存储
- 权限查询
- 中间件
- 权限扫描
- 会话统计

但这些能力之间没有形成清晰分层，导致：

- 可读性差
- 横向复用困难
- 插件层又补了一层 session 编排逻辑

### 3. 运行时权限查询依赖 provider

当前权限和角色查询主要是运行时按请求实时走 provider，这会带来：

- 请求路径性能开销
- provider 职责偏重
- 中间件链路复杂

### 4. session 管理分散

当前会话相关能力分散在：

- `sdk/auth/session_index.go`
- `plugins/plugin-sys/session/service.go`
- `plugins/plugin-client/session/service.go`

导致：

- 核心逻辑和展示逻辑混在一起
- 后台会话页难以统一优化
- 边界不清

### 5. 权限缓存结构不统一

权限扫描和权限展示读取的目标结构不同，容易出现缓存结构不一致问题。

---

## 三、最终设计原则

### 1. 对外 API 使用 realm 对象

统一采用：

```go
auth.Business.xxx
auth.Consumer.xxx
```

而不是：

- `auth.Realms.Business.xxx`
- `auth.For("BUSINESS").xxx`
- `auth.Login(..., "BUSINESS")`

### 2. provider 只负责权限和角色数据

provider 对外只暴露：

```go
type PermissionProvider interface {
    GetPermissionList(ctx context.Context, realmID RealmID, userID string) ([]string, error)
    GetRoleList(ctx context.Context, realmID RealmID, userID string) ([]string, error)
}
```

不承担：

- Gin 相关逻辑
- session 构建
- middleware 判断
- 请求时 scope 拼装

### 3. 权限快照写入 session claims

用户登录成功后，auth 内核构建权限快照并写入 session：

- `Permissions`
- `Roles`
- `ScopeMap`

后续请求直接从 session claims 读取，不再每次请求查 provider。

### 4. ScopeMap 是 auth 输出的数据权限描述，不是业务 SQL 拼装器

`ScopeMap` 只回答：

- 用户对某个权限码能看到什么范围

不负责：

- 某张表怎么拼 where 条件
- 某个模块的归属字段是什么

业务模块通过自己的 `ScopePolicy` 把 `ScopeInfo` 应用到查询。

### 5. middleware 只做 HTTP 接入与鉴权

middleware 负责：

- token 解析
- claims 注入
- request cache
- 登录/权限/角色校验

不负责：

- 查询 provider
- realm 分支判断
- 业务 SQL 判断

### 6. session 管理统一内聚到 auth

单 realm 会话能力挂到 realm 上：

```go
auth.Business.Sessions().Page(...)
auth.Business.Sessions().Tokens(...)
```

跨 realm 聚合能力单独提供：

```go
auth.Sessions(auth.Business, auth.Consumer).Stats(ctx)
```

---

## 四、核心类型设计

### 1. RealmID

```go
type RealmID string

const (
    BusinessID  RealmID = "BUSINESS"
    ConsumerID  RealmID = "CONSUMER"
    DeveloperID RealmID = "DEVELOPER"
)
```

### 2. Realm

```go
type Realm struct {
    ID        RealmID
    TokenName string
    Expire    int
}
```

对外暴露：

```go
var Business *Realm
var Consumer *Realm
var Developer *Realm
```

`Realm` 承担：

- 登录
- 登出
- 登录态检查
- token 续期
- 禁用
- 强退
- 权限读取
- session 服务入口

### 3. SessionClaims

```go
type SessionClaims struct {
    UserID    string         `json:"user_id"`
    RealmID   RealmID        `json:"realm_id"`
    CreatedAt string         `json:"created_at"`
    ACL       ACLSnapshot    `json:"acl"`
    Extra     map[string]any `json:"extra"`
}
```

### 4. ACLSnapshot

```go
type ACLSnapshot struct {
    Permissions []string             `json:"permissions"`
    Roles       []string             `json:"roles"`
    ScopeMap    map[string]ScopeInfo `json:"scope_map"`
    Version     string               `json:"version,omitempty"`
}
```

### 5. ScopeInfo

```go
type ScopeInfo struct {
    GroupScope     string   `json:"group_scope"`
    OrgScope       string   `json:"org_scope"`
    CustomGroupIDs []string `json:"custom_group_ids"`
    CustomOrgIDs   []string `json:"custom_org_ids"`
}
```

---

## 五、目录与职责划分

建议将 `sdk/auth` 重构为以下结构：

```text
sdk/auth/
├── realm.go
├── provider.go
├── claims.go
├── token.go
├── permission.go
├── scope.go
├── permission_scan.go
├── session.go
├── session_aggregate.go
└── middleware/
    ├── login.go
    ├── permission.go
    ├── role.go
    └── context.go
```

### 1. `realm.go`

负责：

- `RealmID`
- `Realm`
- `auth.Business`
- `auth.Consumer`
- 内部 registry 初始化

### 2. `provider.go`

负责：

- `PermissionProvider`
- provider 注册
- provider 获取

### 3. `claims.go`

负责：

- `SessionClaims`
- `ACLSnapshot`
- claims 编解码
- request cache

### 4. `token.go`

负责：

- login/logout
- renew
- disable
- kickout
- Redis token/session 读写

### 5. `permission.go`

负责：

- `GetPermissionList`
- `GetRoleList`
- `HasPermission`
- `HasRole`
- `ScopeFor`
- 从 claims 读取 ACL

### 6. `scope.go`

负责：

- scope 常量
- `MergeScope`
- scope helper
- scope 构建入口

### 7. `permission_scan.go`

负责：

- 权限声明注册
- 权限扫描
- 权限缓存写入
- 权限树与 code 列表读取

### 8. `session.go`

负责：

- 单 realm 会话分页
- token 列表
- 强退用户
- 强退 token
- stats
- trend
- session summary

### 9. `session_aggregate.go`

负责：

- 多 realm 聚合统计
- 多 realm 趋势

### 10. `middleware/*`

负责：

- 登录校验
- 权限校验
- 角色校验
- 上下文注入

---

## 六、认证核心设计

### 1. 登录流程

登录成功后执行：

1. 调 provider 获取权限列表
2. 调 provider 获取角色列表
3. 构建 `ScopeMap`
4. 组装 `SessionClaims`
5. 写入 Redis token key
6. 更新 session set
7. 更新 session 索引
8. 更新 session summary

### 2. 登出与强退

支持：

- 当前 token 登出
- 强退用户所有 token
- 强退指定 token

这些能力都挂在 realm 上：

```go
auth.Business.Logout(c)
auth.Business.KickoutUser(ctx, userID)
auth.Business.KickoutToken(ctx, userID, token)
```

### 3. 续期与禁用

realm 自身负责：

- token/session 续期
- 账号禁用
- 禁用时间查询

---

## 七、权限系统设计

### 1. provider 职责

provider 只提供：

- 权限列表
- 角色列表

```go
type PermissionProvider interface {
    GetPermissionList(ctx context.Context, realmID RealmID, userID string) ([]string, error)
    GetRoleList(ctx context.Context, realmID RealmID, userID string) ([]string, error)
}
```

### 2. 运行时权限判定

运行时不再查询 provider，而是从 claims 读取：

- `claims.ACL.Permissions`
- `claims.ACL.Roles`
- `claims.ACL.ScopeMap`

因此：

- `HasPermission` 只读 claims
- `HasRole` 只读 claims
- `ScopeFor` 只读 claims

### 3. 通用 ACL 原则

auth 内核不内建任何“超级管理员”或其他特权角色的专用分支。

在 auth 视角中：

- 角色只是 ACL 数据的一部分
- 权限判断只依据 session claims 中的 `Permissions`、`Roles`、`ScopeMap`
- 是否存在“拥有所有权限的角色”属于业务层规则，不属于 auth 内核职责

如果某业务需要实现特权角色逻辑，应在：

- provider 侧组装完整权限集合
- 或业务层自行放宽策略

而不是在 auth 内核保留角色名称硬编码。

### 4. ScopeMap 的数据来源

`ScopeMap` 不通过 provider 对外接口返回，而是在登录时由 auth 内核自行构建。

当前阶段明确的数据来源：

- `RelUserPermission`
- `RelRolePermission`

构建原则：

1. 登录时先获取用户角色列表
2. 基于角色查询 `RelRolePermission`
3. 基于用户查询 `RelUserPermission`
4. 对同一 `permissionCode` 的多条 scope 记录做合并
5. 角色授权与用户直授权限共存时，按既定优先级合并

这意味着：

- provider 对外仍然只提供权限列表和角色列表
- scope 构建属于 auth 内核登录态初始化的一部分
- 运行时不再查询关系表，而是直接读取 claims 中的 `ScopeMap`

### 5. ScopeMap 构建规则

scope 的合并应基于当前已有语义继续收敛：

- 同一权限码可能同时来自角色授权和用户直授
- 用户直授权限优先级高于角色授权
- group/org scope 取更严格者
- 自定义 group/org ID 集合做去重合并

换句话说：

- `RelRolePermission` 提供基础数据范围
- `RelUserPermission` 可进一步收紧或覆盖

### 6. ScopeMap 更新策略

由于 `ScopeMap` 存储在 session claims 中，它本质上是登录态快照。

当以下数据发生变更时，需要刷新已有 session：

- 用户直授权限变更
- 角色权限变更
- 用户角色关联变更
- scope 范围字段变更

第一阶段采用：

- 变更后强制刷新对应用户 session

刷新策略可分两类：

1. 简单方案
- 踢下线该用户所有相关 realm session
- 用户重新登录后获取新 claims

2. 主动刷新方案
- 后台管理操作后，重建该用户新 claims 并回写现存 token/session

第一阶段建议先采用“踢下线 + 重登”作为默认机制，保证实现简单且边界清晰。  
如果后续希望减少重登影响，再增量补“在线 claims 刷新”能力。

---

## 八、ScopeMap 设计与业务应用

### 1. ScopeMap 的定位

`ScopeMap` 是：

- `permissionCode -> ScopeInfo`

例如：

```go
map[string]ScopeInfo{
    "sys:user:page":   {...},
    "sys:notice:page": {...},
}
```

它表达的是：

- 对某个权限码，当前用户允许访问的数据范围

### 2. 建议的 scope 语义

建议使用统一常量：

- `ALL`
- `SELF`
- `GROUP`
- `GROUP_AND_CHILDREN`
- `ORG`
- `ORG_AND_CHILDREN`
- `CUSTOM`
- `NONE`

### 3. auth 层提供的能力

auth 层提供：

```go
scope, ok := auth.Business.ScopeFor(c, "sys:user:page")
```

以及基础 helper，例如：

- `IsAll`
- `IsSelfOnly`
- `AllowedGroupIDs`
- `AllowedOrgIDs`

### 4. 业务层如何应用

业务层不直接手写散乱的 if 判断。  
建议每个模块提供自己的 `ScopePolicy`：

```go
type ScopePolicy interface {
    Apply(db *gorm.DB, scope auth.ScopeInfo, actor ActorContext) *gorm.DB
    CanAccess(scope auth.ScopeInfo, actor ActorContext, resource ResourceOwnerInfo) bool
}
```

其中：

- `Apply` 用于列表、分页、导出
- `CanAccess` 用于详情、更新、删除

### 5. 为什么不让 auth 拼 SQL

因为 auth 不应该知道：

- 表字段名
- 归属字段名
- 是否需要 join
- 某业务的 owner 规则

所以最合理的边界是：

- auth 给 `ScopeInfo`
- 业务用 policy 落地到自己的查询

### 6. 权限变更一致性

由于 `ScopeMap` 存在 session claims 中，属于登录态快照，因此权限或范围变更后需要刷新 session。

第一阶段策略：

- 权限/角色/scope 变更后，踢下线用户已有 session
- 用户重新登录后获得新 claims

这条策略同样适用于由 `RelUserPermission`、`RelRolePermission` 引起的 scope 变更。

### 7. 默认失败语义

为了避免各模块对 scope 缺失时做出不一致处理，统一约定：

- 未配置数据权限 policy：拒绝访问
- 当前权限码在 `ScopeMap` 中不存在：拒绝访问
- `ScopeInfo` 为 `NONE`：拒绝访问
- 只有显式 `ALL` 才表示全量放行

这样可以避免因实现遗漏导致“默认放大权限”。

---

## 九、middleware 重构设计

中间件统一改为接收 `*auth.Realm`。

### 1. `RequireLogin(realm)`

负责：

- 解析 token
- 读取 claims
- 检查登录状态
- 注入 `login_id`
- 注入 `login_realm`
- 注入 claims 到 request context
- 建立 request cache

### 2. `RequirePermission(realm, permissions...)`

负责：

- 从 `claims.ACL.Permissions` 读取
- 默认 AND 判定

### 3. `RequirePermissionOr(realm, permissions...)`

负责：

- OR 判定

### 4. `RequireRole(realm, roles...)`

负责：

- 从 `claims.ACL.Roles` 读取

### 5. `RequireRoleOr(realm, roles...)`

负责：

- OR 判定

### 6. `AttachContext`

负责：

- 为 GORM hooks、日志、审计注入统一上下文

### 7. 预期效果

重构后中间件层：

- 不再出现 `if consumer else business`
- 不再传 `loginType string`
- 不直接查 provider

---

## 十、权限声明、扫描与缓存设计

### 1. 权限声明

路由注册时继续使用：

```go
registry.Perm("sys:user:page", "用户分页")
```

但底层统一登记到 auth 的权限注册中心。

### 2. 权限元数据结构

```go
type PermissionEntry struct {
    Code   string `json:"code"`
    Module string `json:"module"`
    Name   string `json:"name"`
}
```

### 3. 权限缓存拆成两个视图

#### 视图一：全量权限码

- key: `hei:auth:permission:codes`
- value: `[]string`

用途：

- 超级管理员权限全集
- 快速权限码集合读取

#### 视图二：模块树

- key: `hei:auth:permission:tree`
- value: `map[string]map[string]PermissionEntry`

用途：

- 后台权限展示
- 模块级查询

### 4. 读取接口

提供：

- `GetAllPermissionCodes()`
- `GetPermissionTree()`
- `GetModules()`
- `GetPermissionsByModule(module)`

### 5. 设计原则

不要再让一个 key 同时承担：

- 权限 code 列表
- 权限树

两种用途必须拆开。

---

## 十一、session 子系统设计

### 1. 单 realm session API

```go
auth.Business.Sessions().Page(ctx, query)
auth.Business.Sessions().Tokens(ctx, userID)
auth.Business.Sessions().KickoutUser(ctx, userID)
auth.Business.Sessions().KickoutToken(ctx, userID, token)
auth.Business.Sessions().Stats(ctx)
auth.Business.Sessions().Trend(ctx, days)
```

### 2. 多 realm 聚合 API

```go
auth.Sessions(auth.Business, auth.Consumer).Stats(ctx)
auth.Sessions(auth.Business, auth.Consumer).Trend(ctx, 7)
```

### 3. Redis 索引

保留现有主索引思路：

- `session:index`
- `session:expiry`
- `session:counts`
- `token:created`
- `token:expiry`
- `token:owners`

### 4. 新增 summary

新增：

- `session:summary:{userID}`

建议字段：

- `user_id`
- `token_count`
- `first_login_at`
- `last_login_at`
- `last_seen_at`
- `device_types`
- `username`

### 5. 查询优化策略

#### 分页页

1. zset 取 userIDs
2. pipeline 取 summary
3. 不扫描所有 token claims

#### token 明细页

1. 读 session token set
2. MGET token claims
3. pipeline TTL

#### 统计页

1. 直接读索引
2. 不扫 claims

### 6. 更新时机

#### 登录时

- 写 claims
- 更新 index
- 更新 summary

#### 续期时

- 更新 expiry
- 更新 last seen

#### 强退时

- 删除 token
- 更新 token_count
- 修正索引

#### 过期清理时

- 修正 summary
- 修正索引

---

## 十二、并发、一致性与稳健性设计

这部分是本次重构必须补齐的核心约束。  
auth 的主要并发风险不在 Go 进程内锁，而在 Redis 多 key 更新、跨实例并发修改和 ACL 快照失效。

### 1. 真相源与派生数据

必须明确：

- token claims 是登录态真相源
- user session set 是用户 token 集合真相源
- session index / expiry index / summary 都是派生数据

这意味着：

- 派生数据允许短暂不一致
- 派生数据必须支持修复
- 任何需要强一致判断的逻辑，优先以 token claims 和 token 实际存在性为准

### 2. 多 key 更新策略

登录、续期、强退、清理都会同时写多个 Redis key：

- token key
- session set
- token owner
- session index
- expiry index
- summary

要求：

1. 同一操作的相关写入必须使用 pipeline 批量提交
2. 对关键路径需要尽量保持幂等
3. summary/index 写失败不应影响 token 主写入成功
4. summary/index 后续必须可修复

建议原则：

- token/session 主写入优先级最高
- index/summary 视作附属更新

### 3. 幂等语义

所有会话操作都要按幂等方式设计：

#### 登录

- 同一 token 只会创建一次
- 重复写索引不会造成脏数据扩大

#### 续期

- 必须先确认 token 仍存在
- 已被踢下线或已过期的 token 不允许被“复活”

#### 强退 token

- 多次执行结果应一致
- token 不存在时直接视为成功

#### 强退用户

- 删除全部 token 后，重复执行仍然安全

#### 过期清理

- 清理过程中发现 token 已不存在，直接修正索引
- 不要求与其他写操作强事务一致，但要求最终一致

### 4. 操作优先级

并发场景下，建议定义以下优先级：

1. 强退
2. 禁用
3. 过期清理
4. 续期
5. 普通读取

重点要求：

- 强退后的 token 不允许被续期重新挂回索引
- 续期前必须检查 token 是否仍存在
- 清理逻辑发现 token 被强退后，应只做索引修正，不做回填

### 5. session summary 修复策略

`session:summary:{userID}` 不是强真相，只是列表页性能优化结构。

因此：

- summary 缺失时可根据 token 集合重建
- token_count 不一致时可按真实 live tokens 修正
- last_login_at 不一致时以 live token 中最大创建时间修正

修复触发点：

- 分页读取发现 summary 缺失
- token 明细读取发现脏数据
- 后台定时清理任务

### 6. ACL 快照一致性策略

`Permissions`、`Roles`、`ScopeMap` 都是 session claims 的 ACL 快照。

快照机制带来的已知特性：

- 运行时鉴权快
- 权限变更存在旧 session 窗口

第一阶段控制策略：

- 权限或 scope 变更后，踢下线相关用户 session

后续可扩展：

- `ACLVersion`
- 敏感接口版本强校验

但第一阶段不把版本校验作为主路径，以控制实现复杂度。

### 7. 多实例场景稳健性要求

系统必须默认运行在多实例环境下可工作。

要求：

- 所有 Redis key 规则以 realm 为隔离前缀
- 不依赖单机内存保存 session 真相
- 清理逻辑可在任意实例执行
- 重复执行清理和强退不会破坏状态

### 8. Scope 构建并发约束

由于 `ScopeMap` 来源于：

- `RelUserPermission`
- `RelRolePermission`

因此要接受以下现实：

- 登录过程中，如果后台正在修改权限关系，用户拿到的是登录时刻快照
- 该快照在权限变更后需要失效

这不是 bug，而是快照模型本身的语义。  
第一阶段以“变更即刷新 session”作为闭环。

---

## 十三、插件层职责调整

### 1. `plugin-sys/provider`

保留：

- provider 实现

改造：

- 签名改为 `GetPermissionList(ctx, realmID, userID)`
- 签名改为 `GetRoleList(ctx, realmID, userID)`

### 2. `plugin-sys/session`

迁回 auth 的：

- session 分页核心
- token 列表核心
- stats/trend 核心
- kickout 核心
- 多 realm 聚合

保留在插件层的：

- handler
- VO/DTO
- 用户资料补充
- 关键词查用户 ID
- 权限码声明

### 3. `plugin-client/session`

同步迁移到新的 realm session API。

---

## 十四、实施步骤

### 阶段一：核心骨架重建

目标：

- 建立新 API 形态
- 建立新核心类型

工作项：

1. 新建 `RealmID`
2. 新建 `Realm`
3. 暴露 `auth.Business` / `auth.Consumer`
4. 建立内部 registry
5. 定义新 provider 接口
6. 新建 `SessionClaims` / `ACLSnapshot`

### 阶段二：认证核心迁移

目标：

- 让登录态全部基于 claims

工作项：

1. 重构登录流程
2. 重构 token/session Redis 存储
3. 重构 logout/renew/disable/kickout
4. 删除旧默认 BUSINESS 包级入口

### 阶段三：权限链路迁移

目标：

- 运行时权限链路只读 claims

工作项：

1. realm 增加权限/角色方法
2. 重构 `HasPermission` / `HasRole`
3. 增加 `ScopeFor`
4. 接入 `RelUserPermission` / `RelRolePermission` 的登录时 scope 构建
5. 权限变更后踢下线机制接入

### 阶段四：middleware 迁移

目标：

- HTTP 接入统一收口

工作项：

1. 改 `RequireLogin`
2. 改 `RequirePermission`
3. 改 `RequireRole`
4. 改 context 注入
5. 删除 realm 分支代码

### 阶段五：权限扫描与缓存迁移

目标：

- 权限元数据标准化

工作项：

1. 权限注册中心重构
2. 生成 `permission:codes`
3. 生成 `permission:tree`
4. 改后台权限读取逻辑

### 阶段六：session 子系统迁移

目标：

- session 核心能力回收至 auth

工作项：

1. realm 增加 `Sessions()`
2. 增加 summary
3. 实现 Page/Tokens/Kickout/Stats/Trend
4. 实现多 realm 聚合器
5. 增加 summary/index 修复逻辑

### 阶段七：插件与调用点迁移

目标：

- 全仓使用统一新 API

工作项：

1. 改 `plugin-sys/provider`
2. 改 `plugin-sys/session`
3. 改 `plugin-client/session`
4. 改所有 middleware 调用点
5. 改所有旧 auth API 调用点

### 阶段八：清理与验证

目标：

- 移除旧实现，完成验收

工作项：

1. 删除 `ToolForLoginType`
2. 删除旧 `loginType string` 型 API
3. 删除旧权限缓存不一致逻辑
4. 删除旧 middleware 分支
5. 补测试并跑回归

---

## 十五、验收标准

重构完成后，必须满足：

### 1. API 形态

- 全仓 auth 主路径统一为 `auth.Business.xxx`
- 不再使用裸 `loginType string`

### 2. middleware

- 不再存在 `BUSINESS/CONSUMER` 条件分支
- middleware 统一只吃 `*auth.Realm`

### 3. provider

- provider 仅暴露 `GetPermissionList` / `GetRoleList`
- provider 不再依赖 Gin

### 4. 运行时鉴权

- 权限、角色、scope 判断都从 session claims 读取
- 请求路径不重复查 provider
- scope 明确来源于登录时查询的 `RelUserPermission` / `RelRolePermission`

### 5. 权限缓存

- `permission:codes` 与 `permission:tree` 两个缓存视图分离
- 权限 code 列表与权限树各自独立读取

### 6. session 管理

- `plugin-sys/session` 不再持有会话核心编排
- session 分页页不再依赖全量 token payload 扫描
- summary/index 出现脏数据时可自动或惰性修复

### 7. ScopeMap 应用

- auth 只输出 `ScopeInfo`
- 业务模块通过独立 `ScopePolicy` 应用到查询
- auth 不拼业务 SQL
- scope 缺失时默认拒绝，不做隐式放大授权

---

## 十六、后续扩展建议

重构完成后，后续可平滑扩展：

1. 新 realm
- 例如 `auth.Developer`

2. ACL 版本控制
- 权限变更后不一定踢下线，可改成 ACL 版本失效

3. 更细粒度 session 聚合
- 多 realm 分页聚合
- 分设备统计

4. 业务侧统一 ScopePolicy 规范
- 为各模块提供统一模板与示例

---

## 十七、结论

这次重构完成后：

- `sdk/auth` 会从“能用但耦合偏重的认证工具包”变成“结构明确的认证鉴权与会话治理内核”
- 权限元数据、运行时鉴权、session 管理三条主线会彻底分开
- 业务侧对数据权限的使用也会收敛为统一 policy 模式

这份文档作为后续实施与验收基线使用。
