/** Author: Charlie */

import { useEffect, useMemo, useRef, useState } from 'react'
import { Button, ConfigProvider, Form, Input, Result, Space, Tabs, message } from 'antd'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '@/api'
import { CaptchaInput, type CaptchaInputHandle } from '@/components/common/CaptchaInput'
import { PasswordStrength } from '@/components/common/PasswordStrength'
import { encryptPasswords } from '@/utils/security'
import { isValidEmail, isValidPhone } from '@/utils/validate'
import { wireBool } from '@/utils/wire'

const OTP_COOLDOWN_SECONDS = 60

type LoginType = 'ACCOUNT' | 'EMAIL' | 'PHONE'

type RegisterFormValues = {
  account?: string
  email?: string
  phone?: string
  otp_code?: string
  password: string
  confirmPassword: string
  captcha_id: string
  captcha_value: string
}

const registerTabItems = [
  { key: 'ACCOUNT', label: '用户名', placeholder: '请输入用户名' },
  { key: 'EMAIL', label: '邮箱', placeholder: '请输入注册邮箱' },
  { key: 'PHONE', label: '手机号', placeholder: '请输入注册手机号' },
]

export function AuthRegisterForm() {
  const navigate = useNavigate()
  const [form] = Form.useForm<RegisterFormValues>()
  const [registerChannel, setRegisterChannel] = useState<LoginType>('ACCOUNT')
  const [loading, setLoading] = useState(false)
  const [sendingCode, setSendingCode] = useState(false)
  const [otpCooldown, setOtpCooldown] = useState(0)
  const [registerEnabled, setRegisterEnabled] = useState(true)
  const [options, setOptions] = useState({
    register_allow_account: true,
    register_allow_email: true,
    register_allow_phone: false,
    register_require_email: false,
    register_require_phone: false,
  })
  const captchaRef = useRef<CaptchaInputHandle>(null)
  const captchaId = Form.useWatch('captcha_id', form) || ''
  const captchaValue = Form.useWatch('captcha_value', form) || ''
  const password = Form.useWatch('password', form) || ''

  const registerTabs = useMemo(
    () =>
      registerTabItems.filter((item) => {
        if (item.key === 'ACCOUNT') return options.register_allow_account
        if (item.key === 'EMAIL') return options.register_allow_email
        return options.register_allow_phone
      }),
    [options],
  )

  const resolvedChannel = registerTabs.some((item) => item.key === registerChannel)
    ? registerChannel
    : (registerTabs[0]?.key as LoginType) || 'ACCOUNT'

  useEffect(() => {
    void authApi
      .authOptions()
      .then((res) => {
        const data = res?.data || {}
        setRegisterEnabled(wireBool(data.register_enabled ?? false))
        setOptions({
          register_allow_account: wireBool(data.register_allow_account ?? true),
          register_allow_email: wireBool(data.register_allow_email ?? true),
          register_allow_phone: wireBool(data.register_allow_phone ?? false),
          register_require_email: wireBool(data.register_require_email ?? false),
          register_require_phone: wireBool(data.register_require_phone ?? false),
        })
      })
      .catch(() => undefined)
  }, [])

  useEffect(() => {
    if (otpCooldown <= 0) return
    const timer = window.setTimeout(() => setOtpCooldown((v) => v - 1), 1000)
    return () => window.clearTimeout(timer)
  }, [otpCooldown])

  async function onSendCode() {
    if (otpCooldown > 0 || sendingCode) return
    if (resolvedChannel !== 'EMAIL' && resolvedChannel !== 'PHONE') return
    const identity =
      resolvedChannel === 'EMAIL'
        ? form.getFieldValue('email')?.trim()
        : form.getFieldValue('phone')?.trim()
    if (!identity) {
      message.warning(`请输入${registerTabs.find((t) => t.key === resolvedChannel)?.label}`)
      return
    }
    if (resolvedChannel === 'EMAIL' && !isValidEmail(identity)) {
      message.warning('请输入有效邮箱')
      return
    }
    if (resolvedChannel === 'PHONE' && !isValidPhone(identity)) {
      message.warning('请输入有效手机号')
      return
    }
    if (!captchaValue.trim()) {
      message.warning('请输入图形验证码')
      return
    }
    setSendingCode(true)
    try {
      await authApi.sendRegisterCode({
        target: identity,
        channel: resolvedChannel,
        captcha_id: captchaId,
        captcha_value: captchaValue,
      })
      message.success('验证码已发送，请查收后填写')
      setOtpCooldown(OTP_COOLDOWN_SECONDS)
      await captchaRef.current?.refresh()
    } catch {
      await captchaRef.current?.refresh()
    } finally {
      setSendingCode(false)
    }
  }

  async function onFinish(values: RegisterFormValues) {
    const encryptedPayload: Record<string, unknown> = {
      register_channel: resolvedChannel,
      captcha_id: values.captcha_id,
      captcha_value: values.captcha_value,
    }
    if (resolvedChannel === 'ACCOUNT') {
      const account = (values.account || '').trim()
      if (account.length < 3 || account.length > 64) {
        message.warning('用户名需 3-64 个字符')
        return
      }
      encryptedPayload.account = account
      if (options.register_require_email) {
        const email = (values.email || '').trim()
        if (!isValidEmail(email) || email.length > 128) {
          message.warning('邮箱格式不正确')
          return
        }
        encryptedPayload.email = email
      }
      if (options.register_require_phone) {
        const phone = (values.phone || '').trim()
        if (!isValidPhone(phone)) {
          message.warning('请输入有效手机号')
          return
        }
        encryptedPayload.phone = phone
      }
    } else if (resolvedChannel === 'EMAIL') {
      const email = (values.email || '').trim()
      if (!isValidEmail(email) || email.length > 128) {
        message.warning('邮箱格式不正确')
        return
      }
      if (!values.otp_code?.trim()) {
        message.warning('请输入邮箱验证码')
        return
      }
      encryptedPayload.email = email
      encryptedPayload.otp_code = values.otp_code.trim()
    } else {
      const phone = (values.phone || '').trim()
      if (!isValidPhone(phone)) {
        message.warning('请输入有效手机号')
        return
      }
      if (!values.otp_code?.trim()) {
        message.warning('请输入手机验证码')
        return
      }
      encryptedPayload.phone = phone
      encryptedPayload.otp_code = values.otp_code.trim()
    }

    setLoading(true)
    try {
      const encrypted = await encryptPasswords({ password: values.password })
      await authApi.register({
        ...encryptedPayload,
        password: encrypted.values.password || '',
        password_key_id: encrypted.password_key_id,
      })
      message.success('注册成功，请登录')
      navigate('/auth/login')
    } catch {
      await captchaRef.current?.refresh()
    } finally {
      setLoading(false)
    }
  }

  if (!registerEnabled || registerTabs.length === 0) {
    return (
      <Result
        status="info"
        title="暂未开放注册"
        extra={
          <Button type="primary" onClick={() => navigate('/auth/login')}>
            返回登录
          </Button>
        }
      />
    )
  }

  return (
    <ConfigProvider componentSize="large">
      <Form
        form={form}
        layout="vertical"
        requiredMark={false}
        initialValues={{ captcha_id: '', captcha_value: '' }}
        onFinish={(v) => void onFinish(v)}
      >
        <Tabs
          activeKey={resolvedChannel}
          items={registerTabs.map((item) => ({ key: item.key, label: item.label }))}
          onChange={(key) => setRegisterChannel(key as LoginType)}
        />

        {resolvedChannel === 'ACCOUNT' ? (
          <>
            <Form.Item
              name="account"
              rules={[
                { required: true, message: '请输入用户名' },
                { min: 3, max: 64, message: '用户名需 3-64 个字符' },
              ]}
            >
              <Input placeholder="用户名" allowClear />
            </Form.Item>
            {options.register_require_email ? (
              <Form.Item
                name="email"
                rules={[
                  { required: true, message: '请输入邮箱' },
                  { type: 'email', message: '邮箱格式不正确' },
                  { max: 128, message: '邮箱最多 128 个字符' },
                ]}
              >
                <Input placeholder="邮箱（必填）" allowClear />
              </Form.Item>
            ) : null}
            {options.register_require_phone ? (
              <Form.Item
                name="phone"
                rules={[{ required: true, message: '请输入手机号' }]}
              >
                <Input placeholder="手机号（必填）" allowClear />
              </Form.Item>
            ) : null}
          </>
        ) : null}

        {resolvedChannel === 'EMAIL' ? (
          <>
            <Form.Item
              name="email"
              rules={[
                { required: true, message: '请输入邮箱' },
                { type: 'email', message: '邮箱格式不正确' },
                { max: 128, message: '邮箱最多 128 个字符' },
              ]}
            >
              <Input placeholder="邮箱" allowClear />
            </Form.Item>
            <Form.Item>
              <Space.Compact block>
                <Form.Item
                  name="otp_code"
                  noStyle
                  rules={[{ required: true, message: '请输入邮箱验证码' }]}
                >
                  <Input placeholder="邮箱验证码" />
                </Form.Item>
                <Button
                  loading={sendingCode}
                  disabled={otpCooldown > 0}
                  onClick={() => void onSendCode()}
                >
                  {otpCooldown > 0 ? `${otpCooldown}s 后重发` : '发送验证码'}
                </Button>
              </Space.Compact>
            </Form.Item>
          </>
        ) : null}

        {resolvedChannel === 'PHONE' ? (
          <>
            <Form.Item name="phone" rules={[{ required: true, message: '请输入手机号' }]}>
              <Input placeholder="手机号" allowClear />
            </Form.Item>
            <Form.Item>
              <Space.Compact block>
                <Form.Item
                  name="otp_code"
                  noStyle
                  rules={[{ required: true, message: '请输入手机验证码' }]}
                >
                  <Input placeholder="手机验证码" />
                </Form.Item>
                <Button
                  loading={sendingCode}
                  disabled={otpCooldown > 0}
                  onClick={() => void onSendCode()}
                >
                  {otpCooldown > 0 ? `${otpCooldown}s 后重发` : '发送验证码'}
                </Button>
              </Space.Compact>
            </Form.Item>
          </>
        ) : null}

        <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
          <Input.Password placeholder="密码" />
        </Form.Item>
        <PasswordStrength password={password} />

        <Form.Item
          name="confirmPassword"
          dependencies={['password']}
          rules={[
            { required: true, message: '请确认密码' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve()
                }
                return Promise.reject(new Error('两次密码输入不一致'))
              },
            }),
          ]}
        >
          <Input.Password placeholder="确认密码" />
        </Form.Item>

        <Form.Item name="captcha_value" rules={[{ required: true, message: '请输入验证码' }]}>
          <CaptchaInput
            ref={captchaRef}
            size="large"
            onCaptchaIdChange={(v) => form.setFieldValue('captcha_id', v)}
          />
        </Form.Item>

        <Form.Item name="captcha_id" hidden>
          <Input />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading} className="auth-submit">
            立即注册
          </Button>
        </Form.Item>

        <p className="portal-auth__switch">
          已有账号？<Link to="/auth/login">去登录</Link>
        </p>
      </Form>
    </ConfigProvider>
  )
}
