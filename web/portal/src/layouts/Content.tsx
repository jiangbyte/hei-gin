/** Author: Charlie */

import { Outlet, useLocation } from 'react-router-dom'

export function Content() {
  const { pathname } = useLocation()
  const isAuthPage = pathname.startsWith('/auth/')
  const isProfilePage = pathname === '/profile' || pathname.startsWith('/profile/')

  const mainClass = isAuthPage
    ? 'w-full min-h-[calc(100vh-64px-72px)] bg-[color-mix(in_srgb,var(--ant-color-fill-quaternary)_85%,#f5f7fa)] px-0 py-0'
    : isProfilePage
      ? 'w-full min-h-[calc(100vh-64px-72px)] px-0 py-0'
      : 'w-full min-h-[calc(100vh-64px-72px)] px-6 py-6'

  return (
    <main className={mainClass}>
      <Outlet />
    </main>
  )
}
