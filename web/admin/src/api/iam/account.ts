/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const accountPrefix = `${API_PREFIX}/sys/accounts`

export function page(params: any) {
  return http.get<any>(`${accountPrefix}/page`, {
    params,
  })
}

export function detail(params: any) {
  return http.get<any>(`${accountPrefix}/detail`, {
    params,
  })
}

export function create(data: any) {
  return http.post<any>(`${accountPrefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${accountPrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${accountPrefix}/delete`, data)
}

export function ownRoles(accountId: string) {
  return http.get<any>(`${accountPrefix}/own-role`, { params: { id: accountId } })
}

export function grantRoles(data: any) {
  return http.post<any>(`${accountPrefix}/grant-role`, data)
}

export function ownGroups(accountId: string) {
  return http.get<any>(`${accountPrefix}/own-group`, { params: { id: accountId } })
}

export function grantGroups(data: any) {
  return http.post<any>(`${accountPrefix}/grant-group`, data)
}

export function ownDepts(accountId: string) {
  return http.get<any>(`${accountPrefix}/own-dept`, { params: { id: accountId } })
}

export function grantDepts(data: any) {
  return http.post<any>(`${accountPrefix}/grant-dept`, data)
}

export function ownResources(accountId: string) {
  return http.get<any>(`${accountPrefix}/own-resource`, { params: { id: accountId } })
}

export function grantResources(data: any) {
  return http.post<any>(`${accountPrefix}/grant-resource`, data)
}

export function ownClientResources(accountId: string) {
  return http.get<any>(`${accountPrefix}/own-client-resource`, { params: { id: accountId } })
}

export function grantClientResources(data: any) {
  return http.post<any>(`${accountPrefix}/grant-client-resource`, data)
}
