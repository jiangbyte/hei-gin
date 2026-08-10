/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/dashboard`

export function overview() {
  return http.get<any>(`${prefix}/overview`)
}
