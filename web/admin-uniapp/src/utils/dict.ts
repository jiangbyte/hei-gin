/** Author: Charlie */

import { shallowRef } from 'vue'

const DICT_TREE_STORAGE_KEY = 'hei:admin-uniapp:dict-tree'

export const dictTreeState = shallowRef<any[]>([])

let refreshDictPromise: Promise<void> | null = null

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
        const {API_PREFIX} = await import('@/constants/api')
      const { http } = await import('@/utils/request')
        const data = await http.get<any[]>(`${API_PREFIX}/sys/dicts/tree`)
      setDictTree(data ?? [])
    } finally {
      refreshDictPromise = null
    }
  })()

  return refreshDictPromise
}

export function clearDict() {
  dictTreeState.value = []
  uni.removeStorageSync(DICT_TREE_STORAGE_KEY)
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
  tree = dictDataAll()
) {
  const dict = findDictItem(dictCode, value, tree)
  return dict ? getDictLabel(dict) : ''
}

export function dictTypeColor(
  dictCode: string,
  value?: string | number | null,
  tree = dictDataAll()
) {
  const dict = findDictItem(dictCode, value, tree)
  return dict?.color || ''
}

export function translateDictTree(
  dictCode: string,
  value?: string | number | null,
  tree = dictDataAll()
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
  return (
    item.status === undefined ||
    item.status === null ||
    item.status === 'ENABLED'
  )
}

function findDictItem(
  dictCode: string,
  value?: string | number | null,
  tree = dictDataAll()
) {
  if (value === undefined || value === null || value === '') {
    return undefined
  }
  const normalizedValue = String(value)
  return dictTypeList(dictCode, tree).find(
    (item: any) => getDictValue(item) === normalizedValue
  )
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
  uni.setStorageSync(DICT_TREE_STORAGE_KEY, JSON.stringify(dictTreeState.value))
}

function readStoredDictTree() {
  const raw = uni.getStorageSync(DICT_TREE_STORAGE_KEY)
  if (!raw) {
    return []
  }
  try {
    const tree = JSON.parse(String(raw))
    return Array.isArray(tree) ? tree : []
  } catch {
    uni.removeStorageSync(DICT_TREE_STORAGE_KEY)
    return []
  }
}

syncDictTree()
