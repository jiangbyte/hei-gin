/** Author: Charlie */

import type { App } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { setupRouterGuard } from './guard'
import { routes } from './routes.inner'

// VITE_BASE_URL 表示前端应用部署的基础路径，例如部署到 /admin/ 时需要配置为 /admin/。
const { VITE_BASE_URL = '/' } = import.meta.env

/**
 * 全局路由实例。
 *
 * 项目统一使用 HTML5 history 模式，不再支持 hash 路由；生产环境需要服务端配置
 * history fallback，把未知前端路径回退到 index.html。
 */
export const router = createRouter({
  history: createWebHistory(VITE_BASE_URL),
  routes,
})

/**
 * 安装 Vue Router。
 *
 * 先注册全局路由守卫，再挂载 router，最后等待初始路由解析完成，避免应用启动时出现短暂空白状态。
 */
export async function installRouter(app: App) {
  setupRouterGuard(router)
  app.use(router)
  await router.isReady()
}
