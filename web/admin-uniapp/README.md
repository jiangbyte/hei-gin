# HEI Admin uni-app

管理端 uni-app（H5 / 小程序）。账号体系为 **ADMIN**，接口前缀 `/api/v1/admin/*`。会话使用本地 Authorization token。

## 功能

- 登录（账号 / 邮箱 / 手机号）
- Dashboard、工作台、用户中心
- IAM：账号、角色、部门、用户组、岗位、资源
- 系统：字典、Banner、文件
- 在线会话
- 通用资源页（列表、详情、表单、筛选、分页）

## 技术栈

uni-app 3 · Vue 3 · Vite · TypeScript · Pinia · uview-pro · UnoCSS

## 开发

```bash
pnpm install
pnpm dev:h5
```

```env
VITE_APP_TITLE="HEI Admin"
VITE_API_URL="http://127.0.0.1:8000"
VITE_PORT=5174
```

```bash
pnpm dev:mp-weixin
pnpm build:mp-weixin
```

其它平台命令见 `package.json`。

## 命令

```bash
pnpm dev:h5
pnpm build:h5
pnpm type-check
pnpm lint
pnpm format
```

## 构建

```bash
pnpm build:h5
```

生产环境 `VITE_API_URL` 为空时，请求走同源 `/api/`，由网关反代后端。

## 目录

```text
src/
  api/          接口
  components/   组件
  config/       资源页配置
  layouts/      布局
  pages/        页面
  stores/       状态
  utils/        工具
  pages.json
  manifest.json
```

## 说明

- 资源与权限来自 `/api/v1/admin/sys/resources/current`
- 发布前配置 `manifest.json` 中各端 appid，并在小程序后台配置合法域名
