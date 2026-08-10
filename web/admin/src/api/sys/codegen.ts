/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { getFilenameFromContentDisposition, http, saveBlob } from '@/utils'

const prefix = `${API_PREFIX}/sys/codegen`

export function page(params: any) {
  return http.get<any>(`${prefix}/page`, { params })
}

export function detail(params: any) {
  return http.get<any>(`${prefix}/detail`, { params })
}

export function create(data: any) {
  return http.post<any>(`${prefix}/create`, data)
}

export function update(data: any) {
  return http.post<any>(`${prefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${prefix}/delete`, data)
}

export function tables() {
  return http.get<any>(`${prefix}/tables`)
}

export function tableColumns(params: any) {
  return http.get<any>(`${prefix}/table-columns`, { params })
}

export function fields(params: any) {
  return http.get<any>(`${prefix}/fields`, { params })
}

export function updateFieldsBatch(data: any) {
  return http.post<any>(`${prefix}/fields/update-batch`, data)
}

export function parentResources(params?: any) {
  return http.get<any>(`${prefix}/parent-resources`, { params })
}

export function preview(params: any) {
  return http.get<any>(`${prefix}/preview`, { params })
}

export function download(params: any) {
  return http.get<Blob>(`${prefix}/download`, {
    params,
    responseType: 'blob',
  })
}

export async function downloadZip(id: string) {
  const response = await download({ id })
  const disposition = response.headers?.['content-disposition']
  const dispositionText = Array.isArray(disposition) ? disposition[0] : disposition
  const filename = getFilenameFromContentDisposition(dispositionText) || `codegen-${id}.zip`
  saveBlob(response.data, filename)
}
