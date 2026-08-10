/** Author: Charlie */

import type { AxiosInstance, AxiosResponse } from 'axios'
import { AxiosError, isAxiosError } from 'axios'

declare module 'axios' {
  interface AxiosResponse {
    rawData?: unknown
  }
}

export interface SetupResponseInterceptorsOptions {
  unwrapResponseData: (response: AxiosResponse) => unknown
  handleError: (error: AxiosError) => unknown
}

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

function unwrapResponseData(response: AxiosResponse, options: SetupResponseInterceptorsOptions) {
  try {
    response.rawData = response.data
    response.data = options.unwrapResponseData(response)
  } catch (e: unknown) {
    throw toAxiosResponseError(e, response)
  }
}

function handleError(error: unknown, options: SetupResponseInterceptorsOptions) {
  return options.handleError(toAxiosError(error))
}

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

function toAxiosError(error: unknown) {
  if (isAxiosError(error)) {
    return error
  }

  return AxiosError.from(toError(error))
}

function toError(error: unknown) {
  if (error instanceof Error) {
    return error
  }

  return new Error(String(error ?? 'Unknown Error'))
}

function getErrorCustomProps(error: unknown) {
  if (typeof error !== 'object' || error === null) {
    return undefined
  }

  const props = Object.fromEntries(Object.entries(error))
  return Object.keys(props).length > 0 ? props : undefined
}
