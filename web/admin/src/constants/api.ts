/** Author: Charlie */

/**
 * 前端 API 前缀（需与后端装饰器当前约定一致：/api/v1/{client}/...）。
 * 改此处只影响本端拼 URL，不会升级后端路由；后端升版需同步改装饰器。
 */
export const API_ROOT = '/api'
export const API_VERSION = String(import.meta.env.VITE_API_VERSION || 'v1')
export const API_CLIENT = 'admin' as const

export const API_PREFIX = `${API_ROOT}/${API_VERSION}/${API_CLIENT}`

export function buildApiUrl(path: string) {
  const normalized = path.startsWith('/') ? path : `/${path}`
  return `${API_PREFIX}${normalized}`
}
