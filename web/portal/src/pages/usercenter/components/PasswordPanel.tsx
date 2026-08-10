/** Author: Charlie */

import { useEffect, useState } from 'react'
import { Button, Form, Input, message } from 'antd'
import { authApi } from '@/api'
import { encryptPasswords } from '@/utils/security'
import { PasswordStrength } from '@/components/common/PasswordStrength'
import '../usercenter.css'

type VerifyMethod = 'OLD_PASSWORD' | 'EMAIL_CODE' | 'PHONE_CODE'

export function PasswordPanel() {
  const [form] = Form.useForm()
  const [saving, setSaving] = useState(false)
  const [sendingCode, setSendingCode] = useState(false)
  const [verifyMethod, setVerifyMethod] = useState<VerifyMethod>('OLD_PASSWORD')
  const newPassword = Form.useWatch('new_password', form) || ''

  useEffect(() => {
    void authApi
      .authOptions()
      .then((res) => {
        const method = String(
          res?.data?.password_change_verify_method || 'OLD_PASSWORD',
        ).toUpperCase()
        if (method === 'EMAIL_CODE' || method === 'PHONE_CODE' || method === 'OLD_PASSWORD') {
          setVerifyMethod(method)
        }
      })
      .catch(() => setVerifyMethod('OLD_PASSWORD'))
  }, [])

  async function sendCode() {
    setSendingCode(true)
    try {
      await authApi.sendPasswordChangeCode()
      message.success('验证码已发送')
    } finally {
      setSendingCode(false)
    }
  }

  async function savePassword() {
    const values = await form.validateFields()
    if (values.new_password !== values.confirm_password) {
      message.warning('两次输入的新密码不一致')
      return
    }
    setSaving(true)
    try {
      const payload: Record<string, string> = {
        new_password: values.new_password,
      }
      if (verifyMethod === 'OLD_PASSWORD') {
        payload.old_password = values.old_password
      }
      const encrypted = await encryptPasswords(payload)
      await authApi.updateUserCenterPassword({
        old_password: encrypted.values.old_password,
        new_password: encrypted.values.new_password || '',
        password_key_id: encrypted.password_key_id,
        otp_code: verifyMethod === 'OLD_PASSWORD' ? undefined : values.otp_code,
      })
      form.resetFields()
      message.success('密码已更新')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Form
      form={form}
      layout="vertical"
      className="user-center-form user-center-form--narrow w-full min-w-0"
    >
      {verifyMethod === 'OLD_PASSWORD' ? (
        <Form.Item
          name="old_password"
          label="旧密码"
          rules={[{ required: true, message: '请输入旧密码' }]}
        >
          <Input.Password />
        </Form.Item>
      ) : (
        <Form.Item
          name="otp_code"
          label={verifyMethod === 'EMAIL_CODE' ? '邮箱验证码' : '手机验证码'}
          rules={[{ required: true, message: '请输入验证码' }]}
        >
          <Input
            placeholder={
              verifyMethod === 'EMAIL_CODE'
                ? '将向已绑定邮箱发送验证码'
                : '将向已绑定手机号发送验证码'
            }
            addonAfter={
              <Button type="link" loading={sendingCode} onClick={() => void sendCode()}>
                发送验证码
              </Button>
            }
          />
        </Form.Item>
      )}
      <Form.Item
        name="new_password"
        label="新密码"
        rules={[{ required: true, message: '请输入新密码' }]}
      >
        <Input.Password placeholder="请输入新密码" />
      </Form.Item>
      <PasswordStrength password={newPassword} />
      <Form.Item
        name="confirm_password"
        label="确认密码"
        dependencies={['new_password']}
        rules={[
          { required: true, message: '请确认新密码' },
          ({ getFieldValue }) => ({
            validator(_, value) {
              if (!value || getFieldValue('new_password') === value) {
                return Promise.resolve()
              }
              return Promise.reject(new Error('两次输入的新密码不一致'))
            },
          }),
        ]}
      >
        <Input.Password />
      </Form.Item>
      <Form.Item>
        <Button type="primary" loading={saving} onClick={() => void savePassword()}>
          更新密码
        </Button>
      </Form.Item>
    </Form>
  )
}
