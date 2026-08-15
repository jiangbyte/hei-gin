/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/sys/notices`

export function myPage(params?: { current?: number; size?: number; kind?: string }) {
  return http.get<any>(`${prefix}/my-page`, { params })
}

export function myDetail(id: string) {
  return http.get<any>(`${prefix}/my-detail`, { params: { id } })
}

export function unreadCount() {
  return http.get<any>(`${prefix}/unread-count`)
}

export function read(data: { ids: string[] }) {
  return http.post<any>(`${prefix}/read`, data)
}

export function readAll() {
  return http.post<any>(`${prefix}/read-all`)
}
