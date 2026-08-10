/** Author: Charlie */

import type { AxiosError, AxiosResponse } from 'axios'
import { message } from 'antd'
import type { ApiResponse } from '@/typing/api'

let isHandlingUnauthorized = false

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

  void redirectToLogin(error).finally(() => {
    window.setTimeout(() => {
      isHandlingUnauthorized = false
    }, 1000)
  })
}

async function redirectToLogin(error: AxiosError) {
  const { useAuthStore } = await import('@/stores/auth')
  const auth = useAuthStore.getState()
  const loggingOut = auth.loggingOut
  const hadSession = auth.isLogin()
  auth.resetSession()

  // 主动退出 / 本来就没登录：静默清理，不弹错（退出流程自己会开登录框）
  if (loggingOut || !hadSession) {
    return
  }

  const msg = getErrorMessage(error)
  if (msg) {
    message.error(msg)
  }

  const { pathname, search } = window.location
  const { useAuthModalStore } = await import('@/stores/authModal')
  if (pathname.startsWith('/auth/forgot-password')) {
    return
  }
  const redirect =
    pathname.startsWith('/auth') || pathname === '/' ? undefined : `${pathname}${search}`
  useAuthModalStore.getState().open('login', redirect)
}

function showErrorMessage(error: AxiosError) {
  if (error.config?.skipErrorMessage) {
    return
  }

  const msg = getErrorMessage(error)
  if (msg) {
    message.error(msg)
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

  return '网络异常，请稍后重试'
}

function getResponseMessage(data: unknown) {
  if (isRecord(data) && typeof data.message === 'string') {
    return data.message
  }
  return undefined
}
