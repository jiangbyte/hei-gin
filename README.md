# HEI Gin

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.x-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Supported-DC382D?logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

HEI Gin 是 HEI 项目的 Go / Gin 后端模板：模块插件、ADMIN / PORTAL 双端、全局 stringly JSON、RBAC + 数据范围。

**单模块单体**：仓根一个 `go.mod`，`go run ./cmd/api` 直接启动。

1. **一般情况下**：整仓使用，改配置、加业务即可跑。
2. **复杂场景**：**可以改 framework**（会话、中间件、注册表等），不是黑盒。
3. **跟进上游**：用 **Git 合并本仓库代码**（merge / rebase）同步，**不是**把本项目当外部 `go get` 依赖来升级。

HTTP JSON 使用 **全局 stringly**：`boolean` 与数字在线上为字符串，对象与 list 保持结构（见 `internal/framework/core/stringly`，由 `response` / `bind.JSON` 统一挂载）。业务 DTO 写普通 `bool`/`int`，**不要**再引入包裹类型。

文档索引见 [docs/README.md](docs/README.md)。

> **请注意：** 生产仍需自行加固密钥、对象存储、Cookie Secure、TLS 后再上线。

## 姊妹项目

| 项目 | 说明 | 协议 |
| :--- | :--- | :--- |
| [**hei-boot**](https://github.com/jiangbyte/hei-boot) | Spring Boot 工程化脚手架 | Apache License 2.0 |
| [**hei-gin**](https://github.com/jiangbyte/hei-gin) | Go 轻量级后端框架 | MIT |
| [**hei-fastapi**](https://github.com/jiangbyte/hei-fastapi) | FastAPI 原型项目（早期阶段，仅供参考） | MIT |

## 能力一览

| 能力 | 说明 |
|------|------|
| 双端账号 | ADMIN / PORTAL；验证码、RSA 密码传输、登录失败锁定、会话绑定 |
| IAM | 账号、角色、部门、用户组、岗位、资源、权限、客户端、关系；账号/角色/用户组授权（own-*/grant-*）、资源按钮、权限注册表 |
| 系统 | 字典、配置（批量保存/测试推送）、Banner（门户互动）、文件（URL/预签名）、弱口令、审计（含告警规则）、**通知/公告/反馈**、**Go 代码生成** |
| 用户中心 | ADMIN / PORTAL 资料维护、头像、改密、换绑手机/邮箱（含身份登录开关） |
| 认证扩展 | OAuth（GitHub 完整；Gitee / 微信为配置桩）、三方绑定/解绑、会话管理（在线会话统计/强退）、忘记/重置密码、登录验证码 |
| 调度 | **SnailJob** 执行器嵌在 `api` 进程（`module.Job` 注册） |
| 可观测 | Prometheus `/metrics`（可关）、访问日志、安全头 / 可选 HSTS |
| 通知 | 邮件 / 短信 / 推送门面（默认关闭，见 `notify`） |
| 存储 | 本地目录或 S3 兼容；公开路径 `/api/v1/files/**`；Portal 端受限文件访问 |
| 前端同仓 | `web/admin`（Vue，与 hei-boot 前端同源同步） |

内置业务模块由 [`internal/modules/all`](internal/modules/all) blank import 汇总；可用 `modules.disabled` / `modules.enabled` 过滤。

## 仓库结构

```text
go.mod                     # 唯一 Go module：hei-gin
cmd/                       # 可执行入口（main 包）
  api/                     # 唯一运行入口：go run ./cmd/api
  migrate/                 # goose 运维命令（up/down/status）
internal/                  # 私有代码（外部不可导入）
  app/                     # 装配：基础设施 + HTTP + SnailJob
  framework/               # 可改的运行时（core/middleware/platform）
  modules/                 # 业务包（auth/iam/profile/sys/dashboard/health/biz/shared）
    all/                   # 汇总 blank import
scripts/                   # db.sql / migrate.sh / docker / sql
configs/                   # config.example.yaml
web/                       # admin（Vue）
storage/                   # 本地文件存储（.gitignore）
```

目录上的 `internal/modules/*` 只是包边界，**全部属于同一个 module**，本地改完即可被 `go run` 编进单体进程。

## 二次开发

| 诉求 | 做法 |
|------|------|
| 跟进上游 bugfix / 新内置模块 | `git fetch` + **merge/rebase**；保留 `_ "hei-gin/internal/modules/all"`，合并 `all` 即可带上新官方模块 |
| 默认使用 | 在仓根 `go run ./cmd/api`；业务写在 `internal/modules/<name>` |
| 只加自有业务 | 新包 `init` 里 `module.Register`；在 **自己的** `cmd` 里再 `_` import（少改官方 `all`） |
| 关掉某内置 | 配置 `modules.disabled`，不必删代码 |
| 改框架行为 | **直接改本仓 `internal/framework/`** |
| 注册定时任务 | 在模块 `Jobs` 中挂 `module.Job{Name, Run}`，由 SnailJob 客户端注册同名执行器 |

## 快速启动

在**仓库根**（需本机 Postgres / Redis，库名见配置）：

```bash
cp configs/config.example.yaml config.yaml
# CREATE DATABASE hei_gin;
# CREATE DATABASE snail_job;  # 调度中心库（独立）

psql -U postgres -d hei_gin -f scripts/db.sql   # 建表 + seed（与 hei-boot 同构 schema）
# 可选：docker compose -f scripts/docker/docker-compose.snailjob.yml up -d
go run ./cmd/api
```

默认地址：`http://127.0.0.1:8000`  
指标（默认开启）：`http://127.0.0.1:8000/metrics`

### 默认账号

| 账号 | 密码 | 说明 |
|------|------|------|
| `superadmin` | `123456` | seed 超管（含 `*:*:*`） |

### 登录

1. `GET /api/v1/admin/captcha` — 4 位字母验证码；答案 bcrypt 存 Redis（`captcha:{id}`）
2. `GET /api/v1/admin/password-key` — RSA-2048；`public_key` 为 SubjectPublicKeyInfo DER 的 base64；私钥在 Redis（`password:crypto:{id}`）
3. `POST /api/v1/admin/login` — `password` 为 OAEP-SHA256 加密后的 Base64；必须带 `password_key_id`（用后删密钥）

门户同理，前缀换为 `/api/v1/portal/...`。

会话：HttpOnly Cookie `Authorization`（可配 Path / SameSite / Secure）或不透明 Header（**非** `Bearer`）。

OAuth（可选）：配置 `oauth.github` 等后，走 `/api/v1/{admin|portal}/oauth/{provider}/authorize` → `callback` → `exchange`。

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
job.go        # 可选：module.Job 定时任务
```

有写库就必须有 `repo`；无持久化（如部分 health）可不造空 repo。样板见 `iam/account`。

## 代码生成（Go）

`sys/codegen` 管理端提供完整的 **Go 代码生成**（对齐 hei-boot 方案模型，输出 Go 而非 Java）：

- 方案 CRUD（`sys_codegen_plan`）+ 字段配置（`sys_codegen_field`，表列反射同步、控件/字典/查询条件）
- 数据库表 / 列元数据（`information_schema`）
- 四种生成类型：单表 `TABLE`、树表 `TREE`、左树右表 `LEFT_TREE_TABLE`、主子表 `MASTER_DETAIL`
- 预览 / ZIP 下载：Go 后端（`internal/modules/<domain>/<pkg>` 的 model/param/result/repo/service/handler/register 同包分文件）+ Vue 前端（`web/admin`）+ 菜单权限 SQL
- 模板在 `internal/modules/sys/codegen/templates.go`，渲染逻辑在 `internal/modules/sys/codegen/emit.go`；与 hei-boot 的 Java 模板一一对应但改为 Go 风格

## 模块装配

业务包在各自 `init` 中调用 `module.Register`；[`internal/modules/all`](internal/modules/all) 仅作汇总 import。

`cmd/api`：

```go
import _ "hei-gin/internal/modules/all"
```

### 定时任务（SnailJob）

API 进程内嵌 **SnailJob** Go 客户端（`internal/framework/platform/snailjob`）。配置键为 `snail_job`（YAML 未写时用代码默认值）：

```yaml
snail_job:
  enabled: true
  server_host: 127.0.0.1
  server_port: "17888"
  host_ip: 127.0.0.1
  host_port: "17889"
  namespace: c8f1a2b3d4e5461789abcdef01234567
  group_name: hei_gin_admin
  token: SJ_heiGinAdminToken1234567890abcd
```

本地 Server：见 [scripts/docker/README.md](scripts/docker/README.md)（`docker-compose.snailjob.yml` + `snailjob-flyway.sh`）。

内置 JobHandler 示例：

| JobHandler | 模块 |
|------------|------|
| `accountPurgeCancelledJob` | iam/account |
| `bannerStatusJob` | sys/banner |
| `auditAlertJob` | sys/audit（暴力破解/非常时段/敏感操作/批量删除/异地 IP 五类规则） |

数据库初始化用 `scripts/db.sql`（与 hei-boot schema 对齐；权限种子含 `sys:notice:*` / `sys:feedback:*`）。

上游若新增官方模块，通常只改 `all` 包；合并上游后即可自动注册。

## 主要 API 前缀

| 前缀 | 用途 |
|------|------|
| `/api/v1/admin/**` | 管理端 |
| `/api/v1/portal/**` | 门户端 |
| `/api/v1/internal/**` | 内部/健康 |
| `/api/v1/files/**` | 本地文件公开访问 |
| `/metrics` | Prometheus（`metrics.enabled`） |

响应信封：

```json
{ "code": "200", "message": "success", "data": {} }
```