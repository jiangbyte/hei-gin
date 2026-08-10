/** Author: Charlie */

import { useEffect, useMemo, useRef, useState } from 'react'
import {
  Button,
  Checkbox,
  ConfigProvider,
  Form,
  Input,
  Modal,
  Radio,
  Result,
  Tabs,
  message,
} from 'antd'
import { CloseOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '@/api'
import { CaptchaInput, type CaptchaInputHandle } from '@/components/common/CaptchaInput'
import { PasswordStrength } from '@/components/common/PasswordStrength'
import { useAuthModalStore } from '@/stores/authModal'
import { useAuthStore } from '@/stores/auth'
import { encryptPasswords } from '@/utils/security'
import { isValidEmail, isValidPhone } from '@/utils/validate'
import { wireBool } from '@/utils/wire'
import '@/pages/auth/auth-page.css'
import './auth-modal.css'

const brandName = import.meta.env.VITE_APP_TITLE || 'HEI'

const OTP_COOLDOWN_SECONDS = 60

type LoginType = 'ACCOUNT' | 'EMAIL' | 'PHONE'

type LoginFormValues = {
  account?: string
  email?: string
  phone?: string
  password?: string
  otp_code?: string
  captcha_id: string
  captcha_value: string
  remember: boolean
}

type RegisterFormValues = {
  account: string
  email?: string
  phone?: string
  password: string
  confirmPassword: string
  captcha_id: string
  captcha_value: string
}

const allTabItems = [
  { key: 'ACCOUNT', label: '账号', placeholder: '请输入账号' },
  { key: 'EMAIL', label: '邮箱', placeholder: '请输入登录邮箱' },
  { key: 'PHONE', label: '手机号', placeholder: '请输入登录手机号' },
]

export function AuthModal() {
  const mode = useAuthModalStore((s) => s.mode)
  const redirect = useAuthModalStore((s) => s.redirect)
  const close = useAuthModalStore((s) => s.close)
  const switchMode = useAuthModalStore((s) => s.switchMode)
  const navigate = useNavigate()
  const login = useAuthStore((s) => s.login)

  const [loginForm] = Form.useForm<LoginFormValues>()
  const [registerForm] = Form.useForm<RegisterFormValues>()
  const [activeType, setActiveType] = useState<LoginType>('ACCOUNT')
  const [loginMode, setLoginMode] = useState<'PASSWORD' | 'OTP'>('PASSWORD')
  const [loading, setLoading] = useState(false)
  const [sendingCode, setSendingCode] = useState(false)
  const [otpCooldown, setOtpCooldown] = useState(0)
  const [registerEnabled, setRegisterEnabled] = useState(true)
  const [requireEmail, setRequireEmail] = useState(true)
  const [requirePhone, setRequirePhone] = useState(false)
  const [options, setOptions] = useState({
    allow_account: true,
    allow_email: true,
    allow_phone: true,
    allow_otp: true,
  })
  const loginCaptchaRef = useRef<CaptchaInputHandle>(null)
  const registerCaptchaRef = useRef<CaptchaInputHandle>(null)
  const loginCaptchaId = Form.useWatch('captcha_id', loginForm) || ''
  const loginCaptchaValue = Form.useWatch('captcha_value', loginForm) || ''
  const registerCaptchaId = Form.useWatch('captcha_id', registerForm) || ''
  const registerCaptchaValue = Form.useWatch('captcha_value', registerForm) || ''
  const password = Form.useWatch('password', registerForm) || ''

  const open = mode !== null
  const isLogin = mode === 'login'

  const tabItems = useMemo(
    () =>
      allTabItems.filter((item) => {
        if (item.key === 'ACCOUNT') return options.allow_account
        if (item.key === 'EMAIL') return options.allow_email
        return options.allow_phone
      }),
    [options],
  )

  const resolvedActiveType = tabItems.some((item) => item.key === activeType)
    ? activeType
    : (tabItems[0]?.key as LoginType) || 'ACCOUNT'

  const otpAvailable =
    options.allow_otp && (resolvedActiveType === 'EMAIL' || resolvedActiveType === 'PHONE')

  const resolvedLoginMode = otpAvailable ? loginMode : 'PASSWORD'

  useEffect(() => {
    if (!open) return
    void authApi
      .authOptions()
      .then((res) => {
        const data = res?.data || {}
        setOptions({
          allow_account: wireBool(data.allow_account ?? true),
          allow_email: wireBool(data.allow_email ?? true),
          allow_phone: wireBool(data.allow_phone ?? true),
          allow_otp: wireBool(data.allow_otp ?? true),
        })
        setRegisterEnabled(wireBool(data.register_enabled ?? false))
        setRequireEmail(wireBool(data.register_require_email ?? false))
        setRequirePhone(wireBool(data.register_require_phone ?? false))
      })
      .catch(() => undefined)
  }, [open])

  useEffect(() => {
    if (otpCooldown <= 0) return
    const timer = window.setTimeout(() => setOtpCooldown((v) => v - 1), 1000)
    return () => window.clearTimeout(timer)
  }, [otpCooldown])

  function resetForms() {
    loginForm.resetFields()
    registerForm.resetFields()
    setLoading(false)
    setOtpCooldown(0)
    setLoginMode('PASSWORD')
    setActiveType('ACCOUNT')
  }

  async function onSendCode() {
    if (otpCooldown > 0 || sendingCode) return
    const identity =
      resolvedActiveType === 'ACCOUNT'
        ? loginForm.getFieldValue('account')?.trim()
        : resolvedActiveType === 'EMAIL'
          ? loginForm.getFieldValue('email')?.trim()
          : loginForm.getFieldValue('phone')?.trim()
    if (!identity) {
      message.warning(`请输入${tabItems.find((t) => t.key === resolvedActiveType)?.label}`)
      return
    }
    if (resolvedActiveType === 'EMAIL' && !isValidEmail(identity)) {
      message.warning('请输入有效邮箱')
      return
    }
    if (resolvedActiveType === 'PHONE' && !isValidPhone(identity)) {
      message.warning('请输入有效手机号')
      return
    }
    if (!loginCaptchaValue.trim()) {
      message.warning('请输入图形验证码')
      return
    }
    setSendingCode(true)
    try {
      await authApi.sendLoginCode({
        target: identity,
        channel: resolvedActiveType === 'EMAIL' ? 'EMAIL' : 'PHONE',
        captcha_id: loginCaptchaId,
        captcha_value: loginCaptchaValue,
      })
      message.success('验证码已发送，请查收后填写')
      setOtpCooldown(OTP_COOLDOWN_SECONDS)
      await loginCaptchaRef.current?.refresh()
    } catch {
      await loginCaptchaRef.current?.refresh()
    } finally {
      setSendingCode(false)
    }
  }

  async function onLoginFinish(values: LoginFormValues) {
    const identity =
      resolvedActiveType === 'ACCOUNT'
        ? values.account?.trim()
        : resolvedActiveType === 'EMAIL'
          ? values.email?.trim()
          : values.phone?.trim()

    if (!identity) {
      message.warning(`请输入${tabItems.find((t) => t.key === resolvedActiveType)?.label}`)
      return
    }
    if (resolvedActiveType === 'EMAIL' && !isValidEmail(identity)) {
      message.warning('请输入有效邮箱')
      return
    }
    if (resolvedActiveType === 'PHONE' && !isValidPhone(identity)) {
      message.warning('请输入有效手机号')
      return
    }

    setLoading(true)
    try {
      let password = ''
      let passwordKeyId: string | undefined
      if (resolvedLoginMode === 'PASSWORD') {
        const encrypted = await encryptPasswords({ password: values.password || '' })
        password = encrypted.values.password || ''
        passwordKeyId = encrypted.password_key_id
      }
      const next = await login(identity, password, redirect, values.remember, resolvedActiveType, {
        password_key_id: passwordKeyId,
        captcha_id: values.captcha_id,
        captcha_value: values.captcha_value,
        login_mode: resolvedLoginMode,
        ...(resolvedLoginMode === 'OTP' && values.otp_code?.trim()
          ? { otp_code: values.otp_code.trim() }
          : {}),
      })
      message.success('登录成功')
      close()
      navigate(next)
    } catch {
      await loginCaptchaRef.current?.refresh()
    } finally {
      setLoading(false)
    }
  }

  async function onRegisterFinish(values: RegisterFormValues) {
    const account = values.account.trim()
    const email = values.email?.trim() || ''
    const phone = values.phone?.trim() || ''

    if (account.length < 3 || account.length > 64) {
      message.warning('用户名需 3-64 个字符')
      return
    }
    if (requireEmail || email) {
      if (!isValidEmail(email) || email.length > 128) {
        message.warning('邮箱格式不正确')
        return
      }
    }
    if (requirePhone || phone) {
      if (!isValidPhone(phone)) {
        message.warning('手机号格式不正确')
        return
      }
    }

    setLoading(true)
    try {
      const encrypted = await encryptPasswords({ password: values.password })
      await authApi.register({
        account,
        email: email || undefined,
        phone: phone || undefined,
        password: encrypted.values.password || '',
        password_key_id: encrypted.password_key_id,
        captcha_id: values.captcha_id,
        captcha_value: values.captcha_value,
      })
      message.success('注册成功，请登录')
      registerForm.resetFields()
      switchMode('login')
    } catch {
      await registerCaptchaRef.current?.refresh()
    } finally {
      setLoading(false)
    }
  }

  const activeField = resolvedActiveType.toLowerCase() as 'account' | 'email' | 'phone'
  const title = isLogin ? '欢迎登录' : '注册账号'
  const headerExtra = isLogin ? (
    registerEnabled ? (
      <button type="button" className="linkish" onClick={() => switchMode('register')}>
        没有账号？去注册
      </button>
    ) : null
  ) : (
    <button type="button" className="linkish" onClick={() => switchMode('login')}>
      已有账号？去登录
    </button>
  )

  return (
    <Modal
      open={open}
      onCancel={close}
      afterOpenChange={(next) => {
        if (!next) resetForms()
      }}
      footer={null}
      width={880}
      centered
      destroyOnHidden
      maskClosable
      keyboard
      className="auth-modal"
      closeIcon={<CloseOutlined />}
      styles={{
        container: { padding: 0, overflow: 'hidden', background: 'transparent' },
        mask: { background: 'rgba(15, 23, 42, 0.45)' },
      }}
    >
      <div className="auth-modal__card auth-card">
        <aside className="auth-card__brand">
          <div className="auth-card__brand-deco" aria-hidden />
          <div className="auth-card__brand-inner">
            <Link to="/" className="auth-card__logo" onClick={close}>
              <span className="auth-card__logo-mark">{brandName.slice(0, 1).toUpperCase()}</span>
              <span className="auth-card__logo-text">{brandName}</span>
            </Link>
            <p className="auth-card__eyebrow">Portal</p>
            <h2 className="auth-card__headline">
              {isLogin ? '登录门户畅享更多服务' : '加入门户开启更多能力'}
            </h2>
            <p className="auth-card__lead">登录注册、个人中心与公告，开箱即用。</p>
            <div className="auth-card__brand-foot">
              <Link to="/" className="auth-card__brand-link" onClick={close}>
                进入门户首页
              </Link>
            </div>
          </div>
        </aside>

        <div className="auth-card__form">
          <div className="auth-card__form-head">
            <h1 className="auth-card__title">{title}</h1>
            {headerExtra ? <div className="auth-card__form-extra">{headerExtra}</div> : null}
          </div>

          <div className="auth-card__form-body">
            <ConfigProvider componentSize="large">
              {isLogin ? (
                <Form
                  form={loginForm}
                  layout="vertical"
                  requiredMark={false}
                  initialValues={{ remember: true, captcha_id: '', captcha_value: '' }}
                  onFinish={(v) => void onLoginFinish(v)}
                >
                  <Tabs
                    activeKey={resolvedActiveType}
                    items={tabItems.map((item) => ({ key: item.key, label: item.label }))}
                    onChange={(key) => setActiveType(key as LoginType)}
                  />

                  <Form.Item
                    name={activeField}
                    rules={[{ required: true, message: '请填写登录身份' }]}
                  >
                    <Input
                      placeholder={tabItems.find((t) => t.key === resolvedActiveType)?.placeholder}
                      allowClear
                    />
                  </Form.Item>

                  {otpAvailable ? (
                    <Form.Item>
                      <Radio.Group
                        value={resolvedLoginMode}
                        optionType="button"
                        options={[
                          { label: '密码登录', value: 'PASSWORD' },
                          { label: '验证码登录', value: 'OTP' },
                        ]}
                        onChange={(e) => setLoginMode(e.target.value)}
                      />
                    </Form.Item>
                  ) : null}

                  {resolvedLoginMode === 'PASSWORD' ? (
                    <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                      <Input.Password placeholder="请输入密码" />
                    </Form.Item>
                  ) : (
                    <Form.Item
                      name="otp_code"
                      rules={[{ required: true, message: '请输入登录验证码' }]}
                    >
                      <Input
                        placeholder="请输入登录验证码"
                        addonAfter={
                          <Button
                            type="link"
                            loading={sendingCode}
                            disabled={otpCooldown > 0}
                            onClick={() => void onSendCode()}
                          >
                            {otpCooldown > 0 ? `${otpCooldown}s 后重发` : '发送验证码'}
                          </Button>
                        }
                      />
                    </Form.Item>
                  )}

                  <Form.Item name="captcha_id" hidden>
                    <Input />
                  </Form.Item>

                  <Form.Item
                    name="captcha_value"
                    rules={[{ required: true, message: '请输入验证码' }]}
                  >
                    <CaptchaInput
                      ref={loginCaptchaRef}
                      size="large"
                      captchaId={loginCaptchaId}
                      captchaValue={loginCaptchaValue}
                      onCaptchaIdChange={(v) => loginForm.setFieldValue('captcha_id', v)}
                      onCaptchaValueChange={(v) => loginForm.setFieldValue('captcha_value', v)}
                    />
                  </Form.Item>

                  <Form.Item>
                    <div className="flex items-center justify-between">
                      <Form.Item name="remember" valuePropName="checked" noStyle>
                        <Checkbox>记住我</Checkbox>
                      </Form.Item>
                      <Link to="/auth/forgot-password" onClick={close}>
                        忘记密码？
                      </Link>
                    </div>
                  </Form.Item>

                  <Form.Item>
                    <Button type="primary" htmlType="submit" block loading={loading}>
                      登录
                    </Button>
                  </Form.Item>
                </Form>
              ) : !registerEnabled ? (
                <Result
                  status="info"
                  title="暂未开放注册"
                  extra={
                    <Button type="primary" onClick={() => switchMode('login')}>
                      返回登录
                    </Button>
                  }
                />
              ) : (
                <Form
                  form={registerForm}
                  layout="vertical"
                  requiredMark={false}
                  initialValues={{ captcha_id: '', captcha_value: '' }}
                  onFinish={(v) => void onRegisterFinish(v)}
                >
                  <Form.Item
                    name="account"
                    rules={[
                      { required: true, message: '请输入用户名' },
                      { min: 3, max: 64, message: '用户名需 3-64 个字符' },
                    ]}
                  >
                    <Input placeholder="用户名" allowClear />
                  </Form.Item>

                  <Form.Item
                    name="email"
                    rules={[
                      { required: requireEmail, message: '请输入邮箱' },
                      { type: 'email', message: '邮箱格式不正确' },
                      { max: 128, message: '邮箱最多 128 个字符' },
                    ]}
                  >
                    <Input placeholder="邮箱地址" allowClear />
                  </Form.Item>

                  <Form.Item
                    name="phone"
                    rules={[{ required: requirePhone, message: '请输入手机号' }]}
                  >
                    <Input placeholder="手机号" allowClear />
                  </Form.Item>

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

                  <Form.Item
                    name="captcha_value"
                    rules={[{ required: true, message: '请输入验证码' }]}
                  >
                    <CaptchaInput
                      ref={registerCaptchaRef}
                      size="large"
                      captchaId={registerCaptchaId}
                      captchaValue={registerCaptchaValue}
                      onCaptchaIdChange={(v) => registerForm.setFieldValue('captcha_id', v)}
                      onCaptchaValueChange={(v) => registerForm.setFieldValue('captcha_value', v)}
                    />
                  </Form.Item>

                  <Form.Item name="captcha_id" hidden>
                    <Input />
                  </Form.Item>

                  <Form.Item>
                    <Button type="primary" htmlType="submit" block loading={loading}>
                      立即注册
                    </Button>
                  </Form.Item>
                </Form>
              )}
            </ConfigProvider>
          </div>

          <div className="auth-modal__footer">注册登录即表示同意相关服务条款与隐私政策</div>
        </div>
      </div>
    </Modal>
  )
}
