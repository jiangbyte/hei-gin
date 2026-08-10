/** Author: Charlie */

import type { RouteLocationNormalizedLoaded } from 'vue-router'
import type { StoreDefinition } from 'pinia'
import { defineStore } from 'pinia'
import { router } from '@/router'

// 关闭当前标签页且没有其它可跳转标签时，回退到环境变量配置的首页。
const homePath = import.meta.env.VITE_HOME_PATH

interface TabState {
  // 固定标签页，对应资源字段 is_affix=true，始终展示在普通标签页前面。
  affixTabs: AppTab[]

  // 普通标签页，支持关闭、关闭左侧、关闭右侧、拖拽排序等交互。
  tabs: AppTab[]

  // 当前激活标签页的 fullPath，用于标签高亮和滚动定位。
  currentTabPath: string
}

/**
 * 标签页快照。
 *
 * 这里只保存标签栏渲染和跳转需要的最小路由信息，避免把完整 路由 对象放入持久化状态。
 */
export interface AppTab {
  // Vue Router 的路由名称；当前资源路由使用 module_id + code 作为 route.name。
  name: RouteLocationNormalizedLoaded['name']

  // 不含查询参数的路由路径；全局标签去重键。
  path: string

  // 含 query/hash 的完整路径，作为跳转目标；query 变化时原地更新，不新建标签。
  fullPath: string

  // 路由元信息快照，标题、图标、资源类型、固定标签等都从这里读取。
  meta: RouteLocationNormalizedLoaded['meta']
}

interface TabGetters {
  // 标签栏展示顺序：固定标签在前，普通标签在后。
  allTabs: (state: TabState) => AppTab[]
}

interface TabAction {
  // 根据当前路由添加标签页；非页面资源或已存在的标签不会重复添加。
  addTab: (route: RouteLocationNormalizedLoaded) => void

  // 关闭指定标签页；如果关闭的是当前标签，会自动切换到相邻标签或首页。
  closeTab: (fullPath: string) => Promise<void>

  // 只保留指定标签页，关闭其它普通标签页。
  closeOtherTabs: (fullPath: string) => void

  // 关闭指定标签左侧的普通标签页。
  closeLeftTabs: (fullPath: string) => void

  // 关闭指定标签右侧的普通标签页。
  closeRightTabs: (fullPath: string) => void

  // 关闭全部普通标签页，并跳转到首页。
  closeAllTabs: () => void

  // 清空全部标签页，通常用于退出登录或重置会话。
  clearAllTabs: () => void

  // 判断指定 path 是否已经存在于固定标签或普通标签中。
  hasExistTab: (path: string) => boolean

  // 设置当前激活标签页。
  setCurrentTab: (fullPath: string) => void

  // 获取普通标签页中的索引；固定标签不参与左右关闭和拖拽排序。
  getTabIndex: (fullPath: string) => number
}

/**
 * 标签页状态。
 *
 * 标签页属于会话级 UI 状态，使用 sessionStorage 持久化，刷新页面保留，关闭浏览器会话后清空。
 */
