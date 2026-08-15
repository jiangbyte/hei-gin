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
  ApiOutlined,
} from '@ant-design/icons'
import { useSearchParams } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth'
import { BasicInfoPanel } from './components/BasicInfoPanel'
import { CancelAccountPanel } from './components/CancelAccountPanel'
import { EmailPanel } from './components/EmailPanel'
import { MyMessagesPanel } from './components/MyMessagesPanel'
import { OauthPanel } from './components/OauthPanel'
import { PasswordPanel } from './components/PasswordPanel'
import { PhonePanel } from './components/PhonePanel'
import './usercenter.css'

const NAV_ITEMS = [
  { key: 'basic_info', label: '公开资料' },
  { key: 'my_messages', label: '我的消息' },
  { key: 'password', label: '密码' },
  { key: 'phone', label: '手机号' },
  { key: 'email', label: '邮箱' },
  { key: 'oauth', label: '三方账号' },
  { key: 'cancel_account', label: '账号注销' },
] as const

type TabKey = (typeof NAV_ITEMS)[number]['key']

const PANEL_MAP: Record<TabKey, ReactNode> = {
  basic_info: <BasicInfoPanel />,
  my_messages: <MyMessagesPanel />,
  password: <PasswordPanel />,
  phone: <PhonePanel />,
  email: <EmailPanel />,
  oauth: <OauthPanel />,
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
  const userInfo = useAuthStore((s) => s.userInfo)
  const activeTab = resolveTab(searchParams.get('tab'))
  const activeNav = NAV_ITEMS.find((item) => item.key === activeTab) ?? NAV_ITEMS[0]

  const lockedTabs = useMemo(() => {
    if (userInfo?.passwordExpired) return new Set<TabKey>(['password'])
    const allowed = new Set<TabKey>()
    if (userInfo?.forceBindEmail) allowed.add('email')
    if (userInfo?.forceBindPhone) allowed.add('phone')
    return allowed.size > 0 ? allowed : null
  }, [userInfo])

  const menuItems: MenuProps['items'] = useMemo(
    () => [
      {
        key: 'basic_info',
        icon: <UserOutlined />,
        label: '公开资料',
        disabled: Boolean(lockedTabs && !lockedTabs.has('basic_info')),
      },
      {
        type: 'group',
        label: '消息',
        children: [
          {
            key: 'my_messages',
            icon: <MessageOutlined />,
            label: '我的消息',
            disabled: Boolean(lockedTabs && !lockedTabs.has('my_messages')),
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
            disabled: Boolean(lockedTabs && !lockedTabs.has('password')),
          },
          {
            key: 'phone',
            icon: <MobileOutlined />,
            label: '手机号',
            disabled: Boolean(lockedTabs && !lockedTabs.has('phone')),
          },
          {
            key: 'email',
            icon: <MailOutlined />,
            label: '邮箱',
            disabled: Boolean(lockedTabs && !lockedTabs.has('email')),
          },
          {
            key: 'oauth',
            icon: <ApiOutlined />,
            label: '三方账号',
            disabled: Boolean(lockedTabs && !lockedTabs.has('oauth')),
          },
          {
            key: 'cancel_account',
            icon: <DeleteOutlined />,
            label: '账号注销',
            danger: true,
            disabled: Boolean(lockedTabs && !lockedTabs.has('cancel_account')),
          },
        ],
      },
    ],
    [lockedTabs],
  )

  function selectTab(key: string) {
    if (!key || activeTab === key) return
    if (!NAV_ITEMS.some((item) => item.key === key)) return
    if (lockedTabs && !lockedTabs.has(key as TabKey)) return
    const next = new URLSearchParams(searchParams)
    next.set('tab', key)
    setSearchParams(next, { replace: true })
  }

  return (
    <div className="profile page-shell w-full min-w-0">
      <div className="profile__body">
        <aside className="profile__sidebar">
          <Menu
            mode="inline"
            selectedKeys={[activeTab]}
            items={menuItems}
            onClick={({ key }) => selectTab(String(key))}
          />
        </aside>

        <section className="profile__content">
          <div className="profile__panel">
            <h2
              className={
                activeTab === 'basic_info'
                  ? 'profile__panel-title profile__panel-title--with-tabs'
                  : 'profile__panel-title'
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
