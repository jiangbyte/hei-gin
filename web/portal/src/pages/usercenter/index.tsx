/** Author: Charlie */

import { useMemo, type ReactNode } from 'react'
import { Menu } from 'antd'
import type { MenuProps } from 'antd'
import {
  DeleteOutlined,
  LockOutlined,
  MailOutlined,
  MessageOutlined,
  MobileOutlined,
  UserOutlined,
} from '@ant-design/icons'
import { useSearchParams } from 'react-router-dom'
import { BasicInfoPanel } from './components/BasicInfoPanel'
import { CancelAccountPanel } from './components/CancelAccountPanel'
import { EmailPanel } from './components/EmailPanel'
import { MyMessagesPanel } from './components/MyMessagesPanel'
import { PasswordPanel } from './components/PasswordPanel'
import { PhonePanel } from './components/PhonePanel'
import './usercenter.css'

const NAV_ITEMS = [
  { key: 'basic_info', label: '公开资料' },
  { key: 'my_messages', label: '我的消息' },
  { key: 'password', label: '密码' },
  { key: 'phone', label: '手机号' },
  { key: 'email', label: '邮箱' },
  { key: 'cancel_account', label: '账号注销' },
] as const

type TabKey = (typeof NAV_ITEMS)[number]['key']

const PANEL_MAP: Record<TabKey, ReactNode> = {
  basic_info: <BasicInfoPanel />,
  my_messages: <MyMessagesPanel />,
  password: <PasswordPanel />,
  phone: <PhonePanel />,
  email: <EmailPanel />,
  cancel_account: <CancelAccountPanel />,
}

function resolveTab(tab: string | null): TabKey {
  if (tab && NAV_ITEMS.some((item) => item.key === tab)) {
    return tab as TabKey
  }
  return 'basic_info'
}

export function UserCenterPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const activeTab = resolveTab(searchParams.get('tab'))
  const activeNav = NAV_ITEMS.find((item) => item.key === activeTab) ?? NAV_ITEMS[0]

  const menuItems: MenuProps['items'] = useMemo(
    () => [
      {
        key: 'basic_info',
        icon: <UserOutlined />,
        label: '公开资料',
      },
      {
        type: 'group',
        label: '消息',
        children: [
          {
            key: 'my_messages',
            icon: <MessageOutlined />,
            label: '我的消息',
          },
        ],
      },
      {
        type: 'group',
        label: '访问与安全',
        children: [
          {
            key: 'password',
            icon: <LockOutlined />,
            label: '密码',
          },
          {
            key: 'phone',
            icon: <MobileOutlined />,
            label: '手机号',
          },
          {
            key: 'email',
            icon: <MailOutlined />,
            label: '邮箱',
          },
          {
            key: 'cancel_account',
            icon: <DeleteOutlined />,
            label: '账号注销',
            danger: true,
          },
        ],
      },
    ],
    [],
  )

  function selectTab(key: string) {
    if (!key || activeTab === key) return
    if (!NAV_ITEMS.some((item) => item.key === key)) return
    const next = new URLSearchParams(searchParams)
    next.set('tab', key)
    setSearchParams(next, { replace: true })
  }

  return (
    <div className="user-center page-shell w-full min-w-0">
      <div className="user-center__body">
        <aside className="user-center__sidebar">
          <Menu
            mode="inline"
            selectedKeys={[activeTab]}
            items={menuItems}
            onClick={({ key }) => selectTab(String(key))}
          />
        </aside>

        <section className="user-center__content">
          <div className="user-center__panel">
            <h2
              className={
                activeTab === 'basic_info'
                  ? 'user-center__panel-title user-center__panel-title--with-tabs'
                  : 'user-center__panel-title'
              }
            >
              {activeNav.label}
            </h2>
            {PANEL_MAP[activeTab]}
          </div>
        </section>
      </div>
    </div>
  )
}
