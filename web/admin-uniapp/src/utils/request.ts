/** Author: Charlie */

import {stringifyScalars} from '@/utils/wire'
import { clearSessionStorage, getToken } from './session'

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: Record<string, any>
    attachSession?: boolean
  skipErrorMessage?: boolean
  header?: Record<string, string>
}

export class ApiResponseError extends Error {
  code?: string
  statusCode?: number
  raw?: unknown

  constructor(
    message: string,
    code?: string,
    statusCode?: number,
    raw?: unknown
  ) {
    super(message)
    this.name = 'ApiResponseError'
    this.code = code
    this.statusCode = statusCode
    this.raw = raw
  }
}

function isSuccessCode(code: unknown): boolean {
    return typeof code === 'string' && code === '200'
}

const baseURL = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '')

export function request<T = any>(url: string, options: RequestOptions = {}) {
  return new Promise<T>((resolve, reject) => {
    const token = getToken()
    const header: Record<string, string> = {
      ...(options.header ?? {}),
    }

      if (options.attachSession !== false && token) {
      header.Authorization = token
    }

    uni.request({
      url: `${baseURL}${url}`,
      method: options.method ?? 'GET',
        data: stringifyScalars(cleanData(options.data ?? {})) as Record<string, any>,
      header,
      success(response) {
        const statusCode = response.statusCode
        if (statusCode === 401) {
          clearSessionStorage()
          reject(
            new ApiResponseError(
              '登录已过期',
              undefined,
              statusCode,
              response.data
            )
          )
          return
        }
        if (statusCode < 200 || statusCode >= 300) {
          const raw = response.data as any
          const bodyCode =
              raw && typeof raw === 'object' && typeof raw.code === 'string'
                  ? raw.code
                  : undefined
          const message =
            readMessage(response.data) || `请求失败(${statusCode})`
          showError(message, options.skipErrorMessage)
          reject(
            new ApiResponseError(message, bodyCode, statusCode, response.data)
          )
          return
        }

        const body = response.data as any
        if (body && typeof body === 'object' && 'code' in body) {
            if (!isSuccessCode(body.code)) {
            const message = body.message || '业务处理失败'
            showError(message, options.skipErrorMessage)
            reject(
                new ApiResponseError(
                    message,
                    typeof body.code === 'string' ? body.code : undefined,
                    statusCode,
                    body,
                ),
            )
            return
          }
          resolve(body.data as T)
          return
        }

        resolve(body as T)
      },
      fail(error) {
        const errMsg: string = error.errMsg || ''
        const message =
          errMsg === 'request:fail' ? '网络请求失败' : errMsg || '网络请求失败'
        showError(message, options.skipErrorMessage)
        reject(new ApiResponseError(message, undefined, undefined, error))
      },
    })
  })
}

export const http = {
  get<T = any>(
    url: string,
    data?: Record<string, any>,
    options?: RequestOptions
  ) {
    return request<T>(url, { ...(options ?? {}), method: 'GET', data })
  },
  post<T = any>(
    url: string,
    data?: Record<string, any>,
    options?: RequestOptions
  ) {
    return request<T>(url, { ...(options ?? {}), method: 'POST', data })
  },
  upload<T = any>(url: string, filePath: string) {
    return new Promise<T>((resolve, reject) => {
      const token = getToken()
      uni.uploadFile({
        url: `${baseURL}${url}`,
        filePath,
        name: 'file',
        header: { Authorization: token },
        success(res) {
          try {
            const data = JSON.parse(res.data)
              if (isSuccessCode(data.code)) {
              resolve(data.data as T)
            } else {
              reject(
                new ApiResponseError(
                  data.message || '上传失败',
                    typeof data.code === 'string' ? data.code : undefined,
                  res.statusCode
                )
              )
            }
          } catch {
            reject(new ApiResponseError('解析上传结果失败'))
          }
        },
        fail: reject,
      })
    })
  },
}

function cleanData(data: Record<string, any>) {
  return Object.fromEntries(
    Object.entries(data).filter(
      ([, value]) => value !== undefined && value !== ''
    )
  )
}

function readMessage(data: unknown) {
  if (!data || typeof data !== 'object') {
    return ''
  }
    return String((data as any).message ?? (data as any).detail ?? '')
}

function showError(message: string, skip?: boolean) {
  if (skip) {
    return
  }
  uni.showToast({
    title: message,
    icon: 'none',
  })
}
