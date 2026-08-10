/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const deptPrefix = `${API_PREFIX}/sys/depts`

export function page(params: any) {
  return http.get<any>(`${deptPrefix}/page`, { params })
}

export function tree(params?: any) {
  return http.get<any>(`${deptPrefix}/tree`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${deptPrefix}/detail`, { params })
}

export function create(data: any) {
  return http.post<any>(`${deptPrefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${deptPrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${deptPrefix}/delete`, data)
}
