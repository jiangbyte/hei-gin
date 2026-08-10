/** Author: Charlie */

export function mapNames(items?: { id: string; name: string }[]) {
  return (items ?? [])
    .map((item) => item.name)
    .filter(Boolean)
    .join(' / ')
}

export function displayValue(value?: unknown) {
  if (value === undefined || value === null || value === '') {
    return '-'
  }
  return String(value)
}
