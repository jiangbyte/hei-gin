/** Author: Charlie */

import type { ReactNode } from 'react'
import { Drawer } from 'antd'
import { Logo } from '../header/Logo'

type Props = {
  open: boolean
  onClose: () => void
  children: ReactNode
}

export function MobileDrawer({ open, onClose, children }: Props) {
  return (
    <Drawer
      open={open}
      onClose={onClose}
      placement="left"
      width={280}
      styles={{ body: { padding: 0 } }}
      title={<Logo />}
      destroyOnHidden
    >
      {children}
    </Drawer>
  )
}
