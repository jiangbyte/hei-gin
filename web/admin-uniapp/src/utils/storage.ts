/** Author: Charlie */

export function getStorage<T>(key: string, fallback: T): T {
  const value = uni.getStorageSync(key)
  if (!value) {
    return fallback
  }
  try {
    return JSON.parse(String(value)) as T
  } catch {
    return fallback
  }
}

export function setStorage<T>(key: string, value: T) {
  uni.setStorageSync(key, JSON.stringify(value))
}

export function removeStorage(key: string) {
  uni.removeStorageSync(key)
}
