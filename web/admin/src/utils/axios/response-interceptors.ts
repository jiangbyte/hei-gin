/** Author: Charlie */

import type { AxiosInstance, AxiosResponse } from 'axios'
import { AxiosError, isAxiosError } from 'axios'

declare module 'axios' {
  interface AxiosResponse {
    // 保存被解包前的原始响应数据，方便调试或错误处理时追溯完整返回。
    rawData?: unknown
  }
}

/**
 * 响应拦截器配置。
 *
 * unwrapResponseData 负责把后端响应转换成业务数据；
 * handleError 负责统一处理请求错误和解包过程中抛出的错误。
 */
export interface SetupResponseInterceptorsOptions {
  unwrapResponseData: (response: AxiosResponse) => unknown
  handleError: (error: AxiosError) => unknown
}

/**
 * 注册响应拦截器。
 *
 * 成功响应会先解包 data，再返回原 AxiosResponse；失败响应会统一转换成 AxiosError。
 */
export function setupResponseInterceptors(
  http: AxiosInstance,
  options: SetupResponseInterceptorsOptions,
) {
  http.interceptors.response.use(
    (response) => {
      unwrapResponseData(response, options)
      return response
    },
    (error) => handleError(error, options),
  )
}

/**
 * 解包响应数据。
 *
 * 解包前先把原始 response.data 存入 rawData；如果解包过程抛出普通 Error，
 * 会转换成带 response 上下文的 AxiosError，方便上层统一识别。
 */
function unwrapResponseData(response: AxiosResponse, options: SetupResponseInterceptorsOptions) {
  try {
    response.rawData = response.data
    response.data = options.unwrapResponseData(response)
  } catch (e: unknown) {
    throw toAxiosResponseError(e, response)
  }
}

// 统一把未知错误转换为 AxiosError 后交给外部错误处理函数。
function handleError(error: unknown, options: SetupResponseInterceptorsOptions) {
  return options.handleError(toAxiosError(error))
}

// 把业务解包阶段的错误包装成带响应上下文的 AxiosError。
function toAxiosResponseError(error: unknown, response: AxiosResponse) {
  if (isAxiosError(error)) {
    return error
  }

  return AxiosError.from(
    toError(error),
    AxiosError.ERR_BAD_RESPONSE,
    response.config,
    response.request,
    response,
    getErrorCustomProps(error),
  )
}

// 把请求阶段抛出的任意错误标准化为 AxiosError。
function toAxiosError(error: unknown) {
  if (isAxiosError(error)) {
    return error
  }

  return AxiosError.from(toError(error))
}

// unknown 到 Error 的兜底转换，避免后续错误链路处理非 Error 值。
function toError(error: unknown) {
  if (error instanceof Error) {
    return error
  }

  return new Error(String(error ?? 'Unknown Error'))
}

// AxiosError.from 支持附加自定义属性；这里把对象型错误自身字段透传出去。
function getErrorCustomProps(error: unknown) {
  if (typeof error !== 'object' || error === null) {
    return undefined
  }

  const props = Object.fromEntries(Object.entries(error))
  return Object.keys(props).length > 0 ? props : undefined
}
