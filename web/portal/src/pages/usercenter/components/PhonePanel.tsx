/** Author: Charlie */

import { useCallback, useEffect, useState } from 'react'
import { Button, Form, Input, Modal, Spin, Switch, Typography, message } from 'antd'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { encryptPasswords } from '@/utils/security'
import '../usercenter.css'

export function PhonePanel() {
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
        phone: currentProfile.phone ?? '',
        phone_login_enabled: Boolean(currentProfile.phone_login_enabled),
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

  async function savePhone() {
    await form.validateFields()
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
      await authApi.updateUserCenterPhone({
        password: encrypted.values.password || '',
        password_key_id: encrypted.password_key_id,
        phone: values.phone || null,
        phone_login_enabled: Boolean(values.phone_login_enabled),
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
          <Form.Item name="phone" label="手机号">
            <Input allowClear />
          </Form.Item>
          <Form.Item name="phone_login_enabled" label="启用手机号登录" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item>
            <Button type="primary" loading={saving} onClick={() => void savePhone()}>
              更新手机号
            </Button>
          </Form.Item>
        </Form>
      </Spin>

      <Modal
        open={confirmOpen}
        title="确认更新手机号"
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
        <Typography.Text type="secondary">
          为保障账号安全，修改手机号需验证当前密码。
        </Typography.Text>
      </Modal>
    </>
  )
}
