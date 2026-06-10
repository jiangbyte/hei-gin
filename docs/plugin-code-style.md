# Plugin 代码风格约定

参考 `plugins/plugin-sys/user`、`plugins/plugin-sys/role`、`plugins/plugin-sys/resource` 对齐后的模式。

---

## 文件职责

| 文件 | 职责 |
|---|---|
| `model.go` | 实体结构体 + `TableName()` |
| `params.go` | VO / 分页参数 / 请求参数 — 纯类型，无方法 |
| `mapper.go` | Entity ↔ VO 转换函数 |
| `service.go` | 业务逻辑 |
| `migrate.go` | 注册模型、种子数据 |
| `api/v1/api.go` | 路由注册 + Handler（薄） |

---

## 细则

### 1. api.go

**Handler 命名**
- 统一 `*Handler` 后缀（`pageHandler`、`createHandler`、`detailHandler`）

**`init()` 位置**
- 紧跟在 `RegisterRoutes()` 之后

**参数传递**
- 参数传指针结构体，不在 handler 层拆字段

```go
// ✓ 正确
var param utils.IdsParam
c.ShouldBindJSON(&param)
role.RoleRemove(c, &param)

// ✗ 错误
role.RoleRemove(c, param.IDs)
```

**Query 参数**
- `c.Query()` 直接内联在调用里，不另赋变量

```go
// ✓ 正确
result.Success(c, role.RoleDetail(c, c.Query("id")))

// ✗ 错误
id := c.Query("id")
vo := role.RoleDetail(c, id)
result.Success(c, vo)
```

**权限与登录**
- Handler 不做 `auth.GetLoginIDDefaultNull` — 由 hooks 处理
- 注解如 `registry.Perm`、`log.SysLog`、`middleware.NoRepeat` 放在路由注册处

**校验**
- 参数绑定错误（`ShouldBindJSON` / `ShouldBindQuery` 失败）在 handler 用 `result.Failure` 返回
- 业务校验（空 ID、数据不存在等）在 service 层用 `result.WriteError` 处理

---

### 2. service.go

**函数命名**
- 统一 `{Module}{Action}` 格式（`RolePage`、`UserCreate`、`PositionModify`、`ResourceRemove`）

**分页**
- 手写，不使用 `crud` 包

```go
func RolePage(c *gin.Context, p *RolePageParam) {
    ctx := c.Request.Context()
    if p.Current < 1 {
        p.Current = 1
    }
    if p.Size < 1 {
        p.Size = 10
    }
    if p.Size > 100 {
        p.Size = 100
    }

    q := db.DB.WithContext(ctx).Model(&SysRole{})
    // 条件拼装...
    if p.Keyword != "" {
        q = q.Where("name LIKE ? OR code LIKE ?", like, like)
    }

    var total int64
    q.Count(&total)

    var rows []SysRole
    q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

    vos := make([]*RoleVO, len(rows))
    for i, r := range rows {
        vos[i] = SysRoleToRoleVO(&r)
    }
    result.PageDataResult(c, vos, total, p.Current, p.Size)
}
```

**Detail**
- 不使用 `crud.Detail`
- 查询风格：`db.DB.WithContext(ctx).First(&e, "id = ?", id)`（条件作为参数，不链式 `.Where()`）

```go
func RoleDetail(c *gin.Context, id string) *RoleVO {
    if id == "" {
        result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
        return nil
    }
    ctx := c.Request.Context()
    var e SysRole
    if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
            return nil
        }
        result.WriteError(c, exception.NewBusinessError("查询角色详情失败: "+err.Error(), 500))
        return nil
    }
    return SysRoleToRoleVO(&e)
}
```

**Create**
- 使用 `XxxVOToSysXxx` mapper 转换
- 只覆写业务默认值（如 `e.Status = string(enums.StatusEnabled)`）
- 不手动设 `ID`、`CreatedAt`、`UpdatedAt`、`CreatedBy`、`UpdatedBy` — 由 GORM hooks 处理

```go
func RoleCreate(c *gin.Context, vo *RoleVO) {
    ctx := c.Request.Context()
    e := RoleVOToSysRole(vo)
    e.Status = string(enums.StatusEnabled)
    if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
        result.WriteError(c, exception.NewBusinessError("添加角色失败: "+err.Error(), 500))
        return
    }
}
```

