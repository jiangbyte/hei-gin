/** Author: Charlie */

export const navItems = [
  { key: '/', label: '首页' },
  { key: '/announcements', label: '公告' },
]

export function getSelectedNavKey(pathname: string) {
  return (
    navItems.find((item) => (item.key === '/' ? pathname === '/' : pathname.startsWith(item.key)))
      ?.key || '/'
  )
}
