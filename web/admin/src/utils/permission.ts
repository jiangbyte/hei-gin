/** Author: Charlie */

import { useAuthStore } from '@/stores'

export function hasPermission(permissionKey: string) {
  const authStore = useAuthStore()
  return authStore.hasPermission(permissionKey)
}
