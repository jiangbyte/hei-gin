/** Author: Charlie */

import { useState } from 'react'
import type { ReactNode } from 'react'
import { theme } from 'antd'

export type CustomTabItem = {
  key: string
  label: string
  icon?: ReactNode
  children?: ReactNode
}

type Props = {
  items: CustomTabItem[]
  defaultActiveKey?: string
  activeKey?: string
  onChange?: (key: string) => void
  className?: string
  contentClassName?: string
  fillHeight?: boolean
}

export function CustomTabs({
  items,
  defaultActiveKey,
  activeKey,
  onChange,
  className,
  contentClassName,
  fillHeight = true,
}: Props) {
  const { token } = theme.useToken()
  const [innerActiveKey, setInnerActiveKey] = useState(defaultActiveKey ?? items[0]?.key)
  const current = activeKey ?? innerActiveKey
  const activeItem = items.find((item) => item.key === current) ?? items[0]

  function handleClick(key: string) {
    if (onChange) {
      onChange(key)
    } else {
      setInnerActiveKey(key)
    }
  }

  return (
    <div
      className={`tabs-shell flex min-h-0 flex-col ${fillHeight ? 'h-full' : ''} ${className ?? ''}`}
    >
      <div className="tabs-bar flex shrink-0 items-center gap-1">
        {items.map((item) => (
          <button
            key={item.key}
            type="button"
            className={`tabs-btn flex h-9 cursor-pointer items-center gap-1 border-b-2 px-4 text-sm transition-colors ${
              current === item.key ? 'tabs-btn-active font-medium' : 'border-transparent'
            }`}
            style={
              current === item.key
                ? { borderColor: token.colorPrimary, color: token.colorPrimary }
                : undefined
            }
            onClick={() => handleClick(item.key)}
          >
            {item.icon}
            {item.label}
          </button>
        ))}
      </div>
      <div
        className={`${fillHeight ? 'min-h-0 flex-1 overflow-y-auto' : ''} ${contentClassName ?? ''}`}
      >
        {activeItem?.children}
      </div>
    </div>
  )
}
