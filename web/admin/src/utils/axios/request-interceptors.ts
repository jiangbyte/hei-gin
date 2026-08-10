/** Author: Charlie */

import type { AxiosInstance } from 'axios'
import { getToken } from '@/utils/session'
import { stringifyScalars } from '@/utils/wire'

declare module 'axios' {
  interface AxiosRequestConfig {
    /** 公开接口：401 时不跳转登录。会话：Cookie（withCredentials）和/或 Authorization 头。 */
    public?: boolean
    skipErrorMessage?: boolean
    customErrorMessage?: string
  }
}

/** 按 wire 约定将 JSON body/params 标量序列化为字符串；附带本地会话 token。 */
export function setupRequestInterceptor(http: AxiosInstance) {
  http.interceptors.request.use((config) => {
    if (config.data && !(config.data instanceof FormData)) {
      config.data = stringifyScalars(config.data)
    }
    if (config.params) {
      config.params = stringifyScalars(config.params) as typeof config.params
    }
    const token = getToken()
    if (token) {
      config.headers = config.headers ?? {}
      config.headers.Authorization = token
    }
    return config
  })
}
