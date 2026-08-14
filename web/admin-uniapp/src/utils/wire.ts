/** Author: Charlie */

/** HTTP JSON wire 辅助 — 标量仅为字符串。 */
export function wireBool(value: string): boolean {
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
        const result: Record<string, unknown> = {}
        for (const [key, item] of Object.entries(value as Record<string, unknown>)) {
            result[key] = stringifyScalars(item)
        }
        return result
    }
    return value
}
