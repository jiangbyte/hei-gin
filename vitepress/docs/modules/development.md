# 模块开发规范

Hei Gin 采用 **Go Workspace + 插件化架构** 组织业务模块。每个业务插件是一个独立的 Go 模块，通过 `module.Register()` 自注册路由、权限、中间件、定时任务和种子数据。

## 插件结构约定

每个业务插件遵循统一的结构约定：

```
plugins/<plugin-name>/
├── plugin.go              # 插件入口：实现 api.Plugin 接口 + module.Register()
├── imports.go             # 导入所有子模块触发 init() 自注册
├── persistence.go         # 日志持久化（可选）
├── provider/              # Provider 实现（如权限/用户 Provider）
├── <module>/              # 子模块目录
│   ├── api/v1/
│   │   └── api.go         # 路由注册 + HTTP Handler
│   ├── service.go         # 业务逻辑层
│   ├── params.go          # 请求参数和响应结构体
│   └── model.go           # GORM 数据模型（可选）
└── go.mod / go.sum        # 独立 Go 模块
```

### 文件说明

| 文件 | 必选 | 说明 |
|------|------|------|
| `plugin.go` | 是 | 插件入口，实现 `api.Plugin` 接口，通过 `init()` 调用 `module.Register()` |
| `imports.go` | 是 | 导入所有子模块包，触发其 `init()` 函数执行自注册 |
| `persistence.go` | 否 | 日志持久化实现（plugin-sys 中实现了 `api.LogPersistenceAPI`）|
| `provider/` | 否 | Provider 实现（如权限查询、用户信息 Provider）|
| `params.go` | 是 | 定义请求参数和响应结构体 |
| `service.go` | 是 | 业务逻辑层，使用 GORM 进行数据操作 |
| `api/v1/api.go` | 是 | 路由注册函数和 HTTP Handler |
| `model.go` | 否 | GORM 数据模型定义（含自动迁移注册）|

## 文件模板

### plugin.go - 插件入口

```go
package plugin_sys

import (
    "hei-gin/sdk/plugin"
)

type SysPlugin struct {
    plugin.NoopPlugin
}

func (p *SysPlugin) Info() plugin.PluginInfo {
    return plugin.PluginInfo{
        Name:        "plugin-sys",
        Version:     "1.0.0",
        Description: "System management plugin",
    }
}

func (p *SysPlugin) Name() string { return "plugin-sys" }

func (p *SysPlugin) Init() error {
    // 插件初始化逻辑
    return nil
}

func init() {
    module.Register(&SysPlugin{})
}
```

### model.go - GORM 数据模型

```go
package sysuser

import "time"

type SysUser struct {
    ID        string     `gorm:"primaryKey;size:32" json:"id"`
    Username  string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
    Password  string     `gorm:"size:255;not null" json:"-"`
    RealName  string     `gorm:"size:64" json:"real_name"`
    Email     string     `gorm:"size:128" json:"email"`
    Status    int        `gorm:"default:1" json:"status"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}

func (SysUser) TableName() string { return "sys_user" }
```

### params.go - 参数定义

```go
package sysuser

// UserListReq 用户列表查询参数
type UserListReq struct {
    Page     int    `json:"page" form:"page"`
    PageSize int    `json:"page_size" form:"page_size"`
    Username string `json:"username" form:"username"`
    Status   int    `json:"status" form:"status"`
}

// UserCreateReq 创建用户请求参数
type UserCreateReq struct {
    Username string   `json:"username" binding:"required"`
    Password string   `json:"password" binding:"required"`
    RealName string   `json:"real_name" binding:"required"`
    Email    string   `json:"email"`
    Phone    string   `json:"phone"`
    RoleIDs  []string `json:"role_ids"`
}
```

### service.go - 业务逻辑

```go
package sysuser

import (
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "hei-gin/sdk/db"
    "hei-gin/sdk/exception"
)

type Service struct{}

// List 用户列表查询
func (s *Service) List(req *UserListReq) ([]SysUser, int64, error) {
    query := db.DB.Model(&SysUser{})
    if req.Username != "" {
        query = query.Where("username LIKE ?", "%"+req.Username+"%")
    }
    if req.Status > 0 {
        query = query.Where("status = ?", req.Status)
    }

    var total int64
    query.Count(&total)

    var users []SysUser
    query.Order("created_at DESC").
        Offset((req.Page - 1) * req.PageSize).
        Limit(req.PageSize).
        Find(&users)

    return users, total, nil
}

// Create 创建用户
func (s *Service) Create(req *UserCreateReq) error {
    var count int64
    db.DB.Model(&SysUser{}).Where("username = ?", req.Username).Count(&count)
    if count > 0 {
        return exception.NewBusinessError("用户名已存在", 400)
    }

    hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    user := SysUser{
        Username: req.Username,
        Password: string(hashedPwd),
        RealName: req.RealName,
        Email:    req.Email,
        Status:   1,
    }
    return db.DB.Create(&user).Error
}
```

### api/v1/api.go - 路由注册 + Handler

```go
package v1

import (
    "github.com/gin-gonic/gin"

    "hei-gin/sdk/middleware"
    authMiddleware "hei-gin/sdk/auth/middleware"
    "hei-gin/sdk/result"
    "hei-gin/sdk/log"
    "hei-gin/sdk/registry"
    "hei-gin/sdk/crud"
)

var srv = &Service{}

func RegisterRoutes(r *gin.RouterGroup) {
    h := &Handler{}
    r.GET("/list",
        authMiddleware.HeiCheckLogin(),
        registry.Perm("sys:user:list", "用户列表查询"),
        h.UserList,
    )
    r.POST("/create",
        authMiddleware.HeiCheckLogin(),
        registry.Perm("sys:user:create", "创建用户"),
        log.SysLog("创建用户"),
        h.UserCreate,
    )
}

