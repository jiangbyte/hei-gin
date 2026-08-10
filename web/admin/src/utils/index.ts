/** Author: Charlie */

// 工具模块统一导出入口，业务代码优先从 '@/utils' 引入通用能力。
export * from './icon'
export * from './axios'
export * from './color'
export * from './normalize'
export * from './permission'
export * from './file'
export * from './time'
export * from './wire'
export * from './rsaEncrypt'

import { createHttp } from './axios'

// 默认 HTTP 客户端实例。VITE_API_URL 为空时使用同源请求，便于本地代理或同域部署。
export const http = createHttp({
  baseURL: import.meta.env.VITE_API_URL || '',
  withCredentials: true,
})

export * from './session'
export * from './dict'
