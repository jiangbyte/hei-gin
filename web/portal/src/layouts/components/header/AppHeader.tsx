/** Author: Charlie */

import { useState } from 'react'
import { Button, Grid } from 'antd'
import { MenuOutlined } from '@ant-design/icons'
import { useLocation } from 'react-router-dom'
import { MobileDrawer } from '../common/MobileDrawer'
import { Logo } from './Logo'
import { NavMenu } from './NavMenu'
import { ThemeSwitch } from './ThemeSwitch'
import { Notices } from './Notices'
import { UserCenter } from './UserCenter'
import { useAuthStore } from '@/stores/auth'
import './header.css'

const { useBreakpoint } = Grid

export const HEADER_HEIGHT = 64

export function AppHeader() {
  const screens = useBreakpoint()
  const location = useLocation()
  const isMobile = !screens.md
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [pathname, setPathname] = useState(location.pathname)
  const isLogin = useAuthStore((s) => s.isLogin)
  const loggedIn = isLogin()

  if (location.pathname !== pathname) {
    setPathname(location.pathname)
    if (drawerOpen) setDrawerOpen(false)
  }

  return (
    <>
      <header className="app-header">
        <div className="app-header__inner">
          <div className="flex-y-center gap-3 min-w-0 flex-1">
            {isMobile ? (
              <Button
                type="text"
                icon={<MenuOutlined />}
                aria-label="打开菜单"
                onClick={() => setDrawerOpen(true)}
              />
            ) : null}
            <Logo />
            {!isMobile ? (
              <div className="min-w-0 flex-1 ml-4">
                <NavMenu mode="horizontal" />
              </div>
            ) : null}
          </div>
          <div className="flex shrink-0 items-center gap-2">
            <ThemeSwitch />
            {loggedIn ? <Notices /> : null}
            <UserCenter />
          </div>
        </div>
      </header>

      <MobileDrawer open={drawerOpen} onClose={() => setDrawerOpen(false)}>
        <NavMenu mode="inline" onNavigate={() => setDrawerOpen(false)} />
      </MobileDrawer>
    </>
  )
}
