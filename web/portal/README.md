# HEI Portal

React 门户。账号体系为 **PORTAL**，接口前缀 `/api/v1/portal/*`。

## 功能

- 全页认证：登录 / 注册 / 找回与重置密码（`/auth/*`）；支持三方登录入口
- Cookie 会话（可关，仅 Header）
- 首页、公告、反馈
- 个人主页、账号中心（资料、密码、邮箱、手机、消息、OAuth 绑定）

## 技术栈

React 19 · Vite · TypeScript · Ant Design 6 · React Router · Zustand · axios · UnoCSS

主题配置：`src/theme/tokens.ts`。

## 开发

```bash
pnpm install
pnpm dev
```

```env
VITE_APP_TITLE=HEI
VITE_PORT=5174
VITE_HOME_PATH=/
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

docker build -t hei-boot-portal .
docker run -d \
  -e BACKEND_URL="http://host.docker.internal:8000" \
  -p 8082:80 \
  hei-boot-portal
```

或在仓库根目录：`docker compose --profile portal up -d`。

## 目录

```text
src/
  api/           接口
  components/    组件（含登录/注册弹窗）
  layouts/       布局
  pages/         页面
  router/        路由与守卫
  stores/        状态
  theme/         主题 token
  utils/         工具
```
