/** Author: Charlie */

import { defineStore } from 'pinia'
import { nextTick, ref } from 'vue'
import { useColorMode, useFullscreen, useMediaQuery, usePreferredDark } from '@vueuse/core'

// 全局页面级能力使用 documentElement 作为作用对象，例如全屏。
const docEle = ref(document.documentElement)

// 全屏状态由 @vueuse/core 维护，store 只暴露当前状态和切换动作。
const { isFullscreen, toggle } = useFullscreen(docEle)

// 主题模式会同步到本地存储；emitAuto 开启后可以保留 auto 模式，而不是立即解析成 light/dark。
const colorMode = useColorMode({ emitAuto: true })
const prefersDark = usePreferredDark()

// 布局在小屏幕下会切换为移动端交互，例如隐藏桌面侧边栏控制。
const isMobile = useMediaQuery('(max-width: 700px)')

/**
 * 应用全局状态。
 *
 * 这里集中保存页面刷新开关等运行态设置。
 * 业务数据不要放在这个 store 中，避免和页面模块状态混在一起。
 */
export const useAppStore = defineStore('app-store', {
  state: () => ({
    // 侧边栏是否折叠。
    collapsed: false,

    // 页面内容渲染开关。reloadPage 会短暂关闭再打开，用于触发当前路由视图重载。
    loadFlag: true,
  }),
  getters: {
    // 用户选择的主题模式，可能是 light、dark 或 auto。
    storeColorMode: () => colorMode.value,

    // Naive UI 只需要 light/dark；auto 在这里按当前结果映射给组件库。
    naiveTheme: () =>
      colorMode.value === 'dark' || (colorMode.value === 'auto' && prefersDark.value)
        ? 'dark'
        : 'light',

    // 当前是否处于全屏状态。
    fullScreen: () => isFullscreen.value,

    // 当前是否命中移动端断点。
    isMobile: () => isMobile.value,
  },
  actions: {
    /**
     * 设置颜色模式。
     *
     * light/dark 表示强制指定主题；auto 表示跟随系统主题。
     */
    setColorMode(mode: 'light' | 'dark' | 'auto') {
      colorMode.value = mode
    },

    /**
     * 切换侧边栏折叠状态。
     */
    toggleCollapse() {
      this.collapsed = !this.collapsed
    },

    /**
     * 切换浏览器全屏状态。
     */
    toggleFullScreen() {
      toggle()
    },

    /**
     * 重载当前页面内容区域。
     *
     * 通过短暂关闭 loadFlag 让视图组件卸载，再在下一轮渲染后恢复。
     * delay 为恢复渲染前的等待时间；传 0 可以立即恢复。
     */
    async reloadPage(delay = 600) {
      this.loadFlag = false
      await nextTick()
      if (delay) {
        window.setTimeout(() => {
          this.loadFlag = true
        }, delay)
      } else {
        this.loadFlag = true
      }
    },
  },
  persist: {
    // 应用偏好使用 localStorage 持久化，刷新页面后继续沿用用户设置。
    storage: localStorage,
  },
})
