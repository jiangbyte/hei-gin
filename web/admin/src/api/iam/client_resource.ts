/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const clientResourcePrefix = `${API_PREFIX}/sys/client-resources`

export function tree(params?: any) {
  return http.get<any>(`${clientResourcePrefix}/tree`, { params })
}

export function page(params: any) {
  return http.get<any>(`${clientResourcePrefix}/page`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${clientResourcePrefix}/detail`, { params })
}

export function create(data: any) {
  return http.post<any>(`${clientResourcePrefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${clientResourcePrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${clientResourcePrefix}/delete`, data)
}
