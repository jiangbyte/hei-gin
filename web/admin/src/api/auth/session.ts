/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/auth/sessions`

export function analysis() {
  return http.get<any>(`${prefix}/analysis`)
}

export function page(params?: any) {
  return http.get<any>(`${prefix}/page`, { params })
}

export function tokens(params: any) {
  return http.get<any[]>(`${prefix}/tokens`, { params })
}

export function exit(data: any) {
  return http.post<any>(`${prefix}/exit`, data)
}

export function tokenExit(data: any) {
  return http.post<any>(`${prefix}/token/exit`, data)
}
