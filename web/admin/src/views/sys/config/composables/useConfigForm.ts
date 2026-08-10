/** Author: Charlie */

import { configApi } from '@/api'

export type ConfigValueMap = Record<string, string>

export interface ConfigSaveItem {
  config_key: string
  config_value: string | null
  category?: string | null
  remark?: string | null
  value_type?: string | null
  label?: string | null
  scope?: string | null
  scene?: string | null
  is_builtin?: boolean | null
}

export function parseBool(value: string | null | undefined): boolean {
  return (
    String(value ?? '')
      .trim()
      .toUpperCase() === 'TRUE'
  )
}

export function toBoolStr(value: boolean): string {
  return value ? 'TRUE' : 'FALSE'
}

export function parseNumber(value: string | null | undefined, fallback = 0): number {
  const n = Number(value)
  return Number.isFinite(n) ? n : fallback
}

export function parseTags(raw: string | null | undefined): string[] {
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed.map(String) : [String(raw)]
  } catch {
    return String(raw)
      .split(/[,\n]/)
      .map((s) => s.trim())
      .filter(Boolean)
  }
}

export function tagsToJson(tags: string[]): string {
  return JSON.stringify(tags)
}

/** Load category rows into a key → value map (no ids). */
export async function loadByCategory(category: string): Promise<ConfigValueMap> {
  const res = await configApi.list({ category })
  const map: ConfigValueMap = {}
  for (const row of res.data ?? []) {
    if (row?.config_key) {
      map[row.config_key] = row.config_value ?? ''
      // STORAGE 敏感键脱敏后通过 ext_json.is_set 标记「已配置」
      if (row.ext_json?.is_set) {
        map[`${row.config_key}_SET`] = 'TRUE'
      }
    }
  }
  return map
}

/** Batch upsert by config_key only (no id). */
export async function saveByKeys(items: ConfigSaveItem[]): Promise<void> {
  await configApi.batchSave({
    items: items.map((item) => ({
      config_key: item.config_key,
      config_value: item.config_value,
      ...(item.category != null ? { category: item.category } : {}),
      ...(item.remark != null ? { remark: item.remark } : {}),
      ...(item.value_type != null ? { value_type: item.value_type } : {}),
      ...(item.label != null ? { label: item.label } : {}),
      ...(item.scope != null ? { scope: item.scope } : {}),
      ...(item.scene != null ? { scene: item.scene } : {}),
      ...(item.is_builtin != null ? { is_builtin: item.is_builtin } : {}),
    })),
  })
}

export function pickKeys(map: ConfigValueMap, keys: string[]): ConfigValueMap {
  const out: ConfigValueMap = {}
  for (const key of keys) {
    out[key] = map[key] ?? ''
  }
  return out
}
