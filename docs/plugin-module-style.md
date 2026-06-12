# Plugin 模块分层规范

参考样板：`plugins/plugin-sys/banner`

本文档用于约束 `plugins/plugin-sys` 下的模块组织方式，目标是把仓储访问、业务逻辑、HTTP 协议处理拆开，同时避免引入过重的 Java 风格依赖注入写法。

---

## 设计目标

- `repository` 只负责数据访问
- `service` 只负责业务逻辑
- `api handler` 只负责参数绑定、调用 service、返回响应
- `module` 只负责模块装配

需要避免：

- 所有逻辑堆在 `service.go`
- 包级函数只是简单转发到 `defaultModule.service.xxx`
- handler 里直接查库
- 在循环里逐条查关联数据，导致 N+1

---

## 推荐目录结构

```text
plugins/plugin-sys/<module>/
  model.go
  params.go
  repository.go
  service.go
  module.go
  api/v1/api.go
```

说明：

- `mapper.go` 已并入 `params.go`
- 模块内不再单独放 `migrate.go`
- 迁移注册统一放在 `plugins/plugin-sys/migrate.go`

---

## 文件职责

| 文件 | 职责 |
|---|---|
| `model.go` | 实体结构体、`TableName()`、GORM 相关声明 |
| `params.go` | VO、分页参数、请求参数、Entity/VO 转换函数 |
| `repository.go` | DB/Redis/持久化访问 |
| `service.go` | 业务逻辑 |
| `module.go` | 依赖装配 |
| `api/v1/api.go` | 路由注册、handler、协议转换 |

---

## module.go 规范

`module.go` 负责模块默认依赖装配，不做业务逻辑。

推荐形式：

```go
type Module struct {
	service *Service
}

var DefaultModule = NewModule()

func NewModule() *Module {
	repo := &repository{db: db.DB}
	svc := &Service{repo: repo}
	return &Module{service: svc}
}

func (m *Module) Service() *Service {
	return m.service
}
```

规则：

- 使用显式装配，不使用容器
- 不引入 `interface + impl + factory` 这一整套 Java 风格结构
- 默认实例统一使用 `DefaultModule`
- 需要测试时可以手动构造 `Module` 或 `Service`

---

## repository.go 规范

`repository` 只负责数据访问。

推荐形式：

```go
type repository struct {
	db *gorm.DB
}
```

规则：

- 方法入参使用 `context.Context`
- 不接收 `*gin.Context`
- 不写 `result.Success`、`result.WriteError`
- 不做业务校验
- 不做权限判断
- 不做分页参数兜底
- 不做 VO 拼装

允许内容：

- `Create`
- `FindByID`
- `Page`
- `ListByIDs`
- `UpdateByID`
- `DeleteByIDs`
- 批量关系查询
- Redis/缓存读写

---

## service.go 规范

`service` 负责业务逻辑编排，依赖 `repository`。

推荐形式：

```go
type Service struct {
	repo *repository
}
```

规则：

- service 方法直接使用业务动作名
- 不要保留包级转发函数
- 纯计算函数可以保留为自由函数
- 只要函数依赖 repo，就应改成 `Service` 方法

推荐命名：

- `Page`
- `Detail`
- `Create`
- `Modify`
- `Remove`
- `Options`
- `ListAll`

避免命名：

- `BannerPage`
- `ConfigOptions`
- `ConfigEditBatch`
- `BannerListAll`

说明：

- 当前项目大量 service 仍直接写 HTTP 错误响应，因此 `service` 可以接收 `*gin.Context`
- 这是当前项目的务实方案，不强行拆成纯 application service

### Service 私有辅助方法

可以保留 `Service` 私有辅助方法，但要满足下面条件：

- 该方法依赖 `repo`
- 该方法依赖 `Service` 内部缓存、批量加载、装配状态
- 该方法只服务于当前模块内部流程，不需要对外暴露
- 该方法本质上是某个公开业务动作的子步骤

推荐示例：

```go
func (s *Service) Get(c *gin.Context) *HomeVO {
	userID := auth.GetLoginIDDefaultNull(c)
	vo := &HomeVO{}
	if userID != "" {
		vo.QuickActions = s.findQuickActionsByUserID(c.Request.Context(), userID)
		vo.AvailableResources = s.getAvailableResources(c.Request.Context(), userID)
	}
	vo.Notices = s.getNotices(c.Request.Context())
	return vo
}

func (s *Service) findQuickActionsByUserID(ctx context.Context, userID string) []QuickActionVO {
	actions := s.repo.ListQuickActions(ctx, userID)
	resourceIDs := make([]string, len(actions))
	for i, a := range actions {
		resourceIDs[i] = a.ResourceID
	}
	resources := s.repo.ListResourcesByIDs(ctx, resourceIDs)
	return assembleQuickActions(actions, resources)
}
```

