/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/sys/weak-password`

export function page(params: any) {
  return http.get<any>(`${prefix}/page`, { params })
}

export function list(params?: { keyword?: string }) {
  return http.get<any[]>(`${prefix}/list`, { params })
}

export function detail(params: { id: string }) {
  return http.get<any>(`${prefix}/detail`, { params })
}

export function create(data: { password: string }) {
  return http.post<any>(`${prefix}/create`, data)
}

export function update(data: { id: string; password: string }) {
  return http.post<any>(`${prefix}/update`, data)
}

export function remove(data: { ids: string[] }) {
  return http.post<any>(`${prefix}/delete`, data)
}
