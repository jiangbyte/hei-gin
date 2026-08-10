/** Author: Charlie */

import axios, { type CreateAxiosDefaults } from 'axios'
import { handleHttpError, unwrapResponseData } from './handle'
import { setupRequestInterceptor } from './request-interceptors'
import { setupResponseInterceptors } from './response-interceptors'

export { ApiResponseError } from './handle'

/**
 * 创建项目 HTTP 客户端（Cookie + Authorization Header 双通道会话 + 统一解包/错误处理）。
 */
export function createHttp(config?: CreateAxiosDefaults) {
  const http = axios.create(config)

  setupRequestInterceptor(http)

  setupResponseInterceptors(http, {
    unwrapResponseData,
    handleError: handleHttpError,
  })

  return http
}
