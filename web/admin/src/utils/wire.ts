/** Author: Charlie */

/** HTTP JSON wire 辅助 — 标量仅为字符串（无旧版 number/bool 路径）。 */
export function wireBool(value: string | boolean): boolean {
  if (typeof value === 'boolean') {
    return value
  }
  return value === 'true'
}

export function wireInt(value: string): number {
  const n = Number(value)
  if (!Number.isFinite(n)) {
    throw new Error(`Invalid wire int: ${value}`)
  }
  return n
}

export function wireFloat(value: string): number {
  const n = Number(value)
  if (!Number.isFinite(n)) {
    throw new Error(`Invalid wire float: ${value}`)
  }
  return n
}

/** 将 PageData meta 字符串转为数字，供 UI 分页组件使用。 */
export function readPageMeta(
  data: { current?: string; size?: string; total?: string; pages?: string },
  fallback: { current?: number; size?: number } = {},
): { current: number; size: number; total: number; pages?: number } {
  return {
    current:
      data.current !== undefined && data.current !== ''
        ? wireInt(data.current)
        : (fallback.current ?? 1),
    size: data.size !== undefined && data.size !== '' ? wireInt(data.size) : (fallback.size ?? 20),
    total: data.total !== undefined && data.total !== '' ? wireInt(data.total) : 0,
    pages: data.pages !== undefined && data.pages !== '' ? wireInt(data.pages) : undefined,
  }
}

export function stringifyScalars(value: unknown): unknown {
  if (value === null || value === undefined) {
    return value
  }
  if (typeof value === 'boolean') {
    return value ? 'true' : 'false'
  }
  if (typeof value === 'number') {
    return Number.isFinite(value) ? String(value) : value
  }
  if (typeof value === 'string') {
    return value
  }
  if (Array.isArray(value)) {
    return value.map((item) => stringifyScalars(item))
  }
  if (typeof value === 'object') {
    if (value instanceof FormData || value instanceof Blob || value instanceof ArrayBuffer) {
      return value
    }
    const result: Record<string, unknown> = {}
    for (const [key, item] of Object.entries(value as Record<string, unknown>)) {
      result[key] = stringifyScalars(item)
    }
    return result
  }
  return value
}
