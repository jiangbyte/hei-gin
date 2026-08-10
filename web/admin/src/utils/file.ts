/** Author: Charlie */

import { API_PREFIX, FILES_PUBLIC_PATH } from '@/constants/api'

export type UploadedFileValueType = 'auto' | 'url' | 'object_name'

export interface FileLike {
  id?: string | number | null
  object_name?: string | null
  objectName?: string | null
  original_name?: string | null
  originalName?: string | null
  name?: string | null
  content_type?: string | null
  contentType?: string | null
  type?: string | null
  size?: number | string | null
  url?: string | null
}

export interface NormalizedUploadedFile {
  raw: Record<string, unknown>
  value: string
  url: string
  objectName: string
  name: string
  contentType: string | null
  size: number | null
}

/** 管理端按 id 下载（非公开文件路径）。 */
export function buildAdminFileDownloadUrl(id?: string | number | null) {
  const rawId = String(id ?? '').trim()
  if (!rawId) {
    return undefined
  }
  return `${API_PREFIX}/sys/file/download?id=${encodeURIComponent(rawId)}`
}

/** object_name -> 公开访问路径。 */
export function buildPublicFileUrl(objectName?: string | null) {
  const name = String(objectName ?? '')
    .replace(/^\/+/, '')
    .trim()
  if (!name) {
    return undefined
  }
  return `${FILES_PUBLIC_PATH}/${name}`
}

export function formatFileSize(size?: number | string | null) {
  const value = Number(size ?? 0)
  if (!Number.isFinite(value) || value <= 0) {
    return '0 B'
  }
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let current = value
  let unitIndex = 0
  while (current >= 1024 && unitIndex < units.length - 1) {
    current /= 1024
    unitIndex += 1
  }
  return `${current.toFixed(unitIndex === 0 ? 0 : 2)} ${units[unitIndex]}`
}

export function isImageFile(file?: FileLike | string | null) {
  const contentType = getFileContentType(file)
  if (contentType.startsWith('image/')) {
    return true
  }
  return /\.(apng|avif|bmp|gif|ico|jpe?g|png|svg|webp)$/i.test(getFileName(file))
}

export function isVideoFile(file?: FileLike | string | null) {
  const contentType = getFileContentType(file)
  if (contentType.startsWith('video/')) {
    return true
  }
  return /\.(mp4|m4v|mov|ogg|ogv|webm|mkv|avi)$/i.test(getFileName(file))
}

export function normalizeUploadedFile(
  data: unknown,
  fallback?: File | FileLike | null,
  valueType: UploadedFileValueType = 'auto',
): NormalizedUploadedFile {
  const raw = isRecord(data) ? data : {}
  const fallbackRecord = isRecord(fallback) ? fallback : {}
  const objectName =
    getStringValue(raw, 'object_name') ||
    getStringValue(raw, 'objectName') ||
    getStringValue(fallbackRecord, 'object_name') ||
    getStringValue(fallbackRecord, 'objectName')
  const rawUrl = getStringValue(raw, 'url') || getStringValue(fallbackRecord, 'url')
  const name =
    getStringValue(raw, 'original_name') ||
    getStringValue(raw, 'originalName') ||
    getStringValue(raw, 'name') ||
    getStringValue(fallbackRecord, 'original_name') ||
    getStringValue(fallbackRecord, 'originalName') ||
    getStringValue(fallbackRecord, 'name') ||
    objectName
  const contentType =
    getStringValue(raw, 'content_type') ||
    getStringValue(raw, 'contentType') ||
    getStringValue(fallbackRecord, 'content_type') ||
    getStringValue(fallbackRecord, 'contentType') ||
    getStringValue(fallbackRecord, 'type') ||
    null
  const size = getNumberValue(raw, 'size') ?? getNumberValue(fallbackRecord, 'size')

  return {
    raw,
    value: normalizeFileValue({ url: rawUrl, object_name: objectName }, valueType),
    url: rawUrl || objectName,
    objectName,
    name,
    contentType,
    size,
  }
}

export function normalizeFileValue(
  file: FileLike | Record<string, unknown>,
  valueType: UploadedFileValueType = 'auto',
) {
  const objectName = getStringValue(file, 'object_name') || getStringValue(file, 'objectName')
  const rawUrl = getStringValue(file, 'url')
  if (valueType === 'object_name') {
    return objectName
  }
  if (valueType === 'url') {
    return rawUrl || objectName
  }
  return rawUrl || objectName
}

export function getFilenameFromContentDisposition(value?: string | null) {
  if (!value) {
    return undefined
  }
  const utf8Match = /filename\*=UTF-8''([^;]+)/i.exec(value)
  if (utf8Match?.[1]) {
    try {
      return decodeURIComponent(utf8Match[1])
    } catch {
      return utf8Match[1]
    }
  }
  const filenameMatch = /filename="?([^";]+)"?/i.exec(value)
  return filenameMatch?.[1]
}

export function saveBlob(blob: Blob, filename = 'download') {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.setTimeout(() => URL.revokeObjectURL(url), 1000)
}

function getFileContentType(file?: FileLike | string | null) {
  if (typeof file === 'string' || !file) {
    return ''
  }
  return String(file.content_type || file.contentType || '').toLowerCase()
}

function getFileName(file?: FileLike | string | null) {
  if (typeof file === 'string') {
    return file
  }
  if (!file) {
    return ''
  }
  return String(
    file.original_name ||
      file.originalName ||
      file.name ||
      file.object_name ||
      file.objectName ||
      file.url ||
      '',
  )
}

function getStringValue(data: unknown, key: string) {
  if (!isRecord(data)) {
    return ''
  }
  const value = data[key]
  return typeof value === 'string' ? value : ''
}

function getNumberValue(data: unknown, key: string) {
  if (!isRecord(data)) {
    return null
  }
  const value = data[key]
  if (typeof value === 'number') {
    return Number.isFinite(value) ? value : null
  }
  if (typeof value === 'string' && value.trim()) {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : null
  }
  return null
}

function isRecord(data: unknown): data is Record<string, unknown> {
  return typeof data === 'object' && data !== null
}
