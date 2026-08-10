/** Author: Charlie */

import { useEffect } from 'react'
import { Avatar, Button, Dropdown, Modal, Space, Tooltip, Typography, message } from 'antd'
import {
  FormOutlined,
  HomeOutlined,
  LogoutOutlined,
  SettingOutlined,
  UserOutlined,
} from '@ant-design/icons'
import type { DropdownProps, MenuProps } from 'antd'
import { useLocation, useNavigate } from 'react-router-dom'
import { useAuthModalStore } from '@/stores/authModal'
import { useAuthStore } from '@/stores/auth'

type Props = {
  compact?: boolean
  placement?: DropdownProps['placement']
}

export function UserCenter({ compact = false, placement = 'bottomRight' }: Props) {
  const navigate = useNavigate()
  const location = useLocation()
  const userInfo = useAuthStore((s) => s.userInfo)
  const ensureSession = useAuthStore((s) => s.ensureSession)
  const logout = useAuthStore((s) => s.logout)
  const openAuthModal = useAuthModalStore((s) => s.open)

  const loggedIn = Boolean(userInfo?.accountId)

  useEffect(() => {
    void ensureSession()
  }, [ensureSession])

  if (!loggedIn) {
    if (compact) {
      return (
        <Tooltip title="登录" placement="right">
          <Button
            type="text"
            className="!h-10 !w-10 !px-0"
            icon={<UserOutlined />}
            aria-label="登录"
            onClick={() => openAuthModal('login')}
          />
        </Tooltip>
      )
    }

    return (
      <Space size={8}>
        <Button onClick={() => openAuthModal('register')}>注册</Button>
        <Button type="primary" onClick={() => openAuthModal('login')}>
          登录
        </Button>
      </Space>
    )
  }

  const displayName = userInfo?.nickname || userInfo?.account || '用户'
  const avatarSrc = userInfo?.avatar || undefined

  const items: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人主页',
    },
    {
      key: 'userCenter',
      icon: <SettingOutlined />,
      label: '账号设置',
    },
    {
      key: 'feedback',
      icon: <FormOutlined />,
      label: '我的反馈',
    },
    { type: 'divider' },
    {
      key: 'home',
      icon: <HomeOutlined />,
      label: '首页',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
    },
  ]

  const onClick: MenuProps['onClick'] = ({ key }) => {
    if (key === 'profile') {
      navigate('/profile')
      return
    }
    if (key === 'userCenter') {
      navigate('/usercenter')
      return
    }
    if (key === 'feedback') {
      navigate('/feedback')
      return
    }
    if (key === 'home') {
      navigate('/')
      return
    }
    if (key === 'logout') {
      Modal.confirm({
        title: '退出登录',
        content: '确定退出当前账号？',
        okText: '确认',
        cancelText: '取消',
        onOk: async () => {
          await logout(location.pathname)
          message.success('已退出登录')
        },
      })
    }
  }

  const avatarBtn = (
    <Dropdown menu={{ items, onClick }} trigger={['click']} placement={placement}>
      <Space className="cursor-pointer select-none" size={8}>
        <Avatar src={avatarSrc} icon={<UserOutlined />} size={compact ? 32 : 'default'} />
        {compact ? null : (
          <Typography.Text className="hidden max-w-28 truncate md:inline">
            {displayName}
          </Typography.Text>
        )}
      </Space>
    </Dropdown>
  )

  if (compact) {
    return <div className="flex flex-col items-center gap-2">{avatarBtn}</div>
  }

  return <Space size={8}>{avatarBtn}</Space>
}
