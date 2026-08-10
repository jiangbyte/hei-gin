/** Author: Charlie */

import type { AxiosError, AxiosResponse } from 'axios'

const loginPath = '/auth/login'
let isHandlingUnauthorized = false

/**
 * 后端统一响应结构。
 *
 * 约定 code="200" 表示业务成功，data 才是业务真正需要的数据。
 * JSON 标量一律为字符串。
 */
interface ApiResponse<T = unknown> {
  code: string
  message?: string
  data: T
}

const httpStatusMessageMap: Record<number, string> = {
  400: '请求参数错误',
  401: '登录已过期，请重新登录',
  403: '无权访问',
  404: '资源不存在',
  422: '校验失败',
  500: '服务器错误',
  502: '网关错误',
  503: '服务不可用',
  504: '网关超时',
}

/**
 * 业务响应错误。
 *
 * HTTP 请求本身成功，但后端业务 code 非 "200" 时抛出该错误。
 */
export class ApiResponseError<T = unknown> extends Error {
  readonly apiCode: string
  readonly apiData: T
  readonly rawData: ApiResponse<T>

  constructor(response: ApiResponse<T>) {
    super(response.message || `请求失败，错误码 ${response.code}`)
    this.name = 'ApiResponseError'
    this.apiCode = response.code
    this.apiData = response.data
    this.rawData = response
  }
}

export function unwrapResponseData(response: AxiosResponse) {
  if (isApiResponse(response.data)) {
    if (response.data.code !== '200') {
      throw new ApiResponseError(response.data)
    }
    return response.data.data
  }
  return response.data
}

export function handleHttpError(error: AxiosError) {
  if (isUnauthorizedError(error) && !error.config?.public) {
    handleUnauthorizedError(error)
    return Promise.reject(error)
  }

  showErrorMessage(error)
  return Promise.reject(error)
}

function isApiResponse(data: unknown): data is ApiResponse {
  return isRecord(data) && typeof data.code === 'string'
}

function isRecord(data: unknown): data is Record<string, unknown> {
  return typeof data === 'object' && data !== null
}

function isUnauthorizedError(error: AxiosError) {
  return error.response?.status === 401 || getApiCode(error) === '401'
}

function getApiCode(error: AxiosError) {
  const apiCode = (error as { apiCode?: string }).apiCode
  if (typeof apiCode === 'string') {
    return apiCode
  }

  const responseData = error.response?.data
  if (isRecord(responseData) && typeof responseData.code === 'string') {
    return responseData.code
  }

  const rawData = (error.response as { rawData?: unknown } | undefined)?.rawData
  if (isRecord(rawData) && typeof rawData.code === 'string') {
    return rawData.code
  }

  return undefined
}

function handleUnauthorizedError(error: AxiosError) {
  if (isHandlingUnauthorized) {
    return
  }

  isHandlingUnauthorized = true
  showUnauthorizedMessage(error)

  void redirectToLogin().finally(() => {
    window.setTimeout(() => {
      isHandlingUnauthorized = false
    }, 1000)
  })
}

function showUnauthorizedMessage(error: AxiosError) {
  const message = getErrorMessage(error)
  if (message) {
    window.$message?.error(message)
  }
}

async function redirectToLogin() {
  const [{ useAuthStore }, { router }] = await Promise.all([import('@/stores'), import('@/router')])
  const authStore = useAuthStore()
  const currentRoute = router.currentRoute.value
  const redirect = currentRoute.fullPath

  authStore.resetSession()

  if (currentRoute.path.startsWith('/auth')) {
    if (currentRoute.path !== loginPath) {
      await router.replace(loginPath)
    }
    return
  }

  await router.replace(
    redirect
      ? {
          path: loginPath,
          query: { redirect },
        }
      : loginPath,
  )
}

function showErrorMessage(error: AxiosError) {
  if (error.config?.skipErrorMessage) {
    return
  }

  const message = getErrorMessage(error)
  if (message) {
    window.$message?.error(message)
  }
}

function getErrorMessage(error: AxiosError) {
  const customErrorMessage = error.config?.customErrorMessage
  if (customErrorMessage) {
    return customErrorMessage
  }

  const responseMessage = getResponseMessage(error.response?.data)
  if (responseMessage) {
    return responseMessage
  }

  const status = error.response?.status
  if (status) {
    return httpStatusMessageMap[status] ?? `请求失败(${status})`
  }

  return 'Network error. Please try again later.'
}

function getResponseMessage(data: unknown) {
  if (isRecord(data) && typeof data.message === 'string') {
    return data.message
  }
  return undefined
}
