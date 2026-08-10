/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}`

export function captcha(format: 'svg' | 'png' = 'svg') {
  return http.get<any>(`${prefix}/captcha`, { params: { format }, public: true })
}

export function passwordKey() {
  return http.get<any>(`${prefix}/password-key`, { public: true })
}

export function login(data: any) {
  return http.post<any>(`${prefix}/login`, data, { public: true })
}

export function authOptions() {
  return http.get<any>(`${prefix}/public/auth-options`, { public: true })
}

export function sendLoginCode(data: any) {
  return http.post<any>(`${prefix}/send-login-code`, data, { public: true })
}

export function sendPasswordChangeCode() {
  return http.post<any>(`${prefix}/user-center/password/send-code`)
}

export function register(data: any) {
  return http.post<any>(`${prefix}/register`, data, { public: true })
}

export function forgotPassword(data: any) {
  return http.post<any>(`${prefix}/forgot-password`, data, { public: true })
}

export function resetPassword(data: any) {
  return http.post<any>(`${prefix}/reset-password`, data, { public: true })
}

/** probe：首访/会话探测用，401 不弹窗、不跳登录。 */
export function me(options?: { probe?: boolean }) {
  return http.get<any>(
    `${prefix}/me`,
    options?.probe ? { public: true, skipErrorMessage: true } : undefined,
  )
}

export function logout() {
  // 登出本身失败（含已无 cookie）不应再走全局 401 提示
  return http.post<any>(`${prefix}/logout`, undefined, {
    public: true,
    skipErrorMessage: true,
  })
}

/** 注销当前登录账号（软注销 + 清理关联数据）。 */
export function cancelAccount(data?: { cancel_reason?: string | null }) {
  return http.post<any>(`${prefix}/cancel`, data ?? {})
}

export function updateUserCenterProfile(data: any) {
  return http.post<any>(`${prefix}/user-center/profile/update`, data)
}

export function updateUserCenterPassword(data: any) {
  return http.post<any>(`${prefix}/user-center/password/update`, data)
}

export function updateUserCenterPhone(data: any) {
  return http.post<any>(`${prefix}/user-center/phone/update`, data)
}

export function updateUserCenterEmail(data: any) {
  return http.post<any>(`${prefix}/user-center/email/update`, data)
}

export function uploadUserCenterAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return http.post<any>(`${prefix}/user-center/avatar/upload`, formData)
}

export function getPublicSpace(accountId: string) {
  return http.get<any>(`${prefix}/spaces/detail`, {
    public: true,
    params: { account_id: accountId },
  })
}
