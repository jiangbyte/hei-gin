/** Author: Charlie */

import { useCallback, useEffect, useState } from 'react'
import { Button, Form, Input, Modal, Space, Spin, Switch, Typography, message } from 'antd'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { encryptPasswords } from '@/utils/security'
import { isValidEmail } from '@/utils/validate'
import { wireBool } from '@/utils/wire'
import '../usercenter.css'

const OTP_COOLDOWN_SECONDS = 60

export function EmailPanel() {
  const refreshUserInfo = useAuthStore((s) => s.refreshUserInfo)
  const forceBindEmail = useAuthStore((s) => Boolean(s.userInfo?.forceBindEmail))
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [password, setPassword] = useState('')
  const [otpCode, setOtpCode] = useState('')
  const [confirmLoading, setConfirmLoading] = useState(false)
  const [sendingCode, setSendingCode] = useState(false)
  const [otpCooldown, setOtpCooldown] = useState(0)

  const applyProfile = useCallback(
    (data: any) => {
      const currentProfile = data?.profile ?? {}
      form.setFieldsValue({
        email: currentProfile.email ?? '',
        email_login_enabled: forceBindEmail
          ? true
          : wireBool(currentProfile.email_login_enabled ?? false),
      })
    },
    [form, forceBindEmail],
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

  useEffect(() => {
    if (otpCooldown <= 0) return
    const timer = window.setTimeout(() => setOtpCooldown((v) => v - 1), 1000)
    return () => window.clearTimeout(timer)
  }, [otpCooldown])

  async function saveEmail() {
    const values = await form.validateFields()
    const email = (values.email ?? '').trim()
    if (email && !isValidEmail(email)) {
      message.warning('请输入有效邮箱')
      return
    }
    if (!email && (values.email_login_enabled || forceBindEmail)) {
      message.warning('请输入邮箱')
      return
    }
    setPassword('')
    setOtpCode('')
    setConfirmOpen(true)
  }

  async function sendBindCode() {
    if (otpCooldown > 0 || sendingCode) return
    const email = (form.getFieldValue('email') ?? '').trim()
    if (!email || !isValidEmail(email)) {
      message.warning('请先填写有效邮箱')
      return
    }
    setSendingCode(true)
    try {
      await authApi.sendBindEmailCode({ target: email })
      message.success('验证码已发送')
      setOtpCooldown(OTP_COOLDOWN_SECONDS)
    } finally {
      setSendingCode(false)
    }
  }

  async function confirmBind() {
    if (!password) {
      message.warning('请输入当前密码')
      return
    }
    const email = (form.getFieldsValue().email ?? '').trim()
    if (email && !otpCode.trim()) {
      message.warning('请输入邮箱验证码')
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
        email: email || null,
        email_login_enabled: forceBindEmail
          ? true
          : wireBool(values.email_login_enabled ?? false),
        otp_code: otpCode.trim() || undefined,
      })
      setConfirmOpen(false)
      setPassword('')
      setOtpCode('')
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
        {forceBindEmail ? (
          <Typography.Paragraph type="warning">请先绑定邮箱后才能继续使用系统。</Typography.Paragraph>
        ) : null}
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
            <Switch disabled={forceBindEmail} />
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
          <Form.Item label="邮箱验证码">
            <Space.Compact style={{ width: '100%' }}>
              <Input
                value={otpCode}
                placeholder="请输入验证码"
                onChange={(e) => setOtpCode(e.target.value)}
              />
              <Button loading={sendingCode} disabled={otpCooldown > 0} onClick={() => void sendBindCode()}>
                {otpCooldown > 0 ? `${otpCooldown}s` : '发送验证码'}
              </Button>
            </Space.Compact>
          </Form.Item>
          <Form.Item label="当前密码">
            <Input.Password
              value={password}
              placeholder="请输入当前密码"
              onChange={(e) => setPassword(e.target.value)}
              onPressEnter={() => void confirmBind()}
            />
          </Form.Item>
        </Form>
        <Typography.Text type="secondary">绑定或更换邮箱需验证码与当前密码。</Typography.Text>
      </Modal>
    </>
  )
}
