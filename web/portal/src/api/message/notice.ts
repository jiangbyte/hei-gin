/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = `${API_PREFIX}/message/notices`

/** 公开公告列表（可选登录以带已读状态） */
export function list(params?: { current?: number; size?: number }) {
  return http.get<any>(`${prefix}/list`, {
    params,
    public: true,
  })
}