适合做成 `Service` 私有方法的典型场景：

- 批量 enrich，如用户列表补角色、岗位、组织名
- 树结构组装前的批量数据预取
- service 内部缓存读写
- 某个公开动作下的子流程拆分

更适合保留为自由函数的场景：

- 不依赖 `repo`
- 不依赖 `Service` 状态
- 纯计算、纯格式化、纯树节点转换

推荐保留为自由函数的例子：

```go
func formatTimeout(seconds int) string
func buildUserMenuTree(cm map[string][]rawResource, pid string) []map[string]interface{}
func getParentIDKey(parentID *string) string
```

避免写法：

- 为了“面向对象”把所有纯函数都挂到 `Service` 上
- 把可复用的纯装配逻辑写成公开方法暴露出去
- 私有辅助方法里直接做 HTTP 响应之外，还被多个 handler 绕过公开方法直接调用
- 一个公开方法只剩一行转发，实际业务都堆在名称含糊的私有方法里

---

## api/v1/api.go 规范

handler 直接依赖注入后的 `service`，不要再调模块包里的空转函数。

推荐形式：

```go
type handler struct {
	service *banner.Service
}

var defaultHandler = newHandler(banner.DefaultModule)

func newHandler(module *banner.Module) *handler {
	return &handler{service: module.Service()}
}
```

推荐完整写法：

```go
type handler struct {
	service *banner.Service
}

var defaultHandler = newHandler(banner.DefaultModule)

func newHandler(module *banner.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/sys/banner/page", defaultHandler.page)
	r.POST("/api/v1/sys/banner/create", defaultHandler.create)
}

func (h *handler) page(c *gin.Context) {
	var param banner.BannerPageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Page(c, &param)
}

func (h *handler) create(c *gin.Context) {
	var vo banner.BannerVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Create(c, &vo)
	result.Success(c, nil)
}
```

路由注册直接绑定 handler 方法：

```go
r.GET("/api/v1/sys/banner/page", defaultHandler.page)
```

规则：

- handler 负责 `ShouldBindJSON` / `ShouldBindQuery`
- 参数绑定失败使用 `result.Failure`
- 业务逻辑交给 service
- 成功响应在 handler 层统一包装
- handler 不直接访问 repo
- 路由注册只绑定 handler 方法，不直接绑包级业务函数
- `defaultHandler` 只作为默认装配入口，不承担业务逻辑
- `newHandler` 只做装配，不做条件判断或兜底初始化

避免写法：

- `func Page(c *gin.Context) { defaultModule.service.Page(c) }`
- handler 里直接 `db.DB.Model(...)`
- handler 同时做参数绑定、查库、VO 拼装、响应输出
- 在 `api` 包内维护独立的全局 `repo` / `service`

---

## 命名规范

类型命名：

- `Module`
- `Service`
- `repository`
- `handler`

方法命名：

- service/repository 方法直接使用动作名
- 不重复模块名

示例：

- `Page`
- `Detail`
- `Create`
- `Modify`
- `Remove`
- `Options`

避免：

- `UserPage`
- `RoleDetail`
- `ConfigListByCategory` 这种仅用于包级暴露的重复命名

---

## N+1 规范

以下场景必须优先考虑批量查询：

- 列表页补充名称信息
- 用户角色、权限回填
- 组织、分组、岗位名称装配
- 树结构菜单构造
- 资源与角色关系加载

推荐做法：

- 使用 `IN ?` 批量查询
- 先查关系表，再组装 map
- service 层统一做 enrich
- 缓存放在 service 层，不放在 handler 层

禁止做法：

- 在循环里一条条 `FindByID`
- handler 中边遍历边查库
- 为了“拆层”把本来一次批量查询拆散成多次单条查询

---

## 迁移规范

- 模块内不要定义自己的 `migrate.go`
- 所有模型注册、种子注册统一放到 `plugins/plugin-sys/migrate.go`

这样可以避免：

- 每个模块重复维护迁移入口
- 插件根部装配不清晰

---

## banner 样板结论

`plugins/plugin-sys/banner` 当前作为标准样板：

- `module.go`：只负责装配
- `repository.go`：只负责数据访问
- `service.go`：只保留业务方法
- `api/v1/api.go`：handler 注入 `service`
- 删除了 `BannerPage/BannerDetail/...` 这类包级空转函数

后续 `plugin-sys` 其他模块按这套结构统一推进。
