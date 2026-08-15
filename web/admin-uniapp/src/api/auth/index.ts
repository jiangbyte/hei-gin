/** Author: Charlie */

import {API_PREFIX} from '@/constants/api'
import { http } from '@/utils/request'

const prefix = API_PREFIX

export function captcha(params?: any) {
    return http.get<any>(`${prefix}/captcha`, params, {attachSession: false})
}

export function passwordKey() {
    return http.get<any>(`${prefix}/password-key`, undefined, {attachSession: false})
}

export function login(data: any) {
    return http.post<any>(`${prefix}/login`, data, {attachSession: false})
}

export function logout() {
  return http.post<any>(`${prefix}/logout`)
}

export function me() {
  return http.get<any>(`${prefix}/me`)
}

export function forgotPassword(data: any) {
    return http.post<any>(`${prefix}/forgot-password`, data, {attachSession: false})
}

export function resetPassword(data: any) {
    return http.post<any>(`${prefix}/reset-password`, data, {attachSession: false})
}

export function updateProfile(data: any) {
  return http.post<any>(`${prefix}/profile/update`, data)
}

export function updatePassword(data: any) {
  return http.post<any>(`${prefix}/profile/password/update`, data)
}

export function updatePhone(data: any) {
  return http.post<any>(`${prefix}/profile/phone/update`, data)
}

export function updateEmail(data: any) {
  return http.post<any>(`${prefix}/profile/email/update`, data)
}

export function orgInfo() {
  return http.get<any>(`${prefix}/profile/org-info`)
}

export function uploadAvatar(filePath: string) {
    return http.upload<any>(`${prefix}/profile/avatar/upload`, filePath)
}
