/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const resourcePrefix = `${API_PREFIX}/sys/resources`

export function tree(params?: any) {
  return http.get<any>(`${resourcePrefix}/tree`, { params })
}

export function current() {
  return http.get<any>(`${resourcePrefix}/current`)
}

export function detail(params: any) {
  return http.get<any>(`${resourcePrefix}/detail`, { params })
}

export function create(data: any) {
  return http.post<any>(`${resourcePrefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${resourcePrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${resourcePrefix}/delete`, data)
}

export function permissionRegistry() {
  return http.get<any>(`${API_PREFIX}/permission-registry`)
}

export function buttonPage(params?: any) {
  return http.get<any>(`${API_PREFIX}/sys/resource-buttons/page`, { params })
}

export function buttonCreate(data: any) {
  return http.post<any>(`${API_PREFIX}/sys/resource-buttons/create`, data)
}

export function buttonUpdate(data: any) {
  return http.post<any>(`${API_PREFIX}/sys/resource-buttons/update`, data)
}

export function buttonRemove(data: any) {
  return http.post<any>(`${API_PREFIX}/sys/resource-buttons/delete`, data)
}
