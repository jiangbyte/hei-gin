/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const configPrefix = `${API_PREFIX}/sys/config`

export function page(params: any) {
  return http.get<any>(`${configPrefix}/page`, {
    params,
  })
}

export function list(params: { category?: string; scope?: string }) {
  return http.get<any[]>(`${configPrefix}/list`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${configPrefix}/detail`, {
    params,
  })
}

export function create(data: any) {
  return http.post<any>(`${configPrefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${configPrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${configPrefix}/delete`, data)
}

export function batchSave(data: {
  items: Array<{
    config_key: string
    config_value: string | null
    category?: string | null
    remark?: string | null
  }>
}) {
  return http.post<any>(`${configPrefix}/batch-save`, data)
}

export function testAuditAlertWebhook(data: { webhook_url: string; webhook_secret: string }) {
  return http.post<any>(`${configPrefix}/audit-alert/test-webhook`, data, {
    skipErrorMessage: true,
  })
}

export function testAuditAlertPush() {
  return http.post<any>(
    `${configPrefix}/audit-alert/test-push`,
    {},
    {
      skipErrorMessage: true,
    },
  )
}
