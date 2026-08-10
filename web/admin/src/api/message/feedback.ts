/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/message/feedbacks`

/** 管理端反馈处理列表 */
export function page(params: any) {
  return http.get<any>(`${prefix}/page`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${prefix}/detail`, { params })
}

export function update(data: any) {
  return http.post<any>(`${prefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${prefix}/delete`, data)
}

/** 当前登录用户「我的反馈」 */
export function submit(data: {
  title: string
  content: string
  category: string
  contact?: string | null
  attach_object_names?: string[]
}) {
  return http.post<any>(`${prefix}/submit`, data)
}

export function myPage(params?: { current?: number; size?: number }) {
  return http.get<any>(`${prefix}/my-page`, { params })
}

export function myDetail(id: string) {
  return http.get<any>(`${prefix}/my-detail`, { params: { id } })
}
