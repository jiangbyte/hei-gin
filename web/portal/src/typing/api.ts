/** Author: Charlie */

export interface ApiResponse<T = unknown> {
  code: string
  message?: string
  data: T
}

export interface PageData<T> {
  current: string
  size: string
  total: string
  pages?: string
  records: T[]
}
