/** Author: Charlie */

import { useEffect, useMemo, useRef, useState } from 'react'
import { Button, Checkbox, ConfigProvider, Form, Input, Space, Tabs, message } from 'antd'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { authApi } from '@/api'
import { CaptchaInput, type CaptchaInputHandle } from '@/components/common/CaptchaInput'
import { useAuthStore } from '@/stores/auth'
import { encryptPasswords } from '@/utils/security'
import { isValidEmail, isValidPhone } from '@/utils/validate'
import { wireBool } from '@/utils/wire'
import { OauthProviderButtons } from './OauthProviderButtons'

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

const allTabItems = [
  { key: 'ACCOUNT', label: '账号', placeholder: '请输入账号' },
  { key: 'EMAIL', label: '邮箱', placeholder: '请输入登录邮箱' },
  { key: 'PHONE', label: '手机号', placeholder: '请输入登录手机号' },
]

export function AuthLoginForm() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const redirect = params.get('redirect')
  const login = useAuthStore((s) => s.login)

  const [form] = Form.useForm<LoginFormValues>()
  const [activeType, setActiveType] = useState<LoginType>('ACCOUNT')
  const [loginMode, setLoginMode] = useState<'PASSWORD' | 'OTP'>('PASSWORD')
  const [loading, setLoading] = useState(false)
  const [sendingCode, setSendingCode] = useState(false)
  const [otpCooldown, setOtpCooldown] = useState(0)
  const [options, setOptions] = useState({
    allow_account: true,
    allow_email: true,
    allow_phone: true,
    allow_otp: true,
  })
  const [oauthProviders, setOauthProviders] = useState<
    Array<{ provider: string; label: string }>
  >([])
  const [oauthLoading, setOauthLoading] = useState<string | null>(null)
  const captchaRef = useRef<CaptchaInputHandle>(null)
  const captchaId = Form.useWatch('captcha_id', form) || ''
  const captchaValue = Form.useWatch('captcha_value', form) || ''

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
  const activeField = resolvedActiveType.toLowerCase() as 'account' | 'email' | 'phone'

  useEffect(() => {
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
        const providers = Array.isArray(data.oauth_providers) ? data.oauth_providers : []
        setOauthProviders(
          providers
            .map((item: any) => ({
              provider: String(item.provider || ''),
              label: String(item.label || item.provider || ''),
              enabled: wireBool(item.enabled ?? false),
              web_oauth: wireBool(item.web_oauth ?? true),
            }))
            .filter((item: { provider: string; enabled: boolean; web_oauth: boolean }) =>
              Boolean(item.provider && item.enabled && item.web_oauth),
            )
            .map((item: { provider: string; label: string }) => ({
              provider: item.provider,
              label: item.label,
            })),
        )
      })
      .catch(() => undefined)
  }, [])

  useEffect(() => {
    if (otpCooldown <= 0) return
    const timer = window.setTimeout(() => setOtpCooldown((v) => v - 1), 1000)
    return () => window.clearTimeout(timer)
  }, [otpCooldown])

  async function onOauthLogin(provider: string) {
    if (oauthLoading) return
    setOauthLoading(provider)
    try {
      const res = await authApi.oauthAuthorize(provider, {
        intent: 'LOGIN',
        redirect: redirect || undefined,
      })
      const url = res?.data?.authorize_url
      if (!url) {
        message.error('无法发起三方登录')
        return
      }
      window.location.href = String(url)
    } catch {
      // 全局错误提示
    } finally {
      setOauthLoading(null)
    }
  }

  async function onSendCode() {
    if (otpCooldown > 0 || sendingCode) return
    const identity = form.getFieldValue(activeField)?.trim()
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
    if (!captchaValue.trim()) {
      message.warning('请输入图形验证码')
      return
    }
    setSendingCode(true)
    try {
      await authApi.sendLoginCode({
        target: identity,
        channel: resolvedActiveType === 'EMAIL' ? 'EMAIL' : 'PHONE',
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

  async function onFinish(values: LoginFormValues) {
    const identity = values[activeField]?.trim()
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
      const next = await login(identity, password, redirect || undefined, values.remember, resolvedActiveType, {
        password_key_id: passwordKeyId,
        captcha_id: values.captcha_id,
        captcha_value: values.captcha_value,
        login_mode: resolvedLoginMode,
        ...(resolvedLoginMode === 'OTP' && values.otp_code?.trim()
          ? { otp_code: values.otp_code.trim() }
          : {}),
      })
      message.success('登录成功')
      navigate(next)
    } catch {
      await captchaRef.current?.refresh()
    } finally {
      setLoading(false)
    }
  }

  return (
    <ConfigProvider componentSize="large">
      <Form
        form={form}
        layout="vertical"
        requiredMark={false}
        initialValues={{ remember: true, captcha_id: '', captcha_value: '' }}
        onFinish={(v) => void onFinish(v)}
      >
        <Tabs
          activeKey={resolvedActiveType}
          items={tabItems.map((item) => ({ key: item.key, label: item.label }))}
          onChange={(key) => setActiveType(key as LoginType)}
        />

        <Form.Item name={activeField} rules={[{ required: true, message: '请填写登录身份' }]}>
          <Input
            placeholder={tabItems.find((t) => t.key === resolvedActiveType)?.placeholder}
            allowClear
          />
        </Form.Item>

        {resolvedLoginMode === 'PASSWORD' ? (
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password placeholder="请输入密码" />
          </Form.Item>
        ) : (
          <Form.Item>
            <Space.Compact block>
              <Form.Item
                name="otp_code"
                noStyle
                rules={[{ required: true, message: '请输入登录验证码' }]}
              >
                <Input placeholder="请输入登录验证码" />
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
        )}

        <Form.Item name="captcha_id" hidden>
          <Input />
        </Form.Item>

        <Form.Item name="captcha_value" rules={[{ required: true, message: '请输入验证码' }]}>
          <CaptchaInput
            ref={captchaRef}
            size="large"
            onCaptchaIdChange={(v) => form.setFieldValue('captcha_id', v)}
          />
        </Form.Item>

        <Form.Item>
          <div className="auth-form-row">
            <Form.Item name="remember" valuePropName="checked" noStyle>
              <Checkbox>记住我</Checkbox>
            </Form.Item>
            <div className="auth-form-row__links">
              {otpAvailable ? (
                <>
                  <button
                    type="button"
                    className="auth-mode-link"
                    onClick={() =>
                      setLoginMode(resolvedLoginMode === 'PASSWORD' ? 'OTP' : 'PASSWORD')
                    }
                  >
                    {resolvedLoginMode === 'PASSWORD' ? '验证码登录' : '密码登录'}
                  </button>
                  <span className="auth-form-row__sep">·</span>
                </>
              ) : null}
              <Link to="/auth/forgot-password">忘记密码？</Link>
            </div>
          </div>
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading} className="auth-submit">
            登录
          </Button>
        </Form.Item>

        <OauthProviderButtons
          providers={oauthProviders}
          loadingProvider={oauthLoading}
          onSelect={(provider) => void onOauthLogin(provider)}
        />
      </Form>
    </ConfigProvider>
  )
}