export const useTabStore = defineStore('tab-store', {
  state: (): TabState => ({
    affixTabs: [],
    tabs: [],
    currentTabPath: '',
  }),
  getters: {
    // 统一提供给标签栏、标签下拉菜单使用的完整标签列表。
    allTabs: (state: TabState) => [...state.affixTabs, ...state.tabs],
  },
  actions: {
    /**
     * 根据路由添加标签页。
     *
     * 只有 MENU/PAGE 资源会进入标签栏；CATALOG、BUTTON、ACTION、API_GROUP
     * 以及 404 等内部路由都不会生成标签页。
     *
     * 全局按 path 去重：同一页面仅一个标签；query/hash 变化只更新该标签的 fullPath。
     */
    addTab(route: RouteLocationNormalizedLoaded) {
      if (route.meta.layout === 'fullscreen' || !isTabResource(route.meta.resource_type)) {
        return
      }

      const existing = this.allTabs.find((item) => item.path === route.path)
      if (existing) {
        existing.name = route.name
        existing.fullPath = route.fullPath
        existing.meta = { ...route.meta }
        return
      }

      const tab: AppTab = {
        name: route.name,
        path: route.path,
        fullPath: route.fullPath,
        meta: { ...route.meta },
      }

      // 固定标签和普通标签分开保存，方便渲染顺序、关闭行为和拖拽边界保持稳定。
      if (route.meta.is_affix) {
        this.affixTabs.push(tab)
      } else {
        this.tabs.push(tab)
      }
    },

    /**
     * 关闭指定普通标签页。
     *
     * 如果关闭的是当前激活标签，优先跳转到右侧标签，其次左侧标签，再其次最后一个固定标签，
     * 最后回退到 VITE_HOME_PATH。
     */
    async closeTab(fullPath: string) {
      const index = this.getTabIndex(fullPath)
      const isCurrent = this.currentTabPath === fullPath

      this.tabs = this.tabs.filter((item) => item.fullPath !== fullPath)

      if (isCurrent) {
        const nextTab = this.tabs[index] ?? this.tabs[index - 1] ?? this.affixTabs.at(-1)
        await router.push(nextTab?.fullPath ?? homePath)
      }
    },

    /**
     * 关闭其它普通标签页。
     *
     * 固定标签不受影响；如果传入的是固定标签 fullPath，则普通标签会被全部关闭。
     */
    closeOtherTabs(fullPath: string) {
      this.tabs = this.tabs.filter((item) => item.fullPath === fullPath)
    },

    /**
     * 关闭当前普通标签左侧的标签页。
     *
     * 这里按普通标签数组的索引计算，固定标签天然处于另一组，不参与删除。
     */
    closeLeftTabs(fullPath: string) {
      const index = this.getTabIndex(fullPath)
      this.tabs = this.tabs.filter((_, i) => i >= index)
    },

    /**
     * 关闭当前普通标签右侧的标签页。
     *
     * 与 closeLeftTabs 对称，只保留索引小于等于当前标签的普通标签。
     */
    closeRightTabs(fullPath: string) {
      const index = this.getTabIndex(fullPath)
      this.tabs = this.tabs.filter((_, i) => i <= index)
    },

    /**
     * 关闭全部普通标签页并回到首页。
     *
     * 固定标签保留，避免首页或其它常驻入口从标签栏消失。
     */
    closeAllTabs() {
      this.tabs = []
      router.push(homePath)
    },

    /**
     * 清空固定标签和普通标签。
     *
     * 这是强制重置方法，不做路由跳转，调用方需要自行决定后续导航。
     */
    clearAllTabs() {
      this.affixTabs = []
      this.tabs = []
    },

    /**
     * 检查标签是否已存在（按 path）。
     */
    hasExistTab(path: string) {
      return this.allTabs.some((item) => item.path === path)
    },

    /**
     * 记录当前激活标签。
     *
     * 路由守卫在每次解析完成前写入，标签栏据此高亮当前项。
     */
    setCurrentTab(fullPath: string) {
      this.currentTabPath = fullPath
    },

    /**
     * 获取普通标签页索引。
     *
     * 固定标签不在 tabs 数组中，因此返回 -1 表示目标不是普通标签。
     */
    getTabIndex(fullPath: string) {
      return this.tabs.findIndex((item) => item.fullPath === fullPath)
    },
  },
  persist: {
    // 标签页属于当前浏览器会话状态，使用 sessionStorage 比 localStorage 更符合预期。
    storage: sessionStorage,
  },
}) as StoreDefinition<'tab-store', TabState, TabGetters, TabAction>

// 只有真实页面资源才进入标签栏；目录和权限节点不应该生成可访问标签。
function isTabResource(resourceType?: AppRoute.ResourceType) {
  return resourceType === 'MENU' || resourceType === 'PAGE'
}
