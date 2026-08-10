/// <reference types="vite/client" />
/** Author: Charlie */

declare module '*.vue' {
  import type { DefineComponent } from 'vue'

  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, unknown>
  export default component
}

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_COPYRIGHT_INFO: string
  readonly VITE_API_URL?: string
  readonly VITE_BASE_URL?: string
  readonly VITE_HOME_PATH: string
  readonly VITE_ROUTE_LOAD_MODE?: 'static' | 'dynamic'
}
