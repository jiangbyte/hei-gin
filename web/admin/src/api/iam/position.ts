/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const positionPrefix = `${API_PREFIX}/sys/positions`

export function page(params: any) {
  return http.get<any>(`${positionPrefix}/page`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${positionPrefix}/detail`, { params })
}

export function create(data: any) {
  return http.post<any>(`${positionPrefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${positionPrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${positionPrefix}/delete`, data)
}
