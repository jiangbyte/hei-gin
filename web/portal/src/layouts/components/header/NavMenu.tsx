/** Author: Charlie */

import { Menu } from 'antd'
import { useLocation, useNavigate } from 'react-router-dom'
import { getSelectedNavKey, navItems } from './navItems'

type Props = {
  mode?: 'horizontal' | 'inline'
  onNavigate?: () => void
}

export function NavMenu({ mode = 'horizontal', onNavigate }: Props) {
  const navigate = useNavigate()
  const location = useLocation()
  const selectedKey = getSelectedNavKey(location.pathname)

  return (
    <Menu
      mode={mode}
      selectedKeys={[selectedKey]}
      items={navItems}
      onClick={({ key }) => {
        navigate(key)
        onNavigate?.()
      }}
      className={mode === 'horizontal' ? 'border-0 min-w-0 flex-1' : 'border-0'}
    />
  )
}
