/**
 * 由 HEI 代码生成器生成。
 * Author: Charlie
 * 生成时间：2026-08-08 21:09:52
 */
import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/biz/cg-test-activity`

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
