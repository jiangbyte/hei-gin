/** Author: Charlie */

import axios, { type CreateAxiosDefaults } from 'axios'
import { handleHttpError, unwrapResponseData } from './handle'
import { setupRequestInterceptor } from './request-interceptors'
import { setupResponseInterceptors } from './response-interceptors'

export { ApiResponseError } from './handle'

export function createHttp(config?: CreateAxiosDefaults) {
  const http = axios.create(config)

  setupRequestInterceptor(http)

  setupResponseInterceptors(http, {
    unwrapResponseData,
    handleError: handleHttpError,
  })

  return http
}
