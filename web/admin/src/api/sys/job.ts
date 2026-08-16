/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const jobPrefix = `${API_PREFIX}/sys/jobs`

/** 任务分页 */
export function page(params: any) {
  return http.get<any>(`${jobPrefix}/page`, { params })
}

/** 可用任务处理器列表 */
export function handlers() {
  return http.get<any>(`${jobPrefix}/handlers`)
}

/** 任务详情 */
export function detail(params: any) {
  return http.get<any>(`${jobPrefix}/detail`, { params })
}

/** 新建任务 */
export function create(data: any) {
  return http.post<any>(`${jobPrefix}/create`, data)
}

/** 更新任务 */
export function update(data: any) {
  return http.post<any>(`${jobPrefix}/update`, data)
}

/** 删除任务 */
export function remove(data: any) {
  return http.post<any>(`${jobPrefix}/delete`, data)
}

/** 启停任务 */
export function setStatus(data: any) {
  return http.post<any>(`${jobPrefix}/status`, data)
}

/** 立即触发 */
export function trigger(data: any) {
  return http.post<any>(`${jobPrefix}/trigger`, data)
}

/** 执行日志分页 */
export function logs(params: any) {
  return http.get<any>(`${jobPrefix}/logs`, { params })
}
