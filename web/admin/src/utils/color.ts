/** Author: Charlie */

const hexColorPattern = /^#(?:[\da-f]{3}|[\da-f]{6})$/i
const lightTextColor = '#ffffff'
const darkTextColor = '#000000'

export function isHexColor(value?: string | null) {
  const color = String(value ?? '').trim()
  if (!color) {
    return true
  }

  return hexColorPattern.test(color)
}

export function createTagColor(value?: string | null) {
  const color = String(value ?? '').trim()
  if (!color) {
    return undefined
  }

  const textColor = getReadableTextColor(color)
  return textColor ? { color, textColor } : { color }
}

function getReadableTextColor(color: string) {
  const rgb = parseHexColor(color)
  if (!rgb) {
    return undefined
  }

  return getRelativeLuminance(rgb) > 0.5 ? darkTextColor : lightTextColor
}

function parseHexColor(color: string) {
  if (!hexColorPattern.test(color)) {
    return null
  }

  const hex = color.slice(1)
  const normalizedHex =
    hex.length === 3
      ? hex
          .split('')
          .map((char) => char + char)
          .join('')
      : hex

  return {
    r: Number.parseInt(normalizedHex.slice(0, 2), 16),
    g: Number.parseInt(normalizedHex.slice(2, 4), 16),
    b: Number.parseInt(normalizedHex.slice(4, 6), 16),
  }
}

function getRelativeLuminance(rgb: { r: number; g: number; b: number }) {
  const [r, g, b] = [rgb.r, rgb.g, rgb.b].map((channel) => {
    const value = channel / 255
    return value <= 0.03928 ? value / 12.92 : ((value + 0.055) / 1.055) ** 2.4
  })

  return 0.2126 * r + 0.7152 * g + 0.0722 * b
}
