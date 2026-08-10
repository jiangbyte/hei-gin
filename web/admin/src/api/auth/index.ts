/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const authPrefix = `${API_PREFIX}`

export function login(data: any) {
  return http.post<any>(`${authPrefix}/login`, data, {
    public: true,
  })
}

export function authOptions() {
  return http.get<any>(`${authPrefix}/public/auth-options`, {
    public: true,
  })
}

export function sendLoginCode(data: any) {
  return http.post<any>(`${authPrefix}/send-login-code`, data, {
    public: true,
  })
}

export function sendPasswordChangeCode() {
  return http.post<any>(`${authPrefix}/user-center/password/send-code`)
}

export function captcha() {
  return http.get<any>(`${authPrefix}/captcha`, {
    public: true,
  })
}

export function passwordKey() {
  return http.get<any>(`${authPrefix}/password-key`, {
    public: true,
  })
}

export function forgotPassword(data: any) {
  return http.post<any>(`${authPrefix}/forgot-password`, data, {
    public: true,
  })
}

export function resetPassword(data: any) {
  return http.post<any>(`${authPrefix}/reset-password`, data, {
    public: true,
  })
}

export function logout() {
  return http.post<any>(`${authPrefix}/logout`)
}

/** 注销当前登录账号（软注销 + 清理关联数据）。 */
export function cancelAccount(data?: { cancel_reason?: string | null }) {
  return http.post<any>(`${authPrefix}/cancel`, data ?? {})
}

export function me() {
  return http.get<any>(`${authPrefix}/me`)
}

export function updateUserCenterProfile(data: any) {
  return http.post<any>(`${authPrefix}/user-center/profile/update`, data)
}

export function uploadUserCenterAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return http.post<any>(`${authPrefix}/user-center/avatar/upload`, formData)
}

export function updateUserCenterPassword(data: any) {
  return http.post<any>(`${authPrefix}/user-center/password/update`, data)
}

export function updateUserCenterPhone(data: any) {
  return http.post<any>(`${authPrefix}/user-center/phone/update`, data)
}

export function updateUserCenterEmail(data: any) {
  return http.post<any>(`${authPrefix}/user-center/email/update`, data)
}