**Modify**
- 使用 `map[string]interface{}{}` 字面量初始化必有字段，条件添加可选字段
- 没有 `else { nil }` 分支 — 不传入的字段不更新，保留 DB 原值
- 不手动设 `updated_at`、`updated_by` — hooks 处理

```go
func RoleModify(c *gin.Context, vo *RoleVO) {
    ctx := c.Request.Context()
    if vo.ID == "" {
        result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
        return
    }

    var e SysRole
    if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
            return
        }
        result.WriteError(c, exception.NewBusinessError("查询角色失败: "+err.Error(), 500))
        return
    }

    up := map[string]interface{}{
        "code": vo.Code, "name": vo.Name, "category": vo.Category,
        "sort_code": vo.SortCode,
    }
    if vo.Description != nil {
        up["description"] = *vo.Description
    }
    if vo.Status != "" {
        up["status"] = vo.Status
    }
    if err := db.DB.WithContext(ctx).Model(&SysRole{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
        result.WriteError(c, exception.NewBusinessError("编辑角色失败: "+err.Error(), 500))
        return
    }
}
```

**Remove**
- 参数 `param *utils.IdsParam`（指针结构体）
- 用事务时注意 `tx.Rollback()` + `result.WriteError` + `return`

```go
func RoleRemove(c *gin.Context, param *utils.IdsParam) {
    ids := param.IDs
    if len(ids) == 0 {
        return
    }
    ...
}
```

**其他规范**
- 不依赖 `sdk/auth` 包
- 不使用 `crud` 包的任何函数
- 不使用 `time` 包（hooks 处理时间）
- 状态值使用 `enums.StatusEnabled` 等常量，不硬编码字符串
- Entity ↔ VO 转换统一走 mapper 函数

---

### 3. params.go

- 纯结构体定义，不包含辅助方法
- 无 `toVO()`、`GetCurrent()` / `GetSize()` 等

---

### 4. mapper.go

- `SysXxxToXxxVO` — Entity → VO，处理 `*time.Time → string` 转换

```go
func SysRoleToRoleVO(src *SysRole) *RoleVO {
    if src == nil { return nil }
    dst := &RoleVO{}
    dst.ID = src.ID
    // ... 逐字段赋值
    dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
    dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
    return dst
}
```

- `XxxVOToSysXxx` — VO → Entity，处理 `string → *time.Time` 转换

```go
func RoleVOToSysRole(src *RoleVO) *SysRole {
    if src == nil { return nil }
    dst := &SysRole{}
    dst.ID = src.ID
    // ... 逐字段赋值
    dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
    dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
    return dst
}
```

---

### 5. model.go

- 实体结构体，gorm tag 完整
- `TableName()` 方法

---

### 6. migrate.go

- `init()` 中通过 `db.RegisterModel()` 注册模型
- 需要种子数据时通过 `db.RegisterSeed()` 注册

---

## 错误处理

### 原则

不使用 `panic`，改用 **显式错误响应 + `c.Abort()`** 模式。

### service 层

发生错误时调用 `result.WriteError(c, err)` 然后 `return`：

```go
func RoleCreate(c *gin.Context, vo *RoleVO) {
    if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
        result.WriteError(c, exception.NewBusinessError("添加角色失败: "+err.Error(), 500))
        return
    }
}
```

空值校验也用 WriteError，不用 return nil：

```go
func RoleDetail(c *gin.Context, id string) *RoleVO {
    if id == "" {
        result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
        return nil   // 值返回函数需要 return 零值
    }
    ...
}
```

### handler 层

handler 不需要检查错误，`result.WriteError` 已写入响应并 `c.Abort()`，后续 `result.Success` 在已 abort 状态下自动成为 no-op：

```go
func detailHandler(c *gin.Context) {
    vo := role.RoleDetail(c, c.Query("id"))  // 出错时内部已 WriteError + Abort
    result.Success(c, vo)                     // 已 abort 时 no-op
}
```

### result 包行为

- `result.WriteError` — 写入错误 JSON 响应 + 调用 `c.Abort()`
- `result.Success` / `result.Failure` / `PageDataResult` — 如果 `c.IsAborted()`，直接 return，不会二次写入
- `Recovery` middleware 仍作为兜底，捕获未处理的 panic
