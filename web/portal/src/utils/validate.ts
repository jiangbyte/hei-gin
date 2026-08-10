/** Author: Charlie */

export function isValidEmail(value: string) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value.trim())
}

export function isValidPhone(value: string) {
  return /^1\d{10}$/.test(value.trim())
}

export function getSafeRedirect(redirect?: string | null) {
  if (!redirect || redirect.startsWith('/auth')) {
    return import.meta.env.VITE_HOME_PATH || '/'
  }
  return redirect
}
