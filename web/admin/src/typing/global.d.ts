/** Author: Charlie
 *
 *
 * 挂载在 window 上的 Naive UI 全局 API。
 *
 * 这些实例通常由 NaiveProvider 在应用启动时注入，业务代码可以通过 window.$message、
 * window.$dialog 等方式在非组件上下文中触发反馈。
 */
interface Window {
  // 顶部加载条，用于路由切换和异步流程状态提示。
  $loadingBar: import('naive-ui').LoadingBarApi

  // 全局对话框 API。
  $dialog: import('naive-ui').DialogApi

  // 全局消息提示 API。
  $message: import('naive-ui').MessageApi

  // 全局通知 API。
  $notification: import('naive-ui').NotificationApi
}
