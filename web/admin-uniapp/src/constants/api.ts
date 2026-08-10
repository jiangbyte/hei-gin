/** Author: Charlie */

export const API_ROOT = '/api'
export const API_VERSION = String(import.meta.env.VITE_API_VERSION || 'v1')
export const API_CLIENT = 'admin' as const

export const API_PREFIX = `${API_ROOT}/${API_VERSION}/${API_CLIENT}`
export const FILES_PUBLIC_PATH = `${API_ROOT}/${API_VERSION}/files`
