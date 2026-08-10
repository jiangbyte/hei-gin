/** Author: Charlie */

type SearchValueNormalizer<T extends object> = Partial<{
  [K in keyof T]: (value: T[K]) => T[K]
}>

export function normalizeSearchValues<T extends object>(
  values: T,
  normalizers: SearchValueNormalizer<T> = {},
) {
  const result = {} as Partial<T>

  Object.entries(values as Record<string, unknown>).forEach(([key, value]) => {
    if (value === undefined || value === '') {
      return
    }

    const field = key as keyof T
    const normalizeValue = normalizers[field]
    const normalizedValue = normalizeValue
      ? normalizeValue(value as T[typeof field])
      : (value as T[typeof field])
    result[field] = normalizedValue
  })

  return result
}

export function toNullableString(value: unknown) {
  if (value === undefined || value === null) {
    return null
  }

  const text = String(value).trim()
  return text || null
}

export function displayValue(value?: string | number | null) {
  return value === undefined || value === null || value === '' ? '-' : String(value)
}

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function isValidEmail(value: unknown) {
  const text = String(value ?? '').trim()
  return EMAIL_PATTERN.test(text)
}

export function createRequiredRule(field: string, trigger: 'input' | 'change') {
  return {
    required: true,
    message: trigger === 'change' ? `请选择${field}` : `请输入${field}`,
    trigger,
  }
}

export function createRequiredArrayRule(field: string) {
  return {
    type: 'array' as const,
    required: true,
    trigger: 'change' as const,
    min: 1,
    message: `请选择${field}`,
  }
}
