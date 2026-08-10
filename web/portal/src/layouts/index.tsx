/** Author: Charlie */

import { useEffect } from 'react'
import { AuthModal } from '@/components/auth/AuthModal'
import { refreshDict, syncDictTree } from '@/utils/dict'
import { AppFooter, AppHeader, HEADER_HEIGHT } from './components'
import { Content } from './Content'

export function MainLayout() {
  useEffect(() => {
    syncDictTree()
    void refreshDict()
  }, [])

  return (
    <div className="flex min-h-screen flex-col bg-[var(--ant-color-bg-layout)] text-[var(--ant-color-text)]">
      <AppHeader />
      <div className="flex flex-1 flex-col" style={{ paddingTop: HEADER_HEIGHT }}>
        <Content />
        <AppFooter />
      </div>
      <AuthModal />
    </div>
  )
}
