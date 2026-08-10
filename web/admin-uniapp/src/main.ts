/** Author: Charlie */

import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import uViewPro from 'uview-pro'
import 'uno.css'
// #ifdef H5
import '@unocss/reset/tailwind.css'
// #endif
export function createApp() {
  const app = createSSRApp(App)
  const pinia = createPinia()
  app.use(pinia)
  app.use(uViewPro)
  return {
    app,
  }
}
