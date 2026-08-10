/** Author: Charlie */

export {}

declare module 'vue' {
  type Hooks = App.AppInstance & Page.PageInstance
  interface ComponentCustomOptions {
    onLaunch?: () => void
    onShow?: () => void
    onHide?: () => void
  }
}
