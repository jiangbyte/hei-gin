/** Author: Charlie */

import { create } from 'zustand'

export type ColorMode = 'light' | 'dark' | 'auto'
export type ResolvedTheme = 'light' | 'dark'

const COLOR_MODE_KEY = 'portalColorMode'

function getStoredColorMode(): ColorMode {
  const value = localStorage.getItem(COLOR_MODE_KEY)
  return value === 'light' || value === 'dark' || value === 'auto' ? value : 'auto'
}

function getSystemDark() {
  return window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
}

function resolveTheme(colorMode: ColorMode, systemDark: boolean): ResolvedTheme {
  if (colorMode === 'auto') {
    return systemDark ? 'dark' : 'light'
  }
  return colorMode
}

type AppState = {
  colorMode: ColorMode
  systemDark: boolean
  resolvedTheme: ResolvedTheme
  setColorMode: (mode: ColorMode) => void
  setSystemDark: (dark: boolean) => void
}

const initialColorMode = getStoredColorMode()
const initialSystemDark = getSystemDark()

export const useAppStore = create<AppState>((set, get) => ({
  colorMode: initialColorMode,
  systemDark: initialSystemDark,
  resolvedTheme: resolveTheme(initialColorMode, initialSystemDark),

  setColorMode: (mode) => {
    localStorage.setItem(COLOR_MODE_KEY, mode)
    set({
      colorMode: mode,
      resolvedTheme: resolveTheme(mode, get().systemDark),
    })
  },

  setSystemDark: (dark) => {
    set({
      systemDark: dark,
      resolvedTheme: resolveTheme(get().colorMode, dark),
    })
  },
}))
