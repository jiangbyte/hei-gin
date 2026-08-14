/** Author: Charlie */

import type { ThemeConfig } from 'antd'

/**
 * Portal 主题：主色对齐 admin `stores/app/theme.json`（#1677FF）。
 */
export const portalSeedToken = {
  colorPrimary: '#1677FF',
  borderRadius: 5,
  colorBgLayout: '#f5f7fb',
  colorBorder: '#e6ebf2',
  colorBorderSecondary: '#eef2f7',
  colorSplit: '#eef2f7',
  fontFamily:
    "Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
} as const

export const portalComponentToken: ThemeConfig['components'] = {
  Table: {
    headerBg: 'transparent',
    borderColor: '#eef2f7',
  },
  Card: {
    paddingLG: 20,
  },
  Tabs: {
    horizontalItemPadding: '10px 16px',
  },
}
