/** Author: Charlie */

import { useEffect, useMemo } from 'react'
import { App as AntApp, ConfigProvider, theme as antTheme } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import { AppRouter } from '@/router'
import { useAppStore } from '@/stores/app'
import { syncDictTree, refreshDict } from '@/utils/dict'
import { portalComponentToken, portalSeedToken } from '@/theme/tokens'

dayjs.locale('zh-cn')

export default function App() {
  const resolvedTheme = useAppStore((s) => s.resolvedTheme)
  const setSystemDark = useAppStore((s) => s.setSystemDark)

  useEffect(() => {
    // 门户字典公开可读：启动时 sync + refresh（对齐 admin 工具用法）
    syncDictTree()
    void refreshDict()
  }, [])

  useEffect(() => {
    const media = window.matchMedia('(prefers-color-scheme: dark)')
    const handleChange = () => setSystemDark(media.matches)
    handleChange()
    media.addEventListener('change', handleChange)
    return () => media.removeEventListener('change', handleChange)
  }, [setSystemDark])

  useEffect(() => {
    const root = document.documentElement
    root.dataset.theme = resolvedTheme
    root.classList.toggle('dark', resolvedTheme === 'dark')
  }, [resolvedTheme])

  const theme = useMemo(
    () => ({
      cssVar: { key: 'portal' },
      algorithm: resolvedTheme === 'dark' ? antTheme.darkAlgorithm : antTheme.defaultAlgorithm,
      token: {
        ...portalSeedToken,
        ...(resolvedTheme === 'dark'
          ? {
              colorBgLayout: '#0f1117',
              colorBorder: '#2a2f3a',
              colorBorderSecondary: '#242933',
              colorSplit: '#242933',
            }
          : null),
      },
      components: portalComponentToken,
    }),
    [resolvedTheme],
  )

  return (
    <ConfigProvider locale={zhCN} theme={theme}>
      <AntApp>
        <div className="min-h-full bg-[var(--ant-color-bg-layout)] text-[var(--ant-color-text)]">
          <AppRouter />
        </div>
      </AntApp>
    </ConfigProvider>
  )
}
