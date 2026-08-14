/** Author: Charlie */

import type { LoaderFunctionArgs } from 'react-router-dom'
import { redirect } from 'react-router-dom'
import { isAllowedUnderSecurityWall, resolveSecurityWallPath, useAuthStore } from '@/stores/auth'
import { refreshDict, syncDictTree } from '@/utils/dict'
import { getSafeRedirect } from '@/utils/validate'

const publicPrefixes = ['/auth']

export function isPublicPath(pathname: string) {
  return publicPrefixes.some((prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`))
}

export async function requireAuth({ request }: LoaderFunctionArgs) {
  syncDictTree()
  void refreshDict()
  const ok = await useAuthStore.getState().ensureSession()
  if (!ok) {
    const url = new URL(request.url)
    const redirectTo = `${url.pathname}${url.search}`
    const search = redirectTo ? `?redirect=${encodeURIComponent(redirectTo)}` : ''
    throw redirect(`/auth/login${search}`)
  }

  const user = useAuthStore.getState().userInfo
  const url = new URL(request.url)
  if (!isAllowedUnderSecurityWall(url.pathname, url.search, user)) {
    const wall = resolveSecurityWallPath(user)
    if (wall) throw redirect(wall)
  }
  return null
}

export async function guestOnly({ request }: LoaderFunctionArgs) {
  syncDictTree()
  void refreshDict()
  const ok = await useAuthStore.getState().ensureSession()
  if (ok) {
    const url = new URL(request.url)
    const wall = resolveSecurityWallPath(useAuthStore.getState().userInfo)
    throw redirect(wall || getSafeRedirect(url.searchParams.get('redirect')))
  }
  return null
}
