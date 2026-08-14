# HEI Admin

Vue 3 管理端。账号体系为 **ADMIN**，接口前缀 `/api/v1/admin/*`。

## 功能

- 认证与 Cookie 会话；忘记 / 重置密码；三方登录（JustAuth，可配置）
- 用户中心：头像、资料、密码、联系方式、OAuth 绑定
- 动态路由与菜单（`/sys/resources/current`）
- IAM：账号、角色、部门、用户组、岗位、资源、授权
- 系统：字典、配置、Banner、文件、弱口令、审计、代码生成、OAuth 配置
- 消息：通知、公告、反馈
- 在线会话管理、Dashboard

## 技术栈

Vue 3 · Vite · TypeScript · Naive UI / Pro Naive UI · Pinia · Vue Router · axios · UnoCSS

## 开发

```bash
pnpm install
pnpm dev
```

```env
VITE_PORT=5173
VITE_HOME_PATH="/dashboard"
VITE_ROUTE_LOAD_MODE="dynamic"
VITE_API_URL=
VITE_API_PROXY_TARGET=http://127.0.0.1:8000
```

`VITE_API_URL` 留空时，请求走同源 `/api`，由 Vite 代理到后端。

## 命令

```bash
pnpm dev
pnpm build
pnpm preview
pnpm lint
pnpm format
```

## Docker

```bash
pnpm build

docker build -t hei-boot-admin .
docker run -d \
  -e BACKEND_URL="http://host.docker.internal:8000" \
  -p 8081:81 \
  hei-boot-admin
```

或在仓库根目录：`docker compose --profile admin up -d`。

环境变量：`BACKEND_URL`、`CLIENT_MAX_BODY_SIZE`（默认 `10m`）。

## 目录

```text
src/
  api/         接口
  components/  组件
  layouts/     布局
  router/      路由与守卫
  stores/      状态
  views/       页面
  utils/       工具
nginx/         生产 nginx 模板
```

## 说明

- `VITE_ROUTE_LOAD_MODE=dynamic`：菜单与路由来自后端资源
- 页面与按钮权限来自资源树与 permission key
