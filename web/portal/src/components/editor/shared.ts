/** Author: Charlie */

export function toCssSize(value: string | number) {
  return typeof value === 'number' ? `${value}px` : value
}
