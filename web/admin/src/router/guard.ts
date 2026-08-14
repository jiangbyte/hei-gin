/** Author: Charlie */

import type { RouteLocationNormalized, Router } from 'vue-router'
import { useAuthStore, useRouteStore, useTabStore } from '@/stores'
import { isAllowedUnderSecurityWall, resolveSecurityWallPath } from '@/stores/auth'
import { getRouteTitle } from '@/stores/route'
import { isDictLoaded, refreshDict, syncDictTree } from '@/utils/dict'

// 浏览器标题后缀，来自应用环境配置。
const appTitle = import.meta.env.VITE_APP_TITLE

// 首页路径只用于跳转和默认重定向，不用于覆盖资源表里的 path 字段。
const homePath = import.meta.env.VITE_HOME_PATH
const loginPath = '/auth/login'

/**
 * 注册全局路由守卫。
 *
 * 这里集中处理外链跳转、授权路由初始化、标签页同步、菜单高亮和页面标题等路由副作用。
 */
export function setupRouterGuard(router: Router) {
  router.beforeEach(async (to) => {
    const authStore = useAuthStore()
    const routeStore = useRouteStore()
    syncDictTree()

    // 资源配置了 href 时视为外链。打开新窗口后阻止当前路由继续跳转。
    if (to.meta.href) {
      window.open(to.meta.href)
      return false
    }

    window.$loadingBar?.start()

    // 当 userInfo 过期或缺失时，从 HttpOnly cookie 恢复登录状态。
    if (!authStore.sessionChecked) {
      await authStore.ensureSession()
    }

    const isLogin = authStore.isLogin

    // 根路径只做入口转发：已登录进入首页，未登录进入登录页。
    if (to.name === 'root') {
      return { path: isLogin ? homePath : loginPath, replace: true }
    }

    // 已登录用户访问认证页时，直接回到 redirect 或首页（OAuth 回调需放行，含绑定场景）。
    if (isAuthRoute(to) && isLogin && to.name !== 'auth-oauth-callback') {
      const wall = resolveSecurityWallPath(authStore.userInfo)
      return { path: wall ?? getRedirectPath(to) ?? homePath, replace: true }
    }

    // 认证页和显式 404 是公开路由，不触发登录态和授权路由初始化。
    if (isPublicRoute(to)) {
      return
    }

    // 未登录访问后台页面时跳转到登录页，并保留原目标地址。
    if (!isLogin) {
      return {
        path: loginPath,
        query: {
          redirect: to.fullPath,
        },
      }
    }

    const tab = typeof to.query.tab === 'string' ? to.query.tab : undefined
    if (!isAllowedUnderSecurityWall(to.path, tab, authStore.userInfo)) {
      const wall = resolveSecurityWallPath(authStore.userInfo)
      if (wall) {
        return wall
      }
    }

    // 登录后首次进入系统时注册授权路由。注册完成后，如果当前命中的是 404 兜底路由，需要重新匹配一次目标路径。
    if (!routeStore.isInitAuthRoute) {
      try {
        await Promise.all([routeStore.initAuthRoute(), refreshDict()])
        if (isFallbackRoute(to)) {
          return {
            path: to.fullPath,
            replace: true,
            query: to.query,
            hash: to.hash,
          }
        }
      } catch {
        authStore.resetSession()
        return {
          path: loginPath,
          replace: true,
          query: {
            redirect: to.fullPath,
          },
        }
      }
    }

    if (!isDictLoaded()) {
      await refreshDict()
    }
  })

  // 路由解析完成前同步 UI 状态：侧边菜单高亮、标签页新增、当前标签记录。
  router.beforeResolve((to) => {
    const routeStore = useRouteStore()
    const tabStore = useTabStore()

    routeStore.setCurrentMenuPath(to.path)
    tabStore.addTab(to)
    tabStore.setCurrentTab(to.fullPath)
  })

  // 路由完成后更新浏览器标题并结束顶部加载条。
  router.afterEach((to) => {
    document.title = to.meta.name ? `${getRouteTitle(to)} - ${appTitle}` : appTitle
    window.$loadingBar?.finish()
  })

  // 路由解析异常时给出加载条错误状态，避免 loading 一直停留。
  router.onError(() => {
    window.$loadingBar?.error()
  })
}

// 判断当前路由是否命中了 pathMatch 兜底路由。
function isFallbackRoute(route: RouteLocationNormalized) {
  return route.matched.some((item) => item.path.includes(':pathMatch'))
}

function isAuthRoute(route: RouteLocationNormalized) {
  return route.path.startsWith('/auth')
}

function isPublicRoute(route: RouteLocationNormalized) {
  return isAuthRoute(route) || route.name === 'not-found'
}

function getRedirectPath(route: RouteLocationNormalized) {
  const redirect = route.query.redirect
  if (typeof redirect !== 'string' || redirect.startsWith('/auth')) {
    return undefined
  }
  return redirect
}
