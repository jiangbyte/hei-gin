# HEI Gin

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Supported-DC382D?logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-Apache_2.0-blue)

**HEI Gin** 是一个 Go / Gin 后端脚手架：一个进程同时提供**管理端（Admin）**与**门户（Portal）**双端 API，覆盖账号认证、组织权限（RBAC）、系统管理、消息反馈与运营工作台等常用能力。配合同仓维护的 Vue 3 / React / uni-app 三套前端，开箱即用、可按需裁剪。

## 特性

- **双端账号体系**：ADMIN / PORTAL 独立会话；密码 RSA 加密传输、验证码登录、登录锁定与限流、OAuth 三方登录
- **RBAC 权限**：账号 / 角色 / 部门 / 用户组 / 岗位，菜单与资源多层授权，数据范围过滤
- **系统管理**：字典、配置（敏感加密）、Banner、文件存储（本地 / S3 兼容）、公告通知、意见反馈、弱口令清单
- **运维能力**：操作审计与告警、运营工作台概览与 7 日趋势、内嵌任务调度（`sys_job` 管理台）
- **代码生成**：单表 / 树表 / 主子表方案，预览与 ZIP 下载（含前端与菜单权限 SQL）
- **三端前端**：`web/admin`（Vue 3 + Naive UI）、`web/portal`（React + Ant Design）、`web/admin-uniapp`（uni-app）

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
    <td width="50%"><img src="docs/images/admin-iam-resource.png" alt="资源授权" /></td>
    <td></td>
  </tr>
  <tr>
    <td align="center">资源授权</td>
    <td></td>
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

## 快速开始

### 环境要求

- Go 1.25+、PostgreSQL、Redis
- Node.js 22+ 与 pnpm 9+（前端）

### 初始化数据库

```bash
createdb -U postgres -h 127.0.0.1 hei_gin
psql -U postgres -h 127.0.0.1 -d hei_gin -f scripts/db.sql
```

### 启动后端

```bash
cp configs/config.example.yaml config.yaml
go run ./cmd/api
```

启动后 API 位于 `http://127.0.0.1:8000`。

### 启动前端

```bash
cd web/admin && pnpm install && pnpm dev   # http://127.0.0.1:5173
cd web/portal && pnpm install && pnpm dev  # http://127.0.0.1:5174
```

### 默认账号

| 端 | 地址 | 账号 | 密码 |
| --- | --- | --- | --- |
| Admin | http://localhost:5173 | `superadmin` | `123456` |
| Portal | http://localhost:5174 | `user` | `123456` |

> 生产环境首次启动后请立即修改默认密码。

## 文档

- [docs/README.md](docs/README.md) — 架构与二次开发指南
- [configs/config.example.yaml](configs/config.example.yaml) — 配置样例
- [scripts/db.sql](scripts/db.sql) — 数据库 schema 与种子数据

## 姊妹项目

| 项目 | 说明 | 协议 |
| --- | --- | --- |
| [**hei-boot**](https://github.com/jiangbyte/hei-boot) | Spring Boot 工程化脚手架 | Apache License 2.0 |
| [**hei-gin**](https://github.com/jiangbyte/hei-gin) | Go 轻量级后端框架 | Apache License 2.0 |
| [**hei-fastapi**](https://github.com/jiangbyte/hei-fastapi) | FastAPI 原型项目（早期阶段，仅供参考） | Apache License 2.0 |

## License

本项目使用 [Apache License 2.0](LICENSE) 开源协议，三个姊妹项目协议一致。完整条款见 [LICENSE](LICENSE)，版权归属声明见 [NOTICE](NOTICE)。
