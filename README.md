# HEI Gin

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Supported-DC382D?logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-Apache_2.0-blue)

HEI Gin 是 HEI 项目的 Go / Gin 轻量级后端脚手架：**一个 Go 单体进程同时提供管理端（Admin）与门户（Portal）两套 API**，配合同仓维护的三个前端工程，覆盖账号认证、组织权限（RBAC）、系统管理、消息反馈与运营工作台等常用能力，开箱即用、可按需裁剪。

- **后端**：Go 1.25 · Gin 1.12 · PostgreSQL · Redis · GORM · 内嵌任务调度器（go-job）· 全局 stringly JSON
- **前端**：`web/admin`（Vue 3 / Naive UI）· `web/portal`（React / Ant Design）· `web/admin-uniapp`（uni-app），均与 [hei-boot](https://github.com/jiangbyte/hei-boot) 前端同源同步
- **数据约定**：对外 JSON 字段使用 `snake_case`；标量（含 boolean / 数字）统一按字符串收发（stringly，见 `internal/framework/core/stringly`，由 `response` / `bind.JSON` 统一挂载），对象与 list 保持结构；业务 DTO 写普通 `bool` / `int`，不要引入包裹类型

> **请注意：** 生产仍需自行加固密钥、对象存储、Cookie Secure、TLS 后再上线。

## 姊妹项目

| 项目 | 说明 | 协议 |
| :--- | :--- | :--- |
| [**hei-boot**](https://github.com/jiangbyte/hei-boot) | Spring Boot 工程化脚手架 | Apache License 2.0 |
| [**hei-gin**](https://github.com/jiangbyte/hei-gin) | Go 轻量级后端框架 | Apache License 2.0 |
| [**hei-fastapi**](https://github.com/jiangbyte/hei-fastapi) | FastAPI 原型项目（早期阶段，仅供参考） | Apache License 2.0 |

## 功能特性

**认证与账号（`modules/auth`）**

- 双端登录：ADMIN / PORTAL 两套独立账号体系与会话（Redis 会话存储，不透明令牌 / HttpOnly Cookie）
- 账号 / 邮箱 / 手机号多种身份登录，密码登录（RSA-2048 OAEP-SHA256 加密传输）与验证码登录
- 图形验证码（4 位字母，答案 bcrypt 存 Redis）、登录失败锁定与限流防护（账号 / IP 双维度）
- 忘记 / 重置密码、注册（门户）
- 三方登录（OAuth）：GitHub 完整实现，Gitee / 微信为配置桩；管理员可绑定 / 解绑
- 会话管理：在线会话统计 / 强制下线、并发会话数限制、会话绑定 IP

**组织与权限（`modules/iam`）**

- 账号、角色、部门、用户组、岗位管理
- 菜单资源、资源模块、客户端资源多层授权（RBAC）+ 权限注册表（PermissionRegistry，启动时同步）
- 数据范围：账号 / 角色 / 用户组授权（`own-*` / `grant-*`），树表按数据范围过滤与断言（对齐 hei-boot）

**系统管理（`modules/sys`）**

- 数据字典、系统配置（管理端配置页全部生效：`sys_config` 权威、yaml 回退，敏感配置 Fernet 加密存储）、Banner（门户互动）
- 文件存储（本地目录 / S3 兼容，支持预签名 URL）、弱口令清单
- 操作审计（Redis Stream 异步落库）+ 审计告警规则、**通知 / 公告 / 反馈**（管理端 + 门户双端）
- **Go 代码生成**：方案 CRUD、表列元数据、四种生成类型、预览 / ZIP 下载（Go 后端 + Vue 前端 + 菜单权限 SQL）

**运营与调度**

- 运营工作台（`modules/dashboard`）：账号、会话、审计、文件等核心指标概览与趋势
- 内嵌任务调度（gojob：`sys_job` 表驱动，`module.Job` 注册，`sys/job` 管理台维护）：注销账号清理、Banner 定时上下架、审计量级告警、互动计数落库

**用户中心（`modules/profile`）**

- ADMIN / PORTAL 资料维护、头像、改密、换绑手机 / 邮箱（含身份登录开关）

## 技术栈

| 分类 | 选型 |
| :--- | :--- |
| 语言 / 框架 | Go 1.25、Gin 1.12、单模块（仓根一个 `go.mod`，module `hei-gin`） |
| 数据 | PostgreSQL、GORM（Postgres 驱动）、gorm.io/datatypes、SQLite 驱动（本地调试可选） |
| 缓存 / 会话 | Redis（go-redis v9）、自研会话存储（不透明令牌） |
| 安全 | 双端会话、登录锁定 / 限流（Redis 固定窗口）、RSA-2048 密码加密、Fernet 敏感配置加密、安全头 / 可选 HSTS、CSRF 头约定（X-HEI-CSRF） |
| 观测 / 运维 | Prometheus `/metrics`（可关）、访问日志、结构化日志（zap / logrus）、操作审计 |
| 任务 | 内嵌调度器（`github.com/cybergarage/go-job` + `robfig/cron`，`sys_job` 表驱动，`sys/job` 管理台） |
| 配置 | viper（yaml + 环境变量）、运行时配置（runtimecfg：DB 权威、yaml 回退） |
| 其他 | snowflake（ID 生成）、AWS SDK v2（S3）、通知门面（邮件 / 短信 / 推送）、OpenTelemetry（可选） |

| 前端 | 技术 |
| :--- | :--- |
| `web/admin` | Vue 3.5、Naive UI 2、Pinia、Vue Router、Vite 8、TypeScript |
| `web/portal` | React 19、Ant Design 6、zustand、Vite 8、TypeScript |
| `web/admin-uniapp` | uni-app 3（H5 / 小程序） |

## 架构

后端只有**一个可运行应用** `cmd/api`，按请求前缀（`/api/v1/admin` / `/api/v1/portal`）区分管理端与门户两套接口，双账号体系会话相互隔离。业务按包划分，`internal/modules/all` blank import 汇总全部内置模块，各业务包在 `init` 中 `module.Register` 自注册；可用 `modules.disabled` / `modules.enabled` 过滤。

**单模块单体**：`internal/modules/*` 只是包边界，全部同属根 module，本地改完即可被 `go run` 编进单体进程。**跟进上游**用 Git 合并本仓库代码（merge / rebase），**不是**把本项目当外部 `go get` 依赖来升级；需要改框架行为时**直接改本仓 `internal/framework/`**（会话、中间件、注册表等，不是黑盒）。

| 分层 | 说明 |
| :--- | :--- |
| `internal/app` | 装配根：基础设施（DB / Redis / 存储 / 审计）+ HTTP + 内嵌任务调度器 |
| `internal/framework` | 可修改的运行时：core（config / security / response / bind / stringly / crypto / logger / schema / errors / context）、middleware（鉴权 / 限流 / 指标 / 安全头 / CORS / 错误处理）、platform（db / cache / audit / gojob / idgen / module / notify / otel / runtimecfg / storage） |
| `internal/modules` | 业务包：auth、iam、profile、sys、dashboard、health，以及样板模块 biz（代码生成样例） |
| `web/*` | 独立前端工程（无共享依赖层） |

## 快速开始

### 环境要求

- Go 1.25+、PostgreSQL、Redis
- Node.js 22+ 与 pnpm 9+（前端）

### 1. 初始化数据库

以 `scripts/db.sql` 为权威建表与种子数据源（与 hei-boot 同构 schema，含全部表结构、菜单、权限、字典、配置与默认账号；**无迁移步骤**，schema / 种子变更以该文件为准，由人工 / 运维执行）。

```bash
# 创建数据库
createdb -U postgres -h 127.0.0.1 hei_gin

# 导入库表与种子数据（也可用 Navicat / DataGrip 等工具直接执行该文件）
psql -U postgres -h 127.0.0.1 -d hei_gin -f scripts/db.sql
```

### 2. 启动后端

开发默认配置见 `configs/config.example.yaml`（复制为 `config.yaml`）：

- 数据库：`postgres://postgres:123456@127.0.0.1:5432/hei_gin`
- Redis：`redis://:123456@127.0.0.1:6379/4`

```bash
cp configs/config.example.yaml config.yaml
go run ./cmd/api
```

启动后可访问：

| 地址 | 说明 |
| :--- | :--- |
| http://127.0.0.1:8000 | Admin / Portal API |
| http://127.0.0.1:8000/metrics | Prometheus 指标（默认开启，`metrics.enabled` 可关） |

> 内嵌任务调度器（gojob）随 `api` 进程启动，按 `sys_job` 表中的启停状态执行定时任务；本地未初始化 `sys_job` 种子时无任务执行，不影响主流程。

### 3. 启动前端

```bash
cd web/admin && pnpm install && pnpm dev    # http://127.0.0.1:5173
cd web/portal && pnpm install && pnpm dev   # http://127.0.0.1:5174
```

前端开发模式通过 Vite 将 `/api` 代理到后端 `http://127.0.0.1:8000`。

### 默认账号

| 端 | 地址 | 账号 | 密码 |
| :--- | :--- | :--- | :--- |
| Admin | http://localhost:5173 | `superadmin` | `123456` |
| Portal | http://localhost:5174 | `user` | `123456` |

> 登录需要图形验证码（4 位字母，答案 bcrypt 存 Redis `captcha:{id}`，TTL 5 分钟）。本地自动化调试可用 `scripts/smoke`（验证码 → RSA 加密 → 登录 → /me）。**生产环境首次启动后请立即修改默认密码。**

### 登录流程（RSA）

1. `GET /api/v1/admin/captcha` — 4 位字母验证码；答案 bcrypt 存 Redis（`captcha:{id}`）
2. `GET /api/v1/admin/password-key` — RSA-2048；`public_key` 为 SubjectPublicKeyInfo DER 的 base64；私钥在 Redis（`password:crypto:{id}`）
3. `POST /api/v1/admin/login` — `password` 为 OAEP-SHA256 加密后的 Base64；必须带 `password_key_id`（用后删密钥）

门户同理，前缀换为 `/api/v1/portal/...`。会话：HttpOnly Cookie `Authorization`（可配 Path / SameSite / Secure）或不透明 Header（**非** `Bearer`）。

OAuth（可选）：配置 `oauth.github` 等后，走 `/api/v1/{admin|portal}/oauth/{provider}/authorize` → `callback` → `exchange`。

## 界面预览

前端与 hei-boot 同源同步，界面一致。

### 门户 Portal

<table>
  <tr>
    <td width="50%"><img src="docs/images/portal-login.png" alt="门户登录" /></td>
    <td width="50%"><img src="docs/images/portal-home.png" alt="门户首页" /></td>
  </tr>
  <tr>
    <td align="center">登录</td>
    <td align="center">首页</td>
  </tr>
</table>

### 管理端 Admin · 登录 / 工作台

<table>
  <tr>
    <td width="50%"><img src="docs/images/admin-login.png" alt="管理端登录" /></td>
    <td width="50%"><img src="docs/images/admin-dashboard.png" alt="运营工作台" /></td>
  </tr>
  <tr>
    <td align="center">登录</td>
    <td align="center">运营工作台</td>
  </tr>
</table>

### 管理端 Admin · 组织权限

<table>
  <tr>
    <td width="50%"><img src="docs/images/admin-iam-account.png" alt="账号管理" /></td>
    <td width="50%"><img src="docs/images/admin-iam-role.png" alt="角色管理" /></td>
  </tr>
  <tr>
    <td align="center">账号管理</td>
    <td align="center">角色管理</td>
  </tr>
  <tr>
    <td width="50%"><img src="docs/images/admin-iam-resource.png" alt="资源授权" /></td>
    <td></td>
  </tr>
  <tr>
    <td align="center">资源授权</td>
    <td></td>
  </tr>
</table>

### 管理端 Admin · 系统运维

<table>
  <tr>
    <td width="50%"><img src="docs/images/admin-sys-config.png" alt="系统配置" /></td>
    <td width="50%"><img src="docs/images/admin-sys-dict.png" alt="字典管理" /></td>
  </tr>
  <tr>
    <td align="center">系统配置</td>
    <td align="center">字典管理</td>
  </tr>
  <tr>
    <td width="50%"><img src="docs/images/admin-sys-audit.png" alt="操作审计" /></td>
    <td width="50%"><img src="docs/images/admin-sys-codegen.png" alt="代码生成" /></td>
  </tr>
  <tr>
    <td align="center">操作审计</td>
    <td align="center">代码生成</td>
  </tr>
</table>

### 管理端 Admin · 消息与文件

<table>
  <tr>
    <td width="50%"><img src="docs/images/admin-sys-banner.png" alt="Banner 管理" /></td>
    <td width="50%"><img src="docs/images/admin-message-notice.png" alt="公告通知" /></td>
  </tr>
  <tr>
    <td align="center">Banner 管理</td>
    <td align="center">公告通知</td>
  </tr>
  <tr>
    <td width="50%"><img src="docs/images/admin-message-feedback.png" alt="意见反馈" /></td>
    <td width="50%"><img src="docs/images/admin-sys-file.png" alt="文件管理" /></td>
  </tr>
  <tr>
    <td align="center">意见反馈</td>
    <td align="center">文件管理</td>
  </tr>
</table>

## 项目结构

```text
hei-gin
├── cmd
│   └── api                       # 唯一可运行入口：go run ./cmd/api
├── internal
│   ├── app                       # 装配根：基础设施 + HTTP + 内嵌任务调度器
│   ├── framework                 # 可修改的运行时（core / middleware / platform）
│   └── modules                   # 业务包（auth / iam / profile / sys / dashboard / health / biz 样板）
│       └── all                   # blank import 汇总全部内置模块
├── web                           # 前端（admin / portal / admin-uniapp）
├── docs                          # 文档与界面截图（文档索引见 docs/README.md）
├── configs
│   └── config.example.yaml       # 配置样例（复制为 config.yaml 使用）
├── scripts
│   ├── db.sql                    # 权威建表 + seed（与 hei-boot schema 对齐；无迁移步骤）
│   └── smoke                     # 开发辅助：登录链路冒烟（验证码 → RSA → 登录 → /me）
├── storage                       # 本地文件存储（.gitignore）
├── LICENSE                       # Apache License 2.0
└── NOTICE                        # 版权归属声明（jiangbyte）
```

## 主要 API

| 前缀 | 用途 |
| :--- | :--- |
| `/api/v1/admin/**` | 管理端接口 |
| `/api/v1/portal/**` | 门户接口 |
| `/api/v1/internal/**` | 内部 / 健康检查 |
| `/api/v1/files/**` | 公开文件读取（可配置） |
| `/metrics` | Prometheus 指标（`metrics.enabled` 可关） |

常用接口：`/api/v1/{admin|portal}/login`、`/captcha`、`/password-key`、`/oauth/**`、`/sys/**`（账号、角色、字典、配置、公告、反馈等）、`/profile/**`、`/dashboard/overview`。

门户公开接口（免登录白名单，`auth.auth_whitelist` 可配）：`/api/v1/portal/sys/banners/list`、`/api/v1/portal/sys/banners/interaction`、`/api/v1/portal/sys/dicts/tree`、`/api/v1/portal/sys/notices/list`。

响应信封：

```json
{ "code": "200", "message": "success", "data": {} }
```

## 配置说明

配置为单个 yaml（`configs/config.example.yaml` → `config.yaml`，亦可通过环境变量覆盖）；应用启动不做迁移，未写出的键使用代码默认值。

| 配置项 | 说明 | 默认（样例） |
| :--- | :--- | :--- |
| `app` | 应用名、监听地址 / 端口、debug、时区 | `127.0.0.1:8000` |
| `db` | PostgreSQL 连接（`postgres://user:pass@host:5432/hei_gin`）与连接池 | 本机 `postgres` / `123456` |
| `redis` | 会话、验证码与审计（`redis://:pass@host:6379/db`） | `127.0.0.1:6379/4` |
| `auth` | 令牌 TTL、登录锁定 / 限流阈值、验证码 TTL、Cookie 参数、白名单、门户注册开关 | 见样例 |
| `storage` | 文件存储：本地目录或 S3 兼容（endpoint / key / bucket / presign） | local |
| `modules` | `disabled` / `enabled` 过滤内置模块 | 空 |
| `audit` | 操作审计队列与 Redis Stream 消费组 | 开启 |
| `oauth` | GitHub / Gitee / 微信三方登录凭据 | 空（桩） |
| `metrics` | Prometheus 指标开关与路径 | `/metrics` |
| `security` | HSTS 开关与时长 | 关 |
| `notify` | 邮件 / 短信 / 推送门面 | 默认关闭 |
| `otel` | OpenTelemetry 开关与 endpoint | 关 |
| `crypto` | 敏感配置 Fernet 密钥（`fernet_key`，示例内置 hei-boot 同款开发密钥）与可选 Vault | 开发内置默认 |

## 生产部署

### 构建后端二进制

```bash
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o hei-gin ./cmd/api
```

### 前端镜像

`web/admin`、`web/portal` 各自提供 Dockerfile（nginx 托管静态资源，`/api` 反向代理到后端，后端地址可配）。

### 生产必改配置

| 变量 / 配置 | 说明 |
| :--- | :--- |
| `db.url` / `redis.url` | 生产数据库 / 缓存连接 |
| `crypto.fernet_key` | 敏感配置 Fernet 密钥（示例为开发缺省，生产必须替换并重加密存量敏感配置） |
| `storage` | 本地目录 → S3 兼容对象存储；收紧 `public_path` |
| `auth.session_cookie_secure` | 仅在 HTTPS 下开启（配合 `hsts_enabled`） |
| `app.debug` | 关闭 debug 日志与 Gin DebugMode |

### 上线检查清单

- 修改 `superadmin` / `user` 默认密码，替换 `crypto.fernet_key` 为生产密钥并重加密存量敏感配置
- 关闭 debug、收紧 CORS 白名单与门户公开白名单
- 仅在可信反向代理 + TLS 后开启 Cookie Secure / HSTS
- 定时任务按 `sys_job` 表启停；生产按需关闭非必要任务，审计告警建议保持开启

## 二次开发

**业务包分层（同包分文件）**：每个领域包（样板见 `internal/modules/iam/account`）按协作边界拆文件，**不**建 `param` / `result` 子 package：

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

有写库就必须有 `repo`；无持久化（如部分 health）可不造空 repo。

**模块装配**：业务包在各自 `init` 中调用 `module.Register`；`internal/modules/all` 仅作汇总 import，`cmd/api` 保留 `_ "hei-gin/internal/modules/all"`。

| 诉求 | 做法 |
|------|------|
| 跟进上游 bugfix / 新内置模块 | `git fetch` + **merge / rebase**；合并 `all` 即可带上新官方模块 |
| 默认使用 | 在仓根 `go run ./cmd/api`；业务写在 `internal/modules/<name>` |
| 只加自有业务 | 新包 `init` 里 `module.Register`；在 **自己的** `cmd` 里再 `_` import（少改官方 `all`） |
| 关掉某内置 | 配置 `modules.disabled`，不必删代码 |
| 改框架行为 | **直接改本仓 `internal/framework/`** |
| 注册定时任务 | 在模块 `Jobs` 中挂 `module.Job{Name, Run}`，gojob 按 `sys_job` 表调度 |

**定时任务（gojob）**：API 进程内嵌任务调度器（`internal/framework/platform/gojob`），由 `sys_job` / `sys_job_log` 表驱动，`sys/job` 管理台维护启停与触发。内置 JobHandler：

| JobHandler | 模块 |
|------------|------|
| `accountPurgeCancelledJob` | iam/account（注销账号过期清理） |
| `bannerStatusJob` | sys/banner（Banner 定时上下架） |
| `bannerFlushInteractions` | sys/banner（互动计数 Redis 增量落库） |
| `auditAlertJob` | sys/audit（暴力破解 / 非常时段 / 敏感操作 / 批量删除 / 异地 IP 五类规则） |
| `sysFileCleanupLocalOrphans` | sys/file（本地孤立文件清理） |

**代码生成（Go）**：`sys/codegen` 管理端提供完整的 **Go 代码生成**（对齐 hei-boot 方案模型，输出 Go 而非 Java）：

- 方案 CRUD（`sys_codegen_plan`）+ 字段配置（`sys_codegen_field`，表列反射同步、控件 / 字典 / 查询条件）
- 数据库表 / 列元数据（`information_schema`）
- 四种生成类型：单表 `TABLE`、树表 `TREE`、左树右表 `LEFT_TREE_TABLE`、主子表 `MASTER_DETAIL`
- 预览 / ZIP 下载：Go 后端（`internal/modules/<domain>/<pkg>` 同包分文件）+ Vue 前端（`web/admin`）+ 菜单权限 SQL
- 模板在 `internal/modules/sys/codegen/templates.go`，渲染逻辑在 `internal/modules/sys/codegen/emit.go`；与 hei-boot 的 Java 模板一一对应但改为 Go 风格

开发约定：业务表继承通用审计字段（id / created_at / updated_at / created_by / updated_by）；领域服务 `XxxService`；权限用 `security.PermissionRegistry` 注册 + `security.HasPermission` 校验（`*:*:*` 通配）；联表查询用 GORM（`Preload` / `Joins`）。

## 代码贡献

欢迎 Issue 与 PR。提交前请确认：

- Controller 入参与出参符合 stringly 线格式约定（标量字符串化）
- 遵守分层与模块边界：`cmd` / `internal/app` / `internal/framework` / `internal/modules` / `web/*`
- 兼容 Go 1.25 工具链（`gofmt` / `go vet` 通过）；敏感配置走环境变量或 `crypto.fernet_key`；文档随行为同步

```bash
git checkout -b feature/your-change
go build ./... && go vet ./...
git commit -m "feat: describe your change"
```

## 开源协议

本项目使用 [Apache License 2.0](LICENSE) 开源协议，三个姊妹项目（hei-boot / hei-gin / hei-fastapi）协议一致。完整条款见 [LICENSE](LICENSE)，版权归属声明见 [NOTICE](NOTICE)。
