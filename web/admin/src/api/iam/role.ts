/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const rolePrefix = `${API_PREFIX}/sys/roles`

export function page(params: any) {
  return http.get<any>(`${rolePrefix}/page`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${rolePrefix}/detail`, { params })
}

export function create(data: any) {
  return http.post<any>(`${rolePrefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${rolePrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${rolePrefix}/delete`, data)
}

export function ownResources(roleId: string, accountType: string) {
  return http.get<any>(`${rolePrefix}/own-resource`, {
    params: { id: roleId, account_type: accountType },
  })
}

export function grantResources(data: any) {
  return http.post<any>(`${rolePrefix}/grant-resource`, data)
}

export function ownClientResources(roleId: string, accountType: string) {
  return http.get<any>(`${rolePrefix}/own-client-resource`, {
    params: { id: roleId, account_type: accountType },
  })
}

export function grantClientResources(data: any) {
  return http.post<any>(`${rolePrefix}/grant-client-resource`, data)
}

export function ownUsers(roleId: string) {
  return http.get<any>(`${rolePrefix}/own-user`, { params: { id: roleId } })
}

export function grantUsers(data: any) {
  return http.post<any>(`${rolePrefix}/grant-user`, data)
}
