/** Author: Charlie */

import { createHttp } from './axios'

export * from './dict'
export * from './session'
export * from './storage'
export * from './time'
export * from './wire'
export * from './validate'
export * from './file'
export * from './normalize'
export * from './color'
export { encryptPasswords } from './security'
export { ApiResponseError, createHttp } from './axios'

export const http = createHttp({
  baseURL: import.meta.env.VITE_API_URL || '',
  timeout: 15000,
  withCredentials: true,
})
