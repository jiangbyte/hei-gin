/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const bannerPrefix = `${API_PREFIX}/sys/banners`

/** 门户公开展示图列表（public：401 不跳登录） */
export function listBanners(params: any) {
  return http.get<any>(`${bannerPrefix}/list`, {
    params,
    public: true,
  })
}

/** 记录展示图点击交互 */
export function recordBannerInteraction(id: string) {
  return http.post<any>(`${bannerPrefix}/interaction`, { id }, { public: true })
}