type Handler struct{}

func (h *Handler) UserList(c *gin.Context) {
    var req UserListReq
    if err := c.ShouldBindQuery(&req); err != nil {
        panic(exception.NewBusinessError("参数错误", 400))
    }
    users, total, err := srv.List(&req)
    if err != nil {
        panic(exception.NewBusinessError(err.Error(), 500))
    }
    c.JSON(200, result.PageDataResult(c, users, total, req.Page, req.PageSize))
}

func (h *Handler) UserCreate(c *gin.Context) {
    var req UserCreateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        panic(exception.NewBusinessError("参数错误", 400))
    }
    if err := srv.Create(&req); err != nil {
        panic(err)
    }
    c.JSON(200, result.Success(c, nil))
}

// 在 init 中注册路由
func init() {
    registry.RegisterRoute(func(r *gin.Engine) {
        group := r.Group("/api/v1/sys/user")
        RegisterRoutes(group)
    })
}
```

## 创建新插件的完整步骤

### 第一步：创建插件目录

```bash
mkdir -p plugins/plugin-<name>/
cd plugins/plugin-<name>
go mod init hei-gin/plugins/plugin-<name>
```

### 第二步：添加依赖

```bash
go mod edit -require hei-gin/sdk@v0.0.0
go mod edit -replace hei-gin/sdk=../../sdk
```

### 第三步：创建子模块目录

```bash
mkdir -p <module>/api/v1
```

### 第四步：实现 GORM 模型

在 `<module>/model.go` 中定义数据模型：

```go
type MyEntity struct {
    ID        string     `gorm:"primaryKey;size:32" json:"id"`
    Name      string     `gorm:"size:64;not null" json:"name"`
    Status    string     `gorm:"size:16;default:ENABLED" json:"status"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}

func (MyEntity) TableName() string { return "sys_<table>" }
```

### 第五步：实现 service.go

使用 GORM 编写业务逻辑。

### 第六步：实现 params.go

定义请求和响应结构体。

### 第七步：实现 api/v1/api.go

实现 HTTP Handler 和路由注册（通过 `registry.RegisterRoute` 自注册）。

### 第八步：实现 plugin.go 和 imports.go

```go
// plugin.go
package plugin_<name>

import "hei-gin/sdk/module"

type MyPlugin struct{ module.NoopModule }
func (p *MyPlugin) Name() string { return "plugin-<name>" }

func init() { module.Register(&MyPlugin{}) }
```

```go
// imports.go
package plugin_<name>

import _ "hei-gin/plugins/plugin-<name>/<module>"
```

### 第九步：注册到 Workspace

编辑根目录 `go.work`：

```
go 1.25.10

use (
    .
    ./sdk
    ./api
    ./plugins/plugin-<name>
)
```

### 第十步：在 app/main.go 中导入

编辑 `app/main.go`：

```go
package main

import (
    "hei-gin/sdk/app"
    _ "hei-gin/plugins/plugin-<name>"
)

func main() {
    app.Run()
}
```

## 权限代码命名规范

权限代码使用统一的命名规范：`<模块>:<操作>`。

### 标准操作

| 权限代码 | 说明 |
|---------|------|
| `<module>:list` | 列表查询 |
| `<module>:create` | 新增 |
| `<module>:update` | 修改 |
| `<module>:delete` | 删除 |
| `<module>:detail` | 详情查询 |
| `<module>:export` | 导出 |
| `<module>:import` | 导入 |

## 技术要点

### 统一响应

```go
import "hei-gin/sdk/result"

// 成功响应
c.JSON(200, result.Success(c, data))

// 分页响应
c.JSON(200, result.PageDataResult(c, records, total, page, size))

// 失败响应
c.JSON(200, result.Failure(c, message, code, data))
```

注意：`result.Failure()` 的参数顺序是 `(c, message, code, data)`。

### 业务异常

```go
import "hei-gin/sdk/exception"

panic(exception.NewBusinessError("用户名已存在", 400))
```

### 获取当前登录用户

```go
import authx "hei-gin/sdk/auth"

loginID := authx.GetLoginID(c)
if loginID == "" {
    panic(exception.NewBusinessError("未登录", 401))
}
```

### 通用 CRUD

`sdk/crud` 提供了通用分页、详情、删除函数：

```go
import "hei-gin/sdk/crud"

// 分页查询
c.JSON(200, crud.Page[MyEntity, *MyListReq](c, &MyEntity{}, param, buildQuery, "created_at DESC", toVO))

// 查询详情
crud.Detail[MyEntity](c, &entity, id, "实体名称")

// 删除
crud.Remove[MyEntity](c, &MyEntity{}, ids)
```

### 定时任务

```go
import "hei-gin/sdk/scheduler"

type CleanupTask struct{}
func (t *CleanupTask) Name() string { return "cleanup" }
func (t *CleanupTask) Run()         { /* ... */ }

func init() {
    scheduler.Register("0 */5 * * * *", &CleanupTask{})  // 每 5 分钟
    scheduler.RegisterInterval(30*time.Second, &CleanupTask{}) // 30 秒间隔
}
```

### DB 迁移

```bash
# 自动迁移所有插件注册的模型 + 种子数据
go run cmd/migrate/main.go

# 仅迁移，跳过种子
go run cmd/migrate/main.go -skip-seed
```

## 现有插件参考

编写新插件时，可以参考以下现有插件：

- **plugin-sys**：系统管理插件，包含认证、用户、角色、组织等完整管理功能
- **plugin-client**：C 端客户端插件，包含认证、会话、用户管理等
- **plugin-im**：WebSocket IM 插件，包含好友、群组、消息、广播等功能
