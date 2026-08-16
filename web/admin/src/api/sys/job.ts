/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const jobPrefix = `${API_PREFIX}/sys/jobs`
const jobLogPrefix = `${API_PREFIX}/sys/job-logs`

/** 任务分页 */
export function page(params: any) {
  return http.get<any>(`${jobPrefix}/page`, { params })
}
export function detail(params: any) {
  return http.get<any>(`${jobPrefix}/detail`, { params })
}
export function create(data: any) {
  return http.post<any>(`${jobPrefix}/create`, data)
}
export function update(data: any) {
  return http.post<any>(`${jobPrefix}/update`, data)
}
export function remove(data: any) {
  return http.post<any>(`${jobPrefix}/delete`, data)
}
/** 启停 */
export function enabled(data: any) {
  return http.post<any>(`${jobPrefix}/enabled`, data)
}
/** 立即执行 */
export function run(data: any) {
  return http.post<any>(`${jobPrefix}/run`, data)
}

/** 执行日志分页 */
export function logPage(params: any) {
  return http.get<any>(`${jobLogPrefix}/page`, { params })
}
