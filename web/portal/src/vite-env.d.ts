/// <reference types="vite/client" />
/** Author: Charlie */

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_API_URL: string
  readonly VITE_PORT: string
  readonly VITE_HOME_PATH: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
