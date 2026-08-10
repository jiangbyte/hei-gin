/** Author: Charlie */

import { useCallback, useEffect, useState } from 'react'
import { Button, Form, Input, Modal, Spin, Switch, Typography, message } from 'antd'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { encryptPasswords } from '@/utils/security'
import { isValidEmail } from '@/utils/validate'
import '../usercenter.css'

export function EmailPanel() {
  const refreshUserInfo = useAuthStore((s) => s.refreshUserInfo)
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [password, setPassword] = useState('')
  const [confirmLoading, setConfirmLoading] = useState(false)

  const applyProfile = useCallback(
    (data: any) => {
      const currentProfile = data?.profile ?? {}
      form.setFieldsValue({
        email: currentProfile.email ?? '',
        email_login_enabled: Boolean(currentProfile.email_login_enabled),
      })
    },
    [form],
  )

  const refresh = useCallback(async () => {
    setLoading(true)
    try {
      applyProfile(await refreshUserInfo())
    } finally {
      setLoading(false)
    }
  }, [applyProfile, refreshUserInfo])

  useEffect(() => {
    let cancelled = false
    void (async () => {
      try {
        const data = await refreshUserInfo()
        if (!cancelled) applyProfile(data)
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [applyProfile, refreshUserInfo])

  async function saveEmail() {
    const values = await form.validateFields()
    const email = (values.email ?? '').trim()
    if (email && !isValidEmail(email)) {
      message.warning('请输入有效邮箱')
      return
    }
    if (!email && values.email_login_enabled) {
      message.warning('请输入邮箱')
      return
    }
    setPassword('')
    setConfirmOpen(true)
  }

  async function confirmBind() {
    if (!password) {
      message.warning('请输入当前密码')
      return
    }
    setConfirmLoading(true)
    setSaving(true)
    try {
      const encrypted = await encryptPasswords({ password })
      const values = form.getFieldsValue()
      await authApi.updateUserCenterEmail({
        password: encrypted.values.password || '',
        password_key_id: encrypted.password_key_id,
        email: (values.email ?? '').trim() || null,
        email_login_enabled: Boolean(values.email_login_enabled),
      })
      setConfirmOpen(false)
      setPassword('')
      await refresh()
      message.success('绑定已更新')
    } finally {
      setConfirmLoading(false)
      setSaving(false)
    }
  }

  return (
    <>
      <Spin spinning={loading}>
        <Form
          form={form}
          layout="vertical"
          className="user-center-form user-center-form--narrow w-full min-w-0"
        >
          <Form.Item
            name="email"
            label="邮箱"
            rules={[{ type: 'email', message: '邮箱格式不正确' }]}
          >
            <Input allowClear placeholder="your@example.com" />
          </Form.Item>
          <Form.Item name="email_login_enabled" label="启用邮箱登录" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item>
            <Button type="primary" loading={saving} onClick={() => void saveEmail()}>
              更新邮箱
            </Button>
          </Form.Item>
        </Form>
      </Spin>

      <Modal
        open={confirmOpen}
        title="确认更新邮箱"
        okText="确认"
        cancelText="取消"
        confirmLoading={confirmLoading}
        maskClosable={false}
        onOk={() => void confirmBind()}
        onCancel={() => setConfirmOpen(false)}
      >
        <Form layout="vertical">
          <Form.Item label="当前密码">
            <Input.Password
              value={password}
              placeholder="请输入当前密码"
              onChange={(e) => setPassword(e.target.value)}
              onPressEnter={() => void confirmBind()}
            />
          </Form.Item>
        </Form>
        <Typography.Text type="secondary">为保障账号安全，修改邮箱需验证当前密码。</Typography.Text>
      </Modal>
    </>
  )
}
