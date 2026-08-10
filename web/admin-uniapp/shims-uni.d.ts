/// <reference types='@dcloudio/types' />
import 'vue'

declare module '@vue/runtime-core' {
  type Hooks = App.AppInstance & Page.PageInstance

  interface ComponentCustomOptions {
    onLaunch?: () => void
    onShow?: () => void
    onHide?: () => void
  }
}
