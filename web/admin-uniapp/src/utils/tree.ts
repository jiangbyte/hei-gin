/** Author: Charlie */

export function arrayToTree<
    T extends { id: string; parent_id?: string | null; children?: T[] },
>(items: T[]) {
    const map = new Map<string, T>()
    const roots: T[] = []

    items.forEach((item) => {
        item.children = []
        map.set(item.id, item)
    })

    items.forEach((item) => {
        if (item.parent_id && map.has(item.parent_id)) {
            map.get(item.parent_id)!.children!.push(item)
        } else {
            roots.push(item)
        }
    })

    return roots
}

export function flattenTree<T extends { children?: T[] }>(items: T[]) {
    const result: T[] = []
    const walk = (nodes: T[]) => {
        nodes.forEach((node) => {
            result.push(node)
            if (node.children?.length) {
                walk(node.children)
            }
        })
    }
    walk(items)
    return result
}
