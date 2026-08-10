/** Author: Charlie */

import type { LoaderFunctionArgs } from 'react-router-dom'
import { redirect } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth'
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
  return null
}

export async function guestOnly({ request }: LoaderFunctionArgs) {
  syncDictTree()
  void refreshDict()
  const ok = await useAuthStore.getState().ensureSession()
  if (ok) {
    const url = new URL(request.url)
    throw redirect(getSafeRedirect(url.searchParams.get('redirect')))
  }
  return null
}
