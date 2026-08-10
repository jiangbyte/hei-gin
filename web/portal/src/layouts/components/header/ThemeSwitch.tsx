/** Author: Charlie */

import { Button, Dropdown, Space } from 'antd'
import { CheckOutlined, DesktopOutlined, MoonOutlined, SunOutlined } from '@ant-design/icons'
import type { DropdownProps, MenuProps } from 'antd'
import type { ReactNode } from 'react'
import { useAppStore, type ColorMode } from '@/stores/app'

const labels: Record<ColorMode, string> = {
  light: '浅色',
  dark: '深色',
  auto: '跟随系统',
}

const icons: Record<ColorMode, ReactNode> = {
  light: <SunOutlined />,
  dark: <MoonOutlined />,
  auto: <DesktopOutlined />,
}

type Props = {
  placement?: DropdownProps['placement']
}

export function ThemeSwitch({ placement = 'bottomRight' }: Props) {
  const colorMode = useAppStore((s) => s.colorMode)
  const setColorMode = useAppStore((s) => s.setColorMode)

  const items: MenuProps['items'] = (['light', 'dark', 'auto'] as ColorMode[]).map((mode) => ({
    key: mode,
    icon: icons[mode],
    label: (
      <Space className="w-24 justify-between">
        <span>{labels[mode]}</span>
        {colorMode === mode ? <CheckOutlined className="text-xs" /> : null}
      </Space>
    ),
  }))

  const onClick: MenuProps['onClick'] = ({ key }) => {
    setColorMode(key as ColorMode)
  }

  return (
    <Dropdown
      menu={{ items, onClick, selectedKeys: [colorMode] }}
      trigger={['click']}
      placement={placement}
    >
      <Button type="text" icon={icons[colorMode]} aria-label="切换主题" />
    </Dropdown>
  )
}
