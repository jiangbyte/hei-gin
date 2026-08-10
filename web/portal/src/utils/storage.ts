/** Author: Charlie */

export const USER_INFO_KEY = 'user_info'

export function getStoredUserInfo<T>() {
  const raw = localStorage.getItem(USER_INFO_KEY)
  if (!raw) {
    return null
  }
  try {
    return JSON.parse(raw) as T
  } catch {
    localStorage.removeItem(USER_INFO_KEY)
    return null
  }
}

export function setStoredUserInfo(userInfo: unknown | null) {
  if (userInfo) {
    localStorage.setItem(USER_INFO_KEY, JSON.stringify(userInfo))
  } else {
    localStorage.removeItem(USER_INFO_KEY)
  }
}

export function clearAuthStorage() {
  localStorage.removeItem(USER_INFO_KEY)
}
