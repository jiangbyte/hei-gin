/** Author: Charlie */

import { useCallback, useEffect, useState } from 'react'
import { Avatar, Button, Descriptions, Form, Input, Spin, Tabs, message } from 'antd'
import { EditOutlined, UserOutlined } from '@ant-design/icons'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { displayValue } from '../composables/useUserCenterProfile'
import { AvatarUploadModal } from './AvatarUploadModal'
import '../usercenter.css'

export function BasicInfoPanel() {
  const refreshUserInfo = useAuthStore((s) => s.refreshUserInfo)
  const [loading, setLoading] = useState(true)
  const [savingProfile, setSavingProfile] = useState(false)
  const [basicInfoTab, setBasicInfoTab] = useState('avatar')
  const [avatarModalShow, setAvatarModalShow] = useState(false)
  const [me, setMe] = useState<any>(null)
  const [profileForm] = Form.useForm()

  const applyMe = useCallback(
    (data: any) => {
      setMe(data)
      const currentProfile = data?.profile ?? {}
      profileForm.setFieldsValue({
        name: data?.name ?? currentProfile.name ?? '',
        nickname: data?.nickname ?? currentProfile.nickname ?? '',
        signature: currentProfile.signature ?? '',
      })
    },
    [profileForm],
  )

  const refresh = useCallback(async () => {
    setLoading(true)
    try {
      applyMe(await refreshUserInfo())
    } finally {
      setLoading(false)
    }
  }, [applyMe, refreshUserInfo])

  useEffect(() => {
    let cancelled = false
    void (async () => {
      try {
        const data = await refreshUserInfo()
        if (!cancelled) applyMe(data)
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [applyMe, refreshUserInfo])

  async function saveProfile() {
    setSavingProfile(true)
    try {
      const values = await profileForm.validateFields()
      await authApi.updateUserCenterProfile({
        name: values.name || null,
        nickname: values.nickname || null,
        signature: values.signature || null,
      })
      await refresh()
      message.success('保存成功')
    } finally {
      setSavingProfile(false)
    }
  }

  const profile = (me?.profile ?? {}) as any
  const avatarUrl = me?.avatar || profile.avatar || undefined
  const nickname = String(me?.nickname ?? '').trim()
  const name = String(me?.name ?? '').trim()
  const displayName =
    nickname && name && nickname !== name ? `${nickname}（${name}）` : nickname || name || '-'
  const contactParts = [profile.phone, profile.email].filter(Boolean)
  const contactText = contactParts.length ? contactParts.join(' / ') : ''

  return (
    <>
      <Spin spinning={loading}>
        <Tabs
          activeKey={basicInfoTab}
          onChange={setBasicInfoTab}
          className="user-center__subtabs"
          items={[
            {
              key: 'avatar',
              label: '头像',
              children: (
                <div className="user-center__avatar-card">
                  <button
                    type="button"
                    className="user-center__avatar-edit"
                    title="更换头像"
                    onClick={() => setAvatarModalShow(true)}
                  >
                    <Avatar size={160} src={avatarUrl} icon={<UserOutlined />} />
                    <span className="user-center__avatar-badge">
                      <EditOutlined />
                      编辑
                    </span>
                  </button>
                  <div className="user-center__avatar-name">{displayName}</div>
                  <div className="user-center__avatar-account">{me?.account || '-'}</div>
                  <Descriptions
                    className="user-center__avatar-desc"
                    column={1}
                    size="small"
                    labelStyle={{ width: 72 }}
                  >
                    <Descriptions.Item label="联系方式">
                      {displayValue(contactText)}
                    </Descriptions.Item>
                  </Descriptions>
                </div>
              ),
            },
            {
              key: 'profile',
              label: '基本信息',
              children: (
                <Form
                  form={profileForm}
                  layout="vertical"
                  className="user-center-form user-center-form--narrow w-full min-w-0"
                >
                  <Form.Item
                    label="账号"
                    extra={<span className="user-center__hint">登录账号不可修改。</span>}
                  >
                    <Input value={me?.account} disabled />
                  </Form.Item>
                  <Form.Item
                    name="name"
                    label="姓名"
                    rules={[{ max: 64, message: '姓名最多 64 个字符' }]}
                    extra={
                      <span className="user-center__hint">姓名可能出现在审批、审计等场景中。</span>
                    }
                  >
                    <Input allowClear />
                  </Form.Item>
                  <Form.Item
                    name="nickname"
                    label="昵称"
                    rules={[{ max: 64, message: '昵称最多 64 个字符' }]}
                  >
                    <Input allowClear />
                  </Form.Item>
                  <Form.Item name="signature" label="个性签名">
                    <Input.TextArea rows={3} placeholder="一句话介绍自己" allowClear />
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      loading={savingProfile}
                      onClick={() => void saveProfile()}
                    >
                      更新资料
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
          ]}
        />
      </Spin>

      <AvatarUploadModal
        open={avatarModalShow}
        avatar={avatarUrl}
        onClose={() => setAvatarModalShow(false)}
        onUploaded={() => void refresh()}
      />
    </>
  )
}
