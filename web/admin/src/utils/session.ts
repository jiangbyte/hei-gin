/** Author: Charlie */

/**
 * Web Admin 会话 token：与 Cookie 双通道。
 * - Cookie 开：服务端 HttpOnly + 本地 token 均可；后端 Cookie 优先。
 * - Cookie 关：仅靠本地 token，以 Authorization 头发送（非 Bearer）。
 * 本地 token 存在 XSS 可读风险，与 uni-app 一致。
 */
export const tokenKey = 'token'

export function getToken(): string {
  const fromLocal = localStorage.getItem(tokenKey)
  if (fromLocal) {
    return fromLocal
  }
  return sessionStorage.getItem(tokenKey) || ''
}

/** rememberMe=true 用 localStorage，否则 sessionStorage（关页即丢）。 */
export function setToken(token: string, rememberMe = true) {
  const value = token?.trim()
  if (!value) {
    clearToken()
    return
  }
  if (rememberMe) {
    localStorage.setItem(tokenKey, value)
    sessionStorage.removeItem(tokenKey)
  } else {
    sessionStorage.setItem(tokenKey, value)
    localStorage.removeItem(tokenKey)
  }
}

export function clearToken() {
  localStorage.removeItem(tokenKey)
  sessionStorage.removeItem(tokenKey)
}
