/** Author: Charlie */

const CHINA_OFFSET_MS = 8 * 60 * 60 * 1000
const LOCAL_DATE_TIME_PATTERN =
  /^(\d{4})-(\d{1,2})-(\d{1,2})(?:[T\s](\d{1,2})(?::(\d{1,2})(?::(\d{1,2})(?:\.\d+)?)?)?)?$/
const TIME_ZONE_PATTERN = /(?:[zZ]|[+-]\d{2}:?\d{2})$/

export function formatDateTime(value: unknown, fallback = '-') {
  if (value === undefined || value === null || value === '') {
    return fallback
  }

  if (value instanceof Date) {
    return Number.isNaN(value.getTime())
      ? fallback
      : formatChinaDate(value.getTime())
  }

  if (typeof value === 'number') {
    return formatTimestamp(value, fallback)
  }

  const text = String(value).trim()
  if (!text) {
    return fallback
  }

  if (/^\d+$/.test(text)) {
    return formatTimestamp(Number(text), fallback)
  }

  const normalized = text.replace(/\//g, '-')
  if (!TIME_ZONE_PATTERN.test(normalized)) {
    const localMatch = normalized.match(LOCAL_DATE_TIME_PATTERN)
    if (localMatch) {
      return formatLocalMatch(localMatch)
    }
  }

  const timestamp = Date.parse(normalized)
  return Number.isNaN(timestamp) ? text : formatChinaDate(timestamp)
}

function formatTimestamp(value: number, fallback: string) {
  if (!Number.isFinite(value)) {
    return fallback
  }
  const timestamp = Math.abs(value) < 1_000_000_000_000 ? value * 1000 : value
  return formatChinaDate(timestamp)
}

function formatLocalMatch(match: RegExpMatchArray) {
  const [, year, month, day, hour = '0', minute = '0', second = '0'] = match
  return `${pad(year, 4)}-${pad(month)}-${pad(day)} ${pad(hour)}:${pad(minute)}:${pad(second)}`
}

function formatChinaDate(timestamp: number) {
  const date = new Date(timestamp + CHINA_OFFSET_MS)
  return (
    [
      date.getUTCFullYear(),
      pad(date.getUTCMonth() + 1),
      pad(date.getUTCDate()),
    ].join('-') +
    ` ${pad(date.getUTCHours())}:${pad(date.getUTCMinutes())}:${pad(date.getUTCSeconds())}`
  )
}

function pad(value: string | number, length = 2) {
  return String(value).padStart(length, '0')
}
