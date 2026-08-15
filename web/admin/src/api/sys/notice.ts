/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/sys/notices`

export function page(params: any) {
  return http.get<any>(`${prefix}/page`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${prefix}/detail`, { params })
}

export function create(data: any) {
  return http.post<any>(`${prefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${prefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${prefix}/delete`, data)
}

export function publish(data: { ids: string[] }) {
  return http.post<any>(`${prefix}/publish`, data)
}

export function revoke(data: { ids: string[] }) {
  return http.post<any>(`${prefix}/revoke`, data)
}

export function pin(data: { id: string; is_pinned: boolean; pinned_until?: string | null }) {
  return http.post<any>(`${prefix}/pin`, data)
}
