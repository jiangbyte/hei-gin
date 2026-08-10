/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const auditPrefix = `${API_PREFIX}/sys/audit`

export function page(params: any) {
  return http.get<any>(`${auditPrefix}/page`, {
    params,
  })
}

export function detail(params: any) {
  return http.get<any>(`${auditPrefix}/detail`, {
    params,
  })
}
