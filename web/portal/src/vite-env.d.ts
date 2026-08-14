/// <reference types="vite/client" />
/** Author: Charlie */

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_API_URL: string
  readonly VITE_PORT: string
  readonly VITE_HOME_PATH: string
  readonly VITE_COPYRIGHT_INFO?: string
  readonly VITE_ICP_FILING?: string
  readonly VITE_SUPPORT_EMAIL?: string
  readonly VITE_SUPPORT_PHONE?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
