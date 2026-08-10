/** Storage provider labels/colors — aligned with StorageProvider enum, not sys_dict. */

export const STORAGE_PROVIDER_OPTIONS = [
  { label: '本地文件', value: 'local', color: '#18a058' },
  { label: '阿里云 OSS', value: 'oss', color: '#f0a020' },
  { label: '腾讯云 COS', value: 's3', color: '#722ed1' },
  { label: 'MinIO', value: 'minio', color: '#2080f0' },
  { label: 'RustFS', value: 'rustfs', color: '#0d9488' },
] as const

const labelMap = Object.fromEntries(
  STORAGE_PROVIDER_OPTIONS.map((item) => [item.value, item.label]),
) as Record<string, string>

const colorMap = Object.fromEntries(
  STORAGE_PROVIDER_OPTIONS.map((item) => [item.value, item.color]),
) as Record<string, string>

export function storageProviderLabel(value?: string | null) {
  if (!value) {
    return ''
  }
  return labelMap[value] || value
}

export function storageProviderColor(value?: string | null) {
  if (!value) {
    return undefined
  }
  return colorMap[value]
}
