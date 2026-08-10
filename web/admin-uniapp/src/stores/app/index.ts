/** Author: Charlie */

import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface AppState {
  collapsed: boolean
  isMobile: boolean
  isDark: boolean
}

export const useAppStore = defineStore('app', () => {
  const collapsed = ref<boolean>(false)
  const isMobile = ref<boolean>(false)
  const isDark = ref<boolean>(false)

  function setCollapsed(value: boolean) {
    collapsed.value = value
  }

  function setIsMobile(value: boolean) {
    isMobile.value = value
  }

  function setIsDark(value: boolean) {
    isDark.value = value
  }

  function toggleTheme() {
    isDark.value = !isDark.value
  }

  return {
    collapsed,
    isMobile,
    isDark,
    setCollapsed,
    setIsMobile,
    setIsDark,
    toggleTheme,
  }
})
