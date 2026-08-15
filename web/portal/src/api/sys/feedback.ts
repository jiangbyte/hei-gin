/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/sys/feedbacks`

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
