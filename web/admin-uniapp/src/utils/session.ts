/** Author: Charlie */

import { removeStorage } from './storage'

/**
 * 原生/uniapp 无法像 Web Admin/Portal 那样使用 HttpOnly Cookie。
 * 会话为 local storage 中的原始 opaque token，以 Authorization 发送（非 Bearer）。
 * H5 构建下该 token 存在 XSS 风险。
 */
export const tokenKey = 'token'
export const userInfoKey = 'user_info'

const clearListeners = new Set<() => void>()

export function getToken() {
  const token = uni.getStorageSync(tokenKey)
  return token ? String(token) : ''
}

export function setToken(token: string) {
  uni.setStorageSync(tokenKey, token)
}

export function clearSessionStorage() {
  removeStorage(tokenKey)
  removeStorage(userInfoKey)
  clearListeners.forEach((listener) => listener())
}

export function onSessionCleared(listener: () => void) {
  clearListeners.add(listener)
  return () => clearListeners.delete(listener)
}
