# HEI Gin

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-4169E1?logo=postgresql&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-Supported-4479A1?logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Supported-DC382D?logo=redis&logoColor=white)
![Vue](https://img.shields.io/badge/Admin-Vue%203-4FC08D?logo=vuedotjs&logoColor=white)
![React](https://img.shields.io/badge/Portal-React-61DAFB?logo=react&logoColor=black)
![License](https://img.shields.io/badge/License-Apache_2.0-blue)
![Version](https://img.shields.io/badge/version-1.0.0--beta-orange)

**HEI Gin** 是一套开箱即用的 Go / Gin 工程化后端脚手架：单个进程同时提供 **Admin** 与 **Portal** 双端 API，覆盖认证授权、组织权限、系统运维、消息通知与工作台等常见后台能力。与 [hei-boot](https://github.com/jiangbyte/hei-boot)、[hei-fastapi](https://github.com/jiangbyte/hei-fastapi) 保持 API 契约一致；前端统一使用姊妹项目 [hei-admin](https://github.com/jiangbyte/hei-admin)、[hei-portal](https://github.com/jiangbyte/hei-portal)。

> 当前版本：`1.0.0-beta` · 协议：[Apache License 2.0](LICENSE)

## 目录

- [界面预览](#界面预览)
- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [工程结构](#工程结构)
- [快速开始](#快速开始)
- [默认账号](#默认账号)
- [相关文档](#相关文档)
- [姊妹项目](#姊妹项目)
- [License](#license)

## 界面预览

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

### 管理端 · 登录 / 工作台

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

### 管理端 · 组织权限

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
    <td width="50%"><img src="docs/images/admin-iam-dept.png" alt="部门管理" /></td>
    <td width="50%"><img src="docs/images/admin-iam-group.png" alt="用户组管理" /></td>
  </tr>
  <tr>
    <td align="center">部门管理</td>
    <td align="center">用户组管理</td>
  </tr>
  <tr>
    <td width="50%"><img src="docs/images/admin-iam-position.png" alt="岗位管理" /></td>
    <td width="50%"><img src="docs/images/admin-iam-resource.png" alt="资源授权" /></td>
  </tr>
  <tr>
    <td align="center">岗位管理</td>
    <td align="center">资源授权</td>
  </tr>
  <tr>
    <td width="50%"><img src="docs/images/admin-iam-resource-module.png" alt="资源模块" /></td>
    <td width="50%"><img src="docs/images/admin-iam-client-resource.png" alt="客户端资源" /></td>
  </tr>
  <tr>
    <td align="center">资源模块</td>
    <td align="center">客户端资源</td>
  </tr>
</table>

### 管理端 · 系统运维

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
  <tr>
    <td width="50%"><img src="docs/images/admin-sys-session.png" alt="在线会话" /></td>
    <td width="50%"><img src="docs/images/admin-sys-login-log.png" alt="登录日志" /></td>
  </tr>
  <tr>
    <td align="center">在线会话</td>
    <td align="center">登录日志</td>
  </tr>
</table>

### 管理端 · 消息与文件

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

### 管理端 · 业务示例

<table>
  <tr>
    <td width="50%"><img src="docs/images/admin-biz-order.png" alt="订单示例" /></td>
    <td></td>
  </tr>
  <tr>
    <td align="center">订单示例</td>
    <td></td>
  </tr>
</table>

## 功能特性

- **双端账号体系**：ADMIN / PORTAL 独立会话（Redis 不透明 Token，键对齐 fastapi `login:*`）；密码 RSA 传输、验证码登录、失败锁定与限流；可配置三方 OAuth 登录
- **RBAC 权限**：账号 / 角色 / 部门 / 用户组 / 岗位；菜单、按钮与 API 资源授权；在线会话踢出
- **系统管理**：字典、动态配置（`sys_config`，敏感项可加密）、Banner、公告 / 通知、意见反馈、弱口令库
- **对象存储**：S3 兼容存储（MinIO / RustFS / 阿里云 OSS 等），引擎与凭证走运行时配置，直链或预签名访问
- **运维能力**：操作审计与告警、登录日志、工作台（快捷应用与近期活动）、内嵌任务调度（`sys_job`；DB 扫描 + Redis 锁 + cron）
- **代码生成**：单表 / 树表 / 主子表方案，预览与 ZIP 下载（含前端与菜单权限 SQL，输出至 `../hei-admin`）
- **实名认证**：工单提交与审核、敏感字段加密存储（对齐 hei-boot）

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 后端 | Go 1.25+ · Gin · 单 module 单体（`cmd/api`） |
| 持久化 | PostgreSQL / MySQL · GORM（按 `db.driver` 或 DSN 二选一） |
| 缓存 / 会话 | Redis（go-redis）· 不透明会话 Token |
| 配置 | Viper（`config.yaml`）+ 运行时 `sys_config` |
| 其他 | AWS SDK v2（S3）· zap · robfig/cron · snowflake |
| 管理端前端 | [hei-admin](../hei-admin)（Vue 3 + Naive UI） |
| 门户前端 | [hei-portal](../hei-portal)（React + Ant Design） |

## 工程结构

```text
hei-gin
├── cmd/api                   # 唯一可启动入口
├── internal/
│   ├── app/                  # 装配根（基础设施 + 模块钩子）
│   ├── framework/            # 可修改运行时（config / security / middleware / storage / gojob …）
│   └── modules/              # 业务模块（auth / iam / sys / profile / workspace / biz …）
├── configs/config.example.yaml
├── scripts/db.sql            # PostgreSQL 结构 + 种子数据
├── scripts/db.mysql.sql      # MySQL 结构 + 种子数据
└── docs/images               # README 截图
```

## 快速开始

### 环境要求

- Go **1.25+**
- PostgreSQL **或** MySQL 8+、Redis
- 前端见 [hei-admin](../hei-admin)、[hei-portal](../hei-portal)（Node.js 22+、pnpm 9+）

### 1. 初始化数据库

**PostgreSQL：**

```bash
createdb -U postgres -h 127.0.0.1 hei_gin
psql -U postgres -h 127.0.0.1 -d hei_gin -f scripts/db.sql
```

**MySQL 8+：**

```bash
mysql -u root -p -e "CREATE DATABASE hei_gin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p hei_gin < scripts/db.mysql.sql
```

在 [`configs/config.example.yaml`](configs/config.example.yaml) 中配置 `db.driver` / `db.url`（也可只写 URL，由 scheme 推断）：

```yaml
db:
  driver: postgres
  url: postgres://postgres:123456@127.0.0.1:5432/hei_gin?sslmode=disable
  # driver: mysql
  # url: mysql://root:123456@127.0.0.1:3306/hei_gin?charset=utf8mb4&parseTime=true&loc=Local
```

> 本地 / 演示环境以对应脚本全量重建库表与种子数据。已有 PostgreSQL 库若列仍为 `jsonb`，可一次性改为 `json`：`ALTER COLUMN ... TYPE json USING col::json`。

### 2. 启动后端

```bash
cp configs/config.example.yaml config.yaml
# 按需修改 db / redis / app 等
go run ./cmd/api
```

| 项 | 地址 |
| --- | --- |
| API | http://127.0.0.1:8000 |

### 3. 启动前端（姊妹项目）

```bash
# 管理端 → http://127.0.0.1:5173
cd ../hei-admin && pnpm install && pnpm dev

# 门户 → http://127.0.0.1:5174
cd ../hei-portal && pnpm install && pnpm dev
```

详见 [hei-admin](../hei-admin/README.md)、[hei-portal](../hei-portal/README.md)。

## 默认账号

| 端 | 地址 | 账号 | 密码 |
| --- | --- | --- | --- |
| Admin | http://localhost:5173 | `superadmin` | `123456` |
| Portal | http://localhost:5174 | `user` | `123456` |

> 仅供本地演示。部署到非本机环境后请立即修改默认密码，并更换配置加密密钥、对象存储凭证等敏感项。

## 相关文档

| 文档 | 说明 |
| --- | --- |
| [`docs/README.md`](docs/README.md) | 架构约定与二次开发索引 |
| [`../hei-admin/README.md`](../hei-admin/README.md) | 管理端前端说明与环境变量 |
| [`../hei-portal/README.md`](../hei-portal/README.md) | 门户前端说明与环境变量 |
| [`configs/config.example.yaml`](configs/config.example.yaml) | 后端配置样例 |
| [`scripts/db.sql`](scripts/db.sql) | PostgreSQL 结构与种子数据 |
| [`scripts/db.mysql.sql`](scripts/db.mysql.sql) | MySQL 结构与种子数据 |

## 姊妹项目

| 项目 | 说明 | 协议 |
| --- | --- | --- |
| [**hei-boot**](../hei-boot) | Spring Boot 工程化脚手架 | Apache License 2.0 |
| [**hei-gin**](https://github.com/jiangbyte/hei-gin) | Go 轻量级后端框架（本仓库） | Apache License 2.0 |
| [**hei-fastapi**](../hei-fastapi) | FastAPI 异步脚手架 | Apache License 2.0 |
| [**hei-admin**](../hei-admin) | Vue 3 管理端 | Apache License 2.0 |
| [**hei-portal**](../hei-portal) | React 门户 | Apache License 2.0 |

## License

本项目基于 [Apache License 2.0](LICENSE) 开源。完整条款见 [LICENSE](LICENSE)，版权声明见 [NOTICE](NOTICE)。
