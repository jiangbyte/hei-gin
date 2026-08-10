/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const clientModulePrefix = `${API_PREFIX}/sys/client-modules`

export function page(params: any) {
  return http.get<any>(`${clientModulePrefix}/page`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${clientModulePrefix}/detail`, { params })
}

export function selector(params?: any) {
  return http.get<any>(`${clientModulePrefix}/selector`, { params })
}

export function create(data: any) {
  return http.post<any>(`${clientModulePrefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${clientModulePrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${clientModulePrefix}/delete`, data)
}
