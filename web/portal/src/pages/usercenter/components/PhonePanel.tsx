/** Author: Charlie */

import { useCallback, useEffect, useState } from 'react'
import { Button, Form, Input, Modal, Space, Spin, Switch, Typography, message } from 'antd'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { encryptPasswords } from '@/utils/security'
import { isValidPhone } from '@/utils/validate'
import { wireBool } from '@/utils/wire'
import '../usercenter.css'

const OTP_COOLDOWN_SECONDS = 60

export function PhonePanel() {
  const refreshUserInfo = useAuthStore((s) => s.refreshUserInfo)
  const forceBindPhone = useAuthStore((s) => Boolean(s.userInfo?.forceBindPhone))
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
        phone: currentProfile.phone ?? '',
        phone_login_enabled: forceBindPhone
          ? true
          : wireBool(currentProfile.phone_login_enabled ?? false),
      })
    },
    [form, forceBindPhone],
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

  async function savePhone() {
    const values = await form.validateFields()
    const phone = (values.phone ?? '').trim()
    if (phone && !isValidPhone(phone)) {
      message.warning('请输入有效手机号')
      return
    }
    if (!phone && (values.phone_login_enabled || forceBindPhone)) {
      message.warning('请输入手机号')
      return
    }
    setPassword('')
    setOtpCode('')
    setConfirmOpen(true)
  }

  async function sendBindCode() {
    if (otpCooldown > 0 || sendingCode) return
    const phone = (form.getFieldValue('phone') ?? '').trim()
    if (!phone || !isValidPhone(phone)) {
      message.warning('请先填写有效手机号')
      return
    }
    setSendingCode(true)
    try {
      await authApi.sendBindPhoneCode({ target: phone })
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
    const phone = (form.getFieldsValue().phone ?? '').trim()
    if (phone && !otpCode.trim()) {
      message.warning('请输入手机验证码')
      return
    }
    setConfirmLoading(true)
    setSaving(true)
    try {
      const encrypted = await encryptPasswords({ password })
      const values = form.getFieldsValue()
      await authApi.updatePhone({
        password: encrypted.values.password || '',
        password_key_id: encrypted.password_key_id,
        phone: phone || null,
        phone_login_enabled: forceBindPhone
          ? true
          : wireBool(values.phone_login_enabled ?? false),
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
        {forceBindPhone ? (
          <Typography.Paragraph type="warning">请先绑定手机号后才能继续使用系统。</Typography.Paragraph>
        ) : null}
        <Form
          form={form}
          layout="vertical"
          className="profile-form profile-form--narrow w-full min-w-0"
        >
          <Form.Item name="phone" label="手机号">
            <Input allowClear />
          </Form.Item>
          <Form.Item name="phone_login_enabled" label="启用手机号登录" valuePropName="checked">
            <Switch disabled={forceBindPhone} />
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
          <Form.Item label="手机验证码">
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
        <Typography.Text type="secondary">绑定或更换手机号需验证码与当前密码。</Typography.Text>
      </Modal>
    </>
  )
}
