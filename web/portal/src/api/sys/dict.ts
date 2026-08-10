/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const dictPrefix = `${API_PREFIX}/sys/dicts`

export function tree(params?: { category?: string }) {
  // 门户字典公开接口（public：401 不跳登录）
  return http.get<any[]>(`${dictPrefix}/tree`, { params, public: true })
}
