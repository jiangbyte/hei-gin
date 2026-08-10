# HEI Gin

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.x-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Supported-DC382D?logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

HEI Gin 是 HEI 项目的 Go / Gin 后端模板。设计思想对齐原型 [hei-fastapi](https://github.com/jiangbyte/hei-fastapi)（模块插件、双端 ADMIN/PORTAL、双配置、wire 字符串 JSON、RBAC + 数据范围）。

工程上用 **`go.work` 多模块工作区**（类比 Maven reactor：多个 `go.mod`，根目录无业务 `go.mod`）——对齐 boot 的工程边界，而不是照抄 Java：

1. **一般情况下**：整仓使用，改配置、加业务即可跑。
2. **复杂场景**：**可以改 framework**（会话、中间件、注册表等），不是黑盒。
3. **跟进上游**：用 **Git 合并本仓库代码**（merge / rebase）同步，**不是**把本项目当外部 `go get` 依赖来升级。

HTTP JSON 对齐 boot 的 **全局 stringly**：`boolean` 与数字在线上为字符串，对象与 list 保持结构（见 `framework/core/stringly`，由 `response` / `bind.JSON` 统一挂载）。业务 DTO 写普通 `bool`/`int`，**不要**再引入包裹类型。

文档索引见 [docs/README.md](docs/README.md)。

> **请注意：** 生产仍需自行加固密钥、对象存储、Cookie Secure、TLS 后再上线。

## 仓库结构（多模块）

根目录只有 `go.work`（类似 Maven 父 POM / reactor），**没有**根 `go.mod`。每个子目录一个独立 module：

```text
go.work                    # 工作区清单（use 下列全部 module）
framework/                 # hei-gin/framework — 可改的运行时
modules/
  shared/                  # 跨业务共享模型/工具
  auth/ iam/ sys/ …        # 粗粒度业务 module（各有 go.mod）
app/                       # hei-gin/app — 组装根（类似 boot 的 admin 应用）
  cmd/{api,worker,migrate}
  internal/app/            # OpenInfra + AttachRegisteredModules
  internal/modules/all/    # 汇总 blank import 全部内置业务 module
migrations/                # goose SQL（仓库根，cwd 从根跑）
web/                       # 前端（admin / portal / admin-uniapp）
config.yaml / scripts/
```

`app` 通过 `require` + `replace` 依赖 `framework` 与各 `modules/*`；本地开发时 `go.work` 把它们绑成一个工作区，改任一 module 立刻生效（无需 publish）。

## 二次开发

| 诉求 | 做法 |
|------|------|
| 跟进上游 bugfix / 新内置模块 | `git fetch` + **merge/rebase**；保留 `_ "hei-gin/app/internal/modules/all"`，合并 `all` 即可带上新官方模块 |
| 默认使用 | 在仓根 `go run ./app/cmd/api`；业务写在 `modules/<name>` 或自有 module |
| 只加自有业务 | 新建 `modules/xxx`（自有 `go.mod`），`init` 里 `module.Register`；在 **自己的** `cmd` 里再 `_` import（少改官方 `all`） |
| 关掉某内置 | 配置 `modules.disabled`，不必删代码 |
| 改框架行为 | **直接改本仓 `framework/`**，再随业务一起合并上游 |

`go.work` / `replace` 只服务 **本仓内** 模块边界与本地开发，不是靠依赖坐标升级 framework。

## 快速启动

在**仓库根**执行（保证读到 `config.yaml` / `migrations/`）：

```bash
cp config.example.yaml config.yaml
# CREATE DATABASE hei_gin;

./scripts/migrate.sh
go run ./app/cmd/api
```

默认地址：`http://127.0.0.1:8000`

### 默认账号

| 账号 | 密码 | 说明 |
|------|------|------|
| `superadmin` | `123456` | seed 超管（含 `*:*:*`） |

### 登录（对齐 fastapi transport）

1. `GET /api/v1/admin/captcha` — 4 位字母验证码；答案 bcrypt 存 Redis（`captcha:{id}`）
2. `GET /api/v1/admin/password-key` — RSA-2048；`public_key` 为 SubjectPublicKeyInfo DER 的 base64；私钥在 Redis（`password:crypto:{id}`）
3. `POST /api/v1/admin/login` — `password` 为 OAEP-SHA256 加密后的 Base64；必须带 `password_key_id`（用后删密钥）

会话：HttpOnly Cookie `Authorization`（Path `/api/v1/{admin|portal}`）或不透明 Header（**非** `Bearer`）。

## 业务包分层（同包分文件）

每个领域包（如 [`modules/iam/account`](modules/iam/account)）按协作边界拆文件，**不**建 `param`/`result` 子 package：

```text
register.go   # init → module.Register
model.go      # GORM entity（原生类型）
param.go      # 入参
result.go     # 出参（需要时）
repo.go       # 持久化；handler 禁止直接碰 DB
service.go    # 业务规则
handler.go    # Bind → Service → response；JSON 用 bind.JSON
```

有写库就必须有 `repo`；无持久化（如部分 health）可不造空 repo。样板见 `iam/account`。

## 模块装配

业务 module 在各自包 `init` 中调用 `module.Register`；[`app/internal/modules/all`](app/internal/modules/all) 仅作汇总 import。

`app/cmd/api` / `worker`：

```go
import _ "hei-gin/app/internal/modules/all"
```

上游若新增官方模块，通常只改 `all` 包；合并上游后即可自动注册。

## 主要 API 前缀

| 前缀 | 用途 |
|------|------|
| `/api/v1/admin/**` | 管理端 |
| `/api/v1/portal/**` | 门户端 |
| `/api/v1/internal/**` | 内部/健康 |
| `/api/v1/files/**` | 本地文件公开访问 |

响应信封：

```json
{ "code": "200", "message": "success", "data": {} }
```

## 与原型 / Boot

| | hei-fastapi | hei-boot | hei-gin |
|--|-------------|----------|---------|
| 上游同步 | 整仓 | 整仓合并 | **整仓 Git 合并**（非依赖升级） |
| 框架边界 | 同仓 platform | common 等 | `framework/` module |
| 多模块 | 包目录 | Maven modules | **`go.work` + 多 `go.mod`** |
| 装配 | ModuleSpec 发现 | 显式 Maven 依赖 | init 自注册 + `app/.../all` |
| 包内分层 | — | param/result/mapper/… | **同包** param/result/repo/service/handler |
| JSON 标量 | 字符串 wire | StringlyTypedJackson | **全局 stringly**（`framework/core/stringly`） |
| 复杂定制 | 改源码 | 改 common | **改 framework** |
