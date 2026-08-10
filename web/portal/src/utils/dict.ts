/** Author: Charlie
 *
 *
 * Portal 字典工具 — API 与 admin/src/utils/dict.ts 对齐。
 * 仅存储 key / 无 Vue 依赖处不同；查询与刷新语义保持一致。
 */
const DICT_TREE_STORAGE_KEY = 'hei:portal:dict-tree'

type DictTreeRef = {
  value: any[]
}

const listeners = new Set<() => void>()

function notifyDictListeners() {
  listeners.forEach((listener) => listener())
}

/** 与 admin `shallowRef` 同形：`.value` 读写；写入时通知 React 订阅方。 */
export const dictTreeState: DictTreeRef = {
  get value() {
    return _dictTree
  },
  set value(next: any[]) {
    _dictTree = Array.isArray(next) ? next : []
    notifyDictListeners()
  },
}

let _dictTree: any[] = []
let refreshDictPromise: Promise<void> | null = null

/** React：订阅字典树变更（对齐 admin 对 shallowRef 的响应式依赖）。 */
export function subscribeDict(listener: () => void) {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}

export function syncDictTree() {
  dictTreeState.value = readStoredDictTree()
  return dictTreeState.value
}

export function isDictLoaded() {
  return dictTreeState.value.length > 0 || syncDictTree().length > 0
}

export async function refreshDict() {
  if (refreshDictPromise) {
    return refreshDictPromise
  }

  refreshDictPromise = (async () => {
    try {
      const dictApi = await import('@/api/sys/dict')
      const response = await dictApi.tree()
      setDictTree(response.data ?? [])
    } finally {
      refreshDictPromise = null
    }
  })()

  return refreshDictPromise
}

export function clearDict() {
  dictTreeState.value = []
  localStorage.removeItem(DICT_TREE_STORAGE_KEY)
}

export function dictDataAll() {
  return dictTreeState.value
}

export function dictTypeList(dictCode: string, tree = dictDataAll()) {
  return findDictRoot(tree, dictCode)?.children ?? []
}

export function dictList(dictCode: any, tree = dictDataAll()) {
  return dictTypeList(dictCode, tree)
    .filter(isEnabledDict)
    .map((item: any) => ({
      label: getDictLabel(item),
      value: getDictValue(item),
    }))
}

export function dictTypeData(
  dictCode: string,
  value?: string | number | null,
  tree = dictDataAll(),
) {
  const dict = findDictItem(dictCode, value, tree)
  return dict ? getDictLabel(dict) : ''
}

export function dictTypeColor(
  dictCode: string,
  value?: string | number | null,
  tree = dictDataAll(),
) {
  const dict = findDictItem(dictCode, value, tree)
  return dict?.color || ''
}

export function translateDictTree(
  dictCode: string,
  value?: string | number | null,
  tree = dictDataAll(),
) {
  const root = findDictRoot(tree, dictCode)
  const dict = findNodeByValue(root, value)
  return dict ? getDictLabel(dict) : ''
}

export function getDictValue(item: any) {
  return item.value || item.code
}

export function getDictLabel(item: any) {
  return item.label || item.code
}

export function isEnabledDict(item: any) {
  return item.status === undefined || item.status === null || item.status === 'ENABLED'
}

function findDictItem(dictCode: string, value?: string | number | null, tree = dictDataAll()) {
  if (value === undefined || value === null || value === '') {
    return undefined
  }
  const normalizedValue = String(value)
  return dictTypeList(dictCode, tree).find((item: any) => getDictValue(item) === normalizedValue)
}

function findDictRoot(tree: any[], dictCode: string) {
  return tree.find((item) => item.code === dictCode)
}

function findNodeByValue(node: any, value?: string | number | null): any {
  if (!node || value === undefined || value === null || value === '') {
    return undefined
  }

  if (getDictValue(node) === String(value)) {
    return node
  }

  for (const child of node.children ?? []) {
    const result = findNodeByValue(child, value)
    if (result) {
      return result
    }
  }

  return undefined
}

function setDictTree(tree: any[]) {
  dictTreeState.value = Array.isArray(tree) ? tree : []
  localStorage.setItem(DICT_TREE_STORAGE_KEY, JSON.stringify(dictTreeState.value))
}

function readStoredDictTree() {
  const raw = localStorage.getItem(DICT_TREE_STORAGE_KEY)
  if (!raw) {
    return []
  }
  try {
    const tree = JSON.parse(raw)
    return Array.isArray(tree) ? tree : []
  } catch {
    localStorage.removeItem(DICT_TREE_STORAGE_KEY)
    return []
  }
}

syncDictTree()
