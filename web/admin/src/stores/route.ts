/** Author: Charlie */

import { defineStore } from 'pinia'
import { router } from '@/router'
import { resourceApi } from '@/api'
import { staticRoutes } from '@/router/routes.static'
import {
  createFullscreenRoutes,
  createMenus,
  createRoutes,
  generateCacheRoutes,
  getActiveMenuPath,
  groupResourcesByModule,
  normalizeRowRoutes,
} from './route/helper'

/**
 * 路由 store 状态。
 */
interface RouteState {
  isInitAuthRoute: boolean

  // 侧边栏菜单数据，由当前活跃模块的资源列表转换而来。
  menus: AppRoute.MenuOption[]

  // 原始资源列表，字段与后端 SysResource 保持一致。
  rowRoutes: AppRoute.RowRoute[]

  // 前端按 module_id 分组后得到的模块列表，供 ModuleSwitch 使用。
  resourceModules: AppRoute.ResourceModule[]

  // 当前侧边栏展示的资源模块 ID。
  activeModuleId: string | null

  // 当前需要高亮的菜单路径。
  currentMenuPath: string | null

  // keep-alive include 列表。
  cacheRoutes: string[]

  // 独立全屏路由名。
  fullscreenRouteNames: string[]
}

type RouteGetters = Record<string, never>

interface RouteActions {
  resetRouteStore: () => void
  resetRoutes: () => void
  setCurrentMenuPath: (path: string) => void
  setActiveModule: (moduleId: string | null) => void
  syncActiveModuleByPath: (path: string) => void
  initRouteInfo: () => Promise<AppRoute.RowRoute[]>
  initAuthRoute: () => Promise<void>
}

export const useRouteStore = defineStore<'route-store', RouteState, RouteGetters, RouteActions>(
  'route-store',
  {
    state: (): RouteState => ({
      isInitAuthRoute: false,
      menus: [],
      rowRoutes: [],
      resourceModules: [],
      activeModuleId: null,
      currentMenuPath: null,
      cacheRoutes: [],
      fullscreenRouteNames: [],
    }),
    actions: {
      resetRouteStore() {
        this.resetRoutes()
        this.$reset()
      },

      resetRoutes() {
        if (router.hasRoute('appRoot')) {
          router.removeRoute('appRoot')
        }
        this.fullscreenRouteNames.forEach((routeName) => {
          if (router.hasRoute(routeName)) {
            router.removeRoute(routeName)
          }
        })
        this.fullscreenRouteNames = []
      },

      setCurrentMenuPath(path: string) {
        this.currentMenuPath = getActiveMenuPath(this.rowRoutes, path)
        this.syncActiveModuleByPath(path)
      },

      setActiveModule(moduleId: string | null) {
        const resourceModule = this.resourceModules.find((item) => item.id === moduleId)
        this.activeModuleId = resourceModule?.id ?? null
        this.menus = resourceModule ? createMenus(resourceModule.resources) : []
      },

      async initRouteInfo() {
        if (import.meta.env.VITE_ROUTE_LOAD_MODE === 'dynamic') {
          return fetchUserRoutes()
        }
        return staticRoutes
      },

      async initAuthRoute() {
        this.isInitAuthRoute = false

        const rowRoutes = await this.initRouteInfo()
        this.rowRoutes = rowRoutes

        // 前端按 module_id 分组，供 ModuleSwitch 模块切换使用
        this.resourceModules = groupResourcesByModule(rowRoutes)

        if (router.hasRoute('appRoot')) {
          router.removeRoute('appRoot')
        }

        router.addRoute(createRoutes(rowRoutes))
        const fullscreenRoutes = createFullscreenRoutes(rowRoutes)
        fullscreenRoutes.forEach((route) => {
          router.addRoute(route)
        })
        this.fullscreenRouteNames = fullscreenRoutes
          .map((route) => route.name)
          .filter((name): name is string => typeof name === 'string')
        this.setActiveModule(this.resourceModules[0]?.id ?? null)
        this.cacheRoutes = generateCacheRoutes(rowRoutes)
        this.isInitAuthRoute = true
      },

      syncActiveModuleByPath(path: string) {
        const activePath = this.currentMenuPath ?? path
        const resource =
          this.rowRoutes.find((item) => item.path === path) ??
          this.rowRoutes.find((item) => item.path === activePath)

        if (
          resource?.module_id &&
          resource.module_id !== this.activeModuleId &&
          this.resourceModules.some((item) => item.id === resource.module_id)
        ) {
          this.setActiveModule(resource.module_id)
        }
      },
    },
  },
)

export function getRouteTitle(route: {
  name?: string | symbol | null
  path?: string
  meta: {
    name?: string
  }
}) {
  return route.meta.name ?? String(route.name ?? route.path)
}

async function fetchUserRoutes(): Promise<AppRoute.RowRoute[]> {
  const response = await resourceApi.current()
  return normalizeRowRoutes(response.data ?? [])
}
