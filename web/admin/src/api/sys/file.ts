/** Author: Charlie */

import { API_PREFIX } from '@/constants/api'
import { downloadRemoteUrl, getFilenameFromContentDisposition, http, saveBlob } from '@/utils'

const filePrefix = `${API_PREFIX}/sys/file`

export interface SysFileItem {
  id: string
  object_name: string
  original_name: string
  storage_provider: string
  bucket?: string | null
  content_type: string
  size: number
  url: string
  created_at?: string
  created_by?: string | null
  updated_at?: string
  updated_by?: string | null
}

export interface FileUrlResponse {
  object_name: string
  url: string
}

export interface FileUploadOptions {
  storage_provider?: string | null
}

export type FileDownloadTarget =
  | string
  | Partial<Pick<SysFileItem, 'id' | 'object_name' | 'original_name' | 'storage_provider'>>

export function page(params: any) {
  return http.get<any>(`${filePrefix}/page`, {
    params,
  })
}

export function detail(params: any) {
  return http.get<SysFileItem>(`${filePrefix}/detail`, {
    params,
  })
}

export function getById(id: string) {
  return detail({ id })
}

export function listByIds(ids: string[]) {
  return http.post<SysFileItem[]>(`${filePrefix}/list_by_ids`, { ids })
}

export function getByIds(ids: string[]) {
  return listByIds(ids)
}

export function update(data: any) {
  return http.post<any>(`${filePrefix}/update`, data)
}

export function remove(data: any) {
  return http.post<any>(`${filePrefix}/delete`, data)
}

export function upload(file: File, options: FileUploadOptions = {}) {
  const data = new FormData()
  data.append('file', file)
  if (options.storage_provider) {
    data.append('storage_provider', options.storage_provider)
  }
  return http.post<SysFileItem>(`${filePrefix}/upload`, data)
}

export function url(objectName: string) {
  return http.post<FileUrlResponse>(`${filePrefix}/url`, {
    object_name: objectName,
  })
}

export function presignedUrl(objectName: string) {
  return http.post<FileUrlResponse>(`${filePrefix}/presigned_url`, {
    object_name: objectName,
  })
}

export function download(id: string) {
  return http.get<Blob>(`${filePrefix}/download`, {
    params: { id },
    responseType: 'blob',
  })
}

export async function downloadFile(target: FileDownloadTarget, fallbackFilename = 'download') {
  const id = getDownloadTargetId(target)
  const objectName = getDownloadTargetObjectName(target)
  const filename = getDownloadTargetFilename(target, fallbackFilename)

  // 有 id 时走后端流式下载（Content-Disposition: attachment），避免跨域直链只「打开」不下载
  if (id) {
    const response = await download(id)
    const disposition = response.headers?.['content-disposition']
    const dispositionText = Array.isArray(disposition)
      ? disposition[0]
      : typeof disposition === 'string'
        ? disposition
        : undefined
    const responseFilename = getFilenameFromContentDisposition(dispositionText) || filename
    saveBlob(response.data, responseFilename)
    return responseFilename
  }

  if (objectName) {
    const response = await url(objectName)
    await downloadRemoteUrl(response.data?.url, filename)
    return filename
  }

  return filename
}

function getDownloadTargetId(target: FileDownloadTarget) {
  return typeof target === 'string' ? target : String(target.id || '')
}

function getDownloadTargetObjectName(target: FileDownloadTarget) {
  return typeof target === 'string' ? '' : String(target.object_name || '')
}

function getDownloadTargetFilename(target: FileDownloadTarget, fallbackFilename: string) {
  if (typeof target === 'string') {
    return fallbackFilename
  }
  return target.original_name || target.object_name || fallbackFilename
}
